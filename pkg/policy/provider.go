package policy

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Provider interface {
	CompiledPolicies(context.Context) ([]PolicyFunc, error)
}

func NewKubeProvider(mgr ctrl.Manager, compiler Compiler) (Provider, error) {
	r := newPolicyReconciler(mgr.GetClient(), compiler)
	if err := ctrl.NewControllerManagedBy(mgr).For(&v1alpha1.AuthorizationPolicy{}).Complete(r); err != nil {
		return nil, fmt.Errorf("failed to construct manager: %w", err)
	}
	return r, nil
}

type policyReconciler struct {
	client       client.Client
	compiler     Compiler
	lock         *sync.RWMutex
	policies     map[string]PolicyFunc
	sortPolicies func() []PolicyFunc
}

func newPolicyReconciler(client client.Client, compiler Compiler) *policyReconciler {
	return &policyReconciler{
		client:   client,
		compiler: compiler,
		lock:     &sync.RWMutex{},
		policies: map[string]PolicyFunc{},
	}
}

func mapToSortedSlice[K cmp.Ordered, V any](in map[K]V) []V {
	if in == nil {
		return nil
	}
	out := make([]V, 0, len(in))
	for _, key := range slices.Sorted(slices.Values(maps.Keys(in))) {
		out = append(out, in[key])
	}
	return out
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy v1alpha1.AuthorizationPolicy
	err := r.client.Get(ctx, req.NamespacedName, &policy)
	if errors.IsNotFound(err) {
		r.lock.Lock()
		defer r.lock.Unlock()
		delete(r.policies, req.NamespacedName.String())
		// r.sortPolicies = mapToSortedSlice(r.policies) // Reset the sorted func on every reconcile so the policies get resorted in next call
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
	r.policies[req.NamespacedName.String()] = compiled
	// r.sortPolicies = mapToSortedSlice(r.policies) // Reset the sorted func on every reconcile so the policies get resorted in next call
	return ctrl.Result{}, nil
}

func (r *policyReconciler) CompiledPolicies(ctx context.Context) ([]PolicyFunc, error) {
	return r.sortPolicies(), nil
}
