# dogbridge

`dogbridge` is a self-hosted Datadog-compatible telemetry ingestion gateway built on OpenTelemetry Collector.

It accepts Datadog-style APM traces, DogStatsD metrics, OTLP telemetry, Prometheus scrapes, and Kubernetes logs, then exports them to open-source observability backends such as Tempo, Jaeger, VictoriaMetrics, Mimir, Prometheus, Loki, OpenSearch, Elasticsearch, and Pyroscope.

This repository provides an opinionated OpenTelemetry Collector distribution, deployment manifests, examples, and migration documentation for teams moving from Datadog-specific clients to OpenTelemetry-native instrumentation.

## Goals

- Keep existing Datadog-style app settings (`DD_TRACE_AGENT_URL`, `DD_AGENT_HOST`, `DD_DOGSTATSD_PORT`) working during migration.
- Convert telemetry through OTel Collector pipelines.
- Export to open-source backends with sensible defaults.
- Encourage a long-term move to native OTel SDKs and OTLP.

## MVPs

1. Datadog traces (`:8126`) to Tempo/Jaeger.
2. DogStatsD metrics (`:8125/udp`) to Prometheus remote write backends.
3. Kubernetes logs (`filelog`) to Loki/OpenSearch/Elasticsearch.

## Repository Layout

```text
cmd/dogbridge/                Collector entrypoint
distro/components.go          Collector component registration
config/examples/              Ready-to-run pipeline examples
helm/dogbridge/               Helm chart for deployment/daemonset modes
docs/                         Architecture, migration, compatibility docs
examples/                     Go examples and local docker-compose demo
tests/e2e/                    End-to-end placeholder tests
```

## Default Ports

| Signal | Protocol | Port |
|---|---:|---:|
| Datadog traces | HTTP | 8126 |
| DogStatsD metrics | UDP | 8125 |
| OTLP gRPC | gRPC | 4317 |
| OTLP HTTP | HTTP | 4318 |
| Health check | HTTP | 13133 |
| Prometheus scrape endpoint | HTTP | 8888 |

## Quick Start

- Use `config/examples/traces-to-tempo.yaml` for MVP 1.
- Use `config/examples/metrics-to-victoriametrics.yaml` for MVP 2.
- Use `config/examples/logs-to-loki.yaml` for MVP 3.
- Use `examples/docker-compose/` for a local all-in-one stack.

## License

Apache-2.0.
