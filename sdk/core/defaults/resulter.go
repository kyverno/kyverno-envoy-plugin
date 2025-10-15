package defaults

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
)

func Resulter[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
]() core.ResulterFactory[POLICY, IN, OUT, DATA, Result[POLICY, IN, OUT, DATA]] {
	return func(ctx context.Context, data DATA, policies []POLICY, err error) core.Resulter[POLICY, IN, OUT, Result[POLICY, IN, OUT, DATA]] {
		return resulters.NewTransformer(
			func(policy POLICY, in IN, out OUT) PolicyResult[POLICY, IN, OUT] {
				return MakePolicyResult(policy, in, out)
			},
			func(results []PolicyResult[POLICY, IN, OUT]) Result[POLICY, IN, OUT, DATA] {
				return Result[POLICY, IN, OUT, DATA]{
					Data:     data,
					Source:   MakeSourceResult(policies, err),
					Policies: results,
				}
			},
			resulters.NewAppender[POLICY, IN, PolicyResult[POLICY, IN, OUT]](),
		)
	}
}
