package resulters

import (
	"context"
)

type first[
	POLICY any,
	IN any,
	OUT any,
] struct {
	predicate func(OUT) bool
	result    *OUT
}

func NewFirst[
	POLICY any,
	IN any,
	OUT any,
](predicate func(OUT) bool) *first[POLICY, IN, OUT] {
	return &first[POLICY, IN, OUT]{
		predicate: predicate,
	}
}

func (r *first[POLICY, IN, OUT]) Collect(_ context.Context, _ POLICY, _ IN, out OUT) {
	if r.result == nil {
		if r.predicate(out) {
			r.result = &out
		}
	}
}

func (r *first[POLICY, IN, OUT]) Result() OUT {
	if r.result == nil {
		var out OUT
		return out
	}
	return *r.result
}
