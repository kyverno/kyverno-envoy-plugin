package core

import "context"

// TODO: should return the resulter ?

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
	IN any,
	OUT any,
	DATA any,
] = func(context.Context, DATA, []POLICY, error) Handler[IN, OUT]

func MakeHandlerFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, DATA, []POLICY, error) Handler[IN, OUT]) HandlerFactory[POLICY, IN, OUT, DATA] {
	return f
}
