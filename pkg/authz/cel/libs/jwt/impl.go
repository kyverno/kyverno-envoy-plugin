package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

type impl struct {
	types.Adapter
}

func (c *impl) decode(token ref.Val, key ref.Val) ref.Val {
	if token, err := utils.ConvertToNative[string](token); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else {
		claimsMap := jwt.MapClaims{}
		parsed, err := jwt.ParseWithClaims(token, claimsMap, func(*jwt.Token) (any, error) {
			return []byte(key), nil
		})
		if err != nil {
			return c.NativeToValue(nil)
		}
		header, err := structpb.NewStruct(parsed.Header)
		if err != nil {
			return types.WrapErr(err)
		}
		claims, err := structpb.NewStruct(claimsMap)
		if err != nil {
			return types.WrapErr(err)
		}
		return c.NativeToValue(
			Token{
				Header: header,
				Claims: claims,
				Valid:  parsed.Valid,
			},
		)
	}
}
