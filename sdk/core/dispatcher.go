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

func (f DispatcherFunc[IN]) Dispatch(ctx context.Context, in IN) {
	f(ctx, in)
}

func MakeDispatcherFunc[
	IN any,
](f func(ctx context.Context, in IN)) DispatcherFunc[IN] {
	return f
}

type DispatcherFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = func(context.Context, FactoryContext[POLICY, DATA, IN], Collector[POLICY, IN, OUT]) Dispatcher[IN]
