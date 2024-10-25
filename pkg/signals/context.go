package signals

import (
	"context"
	"os/signal"
	"syscall"
)

func Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}
