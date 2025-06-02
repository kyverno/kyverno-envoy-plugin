package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/types/known/structpb"
)

type impl struct {
	types.Adapter
}

func (c *impl) decode(token ref.Val, key ref.Val) ref.Val {
	t, ok := token.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(token)
	}
	k, ok := key.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(key)
	}
	claimsMap := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(string(t), claimsMap, func(*jwt.Token) (any, error) {
		return []byte(k), nil
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
