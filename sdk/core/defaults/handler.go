package defaults

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/handlers"
)

func Handler[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](
	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
) core.HandlerFactory[POLICY, IN, Result[POLICY, IN, OUT, DATA], DATA] {
	return handlers.Handler(
		Dispatcher(evaluator),
		Resulter[POLICY, IN, OUT, DATA](),
	)
}
