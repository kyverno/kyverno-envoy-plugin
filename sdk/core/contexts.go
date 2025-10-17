package core

// SourceContext represents the result of loading data from a Source.
//
// It encapsulates both the successfully loaded data and any error encountered
// during the loading process, allowing downstream consumers to handle partial
// success scenarios gracefully.
//
// Type Parameter:
//   - DATA — the type of items returned by the source.
//
// Example:
//
//	src := core.MakeSource(1, 2, 3)
//	data, err := src.Load(ctx)
//	ctx := core.MakeSourceContext(data, err)
//	// ctx.Data = [1,2,3], ctx.Error = err (if any)
type SourceContext[DATA any] struct {
	// Data contains the successfully loaded items from the source.
	Data []DATA

	// Error represents any error encountered while loading the data.
	// It may be nil if the operation completed successfully.
	Error error
}

// MakeSourceContext constructs a new SourceContext instance from the provided
// data and error values.
//
// This is primarily a convenience helper for more concise and explicit
// initialization of SourceContext objects.
//
// Example:
//
//	ctx := core.MakeSourceContext([]string{"a", "b"}, nil)
func MakeSourceContext[DATA any](data []DATA, err error) SourceContext[DATA] {
	return SourceContext[DATA]{
		Data:  data,
		Error: err,
	}
}

// FactoryContext provides a strongly typed structure encapsulating all contextual
// information required when building or evaluating a policy factory.
//
// It combines the following:
//   - The source context for all loaded policies.
//   - The external data context (e.g., configuration, state).
//   - The input being evaluated.
//
// This type is typically passed into Factory functions and EvaluatorFactories.
//
// Type Parameters:
//   - POLICY — the type representing a policy or rule definition.
//   - DATA   — the external data or configuration type used during evaluation.
//   - IN     — the input object or subject being evaluated.
//
// Example:
//
//	srcCtx := core.MakeSourceContext([]Policy{p1, p2}, nil)
//	fctx := core.MakeFactoryContext(srcCtx, envConfig, request)
//	evaluator := myFactory(ctx, fctx)
type FactoryContext[
	POLICY any,
	DATA any,
	IN any,
] struct {
	// Source holds the result of loading the policy set (data + error).
	Source SourceContext[POLICY]

	// Data contains additional contextual information or configuration
	// relevant to the policy evaluation.
	Data DATA

	// Input represents the runtime subject under evaluation, such as a
	// request, resource, or object.
	Input IN
}

// MakeFactoryContext constructs a new FactoryContext using the provided
// source, data, and input values.
//
// This helper improves clarity when building contexts for evaluator factories.
//
// Example:
//
//	fctx := core.MakeFactoryContext(sourceCtx, env, req)
func MakeFactoryContext[
	POLICY any,
	DATA any,
	IN any,
](source SourceContext[POLICY], data DATA, in IN) FactoryContext[POLICY, DATA, IN] {
	return FactoryContext[POLICY, DATA, IN]{
		Source: source,
		Data:   data,
		Input:  in,
	}
}
