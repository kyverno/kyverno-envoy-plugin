package httpauth

import (
	"context"
	"net/http"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"k8s.io/client-go/dynamic"
)

func NewServer(addr string, dyn dynamic.Interface, p engine.HTTPSource, nestedRequest bool) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		a := &authorizer{
			provider:      p,
			dyn:           dyn,
			nestedRequest: nestedRequest,
		}
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
