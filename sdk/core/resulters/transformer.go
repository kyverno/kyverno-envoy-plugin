package resulters

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

type transformer[
	POLICY any,
	IN any,
	OUT any,
	TRANSFO any,
	RESULT any,
	TRANSFORESULT any,
] struct {
	collect func(POLICY, IN, OUT) TRANSFO
	result  func(RESULT) TRANSFORESULT
	inner   core.Resulter[POLICY, IN, TRANSFO, RESULT]
}

func NewTransformer[
	POLICY any,
	IN any,
	OUT any,
	TRANSFO any,
	RESULT any,
	TRANSFORESULT any,
](
	collect func(POLICY, IN, OUT) TRANSFO,
	result func(RESULT) TRANSFORESULT,
	inner core.Resulter[POLICY, IN, TRANSFO, RESULT],
) *transformer[POLICY, IN, OUT, TRANSFO, RESULT, TRANSFORESULT] {
	return &transformer[POLICY, IN, OUT, TRANSFO, RESULT, TRANSFORESULT]{
		collect: collect,
		result:  result,
		inner:   inner,
	}
}

func (r *transformer[POLICY, IN, OUT, TRANFO, RESULT, TRANSFORESULT]) Collect(ctx context.Context, policy POLICY, in IN, out OUT) {
	r.inner.Collect(ctx, policy, in, r.collect(policy, in, out))
}

func (r *transformer[POLICY, IN, OUT, TRANFO, RESULT, TRANSFORESULT]) Result() TRANSFORESULT {
	return r.result(r.inner.Result())
}
