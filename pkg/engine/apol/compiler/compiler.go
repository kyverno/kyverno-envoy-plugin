package compiler

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	envoy "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type Compiler = engine.Compiler[*v1alpha1.AuthorizationPolicy]

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy *v1alpha1.AuthorizationPolicy) (engine.CompiledPolicy, field.ErrorList) {
	var allErrs field.ErrorList
	base, err := authzcel.NewEnv()
	if err != nil {
		return nil, append(allErrs, field.InternalError(nil, err))
	}
	provider := authzcel.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable(ObjectKey, envoy.CheckRequest),
		cel.Variable(VariablesKey, authzcel.VariablesType),
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
			program, errs := compileAuthorization(path, rule, env, envoy.CheckResponse)
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
			program, errs := compileAuthorization(path, rule, env, envoy.CheckResponse)
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
