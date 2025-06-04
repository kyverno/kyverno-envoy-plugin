package jwt

import (
	"time"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	jwklib "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/jwk"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"google.golang.org/protobuf/types/known/structpb"
)

type impl struct {
	types.Adapter
}

func (c *impl) decode_string_string(token ref.Val, key ref.Val) ref.Val {
	if token, err := utils.ConvertToNative[string](token); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else {
		set := jwk.NewSet()
		key, err := jwk.Import([]byte(key))
		if err != nil {
			return types.WrapErr(err)
		}
		if err := set.AddKey(key); err != nil {
			return types.WrapErr(err)
		}
		tok, err := jwt.Parse(
			[]byte(token),
			jwt.WithValidate(false),
			jwt.WithKeySet(
				set,
				jws.WithUseDefault(true),
				jws.WithInferAlgorithmFromKey(true),
			),
		)
		if err != nil {
			return types.WrapErr(err)
		}
		var claims *structpb.Struct
		if keys := tok.Keys(); len(keys) > 0 {
			fields := make(map[string]any, len(keys))
			for _, key := range keys {
				var value any
				err := tok.Get(key, &value)
				if err != nil {
					return types.WrapErr(err)
				}
				switch value := value.(type) {
				case time.Time:
					fields[key] = value.Unix()
				default:
					fields[key] = value
				}
			}
			encoded, err := structpb.NewStruct(fields)
			if err != nil {
				return types.WrapErr(err)
			}
			claims = encoded
		}
		return c.NativeToValue(
			Token{
				Claims: claims,
				Valid:  jwt.Validate(tok) == nil,
			},
		)
	}
}

func (c *impl) decode_string_set(token ref.Val, set ref.Val) ref.Val {
	if token, err := utils.ConvertToNative[string](token); err != nil {
		return types.WrapErr(err)
	} else if set, err := utils.ConvertToNative[jwklib.Set](set); err != nil {
		return types.WrapErr(err)
	} else {
		tok, err := jwt.Parse(
			[]byte(token),
			jwt.WithValidate(false),
			jwt.WithKeySet(
				set.Set,
				jws.WithUseDefault(true),
				jws.WithInferAlgorithmFromKey(true),
			),
		)
		if err != nil {
			return types.WrapErr(err)
		}
		var claims *structpb.Struct
		if keys := tok.Keys(); len(keys) > 0 {
			fields := make(map[string]any, len(keys))
			for _, key := range keys {
				var value any
				err := tok.Get(key, &value)
				if err != nil {
					return types.WrapErr(err)
				}
				switch value := value.(type) {
				case time.Time:
					fields[key] = value.Unix()
				default:
					fields[key] = value
				}
			}
			encoded, err := structpb.NewStruct(fields)
			if err != nil {
				return types.WrapErr(err)
			}
			claims = encoded
		}
		return c.NativeToValue(
			Token{
				Claims: claims,
				Valid:  jwt.Validate(tok) == nil,
			},
		)
	}
}
