# AGENTS.md

## Scope
These instructions apply to the entire repository.

## Project intent
`dogbridge` is an opinionated OpenTelemetry Collector distribution that provides Datadog-compatible ingestion as a migration bridge to open-source observability backends.

## Current maturity
This repository is in an early scaffold phase. Prefer incremental, runnable progress over broad placeholders.

## Contribution rules
1. Keep changes focused and milestone-driven.
2. Prefer concrete, working examples over stub files.
3. When adding config examples, ensure they are valid OTel Collector YAML.
4. Keep docs aligned with implemented behavior; clearly mark aspirational sections.
5. For new Go code, run `go test ./...` before finalizing.
6. For Helm changes, run `helm lint helm/dogbridge` when Helm is available.

## Near-term priorities
- Implement a runnable collector binary using OTel Collector builder/factories.
- Provide one end-to-end local demo (dd-trace-go -> dogbridge -> tempo).
- Replace placeholder docs with migration and compatibility specifics.

## Commit/PR expectations for agents
- Summarize user-visible behavior changes.
- Include exact commands used for validation.
- Call out any remaining placeholders explicitly.
