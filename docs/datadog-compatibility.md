# Datadog Compatibility Matrix

This document summarizes how Datadog client/protocol traffic maps through `dogbridge` today, and where behavior is intentionally different from Datadog-native ingestion.

Status legend:

- **Supported**: implemented and validated in repo examples/tests.
- **Partial**: protocol or data type is accepted, with caveats.
- **Planned**: part of roadmap, not implemented in this repository state yet.
- **Not supported**: currently out of scope.

## Protocol and signal compatibility

| Datadog input | Port / protocol | Signal | dogbridge status | Notes |
|---|---|---|---|---|
| APM traces (`dd-trace-*`) | `8126` HTTP | Traces | **Supported** | Ingested via OTel Datadog receiver and exported as OTLP traces. |
| OTLP (`otel-sdk-*`) | `4317` gRPC, `4318` HTTP | Traces | **Supported** | Native OTel path for migration end-state. |
| DogStatsD metrics | `8125` UDP | Metrics | **Planned** | Phase 2 target; not wired in current default demo config. |
| Datadog logs intake | Datadog Agent/log endpoints | Logs | **Not supported** | Use OTel log collection path (`filelog`/OTLP) when logs phase is implemented. |

## Span/attribute compatibility

| Datadog concept | OTel representation in dogbridge | Status | Notes |
|---|---|---|---|
| `service` | `service.name` resource attribute | **Supported** | Query and filtering should pivot to `service.name`. |
| `env` | `deployment.environment` or retained tag attribute | **Partial** | Standardize to OTel semantic conventions during migration. |
| `version` | `service.version` resource attribute | **Partial** | Client/library dependent; verify in pilot service telemetry. |
| Datadog span tags | OTel span/resource attributes | **Partial** | Most tags pass through; backend indexing and naming differs. |
| 128-bit trace IDs | OTLP trace IDs | **Supported** | Backend UI formatting may differ from Datadog UI. |

## Query and backend behavior differences

1. **UI/query language differs**: Datadog Trace Explorer queries are not portable to Tempo/Grafana directly; rewrite using TraceQL/search-by-tags.
2. **Retention/indexing model differs**: Datadog managed indexing rules and retention tiers do not automatically exist in OSS backends.
3. **Derived telemetry differs**: Datadog-specific features (Watchdog/AI, some inferred services) are backend-specific and not reproduced by dogbridge.

## Known limitations (current phase)

- Metrics migration path (DogStatsD -> Prometheus remote write) is not yet enabled by default.
- Log migration path is not yet enabled in the local demo.
- Compatibility guarantees are currently focused on trace ingestion and trace delivery.

## Recommended pilot validation checklist

Before broad rollout, validate each candidate service:

1. Service appears in backend with expected `service.name`.
2. End-to-end trace completeness matches Datadog baseline for top endpoints.
3. Error traces preserve status and principal tags needed for triage.
4. p95/p99 latency panels can be rebuilt from backend-native queries.
5. Alert parity exists for SLO/error-rate and saturation signals used by on-call.
