package compiler

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	envoy "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	VariablesKey = "variables"
	ObjectKey    = "object"
)

type Compiler = engine.Compiler[*v1alpha1.ValidatingPolicy]

func NewCompiler() Compiler {
	return &compiler{}
}

type compiler struct{}

func (c *compiler) Compile(policy *v1alpha1.ValidatingPolicy) (engine.CompiledPolicy, field.ErrorList) {
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
	var rules []cel.Program
	{
		path := path.Child("validations")
		for i, rule := range policy.Spec.Validations {
			path := path.Index(i)
			program, errs := compileAuthorization(path, rule, env)
			if errs != nil {
				return nil, append(allErrs, errs...)
			}
			rules = append(rules, program)
		}
	}
	return compiledPolicy{
		failurePolicy:   policy.GetFailurePolicy(),
		matchConditions: matchConditions,
		variables:       variables,
		rules:           rules,
	}, nil
}

func compileAuthorization(path *field.Path, rule admissionregistrationv1.Validation, env *cel.Env) (cel.Program, field.ErrorList) {
	var allErrs field.ErrorList
	{
		path := path.Child("expression")
		ast, issues := env.Compile(rule.Expression)
		if err := issues.Err(); err != nil {
			return nil, append(allErrs, field.Invalid(path, rule.Expression, err.Error()))
		}
		if !ast.OutputType().IsExactType(envoy.CheckResponse) && !ast.OutputType().IsExactType(types.NullType) {
			msg := fmt.Sprintf("rule response output is expected to be of type %s", envoy.CheckResponse.TypeName())
			return nil, append(allErrs, field.Invalid(path, rule.Expression, msg))
		}
		prog, err := env.Program(ast)
		if err != nil {
			return nil, append(allErrs, field.Invalid(path, rule.Expression, err.Error()))
		}
		return prog, nil
	}
}
