package providers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type staticProvider struct {
	compiled []engine.CompiledPolicy
	err      error
}

func NewStaticProvider[T any](compiler engine.Compiler[T], policies ...T) engine.Provider {
	var compiled []engine.CompiledPolicy
	for _, policy := range policies {
		policy, err := compiler.Compile(policy)
		if err != nil {
			return &staticProvider{err: err.ToAggregate()}
		}
		compiled = append(compiled, policy)
	}
	return &staticProvider{compiled: compiled}
}

func (p *staticProvider) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	// TODO: sort based on policy names
	return p.compiled, p.err
}
