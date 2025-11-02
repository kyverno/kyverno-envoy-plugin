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

func (*lib) LibraryName() string {
	return "kyverno.authz.http"
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
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"http.Allowed": {
			cel.Overload("http_allowed", []*cel.Type{}, ResponseOkType, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return impl.allowed() })),
		},
		"http.Denied": {
			cel.Overload("http_denied_string", []*cel.Type{cel.StringType}, ResponseDeniedType, cel.UnaryBinding(impl.denied)),
		},
		"Header": {
			cel.MemberOverload("http_get_header_string", []*cel.Type{RequestAttributesType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_header)),
		},
		"QueryParam": {
			cel.MemberOverload("http_get_queryparam_string", []*cel.Type{RequestAttributesType, cel.StringType}, types.NewListType(cel.StringType), cel.BinaryBinding(impl.get_queryparam)),
		},
		"Response": {
			cel.MemberOverload("http_response_ok", []*cel.Type{ResponseOkType}, ResponseType, cel.UnaryBinding(impl.response_ok)),
			cel.MemberOverload("http_response_denied", []*cel.Type{ResponseDeniedType}, ResponseType, cel.UnaryBinding(impl.response_denied)),
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
