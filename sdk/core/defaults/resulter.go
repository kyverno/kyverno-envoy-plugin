package defaults

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
)

func Resulter[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
]() core.ResulterFactory[POLICY, DATA, IN, OUT, Result[POLICY, DATA, IN, OUT]] {
	return func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN]) core.Resulter[POLICY, IN, OUT, Result[POLICY, DATA, IN, OUT]] {
		return resulters.NewTransformer(
			func(policy POLICY, in IN, out OUT) PolicyResult[POLICY, IN, OUT] {
				return MakePolicyResult(policy, in, out)
			},
			func(results []PolicyResult[POLICY, IN, OUT]) Result[POLICY, DATA, IN, OUT] {
				return MakeResult(
					MakeSourceResult(fctx.Source.Data, fctx.Source.Error),
					fctx.Input,
					fctx.Data,
					results,
				)
			},
			resulters.NewAppender[POLICY, IN, PolicyResult[POLICY, IN, OUT]](),
		)
	}
}
