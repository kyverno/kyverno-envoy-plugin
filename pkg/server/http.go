package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/logging"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/wait"
)

var httpLogger = logging.WithName("http-server")

func RunHttp(ctx context.Context, server *http.Server, certFile, keyFile string) error {
	defer httpLogger.Info("HTTP Server stopped")
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
			httpLogger.Info("HTTP Server shutting down...")
			// create a context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// gracefully shutdown server
			shutdownErr = server.Shutdown(ctx)
		})
		serve := func() error {
			httpLogger.Info("HTTP Server starting ...\n", "serverAddress", server.Addr)
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
