package root

import (
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/run"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:          "kyverno-envoy-plugin",
		Short:        "kyverno-envoy-plugin is a plugin for Envoy",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.AddCommand(
		run.Command(),
		serve.Command(),
		version.Command(),
	)
	return root
}
