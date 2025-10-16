package middlewares

// import (
// 	"context"

// 	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
// )

// func StopOnSourceError[
// 	POLICY any,
// 	IN any,
// 	OUT any,
// 	DATA any,
// ]() core.MiddlewareFactory[POLICY, IN, OUT, DATA] {
// 	return func(ctx context.Context, policies []POLICY, err error, next handler) handler {
// 		return sdk.HandlerFunc[int, int, int](func(ctx context.Context, in int, data int) []sdk.Response[int] {
// 			if err != nil {
// 				return nil
// 			}
// 			return next.Handle(ctx, in, data)
// 		})
// 	}

// }
