package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
)

func RunHttp(ctx context.Context, server *http.Server, certFile, keyFile string) error {
	logger := ctrl.LoggerFrom(ctx).WithValues("address", server.Addr)
	defer logger.Info("HTTP Server stopped")
	// track shutdown error
	var shutdownErr error
	// track serve error
	serveErr := func(ctx context.Context) error {
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
			logger.Info("HTTP Server shutting down...")
			// create a context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// gracefully shutdown server
			shutdownErr = server.Shutdown(ctx)
		})
		serve := func() error {
			logrus.Infof("HTTP Server starting at %s...\n", server.Addr)
			if certFile != "" && keyFile != "" {
				// server over https
				return server.ListenAndServeTLS(certFile, keyFile)
			} else {
				// server over http
				return server.ListenAndServe()
			}
		}
		// server closed is not an error
		if err := serve(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}(ctx)
	// return error if any
	return multierr.Combine(serveErr, shutdownErr)
}
