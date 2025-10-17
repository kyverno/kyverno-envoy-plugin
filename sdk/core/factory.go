package core

import "context"

// Factory defines a generic function type used to create or derive an object
// of type OUT based on contextual information about a policy evaluation.
//
// A Factory is typically responsible for constructing an evaluator or
// other higher-level component that depends on both policy metadata and
// runtime data inputs. It can be used to build reusable logic pipelines
// for dynamic policy processing.
//
// Type Parameters:
//   - POLICY — the policy definition or rule being evaluated
//   - DATA   — the external data or context relevant to policy evaluation
//   - IN     — the input object or request subject to policy evaluation
//   - OUT    — the result or constructed object produced by the factory
//
// Signature:
//
//	func(context.Context, FactoryContext[POLICY, DATA, IN]) OUT
//
// The function receives:
//   - context.Context — allows cancellation, deadlines, and contextual values
//   - FactoryContext  — provides the current policy, data, and input context
//
// Example:
//
//	var factory core.Factory[MyPolicy, MyData, MyInput, MyEvaluator]
//	factory = func(ctx context.Context, fctx core.FactoryContext[MyPolicy, MyData, MyInput]) MyEvaluator {
//	    return NewMyEvaluator(fctx.Policy, fctx.Data)
//	}
//
//	evaluator := factory(ctx, myFactoryContext)
//
// This abstraction enables flexible composition of dynamic evaluation logic
// without coupling to specific types or implementations.
type Factory[
	POLICY any,
	DATA any,
	IN any,
	OUT any,
] = func(context.Context, FactoryContext[POLICY, DATA, IN]) OUT
