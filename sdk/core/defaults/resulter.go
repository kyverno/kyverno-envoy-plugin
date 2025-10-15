package defaults

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

// TODO: make standard appender generic enough to not duplicate it here

type appender[
	POLICY any,
	IN any,
	OUT any,
] struct {
	result []PolicyResult[POLICY, IN, OUT]
}

func (r *appender[POLICY, IN, OUT]) Compute(_ context.Context, policy POLICY, in IN, out OUT) {
	r.result = append(r.result, MakePolicyResult(policy, in, out))
}

func (r *appender[POLICY, IN, OUT]) Result() []PolicyResult[POLICY, IN, OUT] {
	return r.result
}

func Resulter[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
]() core.ResulterFactory[POLICY, IN, OUT, DATA, []PolicyResult[POLICY, IN, OUT]] {
	return func(_ context.Context, data DATA) core.Resulter[POLICY, IN, OUT, []PolicyResult[POLICY, IN, OUT]] {
		return &appender[POLICY, IN, OUT]{}
	}
}
