# dogbridge

`dogbridge` is an opinionated OpenTelemetry Collector configuration, demo, and validation kit for teams migrating from Datadog ingestion paths to open backends.

## What works now (May 2026)

- Runnable local demos using upstream `otel/opentelemetry-collector-contrib`.
- Datadog APM trace ingest on `:8126` and OTLP ingest on `:4317/:4318`.
- Trace forwarding to OTLP backends (Tempo, SigNoz, OpenObserve, and ClickStack demos included).
- DogStatsD ingest and Prometheus remote write example in the Tempo compose config.
- OTLP log forwarding in backend demos that support logs, with smoke coverage for trace/log correlation.
- Local smoke tests for traces, metrics, and backend-specific log pipelines.

## Architecture

```mermaid
flowchart LR
    CFG[dogbridge Collector YAML] -. configures .-> R
    DD[Datadog APM clients] --> R[upstream collector-contrib receivers]
    DS[DogStatsD clients] --> R
    OTLP[OTLP clients] --> R
    R --> P[processors: memory_limiter/resource/transform/batch]
    P --> X[backend exporters]
    X --> B[open observability backends]
```

```mermaid
flowchart TB
    subgraph Ingress
      DD[Datadog APM]
      DS[DogStatsD]
      OTLP[OTLP]
    end

    CFG[dogbridge Collector YAML]

    subgraph Collector[upstream collector-contrib container]
      R[Datadog, StatsD, and OTLP receivers]
      P[Memory, resource, transform, and batch processors]
      X[OTLP, OTLP HTTP, and remote-write exporters]
    end

    subgraph Backends
      T[Tempo]
      VM[VictoriaMetrics]
      SN[SigNoz]
      OO[OpenObserve]
      CS[ClickStack / HyperDX]
    end

    CFG -. mounted config .-> R
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
The demo runs upstream `otel/opentelemetry-collector-contrib` with the dogbridge config mounted at `/etc/dogbridge/config.yaml`.

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

## Collector Usage

```bash
docker run --rm \
  -p 8126:8126 \
  -p 8125:8125/udp \
  -p 4317:4317 \
  -p 4318:4318 \
  -v "$PWD/examples/docker-compose/dogbridge.yaml:/etc/dogbridge/config.yaml:ro" \
  otel/opentelemetry-collector-contrib:0.151.0 \
  --config=/etc/dogbridge/config.yaml
```

## Repository layout

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
