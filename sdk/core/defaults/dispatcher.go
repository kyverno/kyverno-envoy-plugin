package defaults

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/breakers"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/dispatchers"
)

func Dispatcher[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](
	evaluator core.EvaluatorFactory[POLICY, IN, OUT, DATA],
) core.DispatcherFactory[POLICY, IN, OUT, DATA] {
	return dispatchers.Sequential(
		evaluator,
		breakers.NeverFactory[POLICY, IN, OUT, DATA](),
	)
}
