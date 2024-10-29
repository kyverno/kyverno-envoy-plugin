package authz

import (
	"context"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewGrpcServer(network, addr string, config *rest.Config) server.ServerFunc {
	return func(ctx context.Context) error {
		// create a server
		s := grpc.NewServer()
		// configure scheme
		scheme := runtime.NewScheme()
		if err := v1alpha1.Install(scheme); err != nil {
			return err
		}
		// create kubernetes client
		client, err := client.New(config, client.Options{
			Scheme: scheme,
		})
		if err != nil {
			return err
		}
		// setup our authorization service
		svc := &service{
			client: client,
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
