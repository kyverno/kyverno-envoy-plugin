package handlers

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
	RESULT any,
](
	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
	resulter core.ResulterFactory[POLICY, IN, OUT, DATA, RESULT],
	breaker core.BreakerFactory[POLICY, IN, RESULT, DATA],
) core.HandlerFactory[POLICY, IN, RESULT, DATA] {
	if breaker == nil {
		breaker = breakers.NeverFactory[POLICY, IN, RESULT, DATA]()
	}
	return core.MakeHandlerFactory(func(ctx context.Context, data DATA, policies []POLICY, _ error) core.Handler[IN, RESULT] {
		evaluator := evaluator(ctx, data)
		resulter := resulter(ctx, data)
		breaker := breaker(ctx, data)
		return core.MakeHandlerFunc(func(ctx context.Context, input IN) RESULT {
			for _, policy := range policies {
				out := evaluator.Evaluate(ctx, policy, input)
				resulter.Compute(ctx, policy, input, out)
				if breaker.Break(ctx, policy, input, resulter.Result()) {
					break
				}
			}
			return resulter.Result()
		})
	})
}
