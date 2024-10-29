package envoy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// envoy types
var (
	CheckRequest       = types.NewObjectType("envoy.service.auth.v3.CheckRequest")
	CheckResponse      = types.NewObjectType("envoy.service.auth.v3.CheckResponse")
	OkHttpResponse     = types.NewObjectType("envoy.service.auth.v3.OkHttpResponse")
	DeniedHttpResponse = types.NewObjectType("envoy.service.auth.v3.DeniedHttpResponse")
	Metadata           = types.NewObjectType("google.protobuf.Struct")
	HeaderValueOption  = types.NewObjectType("envoy.config.core.v3.HeaderValueOption")
	QueryParameter     = types.NewObjectType("envoy.config.core.v3.QueryParameter")
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register envoy protobuf messages
		cel.Types((*authv3.CheckRequest)(nil), (*authv3.CheckResponse)(nil)),
		// extend environment with function overloads
		c.extenEnv,
	}
}

func (_ *lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (_ *lib) extenEnv(env *cel.Env) (*cel.Env, error) {
	// create implementation, recording the envoy types aware adapter
	impl := impl{
		Adapter: env.CELTypeAdapter(),
	}
	// build our function overloads
	libraryDecls := map[string][]cel.FunctionOpt{
		"envoy.Allowed": {
			cel.Overload("allowed", []*cel.Type{}, OkHttpResponse, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return impl.allowed() })),
		},
		"envoy.Denied": {
			cel.Overload("denied", []*cel.Type{types.IntType}, DeniedHttpResponse, cel.UnaryBinding(impl.denied)),
		},
		"envoy.Response": {
			cel.Overload("response_code", []*cel.Type{types.IntType}, CheckResponse, cel.UnaryBinding(impl.response_code)),
			cel.Overload("response_ok", []*cel.Type{OkHttpResponse}, CheckResponse, cel.UnaryBinding(impl.response_ok)),
			cel.Overload("response_denied", []*cel.Type{DeniedHttpResponse}, CheckResponse, cel.UnaryBinding(impl.response_denied)),
		},
		"envoy.Null": {
			cel.Overload("null", []*cel.Type{}, CheckResponse, cel.FunctionBinding(func(values ...ref.Val) ref.Val {
				return impl.Adapter.NativeToValue((*authv3.CheckResponse)(nil))
			})),
		},
		"envoy.Header": {
			cel.Overload("header_key_value", []*cel.Type{types.StringType, types.StringType}, HeaderValueOption, cel.BinaryBinding(impl.header_key_value)),
		},
		"WithBody": {
			cel.MemberOverload("denied_with_body", []*cel.Type{DeniedHttpResponse, types.StringType}, DeniedHttpResponse, cel.BinaryBinding(impl.denied_with_body)),
		},
		"WithHeader": {
			cel.MemberOverload("ok_with_header", []*cel.Type{OkHttpResponse, HeaderValueOption}, OkHttpResponse, cel.BinaryBinding(impl.ok_with_header)),
			cel.MemberOverload("denied_with_header", []*cel.Type{DeniedHttpResponse, HeaderValueOption}, DeniedHttpResponse, cel.BinaryBinding(impl.denied_with_header)),
		},
		"WithoutHeader": {
			cel.MemberOverload("ok_without_header", []*cel.Type{OkHttpResponse, types.StringType}, OkHttpResponse, cel.BinaryBinding(impl.ok_without_header)),
		},
		"WithResponseHeader": {
			cel.MemberOverload("ok_with_response_header", []*cel.Type{OkHttpResponse, HeaderValueOption}, OkHttpResponse, cel.BinaryBinding(impl.ok_with_response_header)),
		},
		"WithQueryParam": {
			cel.MemberOverload("ok_with_query_param", []*cel.Type{OkHttpResponse, QueryParameter}, OkHttpResponse, cel.BinaryBinding(impl.ok_with_query_param)),
		},
		"WithoutQueryParam": {
			cel.MemberOverload("ok_without_query_param", []*cel.Type{OkHttpResponse, types.StringType}, OkHttpResponse, cel.BinaryBinding(impl.ok_without_query_param)),
		},
		"KeepEmptyValue": {
			cel.MemberOverload("header_keep_empty_value", []*cel.Type{HeaderValueOption}, HeaderValueOption, cel.UnaryBinding(impl.header_keep_empty_value)),
			cel.MemberOverload("header_keep_empty_value_bool", []*cel.Type{HeaderValueOption, types.BoolType}, HeaderValueOption, cel.BinaryBinding(impl.header_keep_empty_value_bool)),
		},
		"Response": {
			cel.MemberOverload("ok_response", []*cel.Type{OkHttpResponse}, CheckResponse, cel.UnaryBinding(impl.response_ok)),
			cel.MemberOverload("denied_response", []*cel.Type{DeniedHttpResponse}, CheckResponse, cel.UnaryBinding(impl.response_denied)),
		},
		"WithMessage": {
			cel.MemberOverload("response_with_message", []*cel.Type{CheckResponse, types.StringType}, CheckResponse, cel.BinaryBinding(impl.response_with_message)),
		},
		"WithMetadata": {
			cel.MemberOverload("response_with_metadata", []*cel.Type{CheckResponse, Metadata}, CheckResponse, cel.BinaryBinding(impl.response_with_metadata)),
		},
	}
	// creatte env options corresponding to our function overloads
	options := []cel.EnvOption{}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	// extend environment with our function overloads
	return env.Extend(options...)
}
