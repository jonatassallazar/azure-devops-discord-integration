# CLAUDE.md

Guidance for AI assistants working in this repository.

## What this project is

A small Go HTTP service that receives **Azure DevOps Service Hook** webhooks and forwards
them to one or more chat apps as formatted messages. It covers three event families: pull
requests, pipeline (build) completions, and release deployments.

Inbound (source of events) is Azure DevOps today; outbound delivery goes through a
vendor-neutral `notify.Sink` interface, currently implemented by **Discord** and
**Google Chat**. Either, both, or neither can be configured per event category — a category
with no webhook URL set simply isn't notified for that destination.

The service is stateless: no database, no queue, no auth on the inbound routes. Each request
is parsed, translated into a vendor-neutral `notify.Message`, and fanned out to whichever
sinks are configured for that event category.

Go module name is `azuredevops-notify`, matching the repo name. All internal imports use
the module path.

## Layout

```
cmd/server/
  main.go                    Entry point: load config -> build httpserver.Server -> Init() (blocking)
  main_test.go                 End-to-end route tests via httptest + gin ServeHTTP
  mocks_test.go                 Fixtures (fakePayload*) used only by main_test.go (test-only, not shipped in the binary)

internal/config/
  config.go                    Config struct (nested Azure/Discord/GoogleChat) + .env loading + os.Getenv mapping

internal/notify/
  message.go                   Message, Author, Field, Level (vendor-neutral notification model)
  sink.go                       Sink interface: Send(ctx, Message) error
  dispatcher.go                 Dispatcher: fans one Message out to every Sink configured for an event category
  discord/discord.go            Sink impl: Level -> Discord embed color, POSTs to a Discord incoming webhook
  googlechat/googlechat.go      Sink impl: Level -> header indicator, POSTs a Cards v2 payload to a Google Chat incoming webhook

internal/azuredevops/
  routes.go                    Route path constants
  avatar.go / avatar_test.go     avatarURL(): Azure identity -> Gravatar URL for notify.Author.IconURL
  models.go                     Structs mirroring Azure DevOps webhook JSON (AzureRequest, Resource, ...)
  payload.go                    AzureRequest -> notify.Message converters + label helpers
  errors.go                      Shared respondError() helper (DRYs the repeated bad-request response)
  pullrequest.go / *_test.go     PullRequestHandler: CreatedPR, ReviewedPR, StatusUpdatedPR, processVotes
  pipeline.go / *_test.go        PipelineHandler: PipelineStatusReport, processStatus, getRepository
  release.go / *_test.go         ReleaseHandler: ReleaseStatusReport, processReleaseStatus

internal/httpserver/
  server.go                     Server struct, Init() -> SetupRouter() + r.Run()
  router.go                      SetupRouter(): gin engine, handler + Dispatcher wiring, route table
  health.go                      RouteHealth constant + healthCheck handler (liveness/readiness probe)

Dockerfile              Multi-stage build (golang:1.22-alpine -> alpine:3.19), exposes 8080
docker-compose.yml      Single `app` service, host port 8080

.github/workflows/
  docker-publish.yml    test job (gofmt/vet/test) gating a build+push to ghcr.io (see "CI")
```

## Request flow

1. Azure DevOps POSTs a Service Hook payload to one of the routes.
2. The handler does `c.ShouldBindJSON(&req)` into a local `azuredevops.AzureRequest` — bound
   into a request-scoped local variable, not a shared struct field (handlers hold only their
   `*notify.Dispatcher`, so one handler instance safely serves concurrent requests).
3. A `process*()` method on the handler maps an Azure status/vote to a `(notify.Level, string)`
   pair.
4. If `title == ""` the handler replies `204 No Content` and dispatches **nothing** — this is
   the deliberate "uninteresting event" path (e.g. a PR update that is not
   completed/conflicts, or a review with only neutral votes).
5. Otherwise `AzureRequest.toPRMessage/toPipelineMessage/toReleaseMessage()` builds a
   vendor-neutral `notify.Message`, and `Dispatcher.Send()` fans it out concurrently to every
   configured sink (Discord and/or Google Chat) for that event category.
6. Success replies `200` echoing the parsed Azure payload; a bind/marshal/all-sinks-failed
   error replies `400` with `{"err": ...}` after `fmt.Println(err)`. A **partial** sink
   failure (e.g. Discord ok, Google Chat down) is logged by the Dispatcher but still replies
   `200` — a single outbound hiccup doesn't fail the inbound webhook.

`PipelineHandler.PipelineStatusReport` additionally calls out to the **Azure DevOps REST
API** (`GET {AZURE_ORGANIZATION}/{AZURE_PROJECT}/_apis/git/repositories/{id}?api-version=5.1`)
using Basic auth with a base64-encoded `":" + AZURE_PAT_TOKEN`, to resolve the triggering
repository name from `resource.triggerInfo["ci.triggerRepository"]`.

## Routes

Defined as constants in `internal/azuredevops/routes.go` — **always reference the
constants**, never hardcode path strings (tests depend on this). Route **paths are live
Azure DevOps Service Hook URLs already configured in production — do not rename them**.

| Constant | Path | Handler | Sinks used |
|---|---|---|---|
| `httpserver.RouteHealth` | `GET /health` | `httpserver.healthCheck` | none |
| `RouteCreatedPR` | `POST /pull-request/created` | `PullRequestHandler.CreatedPR` | `Discord.PRWebhookURL`, `GoogleChat.PRWebhookURL` |
| `RouteReviewedPR` | `POST /pull-request/review` | `PullRequestHandler.ReviewedPR` | same as above |
| `RouteStatusUpdatedPR` | `POST /pull-request/status` | `PullRequestHandler.StatusUpdatedPR` | same as above |
| `RoutePipeline` | `POST /pipeline/` | `PipelineHandler.PipelineStatusReport` | `Discord.PipelineWebhookURL`, `GoogleChat.PipelineWebhookURL` |
| `RouteRelease` | `POST /release/` | `ReleaseHandler.ReleaseStatusReport` | `Discord.ReleaseWebhookURL`, `GoogleChat.ReleaseWebhookURL` |

Note the trailing slash on `/pipeline/` and `/release/` — Azure hook URLs must match.

`GET /health` is not an Azure DevOps hook: it's the container liveness/readiness probe, so its
constant lives in `internal/httpserver/health.go` rather than in `internal/azuredevops/routes.go`.
It replies `200 {"status":"ok"}` and contacts nothing — the service is stateless and only talks
to sinks while handling a webhook, so being able to answer *is* the readiness condition.

## Status / vote mapping (the domain rules)

Pull request reviewer votes (`resource.reviewers[].vote`), in `processVotes` precedence order —
any rejection wins, then waiting, then approval:

| Vote | Meaning | `getVoteText()` label |
|---|---|---|
| `10` | Approved | `Aprovado` |
| `5` | Approved with suggestions | `Aprovado com Sugestões` |
| `0` | No vote | `Sem Voto` |
| `-5` | Waiting for author | `Aguardando o Autor` |
| `-10` | Rejected | `Rejeitado` |

PR status (`StatusUpdatedPR`): `completed` -> `notify.LevelCompleted` "Concluído", `conflicts`
-> `notify.LevelFailure` "com Conflito", anything else -> `204`.

Pipeline `resource.result` and release `resource.deployment.deploymentStatus` share the same
mapping: `succeeded` -> `LevelSuccess` "Concluída", `failed` -> `LevelFailure` "Falhada",
`stopped` -> `LevelWarning` "Interrompida", default -> `LevelUnmapped`
`[Status não mapeado: X]` (still notifies, on purpose).

`notify.Level` is vendor-neutral; each sink privately maps it to its own visual language —
Discord to an embed color int, Google Chat to a colored-circle indicator in the card header
(Cards v2 has no per-card background color). See `internal/notify/discord/discord.go` and
`internal/notify/googlechat/googlechat.go` for the exact mappings.

Merge status labels live in `getMergeStatusText()` (`succeeded`, `conflicts`, `queued`,
`rejectedByPolicy`, `failure`).

## User avatars

The `imageUrl` Azure DevOps puts in a webhook payload points at an authenticated endpoint
(`.../_apis/GraphProfile/MemberAvatars/<descriptor>`), and on Azure DevOps Server that host is
frequently internal-only. Discord and Google Chat fetch an author icon anonymously from their
own servers, so handing them the raw URL yields a sign-in page, not an image, and the chat
message shows a blank avatar.

`avatarURL()` in `internal/azuredevops/avatar.go` is the whole answer: it maps the identity's
`uniqueName` (their email/UPN) to `https://www.gravatar.com/avatar/{md5(lowercased address)}`,
a public URL the chat platforms can fetch themselves. **Azure's `imageUrl` is deliberately
unused** — serving or proxying images is not this service's job, so nothing here fetches,
caches or re-hosts an avatar.

- `?s=80&d=identicon`: an address with no Gravatar account still gets a stable per-user
  pattern instead of a blank circle. The size and default live in the two consts at the top of
  the file.
- `uniqueName` is parsed with `net/mail`. Anything that is not an address — an on-premises
  `DOMAIN\user`, a service account, an empty field — yields `""`, and the message goes out
  with no author icon at all. That is the intended fallback: no icon beats a broken one.
- Privacy note worth keeping in mind: this puts an MD5 of each user's email address into
  every notification, which Discord/Google Chat then resolve against gravatar.com. There is no
  opt-out switch today; add one before that becomes a problem for someone.

Sink-side, the Discord payload sends `avatar_url` (the previous `avatarUrl` spelling was
silently ignored by Discord), and `author`/`thumbnail`/`image`/`footer` are pointers omitted
when empty — an embed object carrying an empty `url` is a broken image for Discord to render
rather than an absent one, which matters now that `IconURL` is legitimately empty for
identities without an address.

## Conventions

- **Language split**: identifiers, comments and log output are English; user-facing chat
  strings (titles, field names, labels) are **Brazilian Portuguese**. Keep new copy in pt-BR
  to stay consistent.
- **Handlers** are structs holding only a `*notify.Dispatcher` (and, for
  `PipelineHandler`, a `config.AzureConfig`), instantiated once in `SetupRouter` and bound as
  method handlers. They hold **no per-request state** — the parsed request is always a local
  variable inside the handler method, never a struct field. Adding a feature means: a new
  struct/method in `internal/azuredevops/`, a route constant in `routes.go`, and wiring in
  `internal/httpserver/router.go`.
- **Business logic goes in a pure `process*()` method** taking the parsed `*AzureRequest` as
  a parameter and returning `(notify.Level, string)` — that is what the
  `internal/azuredevops/*_test.go` table-driven tests exercise. Handlers stay thin (bind, map,
  convert to `notify.Message`, dispatch, respond).
- **Message construction lives in `internal/azuredevops/payload.go`** as methods on
  `*AzureRequest` (`toPRMessage`/`toPipelineMessage`/`toReleaseMessage`). Do not build
  vendor-specific (Discord/Google Chat) JSON inside handlers — that belongs in the
  respective `internal/notify/<vendor>` package, keyed off the vendor-neutral
  `notify.Message`.
- **Adding a new outbound sink** (e.g. Slack): implement `notify.Sink` in a new
  `internal/notify/<vendor>` package (`Send(ctx, notify.Message) error`), add its webhook URL
  field(s) to `config.Config`, and wire it into `buildDispatcher()` in
  `internal/httpserver/router.go`. No changes needed to `internal/azuredevops` or to any
  other sink.
- **Adding a new inbound source** (e.g. GitHub): there is currently no `Source` interface —
  with only one real source (Azure DevOps) so far, one hasn't been introduced to avoid
  premature abstraction. Add a new `internal/<source>` package mirroring the shape of
  `internal/azuredevops` (its own models, handlers, `process*()` methods building
  `notify.Message`), with its own routes wired in `internal/httpserver/router.go`.
- **Levels** are the constants in `internal/notify/message.go` (`LevelPending`,
  `LevelSuccess`, `LevelFailure`, `LevelWarning`, `LevelCompleted`, `LevelUnmapped`) — never a
  raw vendor color. Each sink owns its own Level→visual mapping privately.
- **Error handling** for bind/marshal/dispatch failures is uniform and intentionally simple:
  `internal/azuredevops/errors.go`'s `respondError()` does `fmt.Println(err)` then
  `c.JSON(http.StatusBadRequest, gin.H{"err": err})`. This is a DRY helper for the repeated
  block, not a logging framework or middleware — match it rather than introducing either
  unless asked. `notify.Dispatcher` uses the stdlib `log` package (not a framework) to note a
  *partial* sink failure that doesn't fail the whole request.
- Azure model structs (`internal/azuredevops/models.go`) are a permissive superset:
  `Resource` carries PR, build and release fields together, so one `AzureRequest` type binds
  every hook. Add fields to the existing structs rather than creating parallel types.
- Formatting: standard `gofmt` (tabs). Run `gofmt -l .` before committing — it should report
  nothing (fixtures now live in `cmd/server/mocks_test.go`, a `_test.go` file, so they're
  gofmt-checked like everything else and excluded from the production binary).

## Configuration

`config.Config.LoadEnvironment()` loads dotenv files then reads `os.Getenv`, then validates.
Load order (later files do **not** override values already set — `godotenv.Load` is first-wins,
so the first file listed effectively has priority, and a real environment variable always beats
every file):

1. `.env.{APP_ENV}.local`
2. `.env.local` (skipped entirely when `APP_ENV=test`)
3. `.env.{APP_ENV}`
4. `.env`

`APP_ENV` defaults to `development`. **Dotenv files are optional**: they're a local-development
convenience, and a container deployment (Kubernetes ConfigMap/Secret, `docker run --env-file`,
...) has the variables in the process environment with no file on disk. Finding none is logged
(`config: no env file found, ...`) and is not an error. Note dotenv files are resolved relative
to the process's working directory, which for `go test` is the package directory
(`cmd/server/`), not the repo root.

What *is* fatal is a config the service can do nothing with: `validate()` returns an error when
**no** webhook URL is set at all (any single one is enough), since the process would otherwise
accept webhooks and deliver nowhere. `main` logs that error and exits.

| Variable | Purpose |
|---|---|
| `APP_ENV` | Selects the dotenv variant (`development` / `production` / `test`) |
| `GIN_MODE` | Passed to `gin.SetMode` when non-empty (`debug` / `release` / `test`) |
| `DISCORD_PR_URL` / `DISCORD_PIPELINE_URL` / `DISCORD_RELEASE_URL` | Discord webhooks, one per event category (optional — unset disables Discord for that category) |
| `GOOGLE_CHAT_PR_URL` / `GOOGLE_CHAT_PIPELINE_URL` / `GOOGLE_CHAT_RELEASE_URL` | Google Chat webhooks, one per event category (optional — unset disables Google Chat for that category) |
| `AZURE_ORGANIZATION` | Base URL used to build the Azure REST call (used as a URL prefix, not a bare org name) |
| `AZURE_PROJECT` | Project segment of the Azure REST call |
| `AZURE_PAT_TOKEN` | PAT used for Basic auth against the Azure REST API |

At least one webhook URL (Discord and/or Google Chat) must be set overall, and should be set
per event category you want notifications for. The listen address comes from gin's `r.Run()` default: `:8080`,
overridable with the `PORT` env var. All `.env*` files are gitignored — never commit one, and
never echo real webhook URLs or PAT values into logs, commits, or PR descriptions.

## Build, run, test

```bash
go build ./...                        # compile
go run ./cmd/server                    # run (needs a .env file present in the repo root — copy .env.example)
go vet ./...
gofmt -l .                             # should report nothing
go test ./...                          # full suite
go test -cover ./...
go test ./internal/azuredevops/...     # unit tests only — no env or network needed
docker compose up -d                   # container, host port 8080

# cmd/server route tests, selecting which sink(s) get a configured webhook URL:
TEST_SINKS=discord go test ./cmd/server/...     # Discord only
TEST_SINKS=googlechat go test ./cmd/server/...  # Google Chat only
TEST_SINKS=both go test ./cmd/server/...        # both (default when TEST_SINKS is unset)
```

A `Makefile` wraps the commands above (`make build`, `make test`, `make test-unit`,
`make test-discord`, `make test-googlechat`, `make test-e2e`, ...).

### Testing gotchas — read before running `go test ./...`

- `./internal/azuredevops` tests are pure and always pass offline.
- `./internal/config` tests `os.Chdir` into a `t.TempDir()` and clear every config variable, so
  they exercise dotenv-less (container-style) loading without seeing the repo's own `.env`.
- The **`cmd/server` package tests (`main_test.go`) are true end-to-end tests**, but they run
  fully offline: `TestMain` spins up a local `httptest.Server` stub (answers every POST with
  `200`, and every GET — the pipeline route's Azure REST call — with a JSON
  `azuredevops.AzureRepository` body) and writes a scratch `.env` file in `cmd/server/`
  pointing the sink(s) selected by `TEST_SINKS` (`discord` / `googlechat` / unset-or-`both`)
  and `AZURE_ORGANIZATION` at it, restoring whatever `.env` was already there (or removing the
  scratch one) once the suite finishes — see `writeTestEnv()` in `main_test.go`. No manual env
  setup is required to run `go test ./...`.
- `prepareRouter()` still discards the error from `LoadEnvironment()` with `r, _ :=` — a
  `TestMain` that fails to write the scratch `.env` calls `os.Exit(1)` before any test runs,
  rather than letting individual tests panic with a nil-pointer dereference.
- `TestConfigE2ETesting` sets `APP_ENV=test` and unsets it on return, so the remaining tests in
  the file resolve under `development` — `writeTestEnv()` writes a plain `.env` (the final,
  always-attempted fallback per `config.Config.LoadEnvironment`), which covers both cases.
- New fixtures for root tests go in `cmd/server/mocks_test.go` as `fakePayload*` vars.

## CI

`.github/workflows/docker-publish.yml` is the only workflow. It has two jobs:

- **`test`** — the correctness gate: `gofmt -l .` (fails if it reports anything), `go vet ./...`,
  `go test ./...`. The Go version comes from `go.mod` via `go-version-file`, so bumping the
  module's Go version is enough — don't also hardcode it in the workflow. Every test in this
  repo runs offline (see "Testing gotchas"), so the job needs no secrets or services.
- **`build`** — `needs: test`, so nothing reaches the registry unless the gate is green. Builds
  the `Dockerfile` with Buildx, pushes to `ghcr.io/<owner>/<repo>`, and signs the digest with
  cosign.

Triggers:

| Event | Pushes to GHCR? | Package tags |
|---|---|---|
| Push to `main` | yes | `main`, `sha-<short>`, `latest` |
| Push a `v*.*.*` git tag | yes | `1.2.3`, `1.2`, `latest` |
| Pull request against `main` | no (build only) | `pr-<n>` (computed, never published) |

Tags come from `docker/metadata-action`'s `tags:` input. That input **replaces** the action's
default tag set rather than adding to it, so the defaults are written out explicitly there —
when adding a tag rule, add a line, don't shorten the list. `latest` is applied by
`type=raw,value=latest,enable={{is_default_branch}}` for `main` and by the action's default
`latest=auto` flavor for semver tags. There is deliberately no `schedule` trigger: nightly
rebuilds of an unchanged commit only churned the `latest` digest.

## Known rough edges

Do not "fix" these silently as part of an unrelated change; flag them or fix them deliberately.

- `models.AzurePipeline` (now `azuredevops.AzurePipeline`) is only used by test fixtures; the
  live routes bind `AzureRequest`.
- `Resource.CreatedBy.ImageUrl` / `RequestedFor.ImageUrl` are still bound from the payload but
  no longer used for anything (see "User avatars"). They are kept because the model structs
  deliberately mirror the Azure JSON.
- Response bodies from sink HTTP POSTs *are* now read and closed, and a non-2xx response *is*
  now treated as an error (fixed deliberately as part of the multi-vendor rewrite — previously
  neither was true).
- Inbound routes are unauthenticated — anyone who can reach the port can post notifications to
  the configured Discord/Google Chat destinations.
- There is intentionally no `Source` interface for inbound events (see "Adding a new inbound
  source" above) — only add one when a second real source (e.g. GitHub) is actually being
  built, not speculatively.

## Git workflow

- Default branch is `main`; feature work has historically used `feature/<name>` branches
  merged via PR.
- Commit messages are short and imperative, sometimes with a `chore:`/`fix:` prefix.
- Never commit `.env*` files or real tokens/webhook URLs.
