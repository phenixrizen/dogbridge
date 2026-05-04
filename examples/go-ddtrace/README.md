# go-ddtrace

This example emits one Datadog APM trace using `dd-trace-go` to a Datadog-compatible ingestion endpoint.

## Run

```bash
cd examples/go-ddtrace
go run .
```

Environment variables:

- `DD_TRACE_AGENT_URL` (default: `http://localhost:8126`)
