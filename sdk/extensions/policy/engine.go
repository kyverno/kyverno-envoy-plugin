package policy

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/defaults"
)

// NewEngine constructs a new core.Engine instance for evaluating policies.
//
// It connects a source of policies with a standardized evaluation handler built
// from the EvaluatorFactory function. The resulting engine can load, evaluate,
// and aggregate policies in a uniform way, producing results wrapped in a
// defaults.Result structure.
//
// Generic type parameters:
//
//	POLICY — a concrete type implementing Policy[DATA, IN, OUT]
//	DATA   — contextual data shared by all policy evaluations (e.g., configuration)
//	IN     — input type to evaluate (e.g., requests, resources, or events)
//	OUT    — raw result type returned by a policy before being wrapped in Evaluation
//
// Parameters:
//
//	source — a core.Source that provides one or more POLICY instances to the engine.
//
// Returns:
//
//	core.Engine[DATA, IN, defaults.Result[POLICY, DATA, IN, Evaluation[OUT]]]
//	— an evaluation engine ready to load policies, apply them to inputs, and
//	  return structured results containing both data and evaluation metadata.
//
// The resulting engine uses:
//   - EvaluatorFactory[POLICY]() to produce evaluators that wrap Policy.Evaluate
//     outputs into Evaluation[OUT] values.
//   - defaults.Handler(...) to build the default result-handling pipeline.
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
//	// Define a policy source (static for simplicity).
//	src := core.MakeSource(AccessPolicy{})
//
//	// Build an engine for AccessPolicy.
//	eng := policy.NewEngine[AccessPolicy, map[string]bool, string, bool](src)
//
//	// Use the engine to evaluate an input.
//	result := eng.Evaluate(context.Background(), map[string]bool{"alice": true}, "alice")
//	fmt.Println(result.Output.Result) // true
func NewEngine[
	POLICY Policy[DATA, IN, OUT],
	DATA any,
	IN any,
	OUT any,
](
	source core.Source[POLICY],
) core.Engine[DATA, IN, defaults.Result[POLICY, DATA, IN, Evaluation[OUT]]] {
	return core.NewEngine(
		source,
		// Create a default handler using the policy evaluator factory.
		defaults.Handler(EvaluatorFactory[POLICY]()),
	)
}
