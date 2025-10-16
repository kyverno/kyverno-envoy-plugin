package core

import "context"

type Engine[
	IN any,
	OUT any,
	DATA any,
] interface {
	Handle(context.Context, IN, DATA) OUT
}

type engine[
	IN any,
	OUT any,
	DATA any,
] func(context.Context, IN, DATA) OUT

func (e engine[IN, OUT, DATA]) Handle(ctx context.Context, in IN, data DATA) OUT {
	return e(ctx, in, data)
}

func NewEngine[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
](
	source Source[POLICY],
	handler HandlerFactory[POLICY, DATA, IN, OUT],
) engine[IN, OUT, DATA] {
	return func(ctx context.Context, in IN, data DATA) OUT {
		policies, err := source.Load(ctx)
		sctx := MakeSourceContext(policies, err)
		fctx := MakeFactoryContext(sctx, data, in)
		handler := handler(ctx, fctx)
		return handler.Handle(ctx, in)
	}
}
