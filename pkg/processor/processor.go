package processor

import (
	"context"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
)

type policyAccessor[DATA, IN, OUT any] struct {
	compiler     engine.Compiler[DATA, IN, OUT]
	policies     map[string]policy.Policy[DATA, IN, OUT]
	sortPolicies func() []policy.Policy[DATA, IN, OUT]
	sync.Mutex
}

func NewPolicyAccessor[DATA, IN, OUT any](compiler engine.Compiler[DATA, IN, OUT]) *policyAccessor[DATA, IN, OUT] {
	return &policyAccessor[DATA, IN, OUT]{
		Mutex:    sync.Mutex{},
		compiler: compiler,
		policies: make(map[string]policy.Policy[DATA, IN, OUT]),
		sortPolicies: func() []policy.Policy[DATA, IN, OUT] {
			return nil
		},
	}
}

type Processor interface {
	Process(req *protov1alpha1.ValidatingPolicy)
}

func (p *policyAccessor[DATA, IN, OUT]) Process(req *protov1alpha1.ValidatingPolicy) {
	resetSortPolicies := func() {
		p.sortPolicies = sync.OnceValue(func() []policy.Policy[DATA, IN, OUT] {
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

	vpol := protov1alpha1.FromProto(req)
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

func (p *policyAccessor[DATA, IN, OUT]) Load(ctx context.Context) ([]policy.Policy[DATA, IN, OUT], error) {
	return p.sortPolicies(), nil
}
