package compiler

import (
	"errors"
	"net/http"

	"sync"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/utils"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/variables"
	httpreq "github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/imagedata"
	"go.uber.org/multierr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/cel/lazy"
)

type authorizationProgram struct {
	match    cel.Program
	response cel.Program
}

type compiledPolicy struct {
	failurePolicy   admissionregistrationv1.FailurePolicyType
	matchConditions []cel.Program
	variables       map[string]cel.Program
	allow           []authorizationProgram
	deny            []authorizationProgram
}

func (p compiledPolicy) ForHTTP(r *http.Request) engine.RequestFunc {
	return nil
}

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
	variables := sync.OnceValues(func() (map[string]any, error) {
		loader, err := variables.ImageData(nil)
		if err != nil {
			return nil, err
		}
		vars := lazy.NewMapValue(authzcel.VariablesType)
		data := map[string]any{
			HttpKey:      httpreq.Context{ContextInterface: httpreq.NewHTTP(nil)},
			ImageDataKey: imagedata.Context{ContextInterface: loader},
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
		return data, nil
	})
	allow := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data, err := variables()
		if err != nil {
			return nil, err
		}
		for _, rule := range p.allow {
			matched, err := matchRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// if condition is false, continue
			if !matched {
				continue
			}
			// evaluate the rule
			response, err := evaluateRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// no error and evaluation result is not nil, return
			return response, nil
		}
		return nil, nil
	}
	deny := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data, err := variables()
		if err != nil {
			return nil, err
		}
		for _, rule := range p.deny {
			matched, err := matchRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// if condition is false, continue
			if !matched {
				continue
			}
			// evaluate the rule
			response, err := evaluateRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// no error and evaluation result is not nil, return
			return response, nil
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
	return failurePolicy(allow), failurePolicy(deny)
}

func matchRule(rule authorizationProgram, data map[string]any) (bool, error) {
	// if no match clause, consider a match
	if rule.match == nil {
		return true, nil
	}
	// evaluate rule match condition
	out, _, err := rule.match.Eval(data)
	if err != nil {
		return false, err
	}
	// try to convert to a match result
	matched, err := utils.ConvertToNative[bool](out)
	if err != nil {
		return false, err
	}
	return matched, err
}

func evaluateRule(rule authorizationProgram, data map[string]any) (*authv3.CheckResponse, error) {
	out, _, err := rule.response.Eval(data)
	// check error
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}
	value := out.Value()
	if value == nil {
		return nil, nil
	}
	response, ok := value.(*authv3.CheckResponse)
	if !ok {
		return nil, errors.New("rule result is expected to be authv3.CheckResponse")
	}
	return response, nil
}
