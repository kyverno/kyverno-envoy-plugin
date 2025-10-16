package handlers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

func Handler[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
	RESULT any,
](
	dispatcher core.DispatcherFactory[POLICY, DATA, IN, OUT],
	resulter core.ResulterFactory[POLICY, DATA, IN, OUT, RESULT],
) core.HandlerFactory[POLICY, DATA, IN, RESULT] {
	return func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN]) core.Handler[IN, RESULT] {
		resulter := resulter(ctx, fctx)
		dispatcher := dispatcher(ctx, fctx, resulter)
		return core.MakeHandlerFunc(func(ctx context.Context, in IN) RESULT {
			dispatcher.Dispatch(ctx, in)
			return resulter.Result()
		})
	}
}
