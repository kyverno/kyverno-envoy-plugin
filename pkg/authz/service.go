package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	ctrl "sigs.k8s.io/controller-runtime"
)

type service struct {
	provider engine.Provider
}

func (s *service) Check(ctx context.Context, r *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// execute check
	response, err := s.check(ctx, r)
	// log error if any
	if err != nil {
		ctrl.LoggerFrom(ctx).Error(err, "Check failed")
	}
	// return response and error
	return response, err
}

func (s *service) check(ctx context.Context, r *authv3.CheckRequest) (_r *authv3.CheckResponse, _err error) {
	// fetch compiled policies
	policies, err := s.provider.CompiledPolicies(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: eliminate allocations
	allow := make([]engine.PolicyFunc, 0, len(policies))
	deny := make([]engine.PolicyFunc, 0, len(policies))
	// iterate over policies
	for _, policy := range policies {
		// collect allow/deny
		a, d := policy.For(r)
		if a != nil {
			allow = append(allow, a)
		}
		if d != nil {
			deny = append(deny, d)
		}
	}
	// check deny first
	for _, deny := range deny {
		// execute rule
		response, err := deny()
		// return error if any
		if err != nil {
			return nil, err
		}
		// if the reponse returned by the rule evaluation was not nil, return
		if response != nil {
			return response, nil
		}
	}
	// check allow
	for _, allow := range allow {
		// execute rule
		response, err := allow()
		// return error if any
		if err != nil {
			return nil, err
		}
		// if the reponse returned by the rule evaluation was not nil, return
		if response != nil {
			return response, nil
		}
	}
	// we didn't have a response
	return &authv3.CheckResponse{}, nil
}
