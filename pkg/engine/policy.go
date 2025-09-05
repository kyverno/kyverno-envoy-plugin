package engine

import (
	"net/http"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
)

type PolicyFunc func() (*authv3.CheckResponse, error)

type RequestFunc func() (*httpcel.Response, error)

type CompiledPolicy interface {
	ForEnvoy(r *authv3.CheckRequest) (PolicyFunc, PolicyFunc)
	ForHTTP(r *http.Request) RequestFunc
}
