package mutation

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server/handlers"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/sidecar"
	"gomodules.xyz/jsonpatch/v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

func NewSidecarInjectorServer(addr, sidecarImage, certFile, keyFile string) server.ServerFunc {
	return func(ctx context.Context) error {
		// create mux
		mux := http.NewServeMux()
		// register health check
		mux.Handle("/livez", handlers.Healthy(handlers.True))
		// register ready check
		mux.Handle("/readyz", handlers.Ready(handlers.True))
		// register mutation webhook
		mux.Handle("/mutate", handlers.AdmissionReview(func(ctx context.Context, r *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
			var pod corev1.Pod
			if err := json.Unmarshal(r.Object.Raw, &pod); err != nil {
				return handlers.AdmissionResponse(r, err)
			}
			pod = sidecar.Inject(pod, sidecar.Sidecar(sidecarImage))
			if data, err := json.Marshal(&pod); err != nil {
				return handlers.AdmissionResponse(r, err)
			} else if patch, err := jsonpatch.CreatePatch(r.Object.Raw, data); err != nil {
				return handlers.AdmissionResponse(r, err)
			} else {
				return handlers.AdmissionResponse(r, nil, patch...)
			}
		}))
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
			TLSConfig: &tls.Config{
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
			},
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			IdleTimeout:       5 * time.Minute,
		}
		// run server
		return server.RunHttp(ctx, s, certFile, keyFile)
	}
}
