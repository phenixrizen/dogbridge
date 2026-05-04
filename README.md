# dogbridge

`dogbridge` is an opinionated OpenTelemetry Collector distribution for teams migrating off Datadog agents and proprietary backends.

## Current state (May 2026)

Implemented in this scaffold:

- Runnable collector entrypoint in `cmd/dogbridge` wired with component factories.
- Trace ingestion via Datadog APM (`:8126`) and OTLP (`:4317`/`:4318`).
- Trace export to OTLP backends (Tempo demo included).
- Local docker-compose demo with dogbridge + Tempo + Grafana.
- Smoke script that sends one `dd-trace-go` trace and verifies it is queryable in Tempo.

Not implemented yet (planned phases):

- DogStatsD metrics pipeline and VictoriaMetrics flow.
- Kubernetes logs pipeline and Loki/OpenSearch flows.
- Additional Helm hardening for production operations (PDB/HPA/securityContext defaults).

## Quickstart: Datadog traces to Tempo

```bash
cd examples/docker-compose
docker compose up -d
./smoke-traces.sh
```

If the smoke test succeeds, Tempo search API returns at least one trace ID for service `dogbridge-ddtrace-demo`.

## Repository layout

- `cmd/dogbridge`: collector binary entrypoint.
- `distro/components.go`: enabled receivers/processors/exporters and build settings.
- `examples/docker-compose`: local end-to-end demo stack.
- `examples/go-ddtrace`: simple trace emitter.
- `config/examples`: standalone configuration examples.

## License

Apache-2.0.
