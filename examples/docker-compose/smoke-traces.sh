#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

pushd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null

docker compose up -d

pushd "$ROOT_DIR/examples/go-ddtrace" >/dev/null
DD_TRACE_AGENT_URL="http://localhost:8126" go run .
popd >/dev/null

sleep 3

SERVICE_JSON="$(curl -fsS 'http://localhost:3200/api/search?tags=service.name=dogbridge-ddtrace-demo')"
echo "$SERVICE_JSON" | rg -q 'traceID|traces|id'

echo "Smoke test passed: Tempo search API returned trace data."

popd >/dev/null
