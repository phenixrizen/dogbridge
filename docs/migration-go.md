# Go Migration Guide: `dd-trace-go` to OpenTelemetry SDK + OTLP

This guide is a practical migration path for Go services currently instrumented with `dd-trace-go`.

## Goals

- Keep production risk low with staged rollout.
- Preserve trace visibility during migration.
- End with vendor-neutral OTel SDK + OTLP export.

## Prerequisites

- A running dogbridge collector endpoint (`OTLP gRPC :4317` or HTTP :4318).
- Current service inventory (service name, top endpoints, critical alerts).
- Baseline trace/error/latency dashboards for comparison.

## Step A: Baseline with existing Datadog instrumentation

1. Keep existing `dd-trace-go` instrumentation in service.
2. Point Datadog trace traffic at dogbridge (`:8126`) in lower environment.
3. Validate trace parity against Datadog baseline (volume, latency distributions, key error paths).

Use the repository demo as an example of this mode:

- `examples/go-ddtrace`
- `examples/docker-compose/smoke-traces.sh`

## Step B: Add OTel instrumentation alongside legacy code

> Keep this overlap window short. Dual instrumentation can create duplicate spans and higher overhead.

1. Introduce OTel SDK and instrumentation in one low-risk endpoint path.
2. Export OTLP to dogbridge while legacy `dd-trace-go` path remains enabled.
3. Compare payload quality and field conventions (`service.name`, span names, status codes).

## Step C: Cut over to OTel-only

1. Disable `dd-trace-go` initialization and middleware.
2. Keep OTel SDK + OTLP exporter.
3. Run parity checks for at least one release cycle.
4. Remove Datadog-specific env vars from deployment manifests.

---

## Minimal OTel tracing example (Go)

```go
// go.mod (key deps)
// go.opentelemetry.io/otel
// go.opentelemetry.io/otel/sdk
// go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
```

```go
package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func initTracer(ctx context.Context) func(context.Context) error {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("dogbridge:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("checkout"),
			semconv.ServiceVersion("1.12.0"),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func main() {
	ctx := context.Background()
	shutdown := initTracer(ctx)
	defer shutdown(ctx)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	_ = http.ListenAndServe(":8080", otelhttp.NewHandler(h, "http.request"))
}
```

---

## Dashboard migration examples

### Example 1: Request rate panel

- **Datadog concept**: `trace.http.request.hits` by service/env.
- **OTel/backend target**: request count derived from spans or metrics by `service.name` and route.
- **Migration action**: standardize dashboard template variables to OTel attributes (`service.name`, `deployment.environment`).

### Example 2: p95 latency panel

- **Datadog concept**: APM latency percentiles.
- **OTel/backend target**: percentile query over span duration or RED metrics emitted from traces.
- **Migration action**: verify percentile math/windowing differences and annotate dashboards.

## Alert migration examples

### Error-rate alert

- **Datadog**: monitor on APM error percentage.
- **OTel/backend**: alert expression on `error_spans / total_spans` (or equivalent RED metric).
- **Guardrail**: run in shadow mode for one week before paging cutover.

### Latency SLO alert

- **Datadog**: p95 latency above threshold.
- **OTel/backend**: percentile-based or burn-rate alert depending on backend capabilities.
- **Guardrail**: compare false-positive/false-negative rates during overlap.

---

## Staged Rollout Playbook

1. **Single service pilot (Week 1)**
   - Choose one service with moderate traffic.
   - Validate ingestion, dashboards, and alerts in non-prod then prod.
2. **Dual-write window (Week 2-3)**
   - Keep old and new telemetry paths for fast rollback.
   - Track parity scorecard daily (trace volume, latency bands, error ratio).
3. **Cutover (Week 4)**
   - Switch paging/ops runbooks to backend-native dashboards.
   - Disable legacy Datadog SDK path.
4. **Decommission (Week 5+)**
   - Remove old libraries/env vars and archived dashboards.
   - Capture lessons learned and reusable templates.

## Rollback plan

- Keep previous deployment manifest that still enables legacy tracing path.
- Define objective rollback triggers (e.g., >10% trace drop, sustained alert blind spots).
- Timebox investigation before rollback (for example: 30 minutes during incident window).
