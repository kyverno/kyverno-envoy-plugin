package serve

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var httpAddress string
	var grpcAddress string
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the kyverno-envoy-plugin server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			srv := server.NewServers(grpcAddress, httpAddress)
			server.StartServers(ctx, srv)
		},
	}
	command.Flags().StringVar(&httpAddress, "http-address", "", "Address to listen on")
	command.Flags().StringVar(&grpcAddress, "grpc-address", "", "Address to listen on for health checks")
	return command
}
