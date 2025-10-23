package compiler

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/utils"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/variables"
	"github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/imagedata"
	"github.com/kyverno/kyverno/pkg/cel/libs/resource"
	"go.uber.org/multierr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/cel/lazy"
	"k8s.io/client-go/dynamic"
)

type compiledPolicy[DATA dynamic.Interface, IN, OUT any] struct {
	failurePolicy   admissionregistrationv1.FailurePolicyType
	matchConditions []cel.Program
	variables       map[string]cel.Program
	rules           []cel.Program
}

func (p compiledPolicy[DATA, IN, OUT]) Evaluate(ctx context.Context, dynclient DATA, r IN) (OUT, error) {
	var zero OUT // create a zero variable of the output type
	response, err := p.evaluateRules(r, dynclient)
	if err != nil && p.failurePolicy == admissionregistrationv1.Fail {
		return zero, err
	}
	return response, nil
}

func (p compiledPolicy[DATA, IN, OUT]) match(r IN) (bool, error) {
	data := map[string]any{
		ObjectKey: r,
	}
	var errs []error
	for _, matchCondition := range p.matchConditions {
		// evaluate the condition
		out, _, err := matchCondition.Eval(data)
		// check error
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// try to convert to a bool
		result, err := utils.ConvertToNative[bool](out)
		// check error
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// if condition is false, skip
		if !result {
			return false, nil
		}
	}
	return true, multierr.Combine(errs...)
}

func (p compiledPolicy[DATA, IN, OUT]) setupVariables(r IN, d DATA) (map[string]any, error) {
	loader, err := variables.ImageData(nil)
	if err != nil {
		return nil, err
	}
	vars := lazy.NewMapValue(authzcel.VariablesType)
	data := map[string]any{
		HttpKey:      http.Context{ContextInterface: http.NewHTTP(nil)},
		ImageDataKey: imagedata.Context{ContextInterface: loader},
		ObjectKey:    r,
		ResourceKey:  resource.Context{ContextInterface: variables.NewResourceProvider(d)},
		VariablesKey: vars,
	}
	for name, variable := range p.variables {
		vars.Append(name, func(*lazy.MapValue) ref.Val {
			out, _, err := variable.Eval(data)
			if out != nil {
				return out
			}
			if err != nil {
				return types.WrapErr(err)
			}
			return nil
		})
	}
	return data, nil
}

func (p compiledPolicy[DATA, IN, OUT]) evaluateRules(r IN, dynclient DATA) (OUT, error) {
	var zero OUT // create a zero variable of the output type
	if match, err := p.match(r); err != nil {
		return zero, err
	} else if !match {
		return zero, nil
	}
	data, err := p.setupVariables(r, dynclient)
	if err != nil {
		return zero, err
	}
	for _, rule := range p.rules {
		// evaluate the rule
		response, err := evaluateRule(rule, data)
		// check error
		if err != nil {
			return zero, err
		}
		if response != nil {
			// no error and evaluation result is not nil, return
			val, ok := response.(OUT)
			if !ok {
				return zero, fmt.Errorf("rule result is expected to be %T", zero)
			}
			return val, nil
		}
	}
	return zero, nil
}

func evaluateRule(rule cel.Program, data map[string]any) (any, error) {
	out, _, err := rule.Eval(data)
	// check error
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}
	if out == types.NullValue {
		return nil, nil
	}
	value := out.Value()
	if value == nil {
		return nil, nil
	}
	return value, nil
}
