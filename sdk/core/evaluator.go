package core

import "context"

type Evaluator[
	POLICY any,
	IN any,
	OUT any,
] interface {
	Evaluate(context.Context, POLICY, IN) OUT
}

type EvaluatorFunc[
	POLICY any,
	IN any,
	OUT any,
] func(context.Context, POLICY, IN) OUT

func (f EvaluatorFunc[POLICY, IN, OUT]) Evaluate(ctx context.Context, policy POLICY, input IN) OUT {
	return f(ctx, policy, input)
}

func MakeEvaluatorFunc[
	POLICY any,
	IN any,
	OUT any,
](f func(context.Context, POLICY, IN) OUT) EvaluatorFunc[POLICY, IN, OUT] {
	return f
}

type EvaluatorFactory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = Factory[POLICY, DATA, IN, Evaluator[POLICY, IN, OUT]]
