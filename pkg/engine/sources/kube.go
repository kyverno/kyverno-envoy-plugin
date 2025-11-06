package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	controllerruntime "github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/controller-runtime"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

func NewKube[POLICY any](name string, mgr ctrl.Manager, compiler engine.Compiler[POLICY]) (core.Source[POLICY], error) {
	options := controller.Options{
		NeedLeaderElection: ptr.To(false),
	}
	apis, err := controllerruntime.NewApiSource[v1alpha1.ValidatingPolicy](name, mgr, options)
	if err != nil {
		return nil, err

	}
	cache := sources.NewCache(
		apis,
		func(_ context.Context, in *v1alpha1.ValidatingPolicy) (string, error) {
			return in.Name + in.ResourceVersion, nil
		},
		func(_ context.Context, _ string, in *v1alpha1.ValidatingPolicy) (POLICY, error) {
			policy, err := compiler.Compile(in)
			return policy, err.ToAggregate()
		},
	)
	return cache, nil
}
