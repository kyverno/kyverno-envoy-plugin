package assert

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Assertion interface {
	assert(context.Context, *field.Path, any, binding.Bindings, ...template.Option) (field.ErrorList, error)
}
