package providers

import (
	"context"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type onceProvider struct {
	called   bool
	lock     sync.Mutex
	inner    engine.Source
	policies []engine.CompiledPolicy
	err      error
}

func NewOnceProvider(inner engine.Source) engine.Source {
	return &onceProvider{
		inner: inner,
	}
}

func (p *onceProvider) Load(ctx context.Context) ([]engine.CompiledPolicy, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.called {
		p.policies, p.err = p.inner.Load(ctx)
		p.called = true
	}
	return p.policies, p.err
}
