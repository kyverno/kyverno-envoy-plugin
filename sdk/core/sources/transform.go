package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// NewTransform wraps an existing core.Source and returns a new one that applies
// a transformation function to each element, producing a new source of a different
// (or same) type.
//
// It is conceptually equivalent to a functional “map” operation, enabling
// type-safe data conversion within composable pipelines.
//
// Type parameters:
//
//	IN  — the input type of the original Source.
//	OUT — the output type after applying the transformation.
//
// Parameters:
//
//	inner       — the underlying data source that provides values of type IN.
//	transform — a function that maps each IN to an OUT value.
//
// Returns:
//
//	core.Source[OUT] — a new source that yields transformed elements.
//
// Example:
//
//	src := core.MakeSource(1, 2, 3)
//	strs := sources.NewTransform(src, func(n int) string {
//	    return fmt.Sprintf("num=%d", n)
//	})
//
//	values, _ := strs.Load(context.Background())
//	fmt.Println(values) // Output: ["num=1" "num=2" "num=3"]
//
// Notes:
//
//   - All transformation happens eagerly on Load(ctx).
//   - Errors from the inner source are propagated unchanged.
//   - The returned Source is stateless and thread-safe if transform is.
func NewTransform[IN any, OUT any](
	inner core.Source[IN],
	transform func(IN) OUT,
) core.Source[OUT] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]OUT, error) {
		// Load data from the inner source.
		data, err := inner.Load(ctx)

		// Transform each element.
		transformed := make([]OUT, 0, len(data))
		for _, item := range data {
			transformed = append(transformed, transform(item))
		}

		return transformed, err
	})
}

// NewTransformMultiErr wraps an existing core.Source and returns a new source
// that applies a transformation function to each element. Transformation errors
// are collected and returned using multierr, while successfully transformed
// items are included in the output.
//
// Type parameters:
//
//	IN  — the input type of the original Source.
//	OUT — the output type after applying the transformation.
//
// Parameters:
//
//	inner       — the underlying data source that provides values of type IN.
//	transform — a function that maps each IN to (OUT, error). Any errors
//	              returned are aggregated.
//
// Returns:
//
//	core.Source[OUT] — a new source that yields transformed elements and
//	                   aggregated errors.
//
// Example:
//
//	src := core.MakeSource("1", "2", "x", "4")
//	ints := sources.NewTransformMultiErr(src, func(s string) (int, error) {
//	    return strconv.Atoi(s) // may fail per item
//	})
//
//	values, err := ints.Load(context.Background())
//	fmt.Println(values) // contains successfully converted integers
//	fmt.Println(err)    // aggregated errors for failed items (e.g., "x")
func NewTransformErr[IN any, OUT any](
	inner core.Source[IN],
	transform func(IN) (OUT, error),
) core.Source[OUT] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]OUT, error) {
		// Load all data from the inner source and start with any errors from it.
		data, errs := inner.Load(ctx)

		// Transform each element.
		transformed := make([]OUT, 0, len(data))
		for _, item := range data {
			out, err := transform(item)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
			transformed = append(transformed, out)
		}

		return transformed, errs
	})
}
