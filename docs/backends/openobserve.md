# OpenObserve backend

Dogbridge can send Datadog-compatible traces, DogStatsD metrics, and OTLP logs to OpenObserve through the standard OTLP HTTP exporter.

## Pipeline summary

- Receivers: `datadog`, `statsd`, `otlp`
- Processors: `memory_limiter`, `resource`, metric cardinality transforms, `batch`
- Exporter: `otlp_http/openobserve`

## Local compose target

The `examples/docker-compose/openobserve/docker-compose.yaml` stack runs OpenObserve and a dogbridge collector. It pins OpenObserve OSS to `public.ecr.aws/zinclabs/openobserve:v0.80.2`.

```bash
make demo-openobserve-up
make demo-openobserve-smoke
```

OpenObserve is available at `http://localhost:5080` with local demo credentials:

- user: `root@example.com`
- password: `Complexpass#123`

Dogbridge uses non-default host ports in this demo so it can run alongside the SigNoz stack:

- Datadog traces: `http://localhost:9126`
- DogStatsD UDP: `localhost:9125`
- OTLP HTTP logs: `http://localhost:5318/v1/logs`

## Exporter shape

OpenObserve expects the OTLP HTTP exporter endpoint to point at the organization base path without a trailing slash:

```yaml
exporters:
  otlp_http/openobserve:
    endpoint: http://openobserve:5080/api/default
    headers:
      Authorization: Basic cm9vdEBleGFtcGxlLmNvbTpDb21wbGV4cGFzcyMxMjM=
      stream-name: dogbridge_demo
```

The demo intentionally uses hard-coded local credentials. Replace the `Authorization` header and organization path before using the config outside local development.
