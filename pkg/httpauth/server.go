package httpauth

import (
	"context"
	"net/http"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
)

func NewServer(addr string, provider engine.Provider) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		a := Authorizer{
			provider: provider,
		}

		mux.Handle("POST /authorize", http.HandlerFunc(a.NewHandler()))
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, "", "")
	}
}
