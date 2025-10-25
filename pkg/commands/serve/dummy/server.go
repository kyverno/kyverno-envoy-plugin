package dummy

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"k8s.io/apimachinery/pkg/util/wait"
)

type server struct {
	cancel      context.CancelFunc
	group       wait.Group
	grpcNetwork string
	grpcAddress string
}

func (s *server) Start(ctx context.Context) error {
	// envoyCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
	grpc := authz.NewServer(s.grpcNetwork, s.grpcAddress, nil, nil, nil)

	s.group.StartWithContext(ctx, func(ctx context.Context) {
		grpc.Run(ctx)
	})
	return nil
}

func (s *server) Stop() error {
	s.cancel()
	s.group.Wait()
	return nil
}
