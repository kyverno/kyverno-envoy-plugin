package controlplane

import (
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToProto(pol *vpol.ValidatingPolicy) *protov1alpha1.ValidatingPolicy {
	validations := []*protov1alpha1.Validation{}
	for _, v := range pol.Spec.Validations {
		validations = append(validations, &protov1alpha1.Validation{
			Expression: v.Expression,
			Message:    &v.Message,
			Reason:     (*string)(v.Reason),
		})
	}
	variables := []*protov1alpha1.Variable{}
	for _, v := range pol.Spec.Variables {
		variables = append(variables,
			&protov1alpha1.Variable{
				Name:       v.Name,
				Expression: v.Expression,
			},
		)
	}
	matchConds := []*protov1alpha1.MatchCondition{}
	for _, m := range pol.Spec.MatchConditions {
		matchConds = append(matchConds, &protov1alpha1.MatchCondition{
			Name:       m.Name,
			Expression: m.Expression,
		})
	}
	// TODO: check if ignore is the default
	var fp string
	if pol.Spec.FailurePolicy != nil {
		fp = string(*pol.Spec.FailurePolicy)
	} else {
		fp = "Ignore"
	}
	return &protov1alpha1.ValidatingPolicy{
		Name: pol.Name,
		Spec: &protov1alpha1.ValidatingPolicySpec{
			EvaluationMode:  string(pol.Spec.EvaluationMode()),
			Validations:     validations,
			Variables:       variables,
			MatchConditions: matchConds,
			FailurePolicy:   &fp,
		},
	}
}

func FromProto(pol *protov1alpha1.ValidatingPolicy) *vpol.ValidatingPolicy {
	validations := []admissionregistrationv1.Validation{}
	for _, v := range pol.Spec.Validations {
		validations = append(validations, admissionregistrationv1.Validation{
			Expression: v.Expression,
			Message:    *v.Message,
			Reason:     (*metav1.StatusReason)(v.Reason),
		})
	}
	variables := []admissionregistrationv1.Variable{}
	for _, v := range pol.Spec.Variables {
		variables = append(variables,
			admissionregistrationv1.Variable{
				Name:       v.Name,
				Expression: v.Expression,
			},
		)
	}
	matchConds := []admissionregistrationv1.MatchCondition{}
	for _, m := range pol.Spec.MatchConditions {
		matchConds = append(matchConds, admissionregistrationv1.MatchCondition{
			Name:       m.Name,
			Expression: m.Expression,
		})
	}
	var evalMode vpol.EvaluationMode
	switch pol.Spec.EvaluationMode {
	case "Envoy":
		evalMode = v1alpha1.EvaluationModeEnvoy
	case "HTTP":
		evalMode = v1alpha1.EvaluationModeHTTP
	}
	var fp admissionregistrationv1.FailurePolicyType = "Ignore"
	if pol.Spec.FailurePolicy != nil {
		fp = admissionregistrationv1.FailurePolicyType(*pol.Spec.FailurePolicy)
	}

	return &vpol.ValidatingPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: pol.Name,
		},
		Spec: vpol.ValidatingPolicySpec{
			EvaluationConfiguration: &vpol.EvaluationConfiguration{
				Mode: evalMode,
			},
			Validations:     validations,
			Variables:       variables,
			MatchConditions: matchConds,
			FailurePolicy:   &fp,
		},
	}
}
