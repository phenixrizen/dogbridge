# Loki backend

Dogbridge sends logs to Loki using the stable OTLP HTTP exporter and Loki's
native OTLP ingestion endpoint.

Use `config/examples/logs-to-loki.yaml` as the focused example:

```yaml
exporters:
  otlphttp/loki:
    endpoint: http://loki.monitoring.svc.cluster.local:3100/otlp
```

The deprecated OpenTelemetry `lokiexporter` is intentionally not used. Send logs
to Loki through the stable OTLP HTTP exporter and Loki's native OTLP endpoint.
