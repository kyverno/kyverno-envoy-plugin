package defaults

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
)

func Resulter[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
]() core.ResulterFactory[POLICY, IN, OUT, DATA, Result[POLICY, IN, OUT, DATA]] {
	return resulters.Transformer(
		func(policy POLICY, in IN, out OUT) PolicyResult[POLICY, IN, OUT] {
			return MakePolicyResult(policy, in, out)
		},
		func(results []PolicyResult[POLICY, IN, OUT]) Result[POLICY, IN, OUT, DATA] {
			return Result[POLICY, IN, OUT, DATA]{
				Policies: results,
			}
		},
		resulters.Appender[POLICY, IN, PolicyResult[POLICY, IN, OUT], DATA](),
	)
}
