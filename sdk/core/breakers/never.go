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
	IN any,
	OUT any,
	DATA any,
]() core.BreakerFactory[POLICY, IN, OUT, DATA] {
	return func(_ context.Context, data DATA) core.Breaker[POLICY, IN, OUT] {
		return Never[POLICY, IN, OUT]()
	}
}
