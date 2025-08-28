package engine

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type PolicyFunc func() (*authv3.CheckResponse, error)

type CompiledPolicy interface {
	For(r *authv3.CheckRequest) (PolicyFunc, PolicyFunc)
}
