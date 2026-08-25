#!/usr/bin/env bash
# Real end-to-end smoke test: starts the actual server against your real
# .env (whatever webhook URLs you have configured there) and POSTs a
# minimal fake Azure DevOps payload at it, so a configured sink should
# receive a real message. Unlike `go test`, this hits real network
# endpoints — run it manually, never in CI.
#
# Usage:
#   ./scripts/smoke-test.sh                       # POST /pull-request/created
#   ./scripts/smoke-test.sh /pipeline/             # any route constant's path
set -euo pipefail

ROUTE="${1:-/pull-request/created}"
PORT="${PORT:-8080}"

if [ ! -f .env ] && [ ! -f cmd/server/.env ]; then
  echo "No .env found in repo root. Copy .env.example to .env and set a real webhook URL first." >&2
  exit 1
fi

echo "Starting server on :$PORT ..."
go run ./cmd/server &
SERVER_PID=$!
trap 'kill "$SERVER_PID" 2>/dev/null || true' EXIT

for _ in $(seq 1 30); do
  if curl -s -o /dev/null "http://localhost:$PORT/"; then
    break
  fi
  sleep 0.5
done

echo "POSTing fake pull request payload to $ROUTE ..."
curl -sS -X POST "http://localhost:$PORT$ROUTE" \
  -H "Content-Type: application/json" \
  -d '{
    "eventType": "git.pullrequest.created",
    "resource": {
      "pullRequestId": 1,
      "title": "Smoke test PR",
      "sourceRefName": "refs/heads/feature",
      "targetRefName": "refs/heads/main",
      "mergeStatus": "succeeded",
      "createdBy": { "displayName": "Smoke Test" },
      "repository": { "name": "smoke-test-repo", "remoteUrl": "https://example.com" }
    }
  }'
echo
echo "Done. Check your configured chat destination(s) for the message."
