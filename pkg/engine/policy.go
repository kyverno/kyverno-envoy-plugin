package engine

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"k8s.io/client-go/dynamic"
)

type CompiledPolicy interface {
	Evaluate(*authv3.CheckRequest, dynamic.Interface) (*authv3.CheckResponse, error)
}
