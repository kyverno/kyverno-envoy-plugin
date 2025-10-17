package sources

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// composite is a generic type that represents a collection of core.Source[DATA]
// instances. When loaded, it aggregates the results of each contained source.
//
// This type is unexported (lowercase) because users should typically create it
// via the NewComposite constructor.
type composite[DATA any] []core.Source[DATA]

// NewComposite creates a new composite source that combines multiple
// core.Source[DATA] instances into one.
//
// When the composite is loaded, it will sequentially invoke Load on each
// provided source and merge their results into a single slice.
//
// Example:
//
//	src1 := core.MakeSource(1, 2)
//	src2 := core.MakeSource(3, 4)
//	combined := sources.NewComposite(src1, src2)
//	data, err := combined.Load(context.Background())
//	// data = [1 2 3 4]
//
// If any of the underlying sources return an error, all such errors are combined
// using go.uber.org/multierr and returned together.
func NewComposite[DATA any](sources ...core.Source[DATA]) composite[DATA] {
	return composite[DATA](sources)
}

// Load implements the core.Source[DATA] interface for composite.
//
// It iterates through all contained sources, calling Load on each.
// - Successfully loaded items are appended to the output slice.
// - Any errors encountered are aggregated via multierr.Append.
//
// The method never stops on error; it continues loading from all sources to
// ensure maximum data availability.
//
// Example:
//
//	data, err := composite.Load(ctx)
//	// data contains concatenated results from all sub-sources
//	// err may include multiple aggregated errors, or nil if all succeeded.
func (s composite[DATA]) Load(ctx context.Context) ([]DATA, error) {
	var out []DATA
	var errs error
	for _, source := range s {
		items, err := source.Load(ctx)
		if err != nil {
			errs = multierr.Append(errs, err)
		} else {
			out = append(out, items...)
		}
	}
	return out, errs
}
