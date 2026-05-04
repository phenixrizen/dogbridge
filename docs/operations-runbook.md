# Dogbridge operations runbook (Phase 5)

This runbook documents production hardening defaults, backpressure behavior, and SRE operating guidance.

## Recommended baseline config

Use `config/examples/hardened-all-in-one.yaml` as the starting point for production:

- `memory_limiter` constrains process memory and sheds load when limits are reached.
- `batch` increases exporter efficiency and smooths burst traffic.
- Exporter `sending_queue` absorbs short backend interruptions.
- Exporter `retry_on_failure` applies bounded exponential backoff.
- Exporter `timeout` prevents indefinite hangs.

## Backpressure behavior

Dogbridge backpressure for each signal follows this order:

1. Exporter failures fill `sending_queue`.
2. Once queue capacity is exhausted, exporter send failures increase.
3. Pipeline pressure can propagate to receivers, increasing refused telemetry counters.
4. `memory_limiter` rejects data before OOM to preserve collector process health.

## Self-observability dashboards

Scrape collector metrics endpoint (`service.telemetry.metrics.address`, default `:8888`) and chart these by pipeline/signal:

- Accepted vs refused: receiver accepted/refused counters
- Queue utilization: exporter queue size/capacity
- Delivery success/failure: exporter sent/failed counters
- Batch behavior: processor batch send size and timeout flush frequency
- Resource pressure: process RSS/heap and GC pause

Recommended alerting:

- Any sustained non-zero refused telemetry rate for 5+ minutes
- Exporter send-failed ratio > 1% over 10 minutes
- Queue utilization sustained > 80%

## Load test procedure

Run signal-specific scenarios:

```bash
./tests/load/traces-load.sh
./tests/load/metrics-load.sh
./tests/load/logs-load.sh
```

Scale `REQUESTS` and `CONCURRENCY` until you find the first sustained drop point, then keep steady-state load below that threshold with a safety margin.

## Scaling guidance by signal type

- Traces: usually CPU bound due to high span volume and serialization; scale with more replicas and moderate batch sizes.
- Metrics: sensitive to cardinality; enforce tag controls and monitor remote-write queue depth.
- Logs: throughput and payload-size heavy; watch memory and queue growth, and prefer dedicated nodes for high-ingest tiers.

## Security guidance

- Enable mTLS on all receiver and exporter transport paths crossing trust boundaries.
- Store backend credentials and certificates in Kubernetes Secrets; mount as files, not inline values.
- Apply NetworkPolicy to limit ingress to required ports (`8126`, `4317`, `4318`) and egress to approved telemetry backends.
- Run with least privilege service accounts and avoid hostPath unless log collection requires it.
