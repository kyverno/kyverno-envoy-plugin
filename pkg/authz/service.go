package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type service struct {
	engine    core.Engine[dynamic.Interface, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse]]
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
	// invoke engine
	response := s.engine.Handle(ctx, s.dynclient, r)
	if response.Result == nil {
		// we didn't have a response
		return &authv3.CheckResponse{}, response.Error
	}
	return response.Result, nil
}
