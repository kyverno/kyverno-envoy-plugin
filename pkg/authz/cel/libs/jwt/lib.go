package jwt

import (
	"reflect"

	"github.com/golang-jwt/jwt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"google.golang.org/protobuf/types/known/structpb"
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (*lib) LibraryName() string {
	return "kyverno.jwt"
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register token type
		ext.NativeTypes(reflect.TypeFor[Token]()),
		// extend environment with function overloads
		c.extendEnv,
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

var TokenType = types.NewObjectType("jwt.Token")

type Token struct {
	Header *structpb.Struct
	Claims *structpb.Struct
	Valid  bool
}

func (*lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	adapter := env.CELTypeAdapter()
	decode := func(token ref.Val, key ref.Val) ref.Val {
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
			return adapter.NativeToValue(nil)
		}
		header, err := structpb.NewStruct(parsed.Header)
		if err != nil {
			return types.WrapErr(err)
		}
		claims, err := structpb.NewStruct(claimsMap)
		if err != nil {
			return types.WrapErr(err)
		}
		return adapter.NativeToValue(
			Token{
				Header: header,
				Claims: claims,
				Valid:  parsed.Valid,
			},
		)
	}
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"jwt.Decode": {
			cel.Overload("decode_string_string", []*cel.Type{types.StringType, types.StringType}, TokenType, cel.BinaryBinding(decode)),
		},
	}
	// create env options corresponding to our function overloads
	options := []cel.EnvOption{}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	// extend environment with our function overloads
	return env.Extend(options...)
}
