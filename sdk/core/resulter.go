package core

import "context"

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
	IN any,
	OUT any,
	DATA any,
	RESULT any,
] = func(context.Context, DATA, []POLICY, error) Resulter[POLICY, IN, OUT, RESULT]

func MakeResulterFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
	RESULT any,
](f func(context.Context, DATA, []POLICY, error) Resulter[POLICY, IN, OUT, RESULT]) ResulterFactory[POLICY, IN, OUT, DATA, RESULT] {
	return f
}
