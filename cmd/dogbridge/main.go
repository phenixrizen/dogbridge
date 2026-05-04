package main

import (
	"log"
	"os"

	"github.com/your-org/dogbridge/distro"
	"go.opentelemetry.io/collector/otelcol"
)

func main() {
	info := otelcol.CollectorSettings{
		BuildInfo: distro.BuildInfo,
		Factories: distro.Factories,
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: distro.ResolverSettings,
		},
	}

	cmd := otelcol.NewCommand(info)
	if err := cmd.Execute(); err != nil {
		log.Printf("dogbridge collector exited with error: %v", err)
		os.Exit(1)
	}
}
