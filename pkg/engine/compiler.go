package engine

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Compiler[DATA, IN, OUT any] interface {
	Compile(*vpol.ValidatingPolicy) (policy.Policy[DATA, IN, OUT], field.ErrorList)
}
