package distro

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpsprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
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

func Factories() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Receivers, err = otelcol.MakeFactoryMap(
		datadogreceiver.NewFactory(),
		statsdreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		otlpreceiver.NewFactory(),
	)
	if err != nil {
		return factories, err
	}

	factories.Processors, err = otelcol.MakeFactoryMap[processor.Factory](
		memorylimiterprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		k8sattributesprocessor.NewFactory(),
		transformprocessor.NewFactory(),
		batchprocessor.NewFactory(),
	)
	if err != nil {
		return factories, err
	}

	factories.Exporters, err = otelcol.MakeFactoryMap(
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		elasticsearchexporter.NewFactory(),
	)
	if err != nil {
		return factories, err
	}

	return factories, nil
}

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
}
