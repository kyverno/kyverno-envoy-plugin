package authz

import (
	"context"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
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
		compiled, err := compile(policy)
		if err != nil {
			return nil, err
		}
		result, err := compiled(req)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}
