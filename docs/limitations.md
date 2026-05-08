# Current Limitations

This project is currently optimized for local migration pilots and validated backend demos built on upstream OpenTelemetry Collector components.

## Scope limitations

- Backend demos cover Datadog APM traces, DogStatsD metrics, and OTLP logs, but production rollout patterns still need environment-specific validation.
- Datadog logs intake endpoints are not implemented; use OpenTelemetry log collection paths such as OTLP or filelog-style collection.
- Helm chart production hardening is incomplete for large-cluster operational use.

## Datadog compatibility limitations

- Datadog-native query UX is not portable; dashboards and monitors must be rebuilt in backend-native tooling.
- Datadog managed features (certain anomaly/inference products) are not reproduced by OSS backends.
- Tag/attribute naming may require normalization to OTel semantic conventions for consistent cross-service querying.

## Operational limitations

- Capacity guidance and autoscaling defaults are not yet benchmarked for all three signals under sustained production load.
- Security hardening guidance exists but should be adapted per environment (mTLS, secret rotation, NetworkPolicy).
- Formal compatibility certification across all `dd-trace-*` language clients has not yet been completed.

## Documentation maturity notes

- Migration docs are actionable for pilot rollouts, but should be validated against your own service SLIs before broad cutover.
- Where examples are provided, treat them as templates and adjust names, dimensions, and thresholds to your environment.
