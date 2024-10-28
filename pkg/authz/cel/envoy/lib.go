package envoy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type lib struct{}

func Lib() cel.Library {
	return &lib{}
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Types((*authv3.CheckRequest)(nil), (*authv3.CheckResponse)(nil)),
		c.extenEnv,
	}
}

func (_ *lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func (_ *lib) extenEnv(env *cel.Env) (*cel.Env, error) {
	impl := impl{
		Adapter: env.CELTypeAdapter(),
	}
	libraryDecls := map[string][]cel.FunctionOpt{
		"envoy.Allowed": {
			cel.Overload("allowed", []*cel.Type{}, _ok, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return impl.allowed() })),
		},
		"envoy.Denied": {
			cel.Overload("denied", []*cel.Type{types.IntType}, _denied, cel.UnaryBinding(impl.denied)),
		},
		"envoy.Response": {
			cel.Overload("response_ok", []*cel.Type{_ok}, _response, cel.UnaryBinding(impl.response_ok)),
			cel.Overload("response_denied", []*cel.Type{_denied}, _response, cel.UnaryBinding(impl.response_denied)),
		},
		"envoy.Header": {
			cel.Overload("header_key_value", []*cel.Type{types.StringType, types.StringType}, _header, cel.BinaryBinding(impl.header_key_value)),
		},
		"WithBody": {
			cel.MemberOverload("denied_with_body", []*cel.Type{_denied, types.StringType}, _denied, cel.BinaryBinding(impl.denied_with_body)),
		},
		"WithHeader": {
			cel.MemberOverload("ok_with_header", []*cel.Type{_ok, _header}, _ok, cel.BinaryBinding(impl.ok_with_header)),
			cel.MemberOverload("denied_with_header", []*cel.Type{_denied, _header}, _denied, cel.BinaryBinding(impl.denied_with_header)),
		},
		"WithoutHeader": {
			cel.MemberOverload("ok_without_header", []*cel.Type{_ok, types.StringType}, _ok, cel.BinaryBinding(impl.ok_without_header)),
		},
		"WithResponseHeader": {
			cel.MemberOverload("ok_with_response_header", []*cel.Type{_ok, _header}, _ok, cel.BinaryBinding(impl.ok_with_response_header)),
		},
		"KeepEmptyValue": {
			cel.MemberOverload("header_keep_empty_value", []*cel.Type{_header}, _header, cel.UnaryBinding(impl.header_keep_empty_value)),
			cel.MemberOverload("header_keep_empty_value_bool", []*cel.Type{_header, types.BoolType}, _header, cel.BinaryBinding(impl.header_keep_empty_value_bool)),
		},
		"Response": {
			cel.MemberOverload("ok_response", []*cel.Type{_ok}, _response, cel.UnaryBinding(impl.response_ok)),
			cel.MemberOverload("denied_response", []*cel.Type{_denied}, _response, cel.UnaryBinding(impl.response_denied)),
		},
		"WithMetadata": {
			cel.MemberOverload("response_with_metadata", []*cel.Type{_response, _struct}, _ok, cel.BinaryBinding(impl.response_with_metadata)),
		},
	}
	options := []cel.EnvOption{}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	return env.Extend(options...)
}
