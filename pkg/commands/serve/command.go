package serve

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/authz-server"
	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/control-plane"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/dummy"
	sidecarinjector "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-injector"
	validationwebhook "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/validation-webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run Kyverno Envoy Plugin servers",
	}
	command.AddCommand(sidecarinjector.Command())
	command.AddCommand(validationwebhook.Command())
	command.AddCommand(controlplane.Command())
	command.AddCommand(authzserver.Command())
	command.AddCommand(dummy.Command())
	return command
}
