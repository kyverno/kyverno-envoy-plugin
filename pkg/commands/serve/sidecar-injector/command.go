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
	var controlPlaneAddr string
	var controlPlaneReconnectWait, controlPlaneMaxDialInterval, healthCheckInterval string
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Start the Kubernetes mutating webhook injecting Kyverno Authorizer sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// create server
				http := mutation.NewSidecarInjectorServer(address, sidecarImage, controlPlaneAddr, certFile, keyFile,
					controlPlaneReconnectWait, controlPlaneMaxDialInterval, healthCheckInterval)
				// run server
				return http.Run(ctx)
			})
		},
	}
	command.Flags().StringVar(&address, "address", ":9443", "Address to listen on")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	command.Flags().StringVar(&sidecarImage, "sidecar-image", "", "Image to use in sidecar")
	command.Flags().StringVar(&controlPlaneAddr, "control-plane-address", "", "The control plane address to inject into the sidecars")
	command.Flags().StringVar(&controlPlaneReconnectWait, "control-plane-reconnect-wait", "3s", "Duration to wait before retrying connecting to the control plane")
	command.Flags().StringVar(&controlPlaneMaxDialInterval, "control-plane-max-dial-interval", "8s", "Duration to wait before stopping attempts of sending a policy to a client")
	command.Flags().StringVar(&healthCheckInterval, "health-check-interval", "30s", "Interval for sending health checks")

	_ = command.MarkFlagRequired("control-plane-address")
	return command
}
