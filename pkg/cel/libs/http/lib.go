package http

import (
	"reflect"

	"github.com/google/cel-go/cel"
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
		ext.NativeTypes(reflect.TypeFor[Req](), reflect.TypeFor[Resp](), reflect.TypeFor[KV](), ext.ParseStructTags(true)),
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
			cel.Overload("http_response", []*cel.Type{}, ResponseType, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return impl.response() })),
		},
		"WithHeader": {
			cel.MemberOverload("with_header", []*cel.Type{ResponseType, cel.StringType, cel.StringType}, ResponseType, cel.FunctionBinding(impl.with_header)),
		},
		"WithBody": {
			cel.MemberOverload("with_body", []*cel.Type{ResponseType, cel.StringType}, ResponseType, cel.BinaryBinding(impl.with_body)),
		},
		"WithStatus": {
			cel.MemberOverload("with_status", []*cel.Type{ResponseType, cel.IntType}, ResponseType, cel.BinaryBinding(impl.with_status)),
		},
		"Get": {
			cel.MemberOverload("get_header_value", []*cel.Type{KVType, cel.StringType}, cel.StringType, cel.BinaryBinding(impl.get_header_value)),
		},
		"GetAll": {
			cel.MemberOverload("get_header_all", []*cel.Type{KVType, cel.StringType}, cel.ListType(cel.StringType), cel.BinaryBinding(impl.get_header_all)),
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
