package providers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type composite []engine.Source

func NewComposite(providers ...engine.Source) engine.Source {
	return composite(providers)
}

func (p composite) Load(ctx context.Context) ([]engine.CompiledPolicy, error) {
	var out []engine.CompiledPolicy
	for _, provider := range p {
		c, err := provider.Load(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, c...)
	}
	return out, nil
}
