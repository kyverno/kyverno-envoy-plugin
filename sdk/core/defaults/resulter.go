package defaults

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
)

type transformer[
	POLICY any,
	IN any,
	OUT any,
	TRANFO any,
	RESULT any,
] struct {
	transfo func(POLICY, IN, OUT) TRANFO
	inner   core.Resulter[POLICY, IN, TRANFO, RESULT]
}

func (r *transformer[POLICY, IN, OUT, TRANFO, RESULT]) Collect(ctx context.Context, policy POLICY, in IN, out OUT) {
	r.inner.Collect(ctx, policy, in, r.transfo(policy, in, out))
}

func (r *transformer[POLICY, IN, OUT, TRANFO, RESULT]) Result() RESULT {
	return r.inner.Result()
}

func Transformer[
	POLICY any,
	IN any,
	OUT any,
	TRANFO any,
	DATA any,
	RESULT any,
](
	transfo func(POLICY, IN, OUT) TRANFO,
	inner core.ResulterFactory[POLICY, IN, TRANFO, DATA, RESULT],
) core.ResulterFactory[POLICY, IN, OUT, DATA, RESULT] {
	return func(ctx context.Context, data DATA) core.Resulter[POLICY, IN, OUT, RESULT] {
		inner := inner(ctx, data)
		return &transformer[POLICY, IN, OUT, TRANFO, RESULT]{
			transfo: transfo,
			inner:   inner,
		}
	}
}

func Resulter[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
]() core.ResulterFactory[POLICY, IN, OUT, DATA, []PolicyResult[POLICY, IN, OUT]] {
	return Transformer(
		func(policy POLICY, in IN, out OUT) PolicyResult[POLICY, IN, OUT] {
			return MakePolicyResult(policy, in, out)
		},
		resulters.Appender[POLICY, IN, PolicyResult[POLICY, IN, OUT], DATA](),
	)
}
