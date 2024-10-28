package cel

import (
	"github.com/google/cel-go/cel"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/envoy"
)

func NewEnv() (*cel.Env, error) {
	return cel.NewEnv(cel.Lib(envoy.Lib()))
}
