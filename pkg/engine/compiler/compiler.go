package compiler

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	authzcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	envoy "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/envoy"
	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/imagedata"
	"github.com/kyverno/kyverno/pkg/cel/libs/resource"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/dynamic"
)

const (
	HttpKey      = "http"
	ImageDataKey = "image"
	ObjectKey    = "object"
	VariablesKey = "variables"
	ResourceKey  = "resource"
)

func NewCompiler[DATA dynamic.Interface, IN, OUT any]() *compiler[DATA, IN, OUT] {
	return &compiler[DATA, IN, OUT]{}
}

type compiler[DATA dynamic.Interface, IN, OUT any] struct{}

func (c *compiler[DATA, IN, OUT]) Compile(policy *vpol.ValidatingPolicy) (policy.Policy[DATA, IN, OUT], field.ErrorList) {
	matchConditions, variables, rules, err := c.compiledEnvironment(policy)
	if err != nil {
		return compiledPolicy[DATA, IN, OUT]{}, err
	}
	return compiledPolicy[DATA, IN, OUT]{
		failurePolicy:   policy.GetFailurePolicy(),
		variables:       variables,
		matchConditions: matchConditions,
		rules:           rules,
	}, err
}

func (c *compiler[DATA, IN, OUT]) compiledEnvironment(policy *vpol.ValidatingPolicy) ([]cel.Program, map[string]cel.Program, []cel.Program, field.ErrorList) {
	var allErrs field.ErrorList
	base, err := authzcel.NewEnv()
	if err != nil {
		return nil, nil, nil, append(allErrs, field.InternalError(nil, err))
	}
	var objectKey cel.EnvOption

	switch policy.Spec.EvaluationMode() {
	case v1alpha1.EvaluationModeEnvoy:
		objectKey = cel.Variable(ObjectKey, envoy.CheckRequest)
	case v1alpha1.EvaluationModeHTTP:
		objectKey = cel.Variable(ObjectKey, httpauth.RequestType)
	default:
		return nil, nil, nil, append(allErrs, field.InternalError(nil, fmt.Errorf("invalid policy evaluation mode: %s", policy.Spec.EvaluationMode())))
	}
	provider := authzcel.NewVariablesProvider(base.CELTypeProvider())
	env, err := base.Extend(
		cel.Variable(HttpKey, http.ContextType),
		cel.Variable(ImageDataKey, imagedata.ContextType),
		objectKey,
		cel.Variable(VariablesKey, authzcel.VariablesType),
		cel.Variable(ResourceKey, resource.ContextType),
		cel.CustomTypeProvider(provider),
	)
	if err != nil {
		return nil, nil, nil, append(allErrs, field.InternalError(nil, err))
	}
	path := field.NewPath("spec")
	matchConditions := make([]cel.Program, 0, len(policy.Spec.MatchConditions))
	{
		path := path.Child("matchConditions")
		for i, matchCondition := range policy.Spec.MatchConditions {
			path := path.Index(i).Child("expression")
			ast, issues := env.Compile(matchCondition.Expression)
			if err := issues.Err(); err != nil {
				return nil, nil, nil, append(allErrs, field.Invalid(path, matchCondition.Expression, err.Error()))
			}
			if !ast.OutputType().IsExactType(types.BoolType) {
				return nil, nil, nil, append(allErrs, field.Invalid(path, matchCondition.Expression, "matchCondition output is expected to be of type bool"))
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, nil, nil, append(allErrs, field.Invalid(path, matchCondition.Expression, err.Error()))
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
				return nil, nil, nil, append(allErrs, field.Invalid(path, variable.Expression, err.Error()))
			}
			provider.RegisterField(variable.Name, ast.OutputType())
			prog, err := env.Program(ast)
			if err != nil {
				return nil, nil, nil, append(allErrs, field.Invalid(path, variable.Expression, err.Error()))
			}
			variables[variable.Name] = prog
		}
	}
	var rules []cel.Program
	{
		path := path.Child("validations")
		for i, rule := range policy.Spec.Validations {
			path := path.Index(i)
			program, errs := c.compileAuthorization(path, policy.Spec.EvaluationMode(), rule, env)
			if errs != nil {
				return nil, nil, nil, append(allErrs, errs...)
			}
			rules = append(rules, program)
		}
	}
	return matchConditions, variables, rules, nil
}

func (c *compiler[DATA, IN, OUT]) compileAuthorization(path *field.Path, evalMode vpol.EvaluationMode, rule admissionregistrationv1.Validation, env *cel.Env) (cel.Program, field.ErrorList) {
	var allErrs field.ErrorList
	{
		path := path.Child("expression")
		ast, issues := env.Compile(rule.Expression)
		if err := issues.Err(); err != nil {
			return nil, append(allErrs, field.Invalid(path, rule.Expression, err.Error()))
		}
		switch evalMode {
		case v1alpha1.EvaluationModeEnvoy:
			if !ast.OutputType().IsExactType(envoy.CheckResponse) && !ast.OutputType().IsExactType(types.NullType) {
				msg := fmt.Sprintf("rule response output is expected to be of type %s", envoy.CheckResponse.TypeName())
				return nil, append(allErrs, field.Invalid(path, rule.Expression, msg))
			}
		case v1alpha1.EvaluationModeHTTP:
			if !ast.OutputType().IsExactType(httpauth.ResponseType) && !ast.OutputType().IsExactType(types.NullType) {
				msg := fmt.Sprintf("rule response output is expected to be of type %s", envoy.CheckResponse.TypeName())
				return nil, append(allErrs, field.Invalid(path, rule.Expression, msg))
			}
		}
		prog, err := env.Program(ast)
		if err != nil {
			return nil, append(allErrs, field.Invalid(path, rule.Expression, err.Error()))
		}
		return prog, nil
	}
}
