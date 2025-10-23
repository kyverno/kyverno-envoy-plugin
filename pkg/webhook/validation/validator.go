package validation

import (
	"context"
	"fmt"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func NewValidator(compileVpol func(*vpol.ValidatingPolicy) field.ErrorList) *validator {
	return &validator{
		compileVpol: compileVpol,
	}
}

type validator struct {
	compileVpol func(*vpol.ValidatingPolicy) field.ErrorList
}

func (v *validator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	switch obj := obj.(type) {
	case *vpol.ValidatingPolicy:
		return nil, v.validateVpol(obj)
	}
	return nil, fmt.Errorf("expected a ValidatingPolicy object but got %T", obj)
}

func (v *validator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	switch newObj := newObj.(type) {
	case *vpol.ValidatingPolicy:
		return nil, v.validateVpol(newObj)
	}
	return nil, fmt.Errorf("expected a ValidatingPolicy object but got %T", newObj)
}

func (*validator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *validator) validateVpol(policy *vpol.ValidatingPolicy) error {
	if policy.Spec.EvaluationMode() == v1alpha1.EvaluationModeEnvoy ||
		policy.Spec.EvaluationMode() == v1alpha1.EvaluationModeHTTP {
		if allErrs := v.compileVpol(policy); len(allErrs) > 0 {
			return apierrors.NewInvalid(
				vpol.SchemeGroupVersion.WithKind("ValidatingPolicy").GroupKind(),
				policy.Name,
				allErrs,
			)
		}
	}
	return nil
}
