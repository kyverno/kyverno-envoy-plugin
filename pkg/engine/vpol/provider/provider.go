package provider

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubeProvider(mgr ctrl.Manager, compiler compiler.Compiler) (engine.Provider, error) {
	provider := newPolicyReconciler(mgr.GetClient(), compiler)
	builder := ctrl.
		NewControllerManagedBy(mgr).
		For(&vpol.ValidatingPolicy{})
	if err := builder.Complete(provider); err != nil {
		return nil, fmt.Errorf("failed to construct controller: %w", err)
	}
	return provider, nil
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

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Reconcile...")
	var policy vpol.ValidatingPolicy
	// Reset the sorted func on every reconcile so the policies get resorted in next call
	resetSortPolicies := func() {
		r.sortPolicies = sync.OnceValue(func() []engine.CompiledPolicy {
			r.lock.Lock()
			defer r.lock.Unlock()
			return utils.ToSortedSlice(r.policies)
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
	if policy.Spec.EvaluationMode() == v1alpha1.EvaluationModeEnvoy {
		compiled, errs := r.compiler.Compile(&policy)
		if len(errs) > 0 {
			logger.Error(errs.ToAggregate(), "Policy compilation error")
			// No need to retry it
			return ctrl.Result{}, nil
		}
		r.lock.Lock()
		defer r.lock.Unlock()
		r.policies[req.String()] = compiled
		resetSortPolicies()
	}
	return ctrl.Result{}, nil
}

func (r *policyReconciler) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	return slices.Clone(r.sortPolicies()), nil
}
