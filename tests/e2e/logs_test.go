package e2e

import (
	"os"
	"strings"
	"testing"
)

func TestLogsPipelineLokiConfig(t *testing.T) {
	cfgBytes, err := os.ReadFile("../../config/examples/logs-to-loki.yaml")
	if err != nil {
		t.Fatalf("read loki logs config: %v", err)
	}

	cfg := string(cfgBytes)
	requiredSnippets := []string{
		"filelog:",
		"type: trace_parser",
		"k8sattributes:",
		"k8s.namespace.name",
		"otlphttp/loki:",
		"endpoint: http://loki.monitoring.svc.cluster.local:3100/otlp",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(cfg, snippet) {
			t.Fatalf("logs-to-loki config missing required snippet %q", snippet)
		}
	}
}

func TestLogsPipelineOpenSearchConfig(t *testing.T) {
	cfgBytes, err := os.ReadFile("../../config/examples/logs-to-opensearch.yaml")
	if err != nil {
		t.Fatalf("read opensearch logs config: %v", err)
	}

	cfg := string(cfgBytes)
	requiredSnippets := []string{
		"filelog:",
		"type: trace_parser",
		"k8sattributes:",
		"elasticsearch:",
		"logs_index: dogbridge-logs",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(cfg, snippet) {
			t.Fatalf("logs-to-opensearch config missing required snippet %q", snippet)
		}
	}
}
