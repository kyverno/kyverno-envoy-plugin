package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
}
