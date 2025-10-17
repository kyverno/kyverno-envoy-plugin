package engine

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"k8s.io/client-go/dynamic"
)

type Policy = policy.Policy[
	dynamic.Interface,
	*authv3.CheckRequest,
	*authv3.CheckResponse,
]
