package jwk

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (*lib) LibraryName() string {
	return "kyverno.jwk"
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register native types
		ext.NativeTypes(reflect.TypeFor[Set]()),
		// extend environment with function overloads
		c.extendEnv,
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (*lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	// get env type adapter
	adapter := env.CELTypeAdapter()
	// create implementation with adapter
	impl := impl{adapter}
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"jwks.Fetch": {
			cel.Overload("fetch_string", []*cel.Type{types.StringType}, SetType, cel.UnaryBinding(impl.fetch)),
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
