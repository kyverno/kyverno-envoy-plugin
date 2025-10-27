package processor

import (
	"context"
	"sync"

	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils"
)

type policyAccessor[POLICY any] struct {
	compiler     engine.Compiler[POLICY]
	policies     map[string]POLICY
	sortPolicies func() []POLICY
	sync.Mutex
}

func NewPolicyAccessor[POLICY any](compiler engine.Compiler[POLICY]) *policyAccessor[POLICY] {
	return &policyAccessor[POLICY]{
		Mutex:    sync.Mutex{},
		compiler: compiler,
		policies: make(map[string]POLICY),
		sortPolicies: func() []POLICY {
			return nil
		},
	}
}

type Processor interface {
	Process(req *protov1alpha1.ValidatingPolicy)
}

func (p *policyAccessor[POLICY]) Process(req *protov1alpha1.ValidatingPolicy) {
	resetSortPolicies := func() {
		p.sortPolicies = sync.OnceValue(func() []POLICY {
			p.Lock()
			defer p.Unlock()
			return utils.ToSortedSlice(p.policies)
		})
	}
	if req.Delete {
		// p.logger.Info("deleting policy: ", req.Name)
		p.Lock()
		delete(p.policies, req.Name)
		p.Unlock()
		resetSortPolicies()
		return
	}

	vpol := controlplane.FromProto(req)
	compiledPolicy, err := p.compiler.Compile(vpol)
	if err != nil {
		// p.logger.Errorf("failed to compile policy %s: %s", req.Name, err)
		return
	}
	p.Lock()
	defer p.Unlock()
	p.policies[req.Name] = compiledPolicy
	resetSortPolicies()
}

func (p *policyAccessor[POLICY]) Load(ctx context.Context) ([]POLICY, error) {
	return p.sortPolicies(), nil
}
