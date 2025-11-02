package http

import (
	"context"
	"net/http"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"k8s.io/client-go/dynamic"
)

func NewServer(config Config, p engine.HTTPSource, dyn dynamic.Interface) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		a := &authorizer{
			provider:      p,
			dyn:           dyn,
			nestedRequest: config.NestedRequest,
		}
		mux.Handle("POST /", a)
		// create server
		s := &http.Server{
			Addr:    config.Address,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, "", "")
	}
}
