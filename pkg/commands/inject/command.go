package inject

import (
	"fmt"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/httpd"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var httpdConf httpd.SimpleServer
	command := &cobra.Command{
		Use:   "sidecar-injector",
		Short: "Responsible for injecting sidecars into pod containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("SimpleServer starting to listen in port %v", httpdConf.Port)
			return httpdConf.Start()
		},
	}
	command.Flags().IntVar(&httpdConf.Port, "port", 443, "server port.")
	command.Flags().StringVar(&httpdConf.CertFile, "certFile", "/etc/mutator/certs/tls.crt", "File containing tls certificate")
	command.Flags().StringVar(&httpdConf.KeyFile, "keyFile", "/etc/mutator/certs/tls.key", "File containing tls private key")
	command.Flags().BoolVar(&httpdConf.Local, "local", false, "Local run mode")
	command.Flags().StringVar(&(&httpdConf.Patcher).SidecarDataKey, "sidecarDataKey", "sidecars.yaml", "ConfigMap Sidecar Data Key")
	return command
}
