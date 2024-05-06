package main

import (
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/sidecar-injector/pkg/httpd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	httpdConf httpd.SimpleServer
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:   "sidecar-injector",
	Short: "Responsible for injecting sidecars into pod containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		log.Infof("SimpleServer starting to listen in port %v", httpdConf.Port)
		return httpdConf.Start()
	},
}

// Execute Kicks off the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVar(&httpdConf.Port, "port", 443, "server port.")
	rootCmd.Flags().StringVar(&httpdConf.CertFile, "certFile", "/etc/mutator/certs/tls.crt", "File containing tls certificate")
	rootCmd.Flags().StringVar(&httpdConf.KeyFile, "keyFile", "/etc/mutator/certs/tls.key", "File containing tls private key")
	rootCmd.Flags().BoolVar(&httpdConf.Local, "local", false, "Local run mode")
	rootCmd.Flags().StringVar(&(&httpdConf.Patcher).SidecarDataKey, "sidecarDataKey", "sidecars.yaml", "ConfigMap Sidecar Data Key")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "enable debug logs")
}

func main() {
	Execute()
}
