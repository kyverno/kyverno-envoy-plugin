package provider

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubeProvider(mgr ctrl.Manager, compiler compiler.Compiler) (engine.Provider, error) {
	r := newPolicyReconciler(mgr.GetClient(), compiler)
	if err := ctrl.NewControllerManagedBy(mgr).For(&vpol.ValidatingPolicy{}).Complete(r); err != nil {
		return nil, fmt.Errorf("failed to construct manager: %w", err)
	}
	return r, nil
}

type policyReconciler struct {
	client       client.Client
	compiler     compiler.Compiler
	lock         *sync.Mutex
	policies     map[string]engine.CompiledPolicy
	sortPolicies func() []engine.CompiledPolicy
}

func newPolicyReconciler(client client.Client, compiler compiler.Compiler) *policyReconciler {
	return &policyReconciler{
		client:   client,
		compiler: compiler,
		lock:     &sync.Mutex{},
		policies: map[string]engine.CompiledPolicy{},
		sortPolicies: func() []engine.CompiledPolicy {
			return nil
		},
	}
}

func mapToSortedSlice[K cmp.Ordered, V any](in map[K]V) []V {
	if in == nil {
		return nil
	}
	out := make([]V, 0, len(in))
	keys := maps.Keys(in)
	slices.Sort(keys)
	for _, key := range keys {
		out = append(out, in[key])
	}
	return out
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy vpol.ValidatingPolicy
	// Reset the sorted func on every reconcile so the policies get resorted in next call
	resetSortPolicies := func() {
		r.sortPolicies = sync.OnceValue(func() []engine.CompiledPolicy {
			r.lock.Lock()
			defer r.lock.Unlock()
			return mapToSortedSlice(r.policies)
		})
	}
	err := r.client.Get(ctx, req.NamespacedName, &policy)
	if errors.IsNotFound(err) {
		r.lock.Lock()
		defer r.lock.Unlock()
		defer resetSortPolicies()
		delete(r.policies, req.String())
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	compiled, errs := r.compiler.Compile(&policy)
	if len(errs) > 0 {
		fmt.Println(errs)
		// No need to retry it
		return ctrl.Result{}, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.policies[req.String()] = compiled
	resetSortPolicies()
	return ctrl.Result{}, nil
}

func (r *policyReconciler) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	return slices.Clone(r.sortPolicies()), nil
}
