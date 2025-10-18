package httpauth

import (
	"context"
	"net/http"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
)

func NewServer(addr string, a *authorizer) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		mux.Handle("POST /", http.HandlerFunc(a.NewHandler()))
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, "", "")
	}
}
