#!/usr/bin/env bash
set -euo pipefail

ENDPOINT="${ENDPOINT:-http://127.0.0.1:4318/v1/traces}"
REQUESTS="${REQUESTS:-500}"
CONCURRENCY="${CONCURRENCY:-20}"

emit_trace() {
  local id
  id="$(openssl rand -hex 16)"
  local span
  span="$(openssl rand -hex 8)"
  curl -sS -o /dev/null -X POST "$ENDPOINT" \
    -H 'Content-Type: application/json' \
    --data "{\"resourceSpans\":[{\"resource\":{\"attributes\":[{\"key\":\"service.name\",\"value\":{\"stringValue\":\"dogbridge-load-traces\"}}]},\"scopeSpans\":[{\"spans\":[{\"traceId\":\"$id\",\"spanId\":\"$span\",\"name\":\"load-span\",\"kind\":2,\"startTimeUnixNano\":\"1714521600000000000\",\"endTimeUnixNano\":\"1714521600100000000\"}]}]}]}"
}

export -f emit_trace
export ENDPOINT
seq "$REQUESTS" | xargs -P "$CONCURRENCY" -n1 bash -c 'emit_trace'
echo "sent $REQUESTS trace requests to $ENDPOINT"
