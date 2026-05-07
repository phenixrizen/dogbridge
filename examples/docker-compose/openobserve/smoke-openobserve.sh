#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
OPENOBSERVE_URL="${OPENOBSERVE_URL:-http://localhost:5080}"
OPENOBSERVE_ORG="${OPENOBSERVE_ORG:-default}"
OPENOBSERVE_AUTH="${OPENOBSERVE_AUTH:-Basic cm9vdEBleGFtcGxlLmNvbTpDb21wbGV4cGFzcyMxMjM=}"
BATCHES="${DOGBRIDGE_DEMO_BATCHES:-3}"
RUN_ID="${DOGBRIDGE_DEMO_RUN_ID:-dogbridge-openobserve-smoke-$(date +%s%N)}"
START_US="$(($(date +%s%N) / 1000 - 60000000))"

now_us() {
  echo "$(($(date +%s%N) / 1000))"
}

search_logs() {
  local sql="$1"
  curl -fsS \
    -H "Authorization: $OPENOBSERVE_AUTH" \
    -H "Content-Type: application/json" \
    "$OPENOBSERVE_URL/api/$OPENOBSERVE_ORG/_search" \
    -d "{\"query\":{\"sql\":\"$sql\",\"start_time\":$START_US,\"end_time\":$(now_us),\"from\":0,\"size\":100},\"search_type\":\"ui\",\"timeout\":0}"
}

count_log_hits_with_trace_context() {
  jq -r '[.hits[]? | select((.trace_id // "") != "" and (.span_id // "") != "")] | length'
}

count_traces() {
  curl -fsS -G \
    -H "Authorization: $OPENOBSERVE_AUTH" \
    --data-urlencode "filter=service_name='dogbridge-dd-demo' AND demo_run_id='$RUN_ID'" \
    --data-urlencode "start_time=$START_US" \
    --data-urlencode "end_time=$(now_us)" \
    --data-urlencode "from=0" \
    --data-urlencode "size=100" \
    "$OPENOBSERVE_URL/api/$OPENOBSERVE_ORG/dogbridge_demo/traces/latest" |
    jq -r '.total // 0'
}

count_metrics() {
  curl -fsS -G \
    -H "Authorization: $OPENOBSERVE_AUTH" \
    --data-urlencode "query=dogbridge_demo_requests_total{run_id=\"$RUN_ID\"}" \
    "$OPENOBSERVE_URL/api/$OPENOBSERVE_ORG/prometheus/api/v1/query" |
    jq -r '.data.result | length'
}

echo "Running dd-trace smoke emitter against OpenObserve with run_id=$RUN_ID..."
pushd "$ROOT_DIR" >/dev/null
DD_TRACE_AGENT_URL=http://localhost:9126 \
DD_DOGSTATSD_ADDR=localhost:9125 \
DOGBRIDGE_OTLP_LOGS_ENDPOINT=http://localhost:5318/v1/logs \
DOGBRIDGE_DEMO_BATCHES="$BATCHES" \
DOGBRIDGE_DEMO_RUN_ID="$RUN_ID" \
go run ./examples/go-ddtrace
popd >/dev/null

deadline=$((SECONDS + 45))
trace_count=0
log_count=0
metric_count=0

while (( SECONDS < deadline )); do
  trace_count="$(count_traces)"
  log_count="$(search_logs "SELECT * FROM dogbridge_demo WHERE service_name='dogbridge-dd-demo' AND demo_run_id='$RUN_ID' AND body LIKE '%dummy app emitted batch%'" | count_log_hits_with_trace_context)"
  metric_count="$(count_metrics)"
  if (( trace_count >= BATCHES && log_count >= BATCHES && metric_count >= BATCHES )); then
    break
  fi
  sleep 1
done

if (( trace_count < BATCHES )); then
  echo "Expected at least $BATCHES OpenObserve trace rows for run_id=$RUN_ID; got=$trace_count"
  exit 1
fi

if (( log_count < BATCHES )); then
  echo "Expected at least $BATCHES OpenObserve log rows for run_id=$RUN_ID; got=$log_count"
  exit 1
fi

if (( metric_count < BATCHES )); then
  echo "Expected at least $BATCHES OpenObserve metric rows for run_id=$RUN_ID; got=$metric_count"
  exit 1
fi

echo "OpenObserve smoke test passed for run_id=$RUN_ID: trace rows=$trace_count log rows=$log_count metric rows=$metric_count."
