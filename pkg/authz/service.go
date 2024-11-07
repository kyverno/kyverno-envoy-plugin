package authz

import (
	"context"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/policy"
)

type service struct {
	provider policy.Provider
}

func (s *service) Check(ctx context.Context, r *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// execute check
	response, err := s.check(ctx, r)
	// log error if any
	if err != nil {
		fmt.Println(err)
	}
	// return response and error
	return response, err
}

func (s *service) check(ctx context.Context, r *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// fetch compiled policies
	policies, err := s.provider.CompiledPolicies(ctx)
	if err != nil {
		return nil, err
	}
	// iterate over policies
	for _, policy := range policies {
		// execute policy
		response, err := policy(r)
		// return error if any
		if err != nil {
			return nil, err
		}
		// if the reponse returned by the policy evaluation was not nil, return
		if response != nil {
			return response, nil
		}
	}
	// we didn't have a response
	return nil, nil
}
