package main

import (
	"fmt"
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	root := setup()
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setup() *cobra.Command {
	var policies []string
	var address string
	var healthaddress string
	serve := &cobra.Command{
		Use:   "serve",
		Short: "Start the kyverno-envoy-plugin server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewServers(policies, address, healthaddress)
			server.StartServers(srv)
		},
	}
	serve.Flags().StringSliceVar(&policies, "policy", nil, "Path to kyverno-json policies")
	serve.Flags().StringVar(&address, "address", ":9000", "Address to listen on")
	serve.Flags().StringVar(&healthaddress, "healthaddress", ":8181", "Address to listen on for health checks")
	root := &cobra.Command{
		Use:   "kyverno-envoy-plugin",
		Short: "kyverno-envoy-plugin is a plugin for Envoy",
	}
	root.AddCommand(serve)
	return root
}
