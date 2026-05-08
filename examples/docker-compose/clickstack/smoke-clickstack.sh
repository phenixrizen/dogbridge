#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/examples/docker-compose/clickstack/docker-compose.yaml"
BATCHES="${DOGBRIDGE_DEMO_BATCHES:-3}"
RUN_ID="${DOGBRIDGE_DEMO_RUN_ID:-dogbridge-clickstack-smoke-$(date +%s%N)}"

clickhouse_query() {
  docker compose -f "$COMPOSE_FILE" exec -T clickstack clickhouse-client --query "$1"
}

count_traces() {
  clickhouse_query "SELECT count() FROM default.otel_traces WHERE ServiceName = 'dogbridge-dd-demo' AND SpanAttributes['demo.run_id'] = '$RUN_ID'"
}

count_logs_with_trace_context() {
  clickhouse_query "SELECT count() FROM default.otel_logs WHERE ServiceName = 'dogbridge-dd-demo' AND LogAttributes['demo.run_id'] = '$RUN_ID' AND Body LIKE '%dummy app emitted batch%' AND TraceId != '' AND SpanId != ''"
}

count_metrics() {
  clickhouse_query "SELECT count() FROM default.otel_metrics_sum WHERE MetricName = 'dogbridge.demo.requests_total' AND Attributes['run_id'] = '$RUN_ID'"
}

echo "Running dd-trace smoke emitter against ClickStack with run_id=$RUN_ID..."
pushd "$ROOT_DIR" >/dev/null
DD_TRACE_AGENT_URL=http://localhost:9226 \
DD_DOGSTATSD_ADDR=localhost:9225 \
DOGBRIDGE_OTLP_LOGS_ENDPOINT=http://localhost:5418/v1/logs \
DOGBRIDGE_DEMO_BATCHES="$BATCHES" \
DOGBRIDGE_DEMO_RUN_ID="$RUN_ID" \
go run ./examples/go-ddtrace
popd >/dev/null

deadline=$((SECONDS + 60))
trace_count=0
log_count=0
metric_count=0

while (( SECONDS < deadline )); do
  trace_count="$(count_traces)"
  log_count="$(count_logs_with_trace_context)"
  metric_count="$(count_metrics)"
  if (( trace_count >= BATCHES && log_count >= BATCHES && metric_count >= BATCHES )); then
    break
  fi
  sleep 1
done

if (( trace_count < BATCHES )); then
  echo "Expected at least $BATCHES ClickStack trace rows for run_id=$RUN_ID; got=$trace_count"
  exit 1
fi

if (( log_count < BATCHES )); then
  echo "Expected at least $BATCHES ClickStack log rows with trace context for run_id=$RUN_ID; got=$log_count"
  exit 1
fi

if (( metric_count < BATCHES )); then
  echo "Expected at least $BATCHES ClickStack metric rows for run_id=$RUN_ID; got=$metric_count"
  exit 1
fi

echo "ClickStack smoke test passed for run_id=$RUN_ID: trace rows=$trace_count log rows=$log_count metric rows=$metric_count."
