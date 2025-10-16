package defaults

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/handlers"
)

func Handler[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
](
	evaluator core.EvaluatorFactory[POLICY, DATA, IN, OUT],
) core.HandlerFactory[POLICY, DATA, IN, Result[POLICY, DATA, IN, OUT]] {
	return handlers.Handler(
		Dispatcher(evaluator),
		Resulter[POLICY, DATA, IN, OUT](),
	)
}
