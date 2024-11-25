package policy

import (
	"context"
	"fmt"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/validation"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
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
	compileFunc := func(policy *v1alpha1.AuthorizationPolicy) field.ErrorList {
		_, err := compiler.Compile(policy)
		fmt.Println("validating policy", policy.Name, err)
		return err
	}
	if err := ctrl.NewWebhookManagedBy(mgr).For(&v1alpha1.AuthorizationPolicy{}).WithValidator(validation.NewValidator(compileFunc)).Complete(); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}
	return r, nil
}

type policyReconciler struct {
	client   client.Client
	compiler Compiler
	lock     *sync.RWMutex
	policies map[types.NamespacedName]PolicyFunc
}

func newPolicyReconciler(client client.Client, compiler Compiler) *policyReconciler {
	return &policyReconciler{
		client:   client,
		compiler: compiler,
		lock:     &sync.RWMutex{},
		policies: map[types.NamespacedName]PolicyFunc{},
	}
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy v1alpha1.AuthorizationPolicy
	err := r.client.Get(ctx, req.NamespacedName, &policy)
	if errors.IsNotFound(err) {
		r.lock.Lock()
		defer r.lock.Unlock()
		delete(r.policies, req.NamespacedName)
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
	r.policies[req.NamespacedName] = compiled
	return ctrl.Result{}, nil
}

func (r *policyReconciler) CompiledPolicies(ctx context.Context) ([]PolicyFunc, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	out := make([]PolicyFunc, 0, len(r.policies))
	for _, policy := range r.policies {
		out = append(out, policy)
	}
	return out, nil
}
