package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/dogbridge/distro"
	"go.opentelemetry.io/collector/otelcol"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "dogbridge CLI error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var cfgFiles []string
	var setValues []string

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the dogbridge collector",
		RunE: func(cmd *cobra.Command, args []string) error {
			uris := make([]string, 0, len(cfgFiles)+len(setValues))
			for _, c := range cfgFiles {
				uris = append(uris, "file:"+c)
			}
			for _, s := range setValues {
				uris = append(uris, "yaml:"+s)
			}

			resolver := distro.ResolverSettings
			if len(uris) > 0 {
				resolver.URIs = uris
			}

			info := otelcol.CollectorSettings{
				BuildInfo: distro.BuildInfo,
				Factories: distro.Factories,
				ConfigProviderSettings: otelcol.ConfigProviderSettings{
					ResolverSettings: resolver,
				},
			}

			return otelcol.NewCommand(info).Execute()
		},
	}
	runCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", nil, "Collector config file(s)")
	runCmd.Flags().StringSliceVar(&setValues, "set", nil, "Inline config override(s) in key::value YAML form")

	root := &cobra.Command{
		Use:   "dogbridge",
		Short: "Datadog-compatible OpenTelemetry Collector distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmd.RunE(cmd, args)
		},
	}
	root.AddCommand(runCmd)
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print dogbridge version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", distro.BuildInfo.Version)
		},
	})

	return root
}
