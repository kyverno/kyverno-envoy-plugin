package http

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

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register types
		ext.NativeTypes(
			reflect.TypeFor[CheckRequest](),
			reflect.TypeFor[CheckResponse](),
			ext.ParseStructTags(true),
		),
		// extend environment with function overloads
		c.extendEnv,
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (c *lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	impl := impl{
		Adapter: env.CELTypeAdapter(),
	}

	libraryDecls := map[string][]cel.FunctionOpt{
		"http.Response": {
			cel.Overload("http_response", []*cel.Type{cel.IntType}, ResponseType, cel.UnaryBinding(impl.response)),
		},
		"WithHeader": {
			cel.MemberOverload("with_header", []*cel.Type{ResponseType, cel.StringType, cel.StringType}, ResponseType, cel.FunctionBinding(impl.with_header)),
		},
		"WithBody": {
			cel.MemberOverload("with_body", []*cel.Type{ResponseType, cel.StringType}, ResponseType, cel.BinaryBinding(impl.with_body)),
		},
		"Header": {
			cel.MemberOverload("get_header", []*cel.Type{RequestType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_header)),
		},
		"QueryParam": {
			cel.MemberOverload("get_queryparam", []*cel.Type{RequestType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_queryparam)),
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
