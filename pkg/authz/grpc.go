package authz

import (
	"context"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/policy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"google.golang.org/grpc"
)

func NewServer(network, addr string, provider policy.Provider) server.ServerFunc {
	return func(ctx context.Context) error {
		// create a server
		s := grpc.NewServer()
		// setup our authorization service
		svc := &service{
			provider: provider,
		}
		// register our authorization service
		authv3.RegisterAuthorizationServer(s, svc)
		// create a listener
		l, err := net.Listen(network, addr)
		if err != nil {
			return err
		}
		// run server
		return server.RunGrpc(ctx, s, l)
	}
}
