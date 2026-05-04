# Load scenarios (Phase 5)

These scripts exercise dogbridge signal pipelines independently through OTLP/HTTP.

## Prerequisites

- A running dogbridge instance exposing `:4318`.
- Optional: exported telemetry backend and collector self-metrics on `:8888`.

## Scenarios

- `traces-load.sh`: emits OTLP traces in parallel.
- `metrics-load.sh`: emits OTLP monotonic sum metrics in parallel.
- `logs-load.sh`: emits OTLP logs in parallel.

All scripts support:

- `ENDPOINT` (default uses `127.0.0.1:4318`)
- `REQUESTS` (default `500`)
- `CONCURRENCY` (default `20`)

Example:

```bash
REQUESTS=2000 CONCURRENCY=50 ./tests/load/traces-load.sh
```

After execution, check self-observability counters such as:

- `otelcol_receiver_accepted_spans`, `otelcol_exporter_sent_spans`
- `otelcol_receiver_refused_metric_points`, `otelcol_exporter_send_failed_metric_points`
- `otelcol_receiver_refused_log_records`, `otelcol_exporter_send_failed_log_records`

Use these to confirm backpressure impact and dropped telemetry behavior.
