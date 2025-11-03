package serve

import (
	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/control-plane"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/http"
	sidecarinjector "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/sidecar-injector"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run Kyverno Authz servers",
	}
	command.AddCommand(envoy.Command())
	command.AddCommand(http.Command())
	command.AddCommand(controlplane.Command())
	command.AddCommand(sidecarinjector.Command())
	return command
}
