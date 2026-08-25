# CLAUDE.md

Guidance for AI assistants working in this repository.

## What this project is

A small Go HTTP service that receives **Azure DevOps Service Hook** webhooks and forwards
them to **Discord** as formatted webhook messages (embeds). It covers three event families:
pull requests, pipeline (build) completions, and release deployments.

The service is stateless: no database, no queue, no auth on the inbound routes. Each request
is parsed, translated into a Discord payload, and POSTed to the corresponding Discord webhook URL.

Go module name is `discord-azure-integration` (note: it differs from the repo directory name
`azure-devops-discord-integration`). All internal imports use that module path.

## Layout

```
main.go              Entry point: load config -> build Server -> Init() (blocking)
config/config.go     ConfigServer struct + .env loading + os.Getenv mapping
server/server.go     Server struct, Init() -> SetupRouter() + r.Run()
server/router.go     SetupRouter(): gin engine, controller wiring, route table
controllers/
  constants.go       Route path constants + HEADER_APP_JSON
  pull-request.go    PullRequestController: CreatedPR, ReviewedPR, StatusUpdatedPR, processVotes
  pipeline.go        PipelineController: PipelineStatusReport, processStatus, getRepository
  release.go         ReleaseController: ReleaseStatusReport, processReleaseStatus
  *_test.go          Unit tests for the pure process*() helpers (table-driven)
models/
  azure.go           Structs mirroring Azure DevOps webhook JSON (AzureRequest, Resource, ...)
  discord.go         Discord webhook payload structs + color constants
  utils.go           AzureRequest -> DiscordPayload converters + label helpers
mocks.go             package main. Large fixtures (fakePayload*) used only by main_test.go
main_test.go         End-to-end route tests via httptest + gin ServeHTTP
Dockerfile           golang:1.16-alpine build, exposes 8080
docker-compose.yml   Single `app` service, host port 8080
```

## Request flow

1. Azure DevOps POSTs a Service Hook payload to one of the routes.
2. The controller does `c.ShouldBindJSON(&p.Response)` into its `*models.AzureRequest` field.
3. A `process*()` helper maps an Azure status/vote to a `(color int32, title string)` pair.
4. If `title == ""` the handler replies `204 No Content` and sends **nothing** to Discord —
   this is the deliberate "uninteresting event" path (e.g. a PR update that is not
   completed/conflicts, or a review with only neutral votes).
5. Otherwise `models.ConvertToDiscordPayload*()` builds the embed, it is marshalled and
   POSTed to the configured Discord webhook URL.
6. Success replies `200` echoing the parsed Azure payload; any error replies
   `400` with `{"err": ...}` after `fmt.Println(err)`.

`PipelineStatusReport` additionally calls out to the **Azure DevOps REST API**
(`GET {AZURE_ORGANIZATION}/{AZURE_PROJECT}/_apis/git/repositories/{id}?api-version=5.1`)
using Basic auth with a base64-encoded `":" + AZURE_PAT_TOKEN`, to resolve the triggering
repository name from `resource.triggerInfo["ci.triggerRepository"]`.

## Routes

Defined as constants in `controllers/constants.go` — **always reference the constants**,
never hardcode path strings (tests depend on this).

| Constant | Path | Handler | Discord webhook used |
|---|---|---|---|
| `CREATED_ROUTE` | `POST /pull-request/created` | `PullRequestController.CreatedPR` | `DISCORD_PR_URL` |
| `REVIEW_ROUTE` | `POST /pull-request/review` | `PullRequestController.ReviewedPR` | `DISCORD_PR_URL` |
| `STATUS_ROUTE` | `POST /pull-request/status` | `PullRequestController.StatusUpdatedPR` | `DISCORD_PR_URL` |
| `PIPELINE_ROUTE` | `POST /pipeline/` | `PipelineController.PipelineStatusReport` | `DISCORD_PIPELINE_URL` |
| `RELEASE_ROUTE` | `POST /release/` | `ReleaseController.ReleaseStatusReport` | `DISCORD_RELEASE_URL` |

Note the trailing slash on `/pipeline/` and `/release/` — Azure hook URLs must match.

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

PR status (`StatusUpdatedPR`): `completed` -> BLURPLE "Concluído", `conflicts` -> RED
"com Conflito", anything else -> `204`.

Pipeline `resource.result` and release `resource.deployment.deploymentStatus` share the same
mapping: `succeeded` -> GREEN "Concluída", `failed` -> RED "Falhada", `stopped` -> ORANGE
"Interrompida", default -> WHITE `[Status não mapeado: X]` (still notifies, on purpose).

Merge status labels live in `getMergeStatusText()` (`succeeded`, `conflicts`, `queued`,
`rejectedByPolicy`, `failure`).

## Conventions

- **Language split**: identifiers, comments and log output are English; user-facing Discord
  strings (titles, field names, labels) are **Brazilian Portuguese**. Keep new Discord copy
  in pt-BR to stay consistent.
- **Controllers** are structs holding `ConfigServer *config.ConfigServer` and
  `Response *models.AzureRequest`, instantiated once in `SetupRouter` and bound as method
  handlers. Adding a feature means: new struct/method in `controllers/`, a route constant in
  `constants.go`, and wiring in `SetupRouter`.
- **Business logic goes in a pure `process*()` method** returning `(int32, string)` — that is
  what the `controllers/*_test.go` table-driven tests exercise. Handlers stay thin (bind,
  map, convert, POST, respond).
- **Payload construction lives in `models/utils.go`** as methods on `*AzureRequest`
  (`ConvertToDiscordPayloadPR/Pipeline/Release`). Do not build Discord JSON inside controllers.
- **Colors** are the constants in `models/discord.go` (`GREEN`, `RED`, `ORANGE`, `YELLOW`,
  `BLURPLE`, `WHITE`, `GRAY`, `BLACK`) — never raw ints.
- **Error handling** is uniform and intentionally simple: `fmt.Println(err)` then
  `c.JSON(http.StatusBadRequest, gin.H{"err": err})` and `return`. Match it rather than
  introducing a logger or middleware unless asked.
- Azure model structs are a permissive superset: `Resource` carries PR, build and release
  fields together, so one `AzureRequest` type binds every hook. Add fields to the existing
  structs rather than creating parallel types.
- Formatting: standard `gofmt` (tabs). Run `gofmt -l .` before committing; note that
  `mocks.go` is already unformatted on `main`, so ignore that one pre-existing hit.

## Configuration

`config.LoadEnvironment()` loads dotenv files then reads `os.Getenv`. Load order (later
files do **not** override values already set — `godotenv.Load` is first-wins, so the first
file listed effectively has priority):

1. `.env.{APP_ENV}.local`
2. `.env.local` (skipped entirely when `APP_ENV=test`)
3. `.env.{APP_ENV}`
4. `.env`

`APP_ENV` defaults to `development`. **If none of the four files exists, `LoadEnvironment`
returns `errors.New("no env file was loaded")`** and `main` exits — env vars alone are not
enough, a file must exist.

| Variable | Purpose |
|---|---|
| `APP_ENV` | Selects the dotenv variant (`development` / `production` / `test`) |
| `GIN_MODE` | Passed to `gin.SetMode` when non-empty (`debug` / `release` / `test`) |
| `DISCORD_PR_URL` | Discord webhook for all three PR routes |
| `DISCORD_PIPELINE_URL` | Discord webhook for `/pipeline/` |
| `DISCORD_RELEASE_URL` | Discord webhook for `/release/` |
| `AZURE_ORGANIZATION` | Base URL used to build the Azure REST call (used as a URL prefix, not a bare org name) |
| `AZURE_PROJECT` | Project segment of the Azure REST call |
| `AZURE_PAT_TOKEN` | PAT used for Basic auth against the Azure REST API |

The listen address comes from gin's `r.Run()` default: `:8080`, overridable with the `PORT`
env var. All `.env*` files are gitignored — never commit one, and never echo real webhook
URLs or PAT values into logs, commits, or PR descriptions.

## Build, run, test

```bash
go build ./...                 # compile
go run main.go                 # run (needs a .env file present)
go vet ./...
gofmt -l .                     # currently reports mocks.go (pre-existing)
go test ./...                  # full suite
go test -cover ./...
go test ./controllers/...      # unit tests only — no env or network needed
docker compose up -d           # container, host port 8080
```

### Testing gotchas — read before running `go test ./...`

- `./controllers` tests are pure and always pass offline.
- The **root package tests (`main_test.go`) are true end-to-end tests and make real outbound
  HTTP requests**. They POST to whatever `DISCORD_*_URL` points at and (for the pipeline
  route) call `AZURE_ORGANIZATION`. With real webhook URLs configured they will spam a live
  Discord channel.
- `prepareRouter()` ignores the error from `LoadEnvironment()` and `main_test.go` discards it
  with `r, _ :=`. With **no `.env*` file present the engine is nil and the suite panics with a
  nil-pointer dereference** in `gin.(*Engine).ServeHTTP` — that panic means "missing env
  file", not a code regression.
- With env files present but unreachable URLs, the handlers' `http.Post` fails and every
  route test fails on `400 != 200`.
- To run them safely, point the three `DISCORD_*_URL` vars and `AZURE_ORGANIZATION` at a local
  stub returning `200` (the pipeline route needs a JSON body decodable into
  `models.AzureRepository`), e.g. via a `.env.test` file plus `APP_ENV=test`. Note that
  `TestConfigE2ETesting` sets `APP_ENV=test` and unsets it on return, so the remaining tests
  in the file resolve under `development`.
- New fixtures for root tests go in `mocks.go` as `fakePayload*` vars (it is `package main`,
  not a `_test.go` file, so it is compiled into the binary too).

## Known rough edges

Do not "fix" these silently as part of an unrelated change; flag them or fix them deliberately.

- `Dockerfile` pins `golang:1.16-alpine` while `go.mod` declares `go 1.17`.
- Controllers store the bound payload in a **shared struct field** (`Response`), and each
  controller instance is shared across all requests — concurrent webhooks can interleave.
  Any refactor toward per-request locals is a behavioural improvement, not a no-op.
- Response bodies from the Discord `http.Post` calls are never read or closed, and a non-2xx
  Discord response is treated as success.
- `models.AzurePipeline` is only used by test fixtures; the live routes bind `AzureRequest`.
- Inbound routes are unauthenticated — anyone who can reach the port can post to the Discord
  channels.
- There is no CI workflow in this repo; checks are local only.

## Git workflow

- Default branch is `main`; feature work has historically used `feature/<name>` branches
  merged via PR.
- Commit messages are short and imperative, sometimes with a `chore:`/`fix:` prefix.
- Never commit `.env*` files or real tokens/webhook URLs.
