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
	"go.opentelemetry.io/collector/otelcol"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/otlpexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/memorylimiterprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

var BuildInfo = component.BuildInfo{
	Command:     "dogbridge",
	Description: "Datadog-compatible OpenTelemetry Collector distribution",
	Version:     "0.1.0",
}

var Factories, _ = otelcol.MakeFactoryMap(
	datadogreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	transformprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	otlpexporter.NewFactory(),
	prometheusremotewriteexporter.NewFactory(),
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
