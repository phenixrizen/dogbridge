# Architecture

Dogbridge is an opinionated set of OpenTelemetry Collector configs, local demos, and validation scripts. It does not maintain a forked Collector runtime; runnable examples use upstream `otel/opentelemetry-collector-contrib`.

```mermaid
flowchart LR
    CFG[dogbridge Collector YAML] -. configures .-> R
    DD[Datadog APM clients] --> R[collector-contrib receivers]
    DS[DogStatsD clients] --> R
    OTLP[OTLP clients] --> R
    R --> P[processors: memory_limiter/resource/transform/batch]
    P --> X[backend exporters]
    X --> B[Tempo, VictoriaMetrics, SigNoz, OpenObserve, ClickStack]
```

The repository focuses on backend-specific Collector YAML, Docker Compose demos, smoke tests, and migration documentation. Add a forked runtime only if upstream Collector components cannot satisfy a concrete requirement.
