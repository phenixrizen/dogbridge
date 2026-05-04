#!/usr/bin/env bash
set -euo pipefail

ENDPOINT="${ENDPOINT:-http://127.0.0.1:4318/v1/metrics}"
REQUESTS="${REQUESTS:-500}"
CONCURRENCY="${CONCURRENCY:-20}"

emit_metric() {
  curl -sS -o /dev/null -X POST "$ENDPOINT" \
    -H 'Content-Type: application/json' \
    --data '{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"dogbridge-load-metrics"}}]},"scopeMetrics":[{"metrics":[{"name":"dogbridge.load.counter","sum":{"aggregationTemporality":2,"isMonotonic":true,"dataPoints":[{"asInt":"1","timeUnixNano":"1714521600000000000"}]}}]}]}]}'
}

export -f emit_metric
export ENDPOINT
seq "$REQUESTS" | xargs -P "$CONCURRENCY" -n1 bash -c 'emit_metric'
echo "sent $REQUESTS metric requests to $ENDPOINT"
