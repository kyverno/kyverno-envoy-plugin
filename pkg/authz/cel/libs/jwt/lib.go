package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
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
		// extend environment with function overloads
		c.extendEnv,
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (*lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"jwt.Decode": {
			cel.Overload("decode_string_string", []*cel.Type{types.StringType, types.StringType}, types.DynType, cel.BinaryBinding(decode)),
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

func decode(token ref.Val, key ref.Val) ref.Val {
	t, ok := token.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(token)
	}
	k, ok := key.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(key)
	}
	parsed, err := jwt.Parse(string(t), func(*jwt.Token) (any, error) {
		return []byte(k), nil
	})
	if err != nil {
		return types.DefaultTypeAdapter.NativeToValue(nil)
	}
	return types.DefaultTypeAdapter.NativeToValue(map[string]any{
		"header": parsed.Header,
		"claims": parsed.Claims.(jwt.MapClaims),
		"valid":  parsed.Valid,
	})
}
