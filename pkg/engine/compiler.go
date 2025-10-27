package engine

import (
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Compiler[POLICY any] interface {
	Compile(*vpol.ValidatingPolicy) (POLICY, field.ErrorList)
}
