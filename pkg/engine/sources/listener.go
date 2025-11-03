package sources

import (
	"context"
	"hash/fnv"
	"slices"
	"sync"

	"github.com/gogo/protobuf/proto"
	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"golang.org/x/exp/maps"
)

type listener struct {
	lock           *sync.Mutex
	evaluationMode v1alpha1.EvaluationMode
	resources      map[int64]*v1alpha1.ValidatingPolicy
}

func NewListener(evalMode v1alpha1.EvaluationMode) *listener {
	return &listener{
		lock:           &sync.Mutex{},
		evaluationMode: evalMode,
		resources:      map[int64]*v1alpha1.ValidatingPolicy{},
	}
}

func (p *listener) Process(req []*protov1alpha1.ValidatingPolicy) {
	newPolicyMap := map[int64]*v1alpha1.ValidatingPolicy{}
	for _, policy := range req {
		if policy.Spec.EvaluationMode != string(p.evaluationMode) {
			continue
		}
		data, err := proto.Marshal(policy)
		if err != nil {
			// todo: handle this somehow ?
			continue
		}

		h := fnv.New64a()
		h.Write(data)
		newPolicyMap[int64(h.Sum64())] = controlplane.FromProto(policy)
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.resources = newPolicyMap
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
