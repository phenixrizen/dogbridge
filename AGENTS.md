# AGENTS.md

## Scope
These instructions apply to the entire repository.

## Project intent
`dogbridge` is an opinionated OpenTelemetry Collector distribution that provides Datadog-compatible ingestion as a migration bridge to open-source observability backends.

## Current maturity
This repository is in an early scaffold stage. Prefer incremental, runnable progress over broad placeholders.

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
- When importing Collector components directly, depend on the component submodules in `go.mod`; do not rely on umbrella modules to supply component packages.
- Do not reintroduce the deprecated OpenTelemetry `lokiexporter`. Send logs to Loki through the stable `otlphttp` exporter and Loki's OTLP ingestion endpoint.
- After dependency or Collector component changes, run `go mod tidy` and `go test ./...`.

## Near-term priorities
- Implement a runnable collector binary using OTel Collector builder/factories.
- Provide one end-to-end local demo (dd-trace-go -> dogbridge -> tempo).
- Replace placeholder docs with migration and compatibility specifics.

## Commit/PR expectations for agents
- Summarize user-visible behavior changes.
- Include exact commands used for validation.
- Call out any remaining placeholders explicitly.
