package policy

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
)

type staticProvider struct {
	compiled []CompiledPolicy
	err      error
}

func NewStaticProvider(compiler Compiler, policies ...v1alpha1.AuthorizationPolicy) Provider {
	var compiled []CompiledPolicy
	for _, policy := range policies {
		policy, err := compiler.Compile(&policy)
		if err != nil {
			return &staticProvider{err: err.ToAggregate()}
		}
		compiled = append(compiled, policy)
	}
	return &staticProvider{compiled: compiled}
}

func (p *staticProvider) CompiledPolicies(ctx context.Context) ([]CompiledPolicy, error) {
	// TODO: sort based on policy names
	return p.compiled, p.err
}
