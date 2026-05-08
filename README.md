# dogbridge

`dogbridge` is an opinionated OpenTelemetry Collector distribution for teams migrating from Datadog ingestion paths to open backends.

## What works now (May 2026)

- Runnable collector entrypoint in `cmd/dogbridge`.
- Datadog APM trace ingest on `:8126` and OTLP ingest on `:4317/:4318`.
- Trace forwarding to OTLP backends (Tempo, SigNoz, OpenObserve, and ClickStack demos included).
- DogStatsD ingest and Prometheus remote write example in the Tempo compose config.
- OTLP log forwarding in backend demos that support logs, with smoke coverage for trace/log correlation.
- Local smoke tests for traces, metrics, and backend-specific log pipelines.

## Architecture

```mermaid
flowchart LR
    A[dd-trace-go] --> B[dogbridge datadogreceiver :8126]
    C[OTLP clients] --> D[dogbridge otlpreceiver :4317/:4318]
    B --> E[processors: memory_limiter/resource/batch]
    D --> E
    E --> F[OTLP exporter]
    F --> G[Tempo]
```

```mermaid
flowchart TB
    subgraph Ingress
      DD[Datadog APM]
      DS[DogStatsD]
      OTLP[OTLP]
    end

    subgraph dogbridge
      R[Receivers]
      P[Processors]
      X[Exporters]
    end

    subgraph Backends
      T[Tempo]
      VM[VictoriaMetrics]
      SN[SigNoz]
      OO[OpenObserve]
      CS[ClickStack / HyperDX]
    end

    DD --> R
    DS --> R
    OTLP --> R
    R --> P --> X
    X --> T
    X --> VM
    X --> SN
    X --> OO
    X --> CS
```

## Quickstart (Tempo demo)

```bash
make demo-up
make demo-smoke-traces
```

Expected result: smoke test reports a trace id queryable via Tempo APIs.

## Quickstart (SigNoz demo)

```bash
make demo-signoz-up
make demo-signoz-smoke
```

Then open SigNoz at `http://localhost:3301`.

## Quickstart (OpenObserve demo)

```bash
make demo-openobserve-up
make demo-openobserve-smoke
```

Then open OpenObserve at `http://localhost:5080` with `root@example.com / Complexpass#123`.

## Quickstart (ClickStack / HyperDX demo)

```bash
make demo-clickstack-up
make demo-clickstack-smoke
```

Then open the HyperDX UI at `http://localhost:3401`.

## CLI usage

```bash
# default run
./dogbridge

# explicit config
./dogbridge run --config examples/docker-compose/dogbridge.yaml

# inline override
./dogbridge run --config examples/docker-compose/dogbridge.yaml --set exporters::otlp::endpoint:tempo:4317

# version
./dogbridge version
```

## Repository layout

- `cmd/dogbridge`: cobra CLI + collector execution.
- `distro/components.go`: enabled collector factories and resolver settings.
- `examples/docker-compose`: local Tempo demo.
- `examples/docker-compose/signoz`: local SigNoz demo.
- `examples/docker-compose/openobserve`: local OpenObserve demo.
- `examples/docker-compose/clickstack`: local ClickStack / HyperDX demo.
- `config/examples`: standalone OTel config examples.
- `docs/backends`: backend-specific notes and example exporter shapes.

## Next Work

- Kubernetes logs pipeline and operational hardening.
- Compatibility matrix and migration playbooks expansion.

## License

Apache-2.0.
