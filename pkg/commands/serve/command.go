package serve

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Command() *cobra.Command {
	var httpAddress string
	var grpcAddress string
	var grpcNetwork string
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the kyverno-envoy-plugin server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var httpErr, grpcErr error
				func(ctx context.Context) {
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// create http and grpc servers
					http := authz.NewHttpServer(httpAddress)
					grpc := authz.NewGrpcServer(grpcNetwork, grpcAddress)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						httpErr = http.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						grpcErr = grpc.Run(ctx)
					})
				}(ctx)
				return multierr.Combine(httpErr, grpcErr)
			})
		},
	}
	command.Flags().StringVar(&httpAddress, "http-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&grpcAddress, "grpc-address", ":9081", "Address to listen on")
	command.Flags().StringVar(&grpcNetwork, "grpc-network", "tcp", "Network to listen on")
	return command
}
