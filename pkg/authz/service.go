package authz

import (
	"context"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	engine "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type service struct {
	client client.Client
}

func (s *service) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	response, err := s.check(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	return response, err
}

func (s *service) check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// fetch policies
	var policies v1alpha1.AuthorizationPolicyList
	if err := s.client.List(ctx, &policies, &client.ListOptions{}); err != nil {
		return nil, err
	}
	for _, policy := range policies.Items {
		// create cel env
		base, err := engine.NewEnv()
		if err != nil {
			return nil, err
		}
		env, err := base.Extend(
			cel.Variable("input", envoy.CheckRequest),
			cel.Variable("variables", types.DynType),
		)
		if err != nil {
			return nil, err
		}
		variables := map[string]any{}
		data := map[string]any{
			"input":     req,
			"variables": variables,
		}
		for _, variable := range policy.Spec.Variables {
			ast, issues := env.Compile(variable.Expression)
			if err := issues.Err(); err != nil {
				return nil, err
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, err
			}
			out, _, err := prog.Eval(data)
			if err != nil {
				return nil, err
			}
			variables[variable.Name] = out.Value()
		}
		for _, rule := range policy.Spec.Authorizations {
			ast, issues := env.Compile(rule.Expression)
			if err := issues.Err(); err != nil {
				return nil, err
			}
			prog, err := env.Program(ast)
			if err != nil {
				return nil, err
			}
			out, _, err := prog.Eval(data)
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
	}
	return nil, nil
}
