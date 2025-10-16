package engine

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

func EvaluatorFactory[
	POLICY Policy[IN, OUT, DATA],
	DATA any,
	IN any,
	OUT any,
]() core.EvaluatorFactory[POLICY, DATA, IN, Evaluation[OUT]] {
	return func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN]) core.Evaluator[POLICY, IN, Evaluation[OUT]] {
		return core.MakeEvaluatorFunc(func(ctx context.Context, policy POLICY, in IN) Evaluation[OUT] {
			out, err := policy.Evaluate(ctx, in, fctx.Data)
			return MakeEvaluation(out, err)
		})
	}
}
