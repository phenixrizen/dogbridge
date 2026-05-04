package distro

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpsprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/processor/resourceprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
)

var BuildInfo = component.BuildInfo{
	Command:     "dogbridge",
	Description: "Datadog-compatible OpenTelemetry Collector distribution",
	Version:     "0.1.0",
}

var Factories, _ = otelcol.MakeFactoryMap(
	datadogreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	k8sattributesprocessor.NewFactory(),
	transformprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	otlpexporter.NewFactory(),
	prometheusremotewriteexporter.NewFactory(),
	lokiexporter.NewFactory(),
	elasticsearchexporter.NewFactory(),
)

var ResolverSettings = confmap.ResolverSettings{
	URIs: []string{"env:", "file:"},
	ProviderFactories: []confmap.ProviderFactory{
		envprovider.NewFactory(),
		fileprovider.NewFactory(),
		httpprovider.NewFactory(),
		httpsprovider.NewFactory(),
		yamlprovider.NewFactory(),
	},
	ConverterFactories: []confmap.ConverterFactory{},
	DefaultScheme:      "env",
	ProviderSettings:   confmap.ProviderSettings{},
	Watcher:            nil,
	DisableURIWarning:  true,
	ErrorHandler:       consumererror.NewNoop(),
}
