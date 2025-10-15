package engine

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

func EvaluatorFactory[
	POLICY Policy[IN, OUT, DATA],
	IN any,
	OUT any,
	DATA any,
]() core.EvaluatorFactory[POLICY, IN, Evaluation[OUT], DATA] {
	return core.MakeEvaluatorFactory(func(ctx context.Context, data DATA) core.Evaluator[POLICY, IN, Evaluation[OUT]] {
		return core.MakeEvaluatorFunc(func(ctx context.Context, policy POLICY, in IN) Evaluation[OUT] {
			out, err := policy.Evaluate(ctx, in, data)
			return MakeEvaluation(out, err)
		})
	})
}
