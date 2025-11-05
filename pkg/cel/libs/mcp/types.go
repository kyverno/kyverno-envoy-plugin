package mcp

import (
	"strconv"

	"github.com/google/cel-go/common/types"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	MCPType        = types.NewOpaqueType("mcp.MCP")
	MCPRequestType = types.NewObjectType("mcp.MCPRequest")
)

type MCPImpl interface {
	Parse([]byte) (*MCPRequest, error)
}

type MCP struct {
	MCPImpl
}

type MCPRequest struct {
	Method string
	ID     mcp.RequestId

	Paginated *mcp.PaginatedParams // For all list methods

	// Tools
	ToolCall *mcp.CallToolParams

	// Resources
	ResourceRead        *mcp.ReadResourceParams
	ResourceSubscribe   *mcp.SubscribeParams
	ResourceUnsubscribe *mcp.UnsubscribeParams

	// Prompts
	PromptGet *mcp.GetPromptParams

	// Lifecycle
	Initialize *mcp.InitializeParams

	// Sampling
	CreateMessage *mcp.CreateMessageParams

	// Elicitation
	Elicitation *mcp.ElicitationParams

	// Utilities
	Complete    *mcp.CompleteParams
	SetLogLevel *mcp.SetLevelParams
}

// GetArguments returns the Arguments from ToolCall, ResourceRead, or PromptGet as map[string]any
// Priority: ToolCall > ResourceRead > PromptGet
// If none are available or Arguments is not a map, it returns nil
func (r *MCPRequest) GetArguments() map[string]any {
	// Check ToolCall first
	if r.ToolCall != nil {
		if args, ok := r.ToolCall.Arguments.(map[string]any); ok {
			return args
		}
	}
	// Check ResourceRead
	if r.ResourceRead != nil && r.ResourceRead.Arguments != nil {
		return r.ResourceRead.Arguments
	}
	// Check PromptGet (convert map[string]string to map[string]any)
	if r.PromptGet != nil && r.PromptGet.Arguments != nil {
		result := make(map[string]any, len(r.PromptGet.Arguments))
		for k, v := range r.PromptGet.Arguments {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetRawArguments returns the Arguments from ToolCall, ResourceRead, or PromptGet as-is without type conversion
// Priority: ToolCall > ResourceRead > PromptGet
func (r *MCPRequest) GetRawArguments() any {
	if r.ToolCall != nil {
		return r.ToolCall.Arguments
	}
	if r.ResourceRead != nil {
		return r.ResourceRead.Arguments
	}
	if r.PromptGet != nil {
		return r.PromptGet.Arguments
	}
	return nil
}

// GetString returns a string argument by key, or the default value if not found
func (r *MCPRequest) GetStringArgument(key string, defaultValue string) string {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt returns an int argument by key, or the default value if not found
func (r *MCPRequest) GetIntArgument(key string, defaultValue int) int {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// GetFloat returns a float64 argument by key, or the default value if not found
func (r *MCPRequest) GetFloatArgument(key string, defaultValue float64) float64 {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return defaultValue
}

// GetBool returns a bool argument by key, or the default value if not found
func (r *MCPRequest) GetBoolArgument(key string, defaultValue bool) bool {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		case int:
			return v != 0
		case float64:
			return v != 0
		}
	}
	return defaultValue
}

// GetStringSlice returns a string slice argument by key, or the default value if not found
func (r *MCPRequest) GetStringSliceArgument(key string, defaultValue []string) []string {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case []string:
			return v
		case []any:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return defaultValue
}

// GetIntSlice returns an int slice argument by key, or the default value if not found
func (r *MCPRequest) GetIntSliceArgument(key string, defaultValue []int) []int {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case []int:
			return v
		case []any:
			result := make([]int, 0, len(v))
			for _, item := range v {
				switch num := item.(type) {
				case int:
					result = append(result, num)
				case float64:
					result = append(result, int(num))
				case string:
					if i, err := strconv.Atoi(num); err == nil {
						result = append(result, i)
					}
				}
			}
			return result
		}
	}
	return defaultValue
}

// GetFloatSlice returns a float64 slice argument by key, or the default value if not found
func (r *MCPRequest) GetFloatSliceArgument(key string, defaultValue []float64) []float64 {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case []float64:
			return v
		case []any:
			result := make([]float64, 0, len(v))
			for _, item := range v {
				switch num := item.(type) {
				case float64:
					result = append(result, num)
				case int:
					result = append(result, float64(num))
				case string:
					if f, err := strconv.ParseFloat(num, 64); err == nil {
						result = append(result, f)
					}
				}
			}
			return result
		}
	}
	return defaultValue
}

// GetBoolSlice returns a bool slice argument by key, or the default value if not found
func (r *MCPRequest) GetBoolSliceArgument(key string, defaultValue []bool) []bool {
	args := r.GetArguments()
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case []bool:
			return v
		case []any:
			result := make([]bool, 0, len(v))
			for _, item := range v {
				switch b := item.(type) {
				case bool:
					result = append(result, b)
				case string:
					if parsed, err := strconv.ParseBool(b); err == nil {
						result = append(result, parsed)
					}
				case int:
					result = append(result, b != 0)
				case float64:
					result = append(result, b != 0)
				}
			}
			return result
		}
	}
	return defaultValue
}
