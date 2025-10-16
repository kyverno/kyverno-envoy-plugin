package core

import "context"

type MiddlewareFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = func(context.Context, []POLICY, error, Handler[IN, OUT]) Handler[IN, OUT]
