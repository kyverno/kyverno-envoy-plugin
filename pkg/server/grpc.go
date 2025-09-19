package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
)

func RunGrpc(ctx context.Context, server *grpc.Server, listener net.Listener) error {
	defer ctrl.LoggerFrom(ctx).Info("GRPC Server stopped")
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
		ctrl.LoggerFrom(ctx).Info("GRPC Server shutting down...")
		// gracefully shutdown server
		server.GracefulStop()
	})
	ctrl.LoggerFrom(ctx).Info("GRPC Server starting...", "address", listener.Addr())
	// serve
	return server.Serve(listener)
}
