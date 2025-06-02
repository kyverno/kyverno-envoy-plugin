package jwk

import (
	"context"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type impl struct {
	types.Adapter
}

func (c *impl) fetch(from ref.Val) ref.Val {
	if from, err := utils.ConvertToNative[string](from); err != nil {
		return types.WrapErr(err)
	} else {
		set, err := jwk.Fetch(context.Background(), from)
		if err != nil {
			return types.WrapErr(err)
		}
		return c.NativeToValue(set)
	}

}
