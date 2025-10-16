package core

import "context"

type Handler[
	IN any,
	OUT any,
] interface {
	Handle(context.Context, IN) OUT
}

type HandlerFunc[
	IN any,
	OUT any,
] func(context.Context, IN) OUT

func (f HandlerFunc[IN, OUT]) Handle(ctx context.Context, input IN) OUT {
	return f(ctx, input)
}

func MakeHandlerFunc[
	IN any,
	OUT any,
](f func(ctx context.Context, input IN) OUT) HandlerFunc[IN, OUT] {
	return f
}

type HandlerFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = Factory[POLICY, DATA, IN, Handler[IN, OUT]]
