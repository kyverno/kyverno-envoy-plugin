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
	IN any,
	OUT any,
	DATA any,
] = func(context.Context, DATA) Evaluator[POLICY, IN, OUT]

func MakeEvaluatorFactory[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, DATA) Evaluator[POLICY, IN, OUT]) EvaluatorFactory[POLICY, IN, OUT, DATA] {
	return f
}

func MakeEvaluatorFactoryFunc[
	POLICY any,
	IN any,
	OUT any,
	DATA any,
](f func(context.Context, POLICY, IN, DATA) OUT) EvaluatorFactory[POLICY, IN, OUT, DATA] {
	return MakeEvaluatorFactory(func(_ context.Context, data DATA) Evaluator[POLICY, IN, OUT] {
		return MakeEvaluatorFunc(func(ctx context.Context, policy POLICY, input IN) OUT {
			return f(ctx, policy, input, data)
		})
	})
}
