package serve

import (
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
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
	return serve
}
