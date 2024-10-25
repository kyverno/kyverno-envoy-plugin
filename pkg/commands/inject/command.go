package inject

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server/handlers"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/sidecar"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	jsonpatch "gomodules.xyz/jsonpatch/v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Command() *cobra.Command {
	var address string
	var certFile string
	var keyFile string
	var sidecarImage string
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Responsible for injecting sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			mux := setupMux(sidecarImage)
			server := setupServer(address, mux)
			return runServer(context.Background(), server, certFile, keyFile)
		},
	}
	command.Flags().StringVar(&address, "address", ":9443", "Address to listen on")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	command.Flags().StringVar(&sidecarImage, "sidecar-image", "", "Image to use in sidecar")
	return command
}

func setupMux(sidecarImage string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/livez", handlers.Healthy(handlers.True))
	mux.Handle("/readyz", handlers.Ready(handlers.True))
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
	return mux
}

func setupServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
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
}

func runServer(ctx context.Context, server *http.Server, certFile, keyFile string) error {
	var group wait.Group
	err := func() error {
		signalsCtx, signalsCancel := signals.Context(ctx)
		defer signalsCancel()
		var shutdownErr error
		group.StartWithContext(signalsCtx, func(ctx context.Context) {
			<-ctx.Done()
			fmt.Println("Shutting down server...")
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			shutdownErr = server.Shutdown(shutdownCtx)
		})
		fmt.Printf("Starting server at %s...\n", server.Addr)
		var serveErr error
		if certFile != "" && keyFile != "" {
			serveErr = server.ListenAndServeTLS(certFile, keyFile)
		} else {
			serveErr = server.ListenAndServe()
		}
		if errors.Is(serveErr, http.ErrServerClosed) {
			serveErr = nil
		}
		return multierr.Combine(serveErr, shutdownErr)
	}()
	group.Wait()
	fmt.Println("Server stopped")
	return err
}
