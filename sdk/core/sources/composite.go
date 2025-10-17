package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// composite represents a collection of core.Source[DATA] instances.
//
// When loaded, it sequentially calls Load on each contained source and
// merges all results into a single slice. Any errors from the sources are
// aggregated using multierr.
//
// This type is unexported because users should typically create it via
// the NewComposite constructor, which ensures proper usage and type safety.
type composite[DATA any] []core.Source[DATA]

// NewComposite creates a new composite source from multiple core.Source[DATA] instances.
//
// When the composite is loaded, it will sequentially invoke Load on each
// provided source and merge their results into a single slice. Errors from
// any of the underlying sources are combined using multierr.
//
// Example:
//
//	src1 := core.MakeSource(1, 2)
//	src2 := core.MakeSource(3, 4)
//	combined := sources.NewComposite(src1, src2)
//
//	data, err := combined.Load(context.Background())
//	// data = [1 2 3 4]
//	// err may contain multiple aggregated errors if any source failed
func NewComposite[DATA any](sources ...core.Source[DATA]) composite[DATA] {
	return composite[DATA](sources)
}

// Load implements the core.Source[DATA] interface for composite.
//
// It performs the following steps:
//  1. Iterates through all contained sources.
//  2. Calls Load on each source.
//  3. Appends successfully loaded items to the output slice.
//  4. Aggregates all errors using multierr.Append without stopping early.
//
// This approach ensures maximum data availability: even if some sources fail,
// successfully loaded items from other sources are returned.
//
// Parameters:
//
//	ctx — context for cancellation or timeout propagation.
//
// Returns:
//
//	[]DATA — combined items from all underlying sources.
//	error  — aggregated errors from sources that failed, or nil if all succeeded.
//
// Example:
//
//	data, err := composite.Load(ctx)
//	// data contains concatenated results from all sub-sources
//	// err may include multiple aggregated errors, or nil if all succeeded
func (s composite[DATA]) Load(ctx context.Context) ([]DATA, error) {
	var out []DATA // Accumulates all successfully loaded items
	var errs error // Aggregates all errors from sub-sources

	for _, source := range s {
		// Load items from the current source
		items, err := source.Load(ctx)

		// Append any error from this source to the aggregated error
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		// Append successfully loaded items to the output slice
		out = append(out, items...)
	}

	return out, errs
}
