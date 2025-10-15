package core

import "context"

type Dispatcher[
	IN any,
] interface {
	Dispatch(context.Context, IN)
}

type DispatcherFunc[
	IN any,
] func(context.Context, IN)

func (f DispatcherFunc[IN]) Dispatch(ctx context.Context, input IN) {
	f(ctx, input)
}

func MakeDispatcherFunc[
	IN any,
](f func(ctx context.Context, input IN)) DispatcherFunc[IN] {
	return f
}

type DispatcherFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
] = func(context.Context, DATA, []POLICY, Collector[POLICY, IN, OUT]) Dispatcher[IN]

func MakeDispatcherFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, DATA, []POLICY, Collector[POLICY, IN, OUT]) Dispatcher[IN]) DispatcherFactory[POLICY, IN, OUT, DATA] {
	return f
}
