package policy

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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
	client   client.Client
	compiler Compiler
	lock     *sync.RWMutex
	policies []policy
}

type policy struct {
	name       types.NamespacedName
	policyFunc PolicyFunc
}

func (r *policyReconciler) addPolicy(pol policy) {
	cmp := func(current, target policy) int {
		return strings.Compare(current.name.String(), target.name.String())
	}

	if i, found := slices.BinarySearchFunc(r.policies, pol, cmp); found {
		r.policies[i] = pol
	} else {
		slices.Insert(r.policies, i, pol)
	}
}

func newPolicyReconciler(client client.Client, compiler Compiler) *policyReconciler {
	return &policyReconciler{
		client:   client,
		compiler: compiler,
		lock:     &sync.RWMutex{},
		policies: make([]policy, 0),
	}
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pol v1alpha1.AuthorizationPolicy
	err := r.client.Get(ctx, req.NamespacedName, &pol)
	if errors.IsNotFound(err) {
		r.lock.Lock()
		defer r.lock.Unlock()
		slices.DeleteFunc(r.policies, func(p policy) bool {
			return req.NamespacedName == p.name
		})
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	compiled, errs := r.compiler.Compile(&pol)
	if len(errs) > 0 {
		fmt.Println(errs)
		// No need to retry it
		return ctrl.Result{}, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.addPolicy(policy{
		name:       req.NamespacedName,
		policyFunc: compiled,
	})
	return ctrl.Result{}, nil
}

func (r *policyReconciler) CompiledPolicies(ctx context.Context) ([]PolicyFunc, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	out := make([]PolicyFunc, 0, len(r.policies))
	for _, policy := range r.policies {
		out = append(out, policy.policyFunc)
	}
	return out, nil
}
