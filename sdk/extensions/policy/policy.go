package policy

import (
	"context"
)

// Policy represents a generic interface for evaluating data according to
// a specific rule, condition, or policy logic.
//
// The interface is fully generic over three type parameters:
//
//	DATA — static or contextual data the policy depends on (e.g., configuration, ruleset)
//	IN   — the input subject to evaluation (e.g., a request, resource, or event)
//	OUT  — the result of evaluation (e.g., a decision, score, or transformed data)
//
// The Evaluate method takes a context (for cancellation, tracing, and metadata),
// the DATA source, and an input value, returning an output value and an error.
//
// Example:
//
//	type Input struct { Value int }
//	type Output struct { Allowed bool }
//
//	var p policy.Policy[[]int, Input, Output]
//	result, err := p.Evaluate(ctx, []int{1, 2, 3}, Input{Value: 2})
type Policy[
	DATA any,
	IN any,
	OUT any,
] interface {
	// Evaluate executes the policy using the provided context, data, and input.
	// It returns an output value and may return an error if evaluation fails.
	Evaluate(context.Context, DATA, IN) (OUT, error)
}

// PolicyFunc is a function adapter that allows a regular Go function with
// the same signature as Policy.Evaluate to satisfy the Policy interface.
//
// This enables you to define lightweight, inline policies without creating
// new struct types.
//
// Example:
//
//	func CheckEven(ctx context.Context, _ struct{}, n int) (bool, error) {
//	    return n%2 == 0, nil
//	}
//
//	p := policy.PolicyFunc[struct{}, int, bool](CheckEven)
//	result, _ := p.Evaluate(context.Background(), struct{}{}, 4)
//	// result == true
type PolicyFunc[
	DATA any,
	IN any,
	OUT any,
] func(context.Context, DATA, IN) (OUT, error)

// Evaluate implements the Policy interface for PolicyFunc.
//
// It simply invokes the wrapped function f with the given parameters.
// This makes PolicyFunc interchangeable with any other Policy implementation.
func (f PolicyFunc[DATA, IN, OUT]) Evaluate(ctx context.Context, data DATA, in IN) (OUT, error) {
	return f(ctx, data, in)
}

// MakePolicyFunc wraps a plain function with the correct signature and returns
// it as a PolicyFunc. This helper is primarily for readability and discoverability,
// and allows easy creation of inline policies with clear intent.
//
// Example:
//
//	check := func(ctx context.Context, conf map[string]bool, input string) (bool, error) {
//	    return conf[input], nil
//	}
//
//	p := policy.MakePolicyFunc(check)
//	allowed, _ := p.Evaluate(context.Background(), map[string]bool{"read": true}, "read")
//	// allowed == true
func MakePolicyFunc[
	DATA any,
	IN any,
	OUT any,
](f func(context.Context, DATA, IN) (OUT, error)) PolicyFunc[DATA, IN, OUT] {
	return f
}
