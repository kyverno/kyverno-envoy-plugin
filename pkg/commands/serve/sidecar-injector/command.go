package sidecarinjector

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/webhook/mutation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var address string
	var certFile string
	var keyFile string
	var sidecarImage string
	var externalPolicySources []string
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Start the Kubernetes mutating webhook injecting Kyverno Authz Server sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// create server
				http := mutation.NewSidecarInjectorServer(address, sidecarImage, certFile, keyFile, externalPolicySources...)
				// run server
				return http.Run(ctx)
			})
		},
	}
	command.Flags().StringVar(&address, "address", ":9443", "Address to listen on")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	command.Flags().StringVar(&sidecarImage, "sidecar-image", "", "Image to use in sidecar")
	command.Flags().StringArrayVar(&externalPolicySources, "external-policy-source", nil, "External policy sources")
	return command
}
