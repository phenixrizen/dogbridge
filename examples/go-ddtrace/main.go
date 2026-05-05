package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	traceAgentURL := getenvDefault("DD_TRACE_AGENT_URL", "http://localhost:8126")
	dogstatsdAddr := getenvDefault("DD_DOGSTATSD_ADDR", "localhost:8125")
	tick := 5 * time.Second

	tracer.Start(
		tracer.WithAgentAddr(traceAgentURL[len("http://"):]),
		tracer.WithService("dogbridge-dd-demo"),
		tracer.WithEnv("local"),
	)
	defer tracer.Stop()

	fmt.Printf("starting demo emitter: traces=%s dogstatsd=%s interval=%s\n", traceAgentURL, dogstatsdAddr, tick)

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for i := 1; ; i++ {
		emitSpan(i)
		emitDogstatsd(dogstatsdAddr, i)
		fmt.Printf("emitted batch %d\n", i)
		<-ticker.C
	}
}

func emitSpan(seq int) {
	root := tracer.StartSpan("demo.request", tracer.ResourceName("GET /dummy"))
	root.SetTag("component", "go-dd-demo")
	root.SetTag("batch.seq", seq)
	time.Sleep(40 * time.Millisecond)

	child := tracer.StartSpan("demo.work", tracer.ChildOf(root.Context()))
	child.SetTag("worker", "dummy-app")
	time.Sleep(20 * time.Millisecond)
	child.Finish()
	root.Finish()
}

func emitDogstatsd(addr string, seq int) {
	payload := []string{
		fmt.Sprintf("dogbridge.demo.requests_total:1|c|#env:local,service:dogbridge-dd-demo,batch:%d", seq),
		fmt.Sprintf("dogbridge.demo.queue_depth:%d|g|#env:local,service:dogbridge-dd-demo", seq%10),
		fmt.Sprintf("dogbridge.demo.request_ms:%d|h|#env:local,service:dogbridge-dd-demo", 60+(seq%7)),
		fmt.Sprintf("_e{18,26}:dogbridge dummy log|dummy app emitted batch=%d|t:info|#env:local,service:dogbridge-dd-demo", seq),
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

func getenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
