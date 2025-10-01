package providers

import (
	"context"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type onceProvider struct {
	called   bool
	lock     sync.Mutex
	inner    engine.Provider
	policies []engine.CompiledPolicy
	err      error
}

func NewOnceProvider(inner engine.Provider) engine.Provider {
	return &onceProvider{
		inner: inner,
	}
}

func (p *onceProvider) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.called {
		p.policies, p.err = p.inner.CompiledPolicies(ctx)
		p.called = true
	}
	return p.policies, p.err
}
