package breakers

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

type never[
	POLICY any,
	IN any,
	OUT any,
] struct{}

func (b never[POLICY, IN, OUT]) Break(context.Context, POLICY, IN, OUT) bool {
	return false
}

func Never[
	POLICY any,
	IN any,
	OUT any,
]() never[POLICY, IN, OUT] {
	return never[POLICY, IN, OUT]{}
}

func NeverFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
]() core.BreakerFactory[POLICY, DATA, IN, OUT] {
	return func(context.Context, core.FactoryContext[POLICY, DATA, IN]) core.Breaker[POLICY, IN, OUT] {
		return Never[POLICY, IN, OUT]()
	}
}
