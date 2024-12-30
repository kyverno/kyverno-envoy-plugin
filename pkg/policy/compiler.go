package policy

import (
	"fmt"
	"sync"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	engine "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	envoy "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/cel/lazy"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type (
	AllowFunc func() (*authv3.CheckResponse, error)
	DenyFunc  func() (*authv3.CheckResponse, error)
)

type CompiledPolicy interface {
	For(r *authv3.CheckRequest) (AllowFunc, DenyFunc)
}

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

func (p compiledPolicy) For(r *authv3.CheckRequest) (AllowFunc, DenyFunc) {
	match := sync.OnceValues(func() (bool, error) {
		data := map[string]any{
			ObjectKey: r,
		}
		for _, matchCondition := range p.matchConditions {
			// evaluate the condition
			out, _, err := matchCondition.Eval(data)
			// check error
			if err != nil {
				return false, err
			}
			// try to convert to a bool
			result, err := utils.ConvertToNative[bool](out)
			// check error
			if err != nil {
				return false, err
			}
			// if condition is false, skip
			if !result {
				return false, nil
			}
		}
		return true, nil
	})
	variables := sync.OnceValue(func() map[string]any {
		vars := lazy.NewMapValue(engine.VariablesType)
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
	allow := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data := variables()
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
			response, err := evaluateRule[envoy.OkResponse](rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// no error and evaluation result is not nil, return
			return &authv3.CheckResponse{
				Status: response.Status,
				HttpResponse: &authv3.CheckResponse_OkResponse{
					OkResponse: response.OkHttpResponse,
				},
				DynamicMetadata: response.DynamicMetadata,
			}, nil
		}
		return nil, nil
	}
	deny := func() (*authv3.CheckResponse, error) {
		if match, err := match(); err != nil {
			return nil, err
		} else if !match {
			return nil, nil
		}
		data := variables()
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
			response, err := evaluateRule[envoy.DeniedResponse](rule, data)
			// check error
			if err != nil {
				return nil, err
			}
			// no error and evaluation result is not nil, return
			return &authv3.CheckResponse{
				Status: response.Status,
				HttpResponse: &authv3.CheckResponse_DeniedResponse{
					DeniedResponse: response.DeniedHttpResponse,
				},
				DynamicMetadata: response.DynamicMetadata,
			}, nil
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

func evaluateRule[T any](rule authorizationProgram, data map[string]any) (*T, error) {
	out, _, err := rule.response.Eval(data)
	// check error
	if err != nil {
		return nil, err
	}
	response, err := utils.ConvertToNative[T](out)
	// check error
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type Compiler interface {
	Compile(*v1alpha1.AuthorizationPolicy) (CompiledPolicy, field.ErrorList)
}

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy *v1alpha1.AuthorizationPolicy) (CompiledPolicy, field.ErrorList) {
	var allErrs field.ErrorList
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
	matchConditions := make([]cel.Program, 0, len(policy.Spec.MatchConditions))
	{
		path := path.Child("matchConditions")
		for i, matchCondition := range policy.Spec.MatchConditions {
			path := path.Index(i).Child("expression")
			ast, issues := env.Compile(matchCondition.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path, matchCondition.Expression, err.Error()))
			}
			if !ast.OutputType().IsExactType(types.BoolType) {
				return nil, append(allErrs, field.Invalid(path, matchCondition.Expression, "matchCondition output is expected to be of type bool"))
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path, matchCondition.Expression, err.Error()))
			}
			matchConditions = append(matchConditions, prog)
		}
	}
	variables := map[string]cel.Program{}
	{
		path := path.Child("variables")
		for i, variable := range policy.Spec.Variables {
			path := path.Index(i).Child("expression")
			ast, issues := env.Compile(variable.Expression)
			if err := issues.Err(); err != nil {
				return nil, append(allErrs, field.Invalid(path, variable.Expression, err.Error()))
			}
			provider.RegisterField(variable.Name, ast.OutputType())
			prog, err := env.Program(ast)
			if err != nil {
				return nil, append(allErrs, field.Invalid(path, variable.Expression, err.Error()))
			}
			variables[variable.Name] = prog
		}
	}
	var denies []authorizationProgram
	{
		path := path.Child("deny")
		for i, rule := range policy.Spec.Deny {
			path := path.Index(i)
			program, errs := compileAuthorization(path, rule, env, envoy.DeniedResponseType)
			if errs != nil {
				return nil, append(allErrs, errs...)
			}
			denies = append(denies, program)
		}
	}
	var allows []authorizationProgram
	{
		path := path.Child("allow")
		for i, rule := range policy.Spec.Allow {
			path := path.Index(i)
			program, errs := compileAuthorization(path, rule, env, envoy.OkResponseType)
			if errs != nil {
				return nil, append(allErrs, errs...)
			}
			allows = append(allows, program)
		}
	}
	return compiledPolicy{
		failurePolicy:   policy.Spec.GetFailurePolicy(),
		matchConditions: matchConditions,
		variables:       variables,
		allow:           allows,
		deny:            denies,
	}, nil
}

func compileAuthorization(path *field.Path, rule v1alpha1.Authorization, env *cel.Env, output *types.Type) (authorizationProgram, field.ErrorList) {
	var allErrs field.ErrorList
	program := authorizationProgram{}
	if rule.Match != "" {
		path := path.Child("match")
		ast, issues := env.Compile(rule.Match)
		if err := issues.Err(); err != nil {
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Match, err.Error()))
		}
		if !ast.OutputType().IsExactType(types.BoolType) {
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Match, "rule match output is expected to be of type bool"))
		}
		prog, err := env.Program(ast)
		if err != nil {
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Match, err.Error()))
		}
		program.match = prog
	}
	{
		path := path.Child("response")
		ast, issues := env.Compile(rule.Response)
		if err := issues.Err(); err != nil {
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Response, err.Error()))
		}
		if !ast.OutputType().IsExactType(output) {
			msg := fmt.Sprintf("rule response output is expected to be of type %s", output.TypeName())
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Response, msg))
		}
		prog, err := env.Program(ast)
		if err != nil {
			return authorizationProgram{}, append(allErrs, field.Invalid(path, rule.Response, err.Error()))
		}
		program.response = prog
	}
	return program, nil
}
