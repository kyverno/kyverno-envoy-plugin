package engine

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/defaults"
)

func NewEngine[
	POLICY Policy[IN, OUT, DATA],
	DATA any,
	IN any,
	OUT any,
](
	source core.Source[POLICY],
) core.Engine[IN, defaults.Result[POLICY, DATA, IN, Evaluation[OUT]], DATA] {
	return core.NewEngine(
		source,
		defaults.Handler(EvaluatorFactory[POLICY]()),
	)
}
