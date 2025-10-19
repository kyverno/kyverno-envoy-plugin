package sidecarinjector

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/sidecar"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/webhook/mutation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var address string
	var certFile string
	var keyFile string
	var configFile string
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Start the Kubernetes mutating webhook injecting Kyverno Authz Server sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// load sidecar
				sidecar, err := sidecar.Load(configFile)
				if err != nil {
					return err
				}
				// create server
				http := mutation.NewSidecarInjectorServer(address, certFile, keyFile, sidecar)
				// run server
				return http.Run(ctx)
			})
		},
	}
	command.Flags().StringVar(&address, "address", ":9443", "Address to listen on")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	command.Flags().StringVar(&configFile, "config-file", "", "File containing the sidecar config")
	return command
}
