package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
)

func NewControlPlane[POLICY any](
	compiler engine.Compiler[POLICY],
) (core.Source[POLICY], error) {
	listener := NewListener()
	cache := sources.NewCache(
		listener,
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
