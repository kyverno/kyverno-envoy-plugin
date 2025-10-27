package sources

import (
	"context"

	policyapi "github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	controlplane "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/discovery"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/processor"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type policyReconciler struct {
	client     client.Client
	polSender  *discovery.PolicyDiscoveryService
	processors map[v1alpha1.EvaluationMode]processor.Processor
}

func NewPolicyReconciler(client client.Client, discoveryService *discovery.PolicyDiscoveryService, processors map[v1alpha1.EvaluationMode]processor.Processor) *policyReconciler {
	return &policyReconciler{
		client:     client,
		polSender:  discoveryService,
		processors: processors,
	}
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy v1alpha1.ValidatingPolicy
	err := r.client.Get(ctx, req.NamespacedName, &policy)
	if errors.IsNotFound(err) {
		// Policy was deleted - notify discovery service and processors
		if r.polSender != nil {
			if err := r.polSender.DeletePolicy(req.Name); err != nil {
				ctrl.LoggerFrom(ctx).Error(err, "Failed to delete policy from discovery service", "policy", req.Name)
			}
		}
		protoRequest := &protov1alpha1.ValidatingPolicy{
			Name:   req.Name,
			Delete: true,
		}
		// delete this policy from any processor who may have it
		for _, p := range r.processors {
			p.Process(protoRequest)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	if policy.Spec.EvaluationMode() != policyapi.EvaluationModeHTTP &&
		policy.Spec.EvaluationMode() != policyapi.EvaluationModeEnvoy {
		return ctrl.Result{}, nil
	}
	protoPolicy := controlplane.ToProto(&policy)
	if r.polSender != nil {
		// StorePolicy automatically broadcasts to all connected clients
		if err := r.polSender.StorePolicy(protoPolicy); err != nil {
			ctrl.LoggerFrom(ctx).Error(err, "Failed to store policy in discovery service", "policy", policy.Name)
			return ctrl.Result{}, err
		}
	}
	go func() {
		if p, ok := r.processors[policy.Spec.EvaluationMode()]; ok {
			p.Process(protoPolicy)
		}
	}()
	return ctrl.Result{}, nil
}
