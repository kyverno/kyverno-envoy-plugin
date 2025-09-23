package serve

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/authz-server"
	sidecarinjector "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-injector"
	validationwebhook "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/validation-webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "serve",
		Short:        "Run Kyverno Envoy Plugin servers",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(
		authzserver.Command(),
		sidecarinjector.Command(),
		validationwebhook.Command(),
	)
	return command
}
