package controlplane

import (
	"context"

	policyapi "github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/sender"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type policyReconciler struct {
	client    client.Client
	polSender *sender.PolicySender
}

func NewReconciler(client client.Client, sender *sender.PolicySender) *policyReconciler {
	return &policyReconciler{
		client:    client,
		polSender: sender,
	}
}

func (r *policyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy v1alpha1.ValidatingPolicy
	err := r.client.Get(ctx, req.NamespacedName, &policy)
	if errors.IsNotFound(err) {
		if r.polSender != nil {
			r.polSender.DeletePolicy(req.Name)
			go r.polSender.SendPolicy(&protov1alpha1.ValidatingPolicy{
				Name:   req.Name,
				Delete: true,
			})
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
	protoPolicy := ToProto(&policy)
	if r.polSender != nil {
		r.polSender.StorePolicy(protoPolicy)
		go r.polSender.SendPolicy(protoPolicy)
	}
	return ctrl.Result{}, nil
}
