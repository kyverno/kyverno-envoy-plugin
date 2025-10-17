package sources

import (
	"context"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

// once is a concurrency-safe wrapper around a core.Source that ensures
// the underlying data source is loaded at most once.
//
// Subsequent calls to Load return the same cached data and error as the first
// invocation, without reloading or calling the inner source again.
//
// This is useful for expensive or deterministic sources whose data does not
// change across multiple calls, e.g., configuration or static metadata.
//
// The type is unexported; users should construct it via NewOnce.
type once[DATA any] struct {
	called bool
	lock   sync.Mutex
	inner  core.Source[DATA]
	data   []DATA
	err    error
}

// NewOnce wraps a given core.Source[DATA] with a caching layer that only
// executes Load once.
//
// Example:
//
//	src := core.MakeSource(1, 2, 3)
//	onceSrc := sources.NewOnce(src)
//
//	// First call loads from src
//	data, err := onceSrc.Load(context.Background())
//
//	// Subsequent calls return cached result, no re-execution
//	again, _ := onceSrc.Load(context.Background())
//
// This is thread-safe — multiple goroutines can safely call Load concurrently;
// only the first call will invoke the underlying source.
func NewOnce[DATA any](inner core.Source[DATA]) *once[DATA] {
	return &once[DATA]{
		inner: inner,
	}
}

// Load implements the core.Source[DATA] interface for the once type.
//
// It ensures the underlying source is loaded only once, caching both the data
// and the error result. Subsequent calls return the cached values directly.
//
// The method uses a mutex to protect against concurrent calls. If multiple
// goroutines invoke Load simultaneously before the first call finishes, only
// the first one will trigger the actual load; others will block until the
// result is available.
//
// Example:
//
//	onceSrc := sources.NewOnce(expensiveSource)
//	data, err := onceSrc.Load(ctx) // triggers the load
//	again, err := onceSrc.Load(ctx) // returns cached result
func (p *once[DATA]) Load(ctx context.Context) ([]DATA, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.called {
		p.data, p.err = p.inner.Load(ctx)
		p.called = true
	}
	return p.data, p.err
}
