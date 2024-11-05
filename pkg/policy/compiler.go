package policy

import (
	"errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	engine "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apiserver/pkg/cel/lazy"
)

type PolicyFunc func(*authv3.CheckRequest) (*authv3.CheckResponse, error)

type Compiler interface {
	Compile(v1alpha1.AuthorizationPolicy) (PolicyFunc, error)
}

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy v1alpha1.AuthorizationPolicy) (PolicyFunc, error) {
	variables := map[string]cel.Program{}
	var authorizations []cel.Program
	base, err := engine.NewEnv()
	if err != nil {
		return nil, err
	}
	provider := engine.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable("input", envoy.CheckRequest),
		cel.Variable("variables", engine.VariablesType),
		cel.CustomTypeProvider(provider),
	)
	if err != nil {
		return nil, err
	}
	for _, variable := range policy.Spec.Variables {
		ast, issues := env.Compile(variable.Expression)
		if err := issues.Err(); err != nil {
			return nil, err
		}
		provider.RegisterField(variable.Name, ast.OutputType())
		prog, err := env.Program(ast)
		if err != nil {
			return nil, err
		}
		variables[variable.Name] = prog
	}
	for _, rule := range policy.Spec.Authorizations {
		ast, issues := env.Compile(rule.Expression)
		if err := issues.Err(); err != nil {
			return nil, err
		}
		if !ast.OutputType().IsExactType(envoy.CheckResponse) {
			return nil, errors.New("rule output is expected to be of type envoy.service.auth.v3.CheckResponse")
		}
		prog, err := env.Program(ast)
		if err != nil {
			return nil, err
		}
		authorizations = append(authorizations, prog)
	}
	eval := func(req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
		vars := lazy.NewMapValue(engine.VariablesType)
		data := map[string]any{
			"input":     req,
			"variables": vars,
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
			out, _, err := rule.Eval(data)
			if err != nil {
				return nil, err
			}
			response, err := utils.ConvertToNative[*authv3.CheckResponse](out)
			if err != nil {
				return nil, err
			}
			if response != nil {
				return response, nil
			}
		}
		return nil, nil
	}
	return func(req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
		response, err := eval(req)
		if err != nil && policy.Spec.GetFailurePolicy() == admissionregistrationv1.Fail {
			return nil, err
		}
		return response, nil
	}, nil
}
