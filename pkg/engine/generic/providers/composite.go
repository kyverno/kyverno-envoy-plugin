package providers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type composite []engine.Provider

func NewComposite(providers ...engine.Provider) engine.Provider {
	return composite(providers)
}

func (p composite) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	var out []engine.CompiledPolicy
	for _, provider := range p {
		c, err := provider.CompiledPolicies(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, c...)
	}
	// TODO: we probably need to sort policies before returning
	return out, nil
}
