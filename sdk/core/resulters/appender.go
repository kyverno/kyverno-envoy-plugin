package resulters

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

type appender[
	POLICY any,
	IN any,
	OUT any,
] struct {
	result []OUT
}

func (r *appender[POLICY, IN, OUT]) Compute(_ context.Context, _ POLICY, _ IN, out OUT) {
	r.result = append(r.result, out)
}

func (r *appender[POLICY, IN, OUT]) Result() []OUT {
	return r.result
}

func Appender[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
]() core.ResulterFactory[POLICY, IN, OUT, DATA, []OUT] {
	return func(_ context.Context, data DATA) core.Resulter[POLICY, IN, OUT, []OUT] {
		return &appender[POLICY, IN, OUT]{}
	}
}
