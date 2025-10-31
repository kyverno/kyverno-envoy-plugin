package http

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"k8s.io/client-go/dynamic"
)

func NewServer(addr string, dyn dynamic.Interface, p engine.HTTPSource, nestedRequest bool, certFile, keyFile string) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		a := &authorizer{
			provider:      p,
			dyn:           dyn,
			nestedRequest: nestedRequest,
		}
		mux.Handle("POST /", a)
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// serve TLS if a certfile and a keyfile are provided
		if certFile != "" && keyFile != "" {
			s.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					// AEADs w/ ECDHE
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			}
		}
		// run server
		return server.RunHttp(ctx, s, certFile, keyFile)
	}
}
