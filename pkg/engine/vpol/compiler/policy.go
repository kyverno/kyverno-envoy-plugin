package compiler

import (
	"errors"
	"net/http"
	"sync"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	envoy "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"go.uber.org/multierr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/cel/lazy"
)

type compiledPolicy struct {
	failurePolicy   admissionregistrationv1.FailurePolicyType
	matchConditions []cel.Program
	variables       map[string]cel.Program
	rules           []cel.Program
}

func (p compiledPolicy) Ammar(r http.Request) {}

func (p compiledPolicy) For(r *authv3.CheckRequest) (engine.PolicyFunc, engine.PolicyFunc) {
	match := sync.OnceValues(func() (bool, error) {
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
	})
	variables := sync.OnceValue(func() map[string]any {
		vars := lazy.NewMapValue(authzcel.VariablesType)
		data := map[string]any{
			ObjectKey:    r,
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
		return data
	})
	rules := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data := variables()
		for _, rule := range p.rules {
			// evaluate the rule
			response, err := evaluateRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			if response != nil {
				// no error and evaluation result is not nil, return
				return response.ToCheckResponse(), nil
			}
		}
		return nil, nil
	}
	failurePolicy := func(inner func() (*authv3.CheckResponse, error)) func() (*authv3.CheckResponse, error) {
		return func() (*authv3.CheckResponse, error) {
			response, err := inner()
			if err != nil && p.failurePolicy == admissionregistrationv1.Fail {
				return nil, err
			}
			return response, nil
		}
	}
	return failurePolicy(rules), nil
}

func evaluateRule(rule cel.Program, data map[string]any) (envoy.Response, error) {
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
	response, ok := value.(envoy.Response)
	if !ok {
		return nil, errors.New("rule result is expected to be envoy.Response")
	}
	return response, nil
}
