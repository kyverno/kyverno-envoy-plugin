package core

import "context"

type Engine[
	DATA any,
	IN any,
	OUT any,
] interface {
	Handle(context.Context, DATA, IN) OUT
}

type engine[
	DATA any,
	IN any,
	OUT any,
] func(context.Context, DATA, IN) OUT

func (e engine[DATA, IN, OUT]) Handle(ctx context.Context, data DATA, in IN) OUT {
	return e(ctx, data, in)
}

func NewEngine[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
](
	source Source[POLICY],
	handler HandlerFactory[POLICY, DATA, IN, OUT],
) engine[DATA, IN, OUT] {
	return func(ctx context.Context, data DATA, in IN) OUT {
		policies, err := source.Load(ctx)
		sctx := MakeSourceContext(policies, err)
		fctx := MakeFactoryContext(sctx, data, in)
		handler := handler(ctx, fctx)
		return handler.Handle(ctx, in)
	}
}
