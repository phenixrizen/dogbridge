# AGENTS.md

## Scope
These instructions apply to the entire repository.

## Project intent
`dogbridge` is an opinionated OpenTelemetry Collector configuration, demo, and validation kit that provides Datadog-compatible ingestion patterns as a migration bridge to open-source observability backends.

## Current maturity
This repository has runnable local demos and smoke-tested backend integrations. Prefer production-oriented incremental changes, keep examples executable, and avoid adding unvalidated placeholders.

## Contribution rules
1. Keep changes focused and milestone-driven.
2. Prefer concrete, working examples over stub files.
3. When adding config examples, ensure they are valid OTel Collector YAML.
4. Keep docs aligned with implemented behavior; clearly mark aspirational sections.
5. For new Go code, run `go test ./...` before finalizing.
6. For Helm changes, run `helm lint helm/dogbridge` when Helm is available.

## Documentation maintenance
- When adding or changing a runnable backend demo, update `README.md` and the relevant `docs/backends/*.md` page in the same change.
- When adding a new backend demo, add or update a backend doc page that includes the compose target, exposed local ports, exporter shape, validation command, and any local-only credentials or placeholders.
- Keep README quickstarts limited to demos that are implemented and smoke-testable in this repository.
- Do not leave backend docs implying support for traces, metrics, or logs unless the corresponding config and smoke validation cover that signal.

## Dependency and component guidance
- Use the latest stable OpenTelemetry Collector release line unless a task explicitly requires an older pin.
- Do not reintroduce the deprecated OpenTelemetry `lokiexporter`. Send logs to Loki through the stable `otlphttp` exporter and Loki's OTLP ingestion endpoint.
- Prefer upstream `otel/opentelemetry-collector-contrib` for runnable demos; do not add a forked Collector runtime without a concrete requirement.
- After Go dependency changes, run `go mod tidy` and `go test ./...`.

## Near-term priorities
- Validate every example config against the pinned upstream Collector version and its backend smoke test.
- Expand Datadog compatibility coverage with tested traces, DogStatsD metrics, OTLP logs, and explicit unsupported Datadog-only behavior.
- Keep backend, migration, and operations docs aligned with implemented behavior.

## Commit/PR expectations for agents
- Summarize user-visible behavior changes.
- Include exact commands used for validation.
- Call out any remaining placeholders explicitly.
