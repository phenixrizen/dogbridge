#!/usr/bin/env bash
set -euo pipefail

ENDPOINT="${ENDPOINT:-http://127.0.0.1:4318/v1/logs}"
REQUESTS="${REQUESTS:-500}"
CONCURRENCY="${CONCURRENCY:-20}"

emit_log() {
  local msg="load-log-$RANDOM"
  curl -sS -o /dev/null -X POST "$ENDPOINT" \
    -H 'Content-Type: application/json' \
    --data "{\"resourceLogs\":[{\"resource\":{\"attributes\":[{\"key\":\"service.name\",\"value\":{\"stringValue\":\"dogbridge-load-logs\"}}]},\"scopeLogs\":[{\"logRecords\":[{\"timeUnixNano\":\"1714521600000000000\",\"severityText\":\"INFO\",\"body\":{\"stringValue\":\"$msg\"}}]}]}]}"
}

export -f emit_log
export ENDPOINT
seq "$REQUESTS" | xargs -P "$CONCURRENCY" -n1 bash -c 'emit_log'
echo "sent $REQUESTS log requests to $ENDPOINT"
