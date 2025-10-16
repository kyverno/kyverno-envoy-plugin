package core

import "context"

type Factory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = func(context.Context, FactoryContext[POLICY, DATA, IN]) OUT
