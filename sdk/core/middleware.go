package core

import "context"

type MiddlewareFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
] = func(context.Context, []POLICY, error, Handler[IN, OUT]) Handler[IN, OUT]
