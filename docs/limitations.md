# Current Limitations

This project is in early scaffold maturity and currently optimized for trace migration pilots.

## Scope limitations

- End-to-end local demo is trace-focused (Datadog APM -> dogbridge -> Tempo).
- Metrics and logs migration tracks are planned but not yet feature-complete in the default distribution path.
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
