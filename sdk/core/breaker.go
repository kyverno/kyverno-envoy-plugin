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

func (f BreakerFunc[POLICY, IN, OUT]) Break(ctx context.Context, policy POLICY, input IN, out OUT) bool {
	return f(ctx, policy, input, out)
}

func MakeBreakerFunc[
	POLICY any,
	IN any,
	OUT any,
](f func(ctx context.Context, policy POLICY, input IN, out OUT) bool) BreakerFunc[POLICY, IN, OUT] {
	return f
}

type BreakerFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
] = func(context.Context, DATA) Breaker[POLICY, IN, OUT]

func MakeBreakerFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, DATA) Breaker[POLICY, IN, OUT]) BreakerFactory[POLICY, IN, OUT, DATA] {
	return f
}
