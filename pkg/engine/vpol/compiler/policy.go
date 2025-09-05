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

	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/variables"

	httpreq "github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/imagedata"

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

func (p compiledPolicy) ForHTTP(r *http.Request) engine.RequestFunc {
	// ammar: you removed match conditions
	variables := sync.OnceValues(func() (map[string]any, error) {
		vars := lazy.NewMapValue(authzcel.VariablesType)
		req, err := httpauth.NewRequest(r)
		if err != nil {
			return nil, err
		}
		loader, err := variables.ImageData(nil)
		if err != nil {
			return nil, err
		}
		data := map[string]any{
			HttpKey:      httpreq.Context{ContextInterface: httpreq.NewHTTP(nil)},
			ImageDataKey: imagedata.Context{ContextInterface: loader},
			ObjectKey:    req,
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
	rules := func() (*httpauth.Response, error) {
		data, err := variables()
		if err != nil {
			return nil, err
		}
		for _, rule := range p.rules {
			// evaluate the rule
			response, err := evaluateHTTP(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			if response != nil {
				// no error and evaluation result is not nil, return
				return response, nil
			}
		}
		return nil, nil
	}
	failurePolicy := func(inner func() (*httpauth.Response, error)) func() (*httpauth.Response, error) {
		return func() (*httpauth.Response, error) {
			response, err := inner()
			if err != nil && p.failurePolicy == admissionregistrationv1.Fail {
				return nil, err
			}
			return response, nil
		}
	}
	return failurePolicy(rules)
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
	rules := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data, err := variables()
		if err != nil {
			return nil, err
		}
		for _, rule := range p.rules {
			// evaluate the rule
			response, err := evaluateRule(rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			if response != nil {
				// no error and evaluation result is not nil, return
				return response, nil
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

func evaluateRule(rule cel.Program, data map[string]any) (*authv3.CheckResponse, error) {
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
	response, ok := value.(*authv3.CheckResponse)
	if !ok {
		return nil, errors.New("rule result is expected to be authv3.CheckResponse")
	}
	return response, nil
}

func evaluateHTTP(rule cel.Program, data map[string]any) (*httpauth.Response, error) {
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
	response, ok := value.(*httpauth.Response)
	if !ok {
		return nil, errors.New("rule result is expected to be http.Response")
	}
	return response, nil
}
