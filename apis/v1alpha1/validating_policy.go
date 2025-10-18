package v1alpha1

import (
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/proto/validatingpolicy/v1alpha1"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToProto(pol *v1alpha1.ValidatingPolicy) *protov1alpha1.ValidatingPolicy {
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
			Validations:     validations,
			Variables:       variables,
			MatchConditions: matchConds,
			FailurePolicy:   &fp,
		},
	}
}

func FromProto(pol *protov1alpha1.ValidatingPolicy) *v1alpha1.ValidatingPolicy {
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
	var evalMode v1alpha1.EvaluationMode
	switch pol.Spec.EvaluationMode {
	case "Envoy":
		evalMode = EvaluationModeEnvoy
	case "HTTP":
		evalMode = EvaluationModeHTTP
	}

	fp := admissionregistrationv1.FailurePolicyType(*pol.Spec.FailurePolicy)
	return &v1alpha1.ValidatingPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: pol.Name,
		},
		Spec: v1alpha1.ValidatingPolicySpec{
			EvaluationConfiguration: &v1alpha1.EvaluationConfiguration{
				Mode: evalMode,
			},
			Validations:     validations,
			Variables:       variables,
			MatchConditions: matchConds,
			FailurePolicy:   &fp,
		},
	}
}
