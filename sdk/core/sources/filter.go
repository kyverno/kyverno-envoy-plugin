package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// NewFilter wraps an existing core.Source and returns a new one that yields only
// the elements matching a given predicate function.
//
// This allows filtering of any data source — whether static, cached, or dynamically
// loaded — without modifying its implementation. It is particularly useful for
// narrowing down datasets before evaluation.
//
// Type parameter:
//
//	DATA — the type of elements produced by the underlying Source.
//
// Parameters:
//
//	inner     — the underlying data source to wrap.
//	predicate — a filtering function that returns true for elements to keep.
//
// Returns:
//
//	core.Source[DATA] — a new source that produces only filtered elements.
//
// Example:
//
//	src := core.MakeSource(1, 2, 3, 4, 5)
//	even := sources.NewFilter(src, func(n int) bool { return n%2 == 0 })
//
//	values, _ := even.Load(context.Background())
//	fmt.Println(values) // Output: [2 4]
//
// Notes:
//
//   - The filter runs eagerly upon Load(ctx), meaning all data from the inner source
//     is first loaded before the predicate is applied.
//   - Errors from the inner source are propagated unchanged to the caller.
//   - The returned Source maintains no additional state.
func NewFilter[DATA any](inner core.Source[DATA], predicate func(DATA) bool) core.Source[DATA] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]DATA, error) {
		// Load all data from the inner source.
		data, err := inner.Load(ctx)

		// Apply filtering logic.
		filtered := make([]DATA, 0, len(data))
		for _, item := range data {
			if predicate(item) {
				filtered = append(filtered, item)
			}
		}

		// Return the filtered slice and any error from the inner source.
		return filtered, err
	})
}

// NewFilterErr wraps an existing core.Source and returns a new source that
// filters elements using a predicate function that may return an error.
// Successfully passing items are returned, and all errors are aggregated using multierr.
//
// Type parameters:
//
//	DATA — the type of elements produced by the underlying Source.
//
// Parameters:
//
//	inner     — the underlying data source providing values of type DATA.
//	predicate — a function that returns (bool, error). True means the item
//	            is included; errors are collected.
//
// Returns:
//
//	core.Source[DATA] — a new source that yields filtered items and aggregated errors.
//
// Example:
//
//	src := core.MakeSource("1", "2", "x", "4")
//	filtered := sources.NewFilterErr(src, func(s string) (bool, error) {
//	    if s == "x" {
//	        return false, fmt.Errorf("invalid value: %q", s)
//	    }
//	    return true, nil
//	})
//
//	values, err := filtered.Load(context.Background())
//	fmt.Println(values) // ["1", "2", "4"]
//	fmt.Println(err)    // aggregated error for "x"
func NewFilterErr[DATA any](
	inner core.Source[DATA],
	predicate func(DATA) (bool, error),
) core.Source[DATA] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]DATA, error) {
		// Load all data from the inner source and start with any errors from it.
		data, errs := inner.Load(ctx)

		// Apply filtering logic.
		filtered := make([]DATA, 0, len(data))
		for _, item := range data {
			ok, err := predicate(item)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
			if ok {
				filtered = append(filtered, item)
			}
		}

		return filtered, errs
	})
}
