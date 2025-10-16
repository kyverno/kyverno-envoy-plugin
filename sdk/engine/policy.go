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

func (f PolicyFunc[DATA, IN, OUT]) Evaluate(ctx context.Context, data DATA, in IN) (OUT, error) {
	return f(ctx, data, in)
}

func MakePolicyFunc[
	DATA any,
	IN any,
	OUT any,
](f func(context.Context, DATA, IN) (OUT, error)) PolicyFunc[DATA, IN, OUT] {
	return f
}
