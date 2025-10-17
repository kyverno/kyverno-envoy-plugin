package sources

import (
	"context"
	"fmt"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// NewCache returns a source that caches items internally using a double-buffered
// strategy. The cache is rebuilt on each Load while reusing existing cached items.
// Any stale entries (no longer produced by the inner source) are automatically discarded.
//
// Both keyFunc and cacheFunc are context-aware and may return errors. All errors
// are aggregated using multierr.
//
// Type parameters:
//
//	DATA — type of elements produced by the inner source
//	KEY  — type of cache keys computed from DATA
//	ITEM — type of cached items returned
//
// Parameters:
//
//	inner      — inner source producing DATA
//	keyFunc    — function to compute cache key from DATA and context; may return error
//	cacheFunc  — function to produce cached ITEM from KEY, DATA and context; may return error
func NewCache[DATA any, KEY comparable, ITEM any](
	inner core.Source[DATA],
	keyFunc func(context.Context, DATA) (KEY, error),
	cacheFunc func(context.Context, KEY, DATA) (ITEM, error),
) core.Source[ITEM] {
	// Mutex to protect concurrent access to the cache
	var lock sync.Mutex

	// Two buffers for double buffering
	cacheA := make(map[KEY]ITEM)
	cacheB := make(map[KEY]ITEM)

	writeBuffer := cacheA
	readBuffer := cacheB

	// Return a core.Source[ITEM] that implements the Load method
	return core.MakeSourceFunc(func(ctx context.Context) ([]ITEM, error) {
		var errs error         // Aggregate all errors
		out := make([]ITEM, 0) // Output slice for cached items

		// Step 1: Load all data from the inner source
		data, innerErr := inner.Load(ctx)
		if innerErr != nil {
			errs = multierr.Append(errs, innerErr)
		}

		// Step 2: Lock for atomic cache update
		lock.Lock()
		defer lock.Unlock() // Ensure lock is released even if panic or continue occurs

		// Step 3: Process each DATA item
		for _, item := range data {
			// Compute the cache key using the context
			key, err := keyFunc(ctx, item)
			if err != nil {
				// Append key computation errors and skip this item
				errs = multierr.Append(errs, fmt.Errorf("key error for item %+v: %w", item, err))
				continue
			}

			// Try to reuse an existing cached item from the previous buffer
			cached, ok := readBuffer[key]
			if !ok {
				// Build a new cached item if it does not exist
				cached, err = cacheFunc(ctx, key, item)
				if err != nil {
					// Append cache function errors and skip this item
					errs = multierr.Append(errs, fmt.Errorf("cache error for key %v: %w", key, err))
					continue
				}
			}

			// Append the cached item to the output and write to the new buffer
			out = append(out, cached)
			writeBuffer[key] = cached
		}

		// Step 4: Swap buffers
		writeBuffer, readBuffer = readBuffer, writeBuffer

		// Step 5: Clear the new write buffer for the next load
		for k := range writeBuffer {
			delete(writeBuffer, k)
		}

		// Return all cached items and aggregated errors
		return out, errs
	})
}
