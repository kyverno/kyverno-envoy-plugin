package server

import (
	"context"
	"net"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/logging"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
)

var grpcLogger = logging.WithName("grpc-server")

func RunGrpc(ctx context.Context, server *grpc.Server, listener net.Listener) error {
	defer grpcLogger.Info("GRPC Server stopped")
	// create a wait group
	var group wait.Group
	// wait all tasks in the group are over
	defer group.Wait()
	// create a cancellable context
	ctx, cancel := context.WithCancel(ctx)
	// cancel context at the end
	defer cancel()
	// shutdown server when context is cancelled
	group.StartWithContext(ctx, func(ctx context.Context) {
		// wait context cancelled
		<-ctx.Done()
		grpcLogger.Info("GRPC Server shutting down")
		// gracefully shutdown server
		server.GracefulStop()
	})
	grpcLogger.Info("GRPC Server started ", "listenerAddress", listener.Addr())
	// serve
	return server.Serve(listener)
}
