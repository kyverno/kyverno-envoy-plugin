package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

// NewFlatten creates a new core.Source[DATA] that flattens a nested source of slices.
//
// This function is useful when an upstream source produces a collection of lists
// (e.g., a source returning `[][]DATA`), but downstream consumers expect a single
// flat list (`[]DATA`). The resulting source merges all inner slices into one.
//
// The flattening process preserves the order of items as emitted by the inner source.
//
// Parameters:
//
//	inner — a core.Source[[]DATA], i.e. a data source whose elements are slices of DATA.
//
// Returns:
//
//	core.Source[DATA] — a derived source that emits all elements of the inner slices
//	                    as a single flattened list.
//
// Behavior:
//   - Calls Load(ctx) on the inner source to obtain all slices of data.
//   - Iterates through each slice and appends its contents into a single output slice.
//   - Returns the flattened result along with any error returned by the inner source.
//   - Does not short-circuit on error: if the inner source returns partial data and an error,
//     all available data is still flattened and returned alongside the error.
//
// Example:
//
//	// Suppose we have a source producing grouped items:
//	grouped := core.MakeSource([][]int{
//	    {1, 2},
//	    {3, 4, 5},
//	})
//
//	// Flatten it into a single list source:
//	flat := sources.NewFlatten[int](grouped)
//
//	// Load all items:
//	items, err := flat.Load(context.Background())
//	// items == []int{1, 2, 3, 4, 5}
//
//	if err != nil {
//	    log.Println("Error:", err)
//	}
func NewFlatten[DATA any](inner core.Source[[]DATA]) core.Source[DATA] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]DATA, error) {
		// Load the nested slices from the inner source.
		data, err := inner.Load(ctx)

		// Prepare an accumulator for flattened output.
		var output []DATA

		// Append the contents of each inner slice in order.
		for _, item := range data {
			output = append(output, item...)
		}

		// Return the flattened list and any error from the inner source.
		return output, err
	})
}
