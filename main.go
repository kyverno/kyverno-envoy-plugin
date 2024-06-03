package main

import (
	"fmt"
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/spf13/cobra"
)

var policies []string

func init() {
	serveCmd.Flags().StringSliceVar(&policies, "policy", nil, "Path to kyverno-json policies")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the kyverno-envoy-plugin server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.NewServers(policies)
		server.StartServers(srv)
	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kyverno-envoy-plugin",
		Short: "kyverno-envoy-plugin is a plugin for Envoy",
	}

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
