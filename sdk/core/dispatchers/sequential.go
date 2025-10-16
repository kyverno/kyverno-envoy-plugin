package dispatchers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/breakers"
)

func Sequential[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
](
	evaluator core.EvaluatorFactory[POLICY, DATA, IN, OUT],
	breaker core.BreakerFactory[POLICY, DATA, IN, OUT],
) core.DispatcherFactory[POLICY, DATA, IN, OUT] {
	if breaker == nil {
		breaker = breakers.NeverFactory[POLICY, DATA, IN, OUT]()
	}
	return func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN], collector core.Collector[POLICY, IN, OUT]) core.Dispatcher[IN] {
		evaluator := evaluator(ctx, fctx)
		breaker := breaker(ctx, fctx)
		return core.MakeDispatcherFunc(func(ctx context.Context, in IN) {
			for _, policy := range fctx.Source.Data {
				out := evaluator.Evaluate(ctx, policy, in)
				collector.Collect(ctx, policy, in, out)
				if breaker.Break(ctx, policy, in, out) {
					break
				}
			}
		})
	}
}
