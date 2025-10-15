package resulters

import (
	"context"
)

type appender[
	POLICY any,
	IN any,
	OUT any,
] struct {
	result []OUT
}

func (r *appender[POLICY, IN, OUT]) Collect(_ context.Context, _ POLICY, _ IN, out OUT) {
	r.result = append(r.result, out)
}

func (r *appender[POLICY, IN, OUT]) Result() []OUT {
	return r.result
}

func NewAppender[
	POLICY any,
	IN any,
	OUT any,
]() *appender[POLICY, IN, OUT] {
	return &appender[POLICY, IN, OUT]{}
}
