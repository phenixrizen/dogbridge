#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/examples/docker-compose/signoz/docker-compose.yaml"
BATCHES="${DOGBRIDGE_DEMO_BATCHES:-3}"
RUN_ID="${DOGBRIDGE_DEMO_RUN_ID:-dogbridge-signoz-smoke-$(date +%s%N)}"
TRACE_COUNT_QUERY="SELECT count() FROM signoz_traces.distributed_signoz_index_v3 WHERE serviceName = 'dogbridge-dd-demo' AND attributes_string['demo.run_id'] = '$RUN_ID'"
LOG_COUNT_QUERY="SELECT count() FROM signoz_logs.distributed_logs_v2 WHERE resources_string['service.name'] = 'dogbridge-dd-demo' AND attributes_string['demo.run_id'] = '$RUN_ID'"
CORRELATED_LOG_COUNT_QUERY="SELECT count() FROM signoz_logs.distributed_logs_v2 AS l GLOBAL ANY INNER JOIN signoz_traces.distributed_signoz_index_v3 AS t ON l.trace_id = toString(t.trace_id) AND l.span_id = t.span_id WHERE l.resources_string['service.name'] = 'dogbridge-dd-demo' AND l.attributes_string['demo.run_id'] = '$RUN_ID'"

query_count() {
  docker compose -f "$COMPOSE_FILE" exec -T clickhouse clickhouse-client --query "$1"
}

echo "Running dd-trace smoke emitter against dogbridge with run_id=$RUN_ID..."
pushd "$ROOT_DIR" >/dev/null
DOGBRIDGE_DEMO_BATCHES="$BATCHES" DOGBRIDGE_DEMO_RUN_ID="$RUN_ID" go run ./examples/go-ddtrace
popd >/dev/null

deadline=$((SECONDS + 30))
trace_count=0
log_count=0
correlated_log_count=0

while (( SECONDS < deadline )); do
  trace_count="$(query_count "$TRACE_COUNT_QUERY")"
  log_count="$(query_count "$LOG_COUNT_QUERY")"
  correlated_log_count="$(query_count "$CORRELATED_LOG_COUNT_QUERY")"
  if (( trace_count >= BATCHES && log_count >= BATCHES && correlated_log_count >= BATCHES )); then
    break
  fi
  sleep 1
done

if (( trace_count < BATCHES )); then
  echo "Expected at least $BATCHES SigNoz trace rows for run_id=$RUN_ID; got=$trace_count"
  exit 1
fi

if (( log_count < BATCHES )); then
  echo "Expected at least $BATCHES SigNoz log rows for run_id=$RUN_ID; got=$log_count"
  exit 1
fi

if (( correlated_log_count < BATCHES )); then
  echo "Expected at least $BATCHES SigNoz logs to correlate to trace/span IDs for run_id=$RUN_ID; got=$correlated_log_count"
  exit 1
fi

echo "SigNoz smoke test passed for run_id=$RUN_ID: trace rows=$trace_count log rows=$log_count correlated logs=$correlated_log_count."
