package sources

import (
	"context"
	"slices"
	"sync"

	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"golang.org/x/exp/maps"
)

type listener struct {
	lock      *sync.Mutex
	resources map[string]*v1alpha1.ValidatingPolicy
}

func NewListener() *listener {
	return &listener{
		lock:      &sync.Mutex{},
		resources: map[string]*v1alpha1.ValidatingPolicy{},
	}
}

func (p *listener) Process(req *protov1alpha1.ValidatingPolicy) {
	if req.Delete {
		p.lock.Lock()
		defer p.lock.Unlock()
		delete(p.resources, req.Name)
		return
	}
	vpol := controlplane.FromProto(req)
	p.lock.Lock()
	defer p.lock.Unlock()
	p.resources[req.Name] = vpol
}

func (r *listener) Load(ctx context.Context) ([]*v1alpha1.ValidatingPolicy, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.resources == nil {
		return nil, nil
	}

	out := make([]*v1alpha1.ValidatingPolicy, 0, len(r.resources))

	// Collect and sort keys for deterministic iteration
	keys := maps.Keys(r.resources)
	slices.Sort(keys)

	for _, key := range keys {
		out = append(out, r.resources[key])
	}
	return out, nil
}
