package envoy

import (
	"context"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/dispatchers"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/handlers"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/client-go/dynamic"
)

func NewServer(network, addr string, source engine.EnvoySource, dynclient dynamic.Interface) server.ServerFunc {
	return func(ctx context.Context) error {
		// create a server
		s := grpc.NewServer()
		// build the engine
		engine := core.NewEngine(
			source,
			handlers.Handler(
				dispatchers.Sequential(
					policy.EvaluatorFactory[engine.EnvoyPolicy](),
					func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Breaker[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse]] {
						return core.MakeBreakerFunc(func(_ context.Context, _ engine.EnvoyPolicy, _ *authv3.CheckRequest, out policy.Evaluation[*authv3.CheckResponse]) bool {
							return out.Result != nil
						})
					},
				),
				func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Resulter[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse], policy.Evaluation[*authv3.CheckResponse]] {
					return resulters.NewFirst[engine.EnvoyPolicy, *authv3.CheckRequest](func(out policy.Evaluation[*authv3.CheckResponse]) bool {
						return out.Result != nil || out.Error != nil
					})
				},
			),
		)
		// setup our authorization service
		svc := &service{
			engine:    engine,
			dynclient: dynclient,
		}
		// register our authorization service
		authv3.RegisterAuthorizationServer(s, svc)
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
