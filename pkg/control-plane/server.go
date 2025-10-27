package controlplane

import (
	"context"
	"net"

	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(network, addr string, srv protov1alpha1.ValidatingPolicyServiceServer) server.ServerFunc {
	return func(ctx context.Context) error {
		// create a server
		s := grpc.NewServer()
		// register validating policy service
		if srv != nil {
			protov1alpha1.RegisterValidatingPolicyServiceServer(s, srv)
		}
		// register reflection service
		reflection.Register(s)
		// create a listener
		l, err := net.Listen(network, addr)
		if err != nil {
			return err
		}
		// run server
		return server.RunGrpc(ctx, s, l)
	}
}
