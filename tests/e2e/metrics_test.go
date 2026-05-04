package e2e

import (
	"os"
	"strings"
	"testing"
)

func TestMetricsPipelineVictoriaMetricsConfig(t *testing.T) {
	cfgBytes, err := os.ReadFile("../../config/examples/metrics-to-victoriametrics.yaml")
	if err != nil {
		t.Fatalf("read metrics example config: %v", err)
	}

	cfg := string(cfgBytes)
	requiredSnippets := []string{
		"statsd:",
		"prometheusremotewrite:",
		"transform/drop_default_high_cardinality:",
		"transform/label_controls:",
		"DOGBRIDGE_METRICS_LABEL_ALLOW_REGEX",
		"DOGBRIDGE_METRICS_LABEL_DENY_REGEX",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(cfg, snippet) {
			t.Fatalf("metrics example missing required snippet %q", snippet)
		}
	}
}
