package engine

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/defaults"
)

func NewEngine[
	POLICY Policy[DATA, IN, OUT],
	DATA any,
	IN any,
	OUT any,
](
	source core.Source[POLICY],
) core.Engine[DATA, IN, defaults.Result[POLICY, DATA, IN, Evaluation[OUT]]] {
	return core.NewEngine(
		source,
		defaults.Handler(EvaluatorFactory[POLICY]()),
	)
}
