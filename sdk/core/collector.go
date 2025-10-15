package core

import "context"

type Collector[
	POLICY any,
	IN any,
	OUT any,
] interface {
	Collect(context.Context, POLICY, IN, OUT)
}
