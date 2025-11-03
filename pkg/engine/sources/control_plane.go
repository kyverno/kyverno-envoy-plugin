package sources

import (
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
)

func NewControlPlane[POLICY any](
	compiler engine.Compiler[POLICY],
) (core.Source[POLICY], error) {
	listener := NewListener("")
	transform := sources.NewTransformErr(listener, func(in *v1alpha1.ValidatingPolicy) (POLICY, error) {
		policy, err := compiler.Compile(in)
		return policy, err.ToAggregate()
	})
	// TODO: cache
	return transform, nil
}
