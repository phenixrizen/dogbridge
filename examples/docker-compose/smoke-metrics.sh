#!/usr/bin/env bash
set -euo pipefail

COMPOSE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

emit() {
  local line="$1"
  echo -n "$line" | nc -4u -w1 127.0.0.1 8125
}

emit 'dogbridge.smoke.counter:3|c|#env:dev,service:smoke,request_id:req-1'
emit 'dogbridge.smoke.gauge:42|g|#env:dev,service:smoke,trace_id:trace-1'
emit 'dogbridge.smoke.hist:15|h|#env:dev,service:smoke,span_id:span-1'

sleep 12

resp="$(curl -fsS 'http://127.0.0.1:8428/api/v1/series?match[]=dogbridge_smoke_counter&start=-5m&end=now')"

if [[ "$resp" != *'"status":"success"'* ]]; then
  echo "VictoriaMetrics query failed"
  exit 1
fi

if [[ "$resp" != *'dogbridge_smoke_counter'* ]]; then
  echo "Expected dogbridge_smoke_counter series not found"
  exit 1
fi

echo "metrics smoke test passed"
