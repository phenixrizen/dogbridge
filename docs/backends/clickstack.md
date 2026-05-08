# ClickStack / HyperDX backend

Dogbridge can send Datadog-compatible traces, DogStatsD metrics, and OTLP logs to ClickStack through the standard OTLP HTTP exporter. ClickStack stores OpenTelemetry data in ClickHouse and presents it through the HyperDX UI.

## Pipeline summary

- Receivers: `datadog`, `statsd`, `otlp`
- Processors: `memory_limiter`, `resource`, metric cardinality transforms, `batch`
- Exporter: `otlp_http/clickstack`

## Local compose target

The `examples/docker-compose/clickstack/docker-compose.yaml` stack runs the ClickStack all-in-one image, the official ClickStack OpenTelemetry Collector image, and a dogbridge collector. It pins both ClickStack images to `2.19.0`.

```bash
make demo-clickstack-up
make demo-clickstack-smoke
```

The HyperDX UI is available at `http://localhost:3401`. The ClickStack OTLP collector is published on `localhost:5517` and `localhost:5518` for this demo.

Dogbridge uses non-default host ports in this demo so it can run alongside the Tempo, SigNoz, and OpenObserve stacks:

- Datadog traces: `http://localhost:9226`
- DogStatsD UDP: `localhost:9225`
- OTLP HTTP logs: `http://localhost:5418/v1/logs`

The smoke test emits a run-scoped trace, log, and metric workload, then verifies the data in ClickStack's ClickHouse tables:

- traces in `default.otel_traces`
- logs in `default.otel_logs`, including non-empty `TraceId` and `SpanId`
- metrics in `default.otel_metrics_sum`

## Exporter shape

The demo sends telemetry from dogbridge to the OTLP HTTP endpoint exposed by ClickStack's bundled OpenTelemetry Collector:

```yaml
exporters:
  otlp_http/clickstack:
    endpoint: http://clickstack-otel-collector:4318
    compression: gzip
```

The all-in-one ClickStack image provides the HyperDX UI, MongoDB, and ClickHouse. This demo runs the official ClickStack OpenTelemetry Collector as a separate service so the OTLP receiver is active immediately and writes into the same ClickHouse tables used by HyperDX.

The local all-in-one demo does not configure an ingestion key. For managed ClickStack or hardened open-source deployments, configure the required authorization header and send dogbridge traffic to the appropriate ClickStack collector endpoint.
