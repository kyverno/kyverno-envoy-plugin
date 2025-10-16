package engine

import (
	"context"
)

type Policy[
	DATA any,
	IN any,
	OUT any,
] interface {
	Evaluate(context.Context, DATA, IN) (OUT, error)
}

type PolicyFunc[
	DATA any,
	IN any,
	OUT any,
] func(context.Context, DATA, IN) (OUT, error)

func (f PolicyFunc[DATA, IN, OUT]) Evaluate(ctx context.Context, runtime DATA, input IN) (OUT, error) {
	return f(ctx, runtime, input)
}

func MakePolicyFunc[
	DATA any,
	IN any,
	OUT any,
](f func(context.Context, DATA, IN) (OUT, error)) PolicyFunc[DATA, IN, OUT] {
	return f
}
