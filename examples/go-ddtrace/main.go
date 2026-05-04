package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	agentURL := os.Getenv("DD_TRACE_AGENT_URL")
	if agentURL == "" {
		agentURL = "http://localhost:8126"
	}

	tracer.Start(
		tracer.WithAgentAddr(agentURL[len("http://"):]),
		tracer.WithService("dogbridge-ddtrace-demo"),
		tracer.WithEnv("local"),
	)
	defer tracer.Stop()

	root := tracer.StartSpan("demo.request")
	root.SetTag("component", "go-ddtrace-example")
	time.Sleep(100 * time.Millisecond)

	child := tracer.StartSpan("demo.http", tracer.ChildOf(root.Context()))
	resp, err := http.Get("https://example.com")
	if err == nil {
		_ = resp.Body.Close()
		child.SetTag("http.status_code", resp.StatusCode)
	}
	child.Finish()
	root.Finish()

	fmt.Println("sent one dd-trace span tree to", agentURL)
}
