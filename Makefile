.PHONY: build run vet fmt test test-cover test-unit test-discord test-googlechat test-e2e

build:
	go build ./...

run:
	go run ./cmd/server

vet:
	go vet ./...

fmt:
	gofmt -l .

# Full suite. cmd/server's tests are end-to-end: they spin up a local stub
# server and point every configured sink (Discord + Google Chat) plus the
# Azure REST call at it, so this runs fully offline.
test:
	go test ./...

test-cover:
	go test -cover ./...

# internal/azuredevops's tests are pure unit tests (no env/network needed).
test-unit:
	go test ./internal/azuredevops/...

# cmd/server's end-to-end route tests, exercising only the Discord sink.
test-discord:
	TEST_SINKS=discord go test ./cmd/server/...

# cmd/server's end-to-end route tests, exercising only the Google Chat sink.
test-googlechat:
	TEST_SINKS=googlechat go test ./cmd/server/...

# cmd/server's end-to-end route tests, exercising both sinks at once.
test-e2e:
	TEST_SINKS=both go test ./cmd/server/...
