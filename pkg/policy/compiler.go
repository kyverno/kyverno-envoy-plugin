package policy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	engine "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/cel/lazy"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type PolicyFunc func(*authv3.CheckRequest) (*authv3.CheckResponse, error)

type Compiler interface {
	Compile(*v1alpha1.AuthorizationPolicy) (PolicyFunc, field.ErrorList)
}

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy *v1alpha1.AuthorizationPolicy) (PolicyFunc, field.ErrorList) {
	var allErrs field.ErrorList
	variables := map[string]cel.Program{}
	var authorizations []cel.Program
	base, err := engine.NewEnv()
	if err != nil {
		return nil, append(allErrs, field.InternalError(nil, err))
	}
	provider := engine.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable(ObjectKey, envoy.CheckRequest),
		cel.Variable(VariablesKey, engine.VariablesType),
		cel.CustomTypeProvider(provider),
	)
	if err != nil {
		return nil, append(allErrs, field.InternalError(nil, err))
	}
	path := field.NewPath("spec")
	{
		path := path.Child("variables")
		for i, variable := range policy.Spec.Variables {
			path := path.Index(i)
			ast, issues := env.Compile(variable.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), variable.Expression, err.Error()))
			}
			provider.RegisterField(variable.Name, ast.OutputType())
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), variable.Expression, err.Error()))
			}
			variables[variable.Name] = prog
		}
	}
	{
		path := path.Child("authorizations")
		for i, rule := range policy.Spec.Authorizations {
			path := path.Index(i)
			ast, issues := env.Compile(rule.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, err.Error()))
			}
			if !ast.OutputType().IsExactType(envoy.CheckResponse) {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, "rule output is expected to be of type envoy.service.auth.v3.CheckResponse"))
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path.Child("expression"), rule.Expression, err.Error()))
			}
			authorizations = append(authorizations, prog)
		}
	}
	eval := func(r *authv3.CheckRequest) (*authv3.CheckResponse, error) {
		vars := lazy.NewMapValue(engine.VariablesType)
		data := map[string]any{
			ObjectKey:    r,
			VariablesKey: vars,
		}
		for name, variable := range variables {
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
		for _, rule := range authorizations {
			// evaluate the rule
			out, _, err := rule.Eval(data)
			// check error
			if err != nil {
				return nil, err
			}
			// evaluation result is nil, continue
			if _, ok := out.(types.Null); ok {
				continue
			}
			// try to convert to a check response
			response, err := utils.ConvertToNative[*authv3.CheckResponse](out)
			// check error
			if err != nil {
				return nil, err
			}
			// evaluation result is nil, continue
			if response == nil {
				continue
			}
			// no error and evaluation result is not nil, return
			return response, nil
		}
		return nil, nil
	}
	return func(r *authv3.CheckRequest) (*authv3.CheckResponse, error) {
		response, err := eval(r)
		if err != nil && policy.Spec.GetFailurePolicy() == admissionregistrationv1.Fail {
			return nil, err
		}
		return response, nil
	}, nil
}
