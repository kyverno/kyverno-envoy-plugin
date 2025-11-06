package mcp

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
	"github.com/mark3labs/mcp-go/mcp"
)

type lib struct {
	mcp MCP
}

func Lib(mcp MCPImpl) cel.EnvOption {
	return cel.Lib(&lib{
		mcp: MCP{mcp},
	})
}

func (*lib) LibraryName() string {
	return "kyverno.mcp"
}

func (l *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("mcp", MCPType),
		// Use native types to register custom structs with the ext native type provider
		ext.NativeTypes(
			reflect.TypeFor[MCP](),
			reflect.TypeFor[MCPRequest](),
			ext.ParseStructTags(true),
		),
		// extend environment with function overloads
		l.extendEnv,
	}
}

func (l *lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Globals(
			map[string]any{
				"mcp": l.mcp,
			},
		),
	}
}

func (l *lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	impl := impl{
		Adapter: env.CELTypeAdapter(),
	}

	mcpConstants := map[string]string{
		"InitializeMethod":             string(mcp.MethodInitialize),
		"PingMethod":                   string(mcp.MethodPing),
		"ResourcesListMethod":          string(mcp.MethodResourcesList),
		"ResourcesTemplatesListMethod": string(mcp.MethodResourcesTemplatesList),
		"ResourcesReadMethod":          string(mcp.MethodResourcesRead),
		"PromptsListMethod":            string(mcp.MethodPromptsList),
		"PromptsGetMethod":             string(mcp.MethodPromptsGet),
		"ToolsListMethod":              string(mcp.MethodToolsList),
		"ToolsCallMethod":              string(mcp.MethodToolsCall),
		"SetLogLevelMethod":            string(mcp.MethodSetLogLevel),
		"ElicitationCreateMethod":      string(mcp.MethodElicitationCreate),
	}

	libraryDecls := map[string][]cel.FunctionOpt{
		"Parse": {
			cel.MemberOverload("mcp_parse_bytes_dyn",
				[]*cel.Type{MCPType, types.DynType},
				MCPRequestType,
				cel.BinaryBinding(impl.mcp_parse),
			),
		},
		"GetStringArgument": {
			cel.MemberOverload("mcp_get_string",
				[]*cel.Type{MCPRequestType, types.StringType, types.StringType},
				types.StringType,
				cel.FunctionBinding(impl.mcp_get_string),
			),
		},
		"GetIntArgument": {
			cel.MemberOverload("mcp_get_int",
				[]*cel.Type{MCPRequestType, types.IntType, types.IntType},
				types.IntType,
				cel.FunctionBinding(impl.mcp_get_int),
			),
		},
		"GetFloatArgument": {
			cel.MemberOverload("mcp_get_float",
				[]*cel.Type{MCPRequestType, types.DoubleType, types.DoubleType},
				types.DoubleType,
				cel.FunctionBinding(impl.mcp_get_float),
			),
		},
		"GetBoolArgument": {
			cel.MemberOverload("mcp_get_bool",
				[]*cel.Type{MCPRequestType, types.BoolType, types.BoolType},
				types.BoolType,
				cel.FunctionBinding(impl.mcp_get_bool),
			),
		},
		"GetStringSliceArgument": {
			cel.MemberOverload("mcp_get_string_slice",
				[]*cel.Type{MCPRequestType, types.StringType, types.StringType},
				cel.ListType(types.StringType),
				cel.FunctionBinding(impl.mcp_get_string_slice),
			),
		},
		"GetIntSliceArgument": {
			cel.MemberOverload("mcp_get_int_slice",
				[]*cel.Type{MCPRequestType, types.IntType, types.IntType},
				cel.ListType(types.IntType),
				cel.FunctionBinding(impl.mcp_get_int_slice),
			),
		},
		"GetFloatSliceArgument": {
			cel.MemberOverload("mcp_get_float_slice",
				[]*cel.Type{MCPRequestType, types.DoubleType, types.DoubleType},
				cel.ListType(types.DoubleType),
				cel.FunctionBinding(impl.mcp_get_float_slice),
			),
		},
		"GetBoolSliceArgument": {
			cel.MemberOverload("mcp_get_bool_slice",
				[]*cel.Type{MCPRequestType, types.BoolType, types.BoolType},
				cel.ListType(types.BoolType),
				cel.FunctionBinding(impl.mcp_get_bool_slice),
			),
		},
	}

	options := []cel.EnvOption{}

	for name, value := range mcpConstants {
		options = append(options, cel.Constant(fmt.Sprintf("mcp.%s", name), types.StringType, types.String(value)))
	}

	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}

	return env.Extend(options...)
}
