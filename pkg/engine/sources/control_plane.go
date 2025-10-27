package sources

import (
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
)

func NewControlPlane[DATA, IN, OUT any](
	compiler engine.Compiler[DATA, IN, OUT],
) (core.Source[policy.Policy[DATA, IN, OUT]], error) {
	listener := NewListener()
	transform := sources.NewTransformErr(listener, func(in *v1alpha1.ValidatingPolicy) (policy.Policy[DATA, IN, OUT], error) {
		policy, err := compiler.Compile(in)
		return policy, err.ToAggregate()
	})
	// TODO: cache
	return transform, nil
}
