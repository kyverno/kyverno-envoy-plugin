package defaults

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

// TODO: make Sequential generic enough to not duplicate it here

func Handler[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](
	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
) core.HandlerFactory[POLICY, IN, Result[POLICY, IN, OUT, DATA], DATA] {
	resulter := Resulter[POLICY, IN, OUT, DATA]()
	return core.MakeHandlerFactory(func(ctx context.Context, data DATA, policies []POLICY, err error) core.Handler[IN, Result[POLICY, IN, OUT, DATA]] {
		evaluator := evaluator(ctx, data)
		resulter := resulter(ctx, data)
		return core.MakeHandlerFunc(func(ctx context.Context, input IN) Result[POLICY, IN, OUT, DATA] {
			for _, policy := range policies {
				out := evaluator.Evaluate(ctx, policy, input)
				resulter.Compute(ctx, policy, input, out)
			}
			return Result[POLICY, IN, OUT, DATA]{
				Input:    input,
				Data:     data,
				Source:   MakeSourceResult(policies, err),
				Policies: resulter.Result(),
			}
		})
	})
}

// func Handler[
// 	POLICY any,
// 	IN any,
// 	OUT any,
// 	DATA any,
// ](
// 	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
// ) core.HandlerFactory[POLICY, IN, []PolicyResult[POLICY, IN, OUT], DATA] {
// 	return handlers.Sequential(
// 		evaluator,
// 		Resulter[POLICY, IN, OUT, DATA](),
// 		nil,
// 	)
// }
