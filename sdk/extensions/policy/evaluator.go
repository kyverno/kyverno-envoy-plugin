package policy

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

// EvaluatorFactory creates a type-safe adapter that integrates policy evaluation
// with the core evaluation framework.
//
// It returns a core.EvaluatorFactory that produces evaluators capable of executing
// a given Policy against an input, using preloaded contextual data. The evaluator
// wraps the result and any error into an Evaluation[OUT] struct for standardized
// downstream handling.
//
// Generic type parameters:
//
//	POLICY — a concrete type that implements Policy[DATA, IN, OUT]
//	DATA   — the type of contextual data shared by the policy (e.g., configuration)
//	IN     — the type of input value being evaluated (e.g., a request or object)
//	OUT    — the type of result produced by the policy evaluation
//
// Returned factory function structure:
//
//	func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN])
//	    core.Evaluator[POLICY, IN, Evaluation[OUT]]
//
// This pattern makes it easy to register generic policies with an evaluation
// engine without manually writing boilerplate for data injection and result wrapping.
//
// Example:
//
//	type AccessPolicy struct{}
//
//	func (p AccessPolicy) Evaluate(ctx context.Context, rules map[string]bool, user string) (bool, error) {
//	    allowed, ok := rules[user]
//	    if !ok {
//	        return false, fmt.Errorf("user %q not found", user)
//	    }
//	    return allowed, nil
//	}
//
//	// Create an evaluator factory for AccessPolicy.
//	factory := policy.EvaluatorFactory[AccessPolicy, map[string]bool, string, bool]()
//
//	// Use the factory to create an evaluator.
//	evaluator := factory(context.Background(), core.FactoryContext[AccessPolicy, map[string]bool, string]{
//	    Data: map[string]bool{"alice": true},
//	})
//
//	// Evaluate a policy instance with an input.
//	result := evaluator.Evaluate(context.Background(), AccessPolicy{}, "alice")
//	fmt.Println(result.Result) // true
func EvaluatorFactory[
	POLICY Policy[DATA, IN, OUT],
	DATA any,
	IN any,
	OUT any,
]() core.EvaluatorFactory[POLICY, DATA, IN, Evaluation[OUT]] {
	return func(ctx context.Context, fctx core.FactoryContext[POLICY, DATA, IN]) core.Evaluator[POLICY, IN, Evaluation[OUT]] {
		return core.MakeEvaluatorFunc(func(ctx context.Context, policy POLICY, in IN) Evaluation[OUT] {
			// Execute the policy’s Evaluate method using the contextual data and input.
			out, err := policy.Evaluate(ctx, fctx.Data, in)

			// Wrap the result and error in a standardized Evaluation struct.
			return MakeEvaluation(out, err)
		})
	}
}
