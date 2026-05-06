package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	traceAgentURL := getenvDefault("DD_TRACE_AGENT_URL", "http://localhost:8126")
	dogstatsdAddr := getenvDefault("DD_DOGSTATSD_ADDR", "localhost:8125")
	otlpLogsEndpoint := getenvDefault("DOGBRIDGE_OTLP_LOGS_ENDPOINT", "http://localhost:4318/v1/logs")
	runID := getenvDefault("DOGBRIDGE_DEMO_RUN_ID", strconv.FormatInt(time.Now().UnixNano(), 10))
	maxBatches := getenvIntDefault("DOGBRIDGE_DEMO_BATCHES", 0)
	tick := 5 * time.Second

	tracer.Start(
		tracer.WithAgentAddr(traceAgentURL[len("http://"):]),
		tracer.WithService("dogbridge-dd-demo"),
		tracer.WithEnv("local"),
	)
	defer tracer.Stop()

	fmt.Printf("starting demo emitter: traces=%s dogstatsd=%s logs=%s interval=%s batches=%s run_id=%s\n", traceAgentURL, dogstatsdAddr, otlpLogsEndpoint, tick, batchLimitLabel(maxBatches), runID)

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for i := 1; maxBatches == 0 || i <= maxBatches; i++ {
		emitSpan(otlpLogsEndpoint, runID, i)
		emitDogstatsd(dogstatsdAddr, runID, i)
		fmt.Printf("emitted batch %d\n", i)
		if maxBatches > 0 && i == maxBatches {
			break
		}
		<-ticker.C
	}
}

func emitSpan(logsEndpoint, runID string, seq int) {
	root := tracer.StartSpan("demo.request", tracer.ResourceName("GET /dummy"))
	root.SetTag("component", "go-dd-demo")
	root.SetTag("batch.seq", seq)
	root.SetTag("demo.run_id", runID)
	time.Sleep(40 * time.Millisecond)

	child := tracer.StartSpan("demo.work", tracer.ChildOf(root.Context()))
	child.SetTag("worker", "dummy-app")
	child.SetTag("demo.run_id", runID)
	time.Sleep(20 * time.Millisecond)
	child.Finish()
	emitOTLPLog(logsEndpoint, root.Context(), runID, seq)
	root.Finish()
}

func emitDogstatsd(addr, runID string, seq int) {
	payload := []string{
		fmt.Sprintf("dogbridge.demo.requests_total:1|c|#env:local,service:dogbridge-dd-demo,batch:%d,run_id:%s", seq, runID),
		fmt.Sprintf("dogbridge.demo.queue_depth:%d|g|#env:local,service:dogbridge-dd-demo,run_id:%s", seq%10, runID),
		fmt.Sprintf("dogbridge.demo.request_ms:%d|h|#env:local,service:dogbridge-dd-demo,run_id:%s", 60+(seq%7), runID),
		fmt.Sprintf("_e{18,26}:dogbridge dummy log|dummy app emitted batch=%d|t:info|#env:local,service:dogbridge-dd-demo,run_id:%s", seq, runID),
	}

	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("dogstatsd dial error: %v\n", err)
		return
	}
	defer conn.Close()

	for _, line := range payload {
		_, _ = conn.Write([]byte(line))
	}
}

func emitOTLPLog(endpoint string, spanContext spanContextIDs, runID string, seq int) {
	now := time.Now().UnixNano()
	body := fmt.Sprintf(`{
		"resourceLogs": [{
			"resource": {
				"attributes": [
					{"key": "service.name", "value": {"stringValue": "dogbridge-dd-demo"}},
					{"key": "deployment.environment", "value": {"stringValue": "local"}},
					{"key": "telemetry.gateway", "value": {"stringValue": "dogbridge"}}
				]
			},
			"scopeLogs": [{
				"scope": {"name": "examples/go-ddtrace"},
				"logRecords": [{
					"timeUnixNano": "%d",
					"observedTimeUnixNano": "%d",
					"traceId": "%s",
					"spanId": "%s",
					"flags": 1,
					"severityText": "INFO",
					"severityNumber": 9,
					"body": {"stringValue": "dummy app emitted batch=%d"},
					"attributes": [
						{"key": "component", "value": {"stringValue": "go-dd-demo"}},
						{"key": "batch.seq", "value": {"intValue": "%d"}},
						{"key": "demo.run_id", "value": {"stringValue": "%s"}},
						{"key": "datadog.trace_id", "value": {"stringValue": "%s"}},
						{"key": "datadog.span_id", "value": {"stringValue": "%s"}}
					]
				}]
			}]
		}]
	}`, now, now, otelTraceID(spanContext), otelSpanID(spanContext), seq, seq, runID, strconv.FormatUint(spanContext.TraceID(), 10), strconv.FormatUint(spanContext.SpanID(), 10))

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(body))
	if err != nil {
		fmt.Printf("otlp log request build error: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("otlp log post error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("otlp log post status: %s\n", resp.Status)
	}
}

type spanContextIDs interface {
	TraceID() uint64
	SpanID() uint64
}

func getenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getenvIntDefault(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

func otelTraceID(spanContext spanContextIDs) string {
	// Match the OpenTelemetry Datadog receiver's default 64-bit Datadog trace ID mapping.
	return fmt.Sprintf("%016x%016x", uint64(0), spanContext.TraceID())
}

func otelSpanID(spanContext spanContextIDs) string {
	return fmt.Sprintf("%016x", spanContext.SpanID())
}

func batchLimitLabel(maxBatches int) string {
	if maxBatches == 0 {
		return "unbounded"
	}
	return strconv.Itoa(maxBatches)
}
