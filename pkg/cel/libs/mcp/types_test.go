package mcp

import (
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// requestWithToolCall creates an MCPRequest with ToolCall arguments
func requestWithToolCall(args any) *MCPRequest {
	return &MCPRequest{
		ToolCall: &mcp.CallToolParams{
			Arguments: args,
		},
	}
}

// requestWithResourceRead creates an MCPRequest with ResourceRead arguments
func requestWithResourceRead(args map[string]any) *MCPRequest {
	return &MCPRequest{
		ResourceRead: &mcp.ReadResourceParams{
			Arguments: args,
		},
	}
}

// requestWithPromptGet creates an MCPRequest with PromptGet arguments
func requestWithPromptGet(args map[string]string) *MCPRequest {
	return &MCPRequest{
		PromptGet: &mcp.GetPromptParams{
			Arguments: args,
		},
	}
}

func TestMCPRequest_GetArguments(t *testing.T) {
	tests := []struct {
		name     string
		request  *MCPRequest
		expected map[string]any
	}{
		{
			name: "ToolCall takes priority over ResourceRead",
			request: func() *MCPRequest {
				req := requestWithToolCall(map[string]any{"tool": "value"})
				req.ResourceRead = &mcp.ReadResourceParams{
					Arguments: map[string]any{"resource": "value"},
				}
				return req
			}(),
			expected: map[string]any{"tool": "value"},
		},
		{
			name:     "ResourceRead used when ToolCall is nil",
			request:  requestWithResourceRead(map[string]any{"resource": "value"}),
			expected: map[string]any{"resource": "value"},
		},
		{
			name:     "PromptGet used when ToolCall and ResourceRead are nil",
			request:  requestWithPromptGet(map[string]string{"prompt": "value"}),
			expected: map[string]any{"prompt": "value"},
		},
		{
			name:     "ToolCall with non-map Arguments returns nil",
			request:  requestWithToolCall("not a map"),
			expected: nil,
		},
		{
			name:     "All nil returns nil",
			request:  &MCPRequest{},
			expected: nil,
		},
		{
			name: "PromptGet converts map[string]string to map[string]any",
			request: requestWithPromptGet(map[string]string{
				"key1": "value1",
				"key2": "value2",
			}),
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetArguments()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetArguments() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetRawArguments(t *testing.T) {
	tests := []struct {
		name     string
		request  *MCPRequest
		expected any
	}{
		{
			name: "ToolCall takes priority",
			request: func() *MCPRequest {
				req := requestWithToolCall(map[string]any{"tool": "value"})
				req.ResourceRead = &mcp.ReadResourceParams{
					Arguments: map[string]any{"resource": "value"},
				}
				return req
			}(),
			expected: map[string]any{"tool": "value"},
		},
		{
			name:     "ResourceRead used when ToolCall is nil",
			request:  requestWithResourceRead(map[string]any{"resource": "value"}),
			expected: map[string]any{"resource": "value"},
		},
		{
			name:     "PromptGet used when others are nil",
			request:  requestWithPromptGet(map[string]string{"prompt": "value"}),
			expected: map[string]string{"prompt": "value"},
		},
		{
			name:     "Returns raw type without conversion",
			request:  requestWithToolCall("raw string"),
			expected: "raw string",
		},
		{
			name:     "All nil returns nil",
			request:  &MCPRequest{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetRawArguments()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetRawArguments() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetStringArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "Returns string value when found",
			request:      requestWithToolCall(map[string]any{"name": "test"}),
			key:          "name",
			defaultValue: "default",
			expected:     "test",
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Returns default when value is not string",
			request:      requestWithToolCall(map[string]any{"name": 123}),
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Returns default when no arguments",
			request:      &MCPRequest{},
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetStringArgument(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetStringArgument(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetIntArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "Returns int value when found",
			request:      requestWithToolCall(map[string]any{"count": 42}),
			key:          "count",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Converts float64 to int",
			request:      requestWithToolCall(map[string]any{"count": 42.7}),
			key:          "count",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Converts valid string to int",
			request:      requestWithToolCall(map[string]any{"count": "42"}),
			key:          "count",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Returns default for invalid string",
			request:      requestWithToolCall(map[string]any{"count": "not a number"}),
			key:          "count",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "count",
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetIntArgument(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetIntArgument(%q, %d) = %d, want %d", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetFloatArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue float64
		expected     float64
	}{
		{
			name:         "Returns float64 value when found",
			request:      requestWithToolCall(map[string]any{"price": 3.14}),
			key:          "price",
			defaultValue: 0.0,
			expected:     3.14,
		},
		{
			name:         "Converts int to float64",
			request:      requestWithToolCall(map[string]any{"price": 42}),
			key:          "price",
			defaultValue: 0.0,
			expected:     42.0,
		},
		{
			name:         "Converts valid string to float64",
			request:      requestWithToolCall(map[string]any{"price": "3.14"}),
			key:          "price",
			defaultValue: 0.0,
			expected:     3.14,
		},
		{
			name:         "Returns default for invalid string",
			request:      requestWithToolCall(map[string]any{"price": "not a number"}),
			key:          "price",
			defaultValue: 1.5,
			expected:     1.5,
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "price",
			defaultValue: 1.5,
			expected:     1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetFloatArgument(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetFloatArgument(%q, %f) = %f, want %f", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetBoolArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "Returns bool value when found",
			request:      requestWithToolCall(map[string]any{"enabled": true}),
			key:          "enabled",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Converts valid string to bool",
			request:      requestWithToolCall(map[string]any{"enabled": "true"}),
			key:          "enabled",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Converts int non-zero to true",
			request:      requestWithToolCall(map[string]any{"enabled": 1}),
			key:          "enabled",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Converts int zero to false",
			request:      requestWithToolCall(map[string]any{"enabled": 0}),
			key:          "enabled",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "Converts float64 non-zero to true",
			request:      requestWithToolCall(map[string]any{"enabled": 0.5}),
			key:          "enabled",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Returns default for invalid string",
			request:      requestWithToolCall(map[string]any{"enabled": "not a bool"}),
			key:          "enabled",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "enabled",
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetBoolArgument(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetBoolArgument(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetStringSliceArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "Returns []string value when found",
			request:      requestWithToolCall(map[string]any{"tags": []string{"a", "b", "c"}}),
			key:          "tags",
			defaultValue: []string{},
			expected:     []string{"a", "b", "c"},
		},
		{
			name:         "Converts []any with strings to []string",
			request:      requestWithToolCall(map[string]any{"tags": []any{"a", "b", "c"}}),
			key:          "tags",
			defaultValue: []string{},
			expected:     []string{"a", "b", "c"},
		},
		{
			name:         "Filters out non-string items from []any",
			request:      requestWithToolCall(map[string]any{"tags": []any{"a", 123, "b"}}),
			key:          "tags",
			defaultValue: []string{},
			expected:     []string{"a", "b"},
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
		{
			name:         "Returns default when value is not a slice",
			request:      requestWithToolCall(map[string]any{"tags": "not a slice"}),
			key:          "tags",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetStringSliceArgument(tt.key, tt.defaultValue)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetStringSliceArgument(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetIntSliceArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue []int
		expected     []int
	}{
		{
			name:         "Returns []int value when found",
			request:      requestWithToolCall(map[string]any{"numbers": []int{1, 2, 3}}),
			key:          "numbers",
			defaultValue: []int{},
			expected:     []int{1, 2, 3},
		},
		{
			name:         "Converts []any with ints to []int",
			request:      requestWithToolCall(map[string]any{"numbers": []any{1, 2, 3}}),
			key:          "numbers",
			defaultValue: []int{},
			expected:     []int{1, 2, 3},
		},
		{
			name:         "Converts []any with float64s to []int",
			request:      requestWithToolCall(map[string]any{"numbers": []any{1.0, 2.5, 3.9}}),
			key:          "numbers",
			defaultValue: []int{},
			expected:     []int{1, 2, 3},
		},
		{
			name:         "Converts []any with valid string numbers to []int",
			request:      requestWithToolCall(map[string]any{"numbers": []any{"1", "2", "3"}}),
			key:          "numbers",
			defaultValue: []int{},
			expected:     []int{1, 2, 3},
		},
		{
			name:         "Filters out invalid values from []any",
			request:      requestWithToolCall(map[string]any{"numbers": []any{1, "invalid", 3}}),
			key:          "numbers",
			defaultValue: []int{},
			expected:     []int{1, 3},
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "numbers",
			defaultValue: []int{0},
			expected:     []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetIntSliceArgument(tt.key, tt.defaultValue)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetIntSliceArgument(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetFloatSliceArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue []float64
		expected     []float64
	}{
		{
			name:         "Returns []float64 value when found",
			request:      requestWithToolCall(map[string]any{"values": []float64{1.1, 2.2, 3.3}}),
			key:          "values",
			defaultValue: []float64{},
			expected:     []float64{1.1, 2.2, 3.3},
		},
		{
			name:         "Converts []any with float64s to []float64",
			request:      requestWithToolCall(map[string]any{"values": []any{1.1, 2.2, 3.3}}),
			key:          "values",
			defaultValue: []float64{},
			expected:     []float64{1.1, 2.2, 3.3},
		},
		{
			name:         "Converts []any with ints to []float64",
			request:      requestWithToolCall(map[string]any{"values": []any{1, 2, 3}}),
			key:          "values",
			defaultValue: []float64{},
			expected:     []float64{1.0, 2.0, 3.0},
		},
		{
			name:         "Converts []any with valid string numbers to []float64",
			request:      requestWithToolCall(map[string]any{"values": []any{"1.1", "2.2", "3.3"}}),
			key:          "values",
			defaultValue: []float64{},
			expected:     []float64{1.1, 2.2, 3.3},
		},
		{
			name:         "Filters out invalid values from []any",
			request:      requestWithToolCall(map[string]any{"values": []any{1.1, "invalid", 3.3}}),
			key:          "values",
			defaultValue: []float64{},
			expected:     []float64{1.1, 3.3},
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "values",
			defaultValue: []float64{0.0},
			expected:     []float64{0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetFloatSliceArgument(tt.key, tt.defaultValue)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetFloatSliceArgument(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMCPRequest_GetBoolSliceArgument(t *testing.T) {
	tests := []struct {
		name         string
		request      *MCPRequest
		key          string
		defaultValue []bool
		expected     []bool
	}{
		{
			name:         "Returns []bool value when found",
			request:      requestWithToolCall(map[string]any{"flags": []bool{true, false, true}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false, true},
		},
		{
			name:         "Converts []any with bools to []bool",
			request:      requestWithToolCall(map[string]any{"flags": []any{true, false, true}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false, true},
		},
		{
			name:         "Converts []any with valid string bools to []bool",
			request:      requestWithToolCall(map[string]any{"flags": []any{"true", "false", "true"}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false, true},
		},
		{
			name:         "Converts []any with ints to []bool (non-zero=true, zero=false)",
			request:      requestWithToolCall(map[string]any{"flags": []any{1, 0, -1}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false, true},
		},
		{
			name:         "Converts []any with float64s to []bool (non-zero=true, zero=false)",
			request:      requestWithToolCall(map[string]any{"flags": []any{1.5, 0.0, -0.5}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false, true},
		},
		{
			name:         "Filters out invalid string values from []any",
			request:      requestWithToolCall(map[string]any{"flags": []any{true, "invalid", false}}),
			key:          "flags",
			defaultValue: []bool{},
			expected:     []bool{true, false},
		},
		{
			name:         "Returns default when key not found",
			request:      requestWithToolCall(map[string]any{"other": "value"}),
			key:          "flags",
			defaultValue: []bool{false},
			expected:     []bool{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetBoolSliceArgument(tt.key, tt.defaultValue)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetBoolSliceArgument(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
