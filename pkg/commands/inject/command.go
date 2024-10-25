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
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	jsonpatch "gomodules.xyz/jsonpatch/v2"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
)

func Command() *cobra.Command {
	var address string
	var certFile string
	var keyFile string
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Responsible for injecting sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(context.Background(), address, certFile, keyFile)
		},
	}
	command.Flags().StringVar(&address, "address", ":9443", "Address to listen on")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	return command
}

func setupMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/livez", handlers.Health())
	mux.Handle("/readyz", handlers.Health())
	mux.Handle("/mutate", handlers.AdmissionReview(func(ctx context.Context, r *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
		response := func(err error, warnings ...string) *admissionv1.AdmissionResponse {
			response := admissionv1.AdmissionResponse{
				Allowed: err == nil,
				UID:     r.UID,
			}
			if err != nil {
				response.Result = &metav1.Status{
					Status:  metav1.StatusFailure,
					Message: err.Error(),
				}
			}
			response.Warnings = warnings
			return &response
		}
		var object unstructured.Unstructured
		if err := object.UnmarshalJSON(r.Object.Raw); err != nil {
			fmt.Println("object.UnmarshalJSON")
			return response(err)
		}
		var pod v1.Pod
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.UnstructuredContent(), &pod); err != nil {
			fmt.Println("FromUnstructured")
			return response(err)
		}
		pod.Spec.Containers = append(pod.Spec.Containers, v1.Container{
			Name:            "kyverno-envoy-plugin",
			ImagePullPolicy: v1.PullIfNotPresent,
			Image:           "foo",
			Args: []string{
				"serve",
			},
		})
		if data, err := json.Marshal(&pod); err != nil {
			fmt.Println("json.Marshal", string(data))
			return response(err)
		} else if patch, err := jsonpatch.CreatePatch(r.Object.Raw, data); err != nil {
			fmt.Println("jsonpatch.CreateMergePatch", string(data))
			return response(err)
		} else if patch, err := json.Marshal(patch); err != nil {
			fmt.Println("jsonpatch.CreateMergePatch", string(data))
			return response(err)
		} else {
			rspn := response(nil)
			rspn.PatchType = ptr.To(admissionv1.PatchTypeJSONPatch)
			rspn.Patch = patch
			fmt.Println("ok", string(patch))
			return rspn
		}
	}))
	return mux
}

func setupServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: setupMux(),
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

func runServer(ctx context.Context, addr, certFile, keyFile string) error {
	var group wait.Group
	server := setupServer(addr)
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
		fmt.Printf("Starting server at %s...\n", addr)
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
