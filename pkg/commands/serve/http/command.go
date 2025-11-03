package http

import (
	authzserver "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/http/authz-server"
	validationwebhook "github.com/kyverno/kyverno-envoy-plugin/pkg/commands/serve/http/validation-webhook"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "http",
		Short: "Run Kyverno HTTP servers",
	}
	command.AddCommand(authzserver.Command())
	command.AddCommand(validationwebhook.Command())
	return command
}
