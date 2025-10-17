package core

import "context"

type Collector[
	POLICY any,
	IN any,
	OUT any,
] interface {
	Collect(context.Context, POLICY, IN, OUT)
}

type Resulter[
	POLICY any,
	IN any,
	OUT any,
	RESULT any,
] interface {
	Collector[POLICY, IN, OUT]
	Result() RESULT
}

type ResulterFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
	RESULT any,
] = Factory[POLICY, DATA, IN, Resulter[POLICY, IN, OUT, RESULT]]
