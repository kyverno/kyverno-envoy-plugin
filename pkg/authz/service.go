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

func (s *service) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	response, err := s.check(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	return response, err
}

func (s *service) check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// fetch policies
	policies, err := s.provider.CompiledPolicies(ctx)
	if err != nil {
		return nil, err
	}
	for _, policy := range policies {
		result, err := policy(req)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}
