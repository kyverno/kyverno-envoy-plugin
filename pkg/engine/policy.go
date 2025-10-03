package engine

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"k8s.io/client-go/dynamic"
)

type PolicyFunc func() (*authv3.CheckResponse, error)

type CompiledPolicy interface {
	For(*authv3.CheckRequest, dynamic.Interface) (PolicyFunc, PolicyFunc)
}
