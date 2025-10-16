package core

import "context"

type Breaker[
	POLICY any,
	IN any,
	OUT any,
] interface {
	Break(context.Context, POLICY, IN, OUT) bool
}

type BreakerFunc[
	POLICY any,
	IN any,
	OUT any,
] func(context.Context, POLICY, IN, OUT) bool

func (f BreakerFunc[POLICY, IN, OUT]) Break(ctx context.Context, policy POLICY, in IN, out OUT) bool {
	return f(ctx, policy, in, out)
}

func MakeBreakerFunc[
	POLICY any,
	IN any,
	OUT any,
](f func(ctx context.Context, policy POLICY, in IN, out OUT) bool) BreakerFunc[POLICY, IN, OUT] {
	return f
}

type BreakerFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = Factory[POLICY, DATA, IN, Breaker[POLICY, IN, OUT]]
