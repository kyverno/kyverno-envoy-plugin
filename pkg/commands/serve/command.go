package serve

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/authz-server"
	sidecarinjector "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-injector"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run Kyverno Envoy Plugin servers",
	}
	command.AddCommand(authzserver.Command())
	command.AddCommand(sidecarinjector.Command())
	return command
}
