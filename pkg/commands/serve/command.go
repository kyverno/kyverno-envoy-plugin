package serve

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var policies []string
	var address string
	var healthaddress string
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the kyverno-envoy-plugin server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			srv := server.NewServers(policies, address, healthaddress)
			server.StartServers(ctx, srv)
		},
	}
	command.Flags().StringSliceVar(&policies, "policy", nil, "Path to kyverno-json policies")
	command.Flags().StringVar(&address, "address", ":9000", "Address to listen on")
	command.Flags().StringVar(&healthaddress, "healthaddress", ":8181", "Address to listen on for health checks")
	return command
}
