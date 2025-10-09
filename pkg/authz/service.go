package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type service struct {
	provider  engine.Provider
	dynclient dynamic.Interface
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
	// check validations
	for _, policy := range policies {
		// execute rule
		response, err := policy.Evaluate(r, s.dynclient)
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
