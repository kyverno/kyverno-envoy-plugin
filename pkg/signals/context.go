package signals

import (
	"context"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/util/wait"
)

func Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}

func Do(ctx context.Context, callback func(context.Context) error) error {
	// create a wait group
	var group wait.Group
	// wait all tasks in the group are over
	defer group.Wait()
	// create a signal aware context
	ctx, stop := Context(ctx)
	// cancel context and restore signals behaviour
	defer stop()
	// wait until context is cancelled or signals are triggered
	group.StartWithContext(ctx, func(ctx context.Context) {
		// restore signals behaviour (context has been cancelled at this point)
		defer stop()
		// wait signals are triggered
		<-ctx.Done()
	})
	// invoke callback with signals aware context
	return callback(ctx)
}
