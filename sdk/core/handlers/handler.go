package handlers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

type appender[
	POLICY any,
	IN any,
	OUT any,
] struct {
	result []OUT
}

func (r *appender[POLICY, IN, OUT]) Collect(_ context.Context, _ POLICY, _ IN, out OUT) {
	r.result = append(r.result, out)
}

func Handler[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
	RESULT any,
](
	dispatcher core.DispatcherFactory[POLICY, IN, OUT, DATA],
	resulter core.ResulterFactory[POLICY, IN, OUT, DATA, RESULT],
) core.HandlerFactory[POLICY, IN, RESULT, DATA] {
	return func(ctx context.Context, data DATA, policies []POLICY, err error) core.Handler[IN, RESULT] {
		resulter := resulter(ctx, data)
		dispatcher := dispatcher(ctx, data, policies, resulter)
		return core.MakeHandlerFunc(func(ctx context.Context, input IN) RESULT {
			dispatcher.Dispatch(ctx, input)
			return resulter.Result()
		})
	}
}
