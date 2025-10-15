package engine

import (
	"context"
)

type Policy[
	IN any,
	OUT any,
	DATA any,
] interface {
	Evaluate(context.Context, IN, DATA) (OUT, error)
}

type PolicyFunc[
	IN any,
	OUT any,
	DATA any,
] func(context.Context, IN, DATA) (OUT, error)

func (f PolicyFunc[IN, OUT, DATA]) Evaluate(ctx context.Context, input IN, runtime DATA) (OUT, error) {
	return f(ctx, input, runtime)
}

func MakePolicyFunc[
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, IN, DATA) (OUT, error)) PolicyFunc[IN, OUT, DATA] {
	return f
}
