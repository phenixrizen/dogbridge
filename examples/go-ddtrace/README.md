# go-ddtrace

This demo emits **Datadog spans, metrics, and correlated OTLP logs** continuously to dogbridge:

- Spans via `dd-trace-go`
- Metrics via DogStatsD metric packets
- Logs via OTLP/HTTP with trace and span IDs from the active Datadog span

## Run

```bash
cd examples/go-ddtrace
go run .
```

Environment variables:

- `DD_TRACE_AGENT_URL` (default: `http://localhost:8126`)
- `DD_DOGSTATSD_ADDR` (default: `localhost:8125`)
- `DOGBRIDGE_OTLP_LOGS_ENDPOINT` (default: `http://localhost:4318/v1/logs`)
- `DOGBRIDGE_DEMO_RUN_ID` (default: current Unix nanoseconds)
- `DOGBRIDGE_DEMO_BATCHES` (default: `0`, emit continuously; set a positive number for smoke tests)
