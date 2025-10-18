package serve

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/authz-server"
	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/control-plane"
	sidecarinjector "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-injector"
	sidecarauthz "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-server"
	validationwebhook "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/validation-webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run Kyverno Envoy Plugin servers",
	}
	command.AddCommand(sidecarauthz.Command())
	command.AddCommand(sidecarinjector.Command())
	command.AddCommand(validationwebhook.Command())
	command.AddCommand(controlplane.Command())
	command.AddCommand(authzserver.Command())
	return command
}
