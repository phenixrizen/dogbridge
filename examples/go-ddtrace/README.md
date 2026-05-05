# go-ddtrace

This demo emits **Datadog spans, metrics, and event-style logs** continuously to dogbridge using Datadog Go libraries:

- Spans via `dd-trace-go`
- Metrics via DogStatsD metric packets
- Log-like events via DogStatsD events

## Run

```bash
cd examples/go-ddtrace
go run .
```

Environment variables:

- `DD_TRACE_AGENT_URL` (default: `http://localhost:8126`)
- `DD_DOGSTATSD_ADDR` (default: `localhost:8125`)
