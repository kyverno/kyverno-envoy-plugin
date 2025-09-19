package utils

import (
	"cmp"
	"slices"

	"golang.org/x/exp/maps"
)

func ToSortedSlice[K cmp.Ordered, V any](in map[K]V) []V {
	if in == nil {
		return nil
	}
	out := make([]V, 0, len(in))
	keys := maps.Keys(in)
	slices.Sort(keys)
	for _, key := range keys {
		out = append(out, in[key])
	}
	return out
}
