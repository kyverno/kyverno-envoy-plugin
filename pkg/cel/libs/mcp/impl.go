package mcp

import (
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno/pkg/cel/utils"
)

type impl struct {
	types.Adapter
}

func (c *impl) mcp_parse(mcp, value ref.Val) ref.Val {
	mcpImpl, err := utils.ConvertToNative[MCP](mcp)
	if err != nil {
		return types.WrapErr(err)
	}

	stringBody, err := utils.ConvertToNative[string](value)
	if err != nil {
		return types.WrapErr(err)
	}

	mcpRequest, err := mcpImpl.Parse([]byte(stringBody))
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest)
}

func (c *impl) mcp_get_string(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[string](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetStringArgument(key, defaultValue))
}

func (c *impl) mcp_get_int(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[int](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetIntArgument(key, defaultValue))
}

func (c *impl) mcp_get_float(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[float64](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetFloatArgument(key, defaultValue))
}

func (c *impl) mcp_get_bool(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[bool](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetBoolArgument(key, defaultValue))
}

func (c *impl) mcp_get_string_slice(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[[]string](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetStringSliceArgument(key, defaultValue))
}

func (c *impl) mcp_get_int_slice(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[[]int](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetIntSliceArgument(key, defaultValue))
}

func (c *impl) mcp_get_float_slice(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[[]float64](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetFloatSliceArgument(key, defaultValue))
}

func (c *impl) mcp_get_bool_slice(args ...ref.Val) ref.Val {
	mcpRequest, err := utils.ConvertToNative[*MCPRequest](args[0])
	if err != nil {
		return types.WrapErr(err)
	}

	key, err := utils.ConvertToNative[string](args[1])
	if err != nil {
		return types.WrapErr(err)
	}

	defaultValue, err := utils.ConvertToNative[[]bool](args[2])
	if err != nil {
		return types.WrapErr(err)
	}

	return c.NativeToValue(mcpRequest.GetBoolSliceArgument(key, defaultValue))
}
