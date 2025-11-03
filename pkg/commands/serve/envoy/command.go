package envoy

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/envoy/authz-server"
	validationwebhook "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/envoy/validation-webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "envoy",
		Short: "Run Kyverno Envoy servers",
	}
	command.AddCommand(authzserver.Command())
	command.AddCommand(validationwebhook.Command())
	return command
}
