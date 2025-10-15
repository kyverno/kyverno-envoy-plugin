package dispatchers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/breakers"
)

func Sequential[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](
	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
	breaker core.BreakerFactory[POLICY, IN, OUT, DATA],
) core.DispatcherFactory[POLICY, IN, OUT, DATA] {
	if breaker == nil {
		breaker = breakers.NeverFactory[POLICY, IN, OUT, DATA]()
	}
	return core.MakeDispatcherFactory(func(ctx context.Context, data DATA, policies []POLICY, collector core.Collector[POLICY, IN, OUT]) core.Dispatcher[IN] {
		evaluator := evaluator(ctx, data)
		breaker := breaker(ctx, data)
		return core.MakeDispatcherFunc(func(ctx context.Context, input IN) {
			for _, policy := range policies {
				out := evaluator.Evaluate(ctx, policy, input)
				collector.Collect(ctx, policy, input, out)
				if breaker.Break(ctx, policy, input, out) {
					break
				}
			}
		})
	})
}
