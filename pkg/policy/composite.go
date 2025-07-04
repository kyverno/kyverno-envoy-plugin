package policy

import "context"

type composite []Provider

func NewComposite(providers ...Provider) Provider {
	return composite(providers)
}

func (p composite) CompiledPolicies(ctx context.Context) ([]CompiledPolicy, error) {
	var out []CompiledPolicy
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
