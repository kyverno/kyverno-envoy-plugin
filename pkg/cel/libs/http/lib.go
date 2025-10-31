package http

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
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
		"http.Allowed": {
			cel.Overload("http_allow", []*cel.Type{}, ResponseType, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return impl.allowed() })),
		},
		"http.Denied": {
			cel.Overload("http_deny", []*cel.Type{cel.StringType}, ResponseType, cel.UnaryBinding(impl.denied)),
		},
		"Header": {
			cel.MemberOverload("get_header", []*cel.Type{RequestType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_header)),
		},
		"QueryParam": {
			cel.MemberOverload("get_queryparam", []*cel.Type{RequestType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_queryparam)),
		},
		"Write": {
			cel.MemberOverload("write", []*cel.Type{ResponseWriterType, cel.BytesType}, ResponseWriterType, cel.BinaryBinding(impl.write)),
		},
		"WriteHeader": {
			cel.MemberOverload("write_header", []*cel.Type{ResponseWriterType, cel.IntType}, ResponseWriterType, cel.BinaryBinding(impl.write_header)),
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
