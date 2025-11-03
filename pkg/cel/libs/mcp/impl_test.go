package mcp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/mark3labs/mcp-go/mcp"
)

// mockMCPImpl is a mock implementation of MCPImpl for tests
type mockMCPImpl struct {
	parseFn func([]byte) (*MCPRequest, error)
}

func (m *mockMCPImpl) Parse(b []byte) (*MCPRequest, error) {
	return m.parseFn(b)
}

// setupTestEnv creates a CEL environment with MCPRequest type registered and returns an impl instance
func setupTestEnv(t *testing.T) *impl {
	env, err := cel.NewEnv(
		ext.NativeTypes(
			reflect.TypeFor[MCPRequest](),
		),
	)
	if err != nil {
		t.Fatalf("failed creating CEL env: %v", err)
	}
	adapter := env.CELTypeAdapter()
	return &impl{adapter}
}

// setupTestEnvWithMCP creates a CEL environment with both MCP and MCPRequest types registered
func setupTestEnvWithMCP(t *testing.T) *impl {
	env, err := cel.NewEnv(
		ext.NativeTypes(
			reflect.TypeFor[MCP](),
			reflect.TypeFor[MCPRequest](),
		),
	)
	if err != nil {
		t.Fatalf("failed creating CEL env: %v", err)
	}
	adapter := env.CELTypeAdapter()
	return &impl{adapter}
}

// runGetterTest is a generic helper to test getter methods with three parameters (request, key, defaultValue)
func runGetterTest[T any](
	t *testing.T,
	impl *impl,
	testName string,
	mcpRequest *MCPRequest,
	key any,
	defaultValue any,
	want any,
	expectError bool,
	wantErrSubstring string,
	executeFunc func(*impl, ref.Val, ref.Val, ref.Val) ref.Val,
) {
	t.Run(testName, func(t *testing.T) {
		var reqVal any
		if mcpRequest != nil {
			reqVal = mcpRequest
		} else {
			reqVal = struct{}{}
		}
		reqRefVal := impl.NativeToValue(reqVal)
		keyRefVal := impl.NativeToValue(key)
		defaultValueRefVal := impl.NativeToValue(defaultValue)

		got := executeFunc(impl, reqRefVal, keyRefVal, defaultValueRefVal)

		if expectError {
			if got.Type() != types.ErrType {
				t.Errorf("Expected error, got %v (type %v)", got, got.Type())
			}
			if wantErrSubstring != "" {
				if s, ok := got.Value().(error); ok {
					if !strings.Contains(s.Error(), wantErrSubstring) {
						t.Errorf("Expected error containing %q, got %v", wantErrSubstring, s)
					}
				}
			}
		} else {
			native, err := got.ConvertToNative(reflect.TypeFor[T]())
			if err != nil {
				t.Errorf("ConvertToNative err: %v", err)
				return
			}
			gotVal, ok := native.(T)
			if !ok {
				t.Fatalf("Failed to convert output to %T: %T", *new(T), native)
			}
			if !reflect.DeepEqual(gotVal, want) {
				t.Errorf("got %v, want %v", gotVal, want)
			}
		}
	})
}

func TestImpl_mcp_parse(t *testing.T) {
	tests := []struct {
		name        string
		mcpImpl     MCPImpl
		value       any
		want        any
		expectError bool
	}{
		{
			name: "successful parse",
			mcpImpl: &mockMCPImpl{
				parseFn: func(b []byte) (*MCPRequest, error) {
					if string(b) != `{"foo":"bar"}` {
						return nil, fmt.Errorf("unexpected body: %s", string(b))
					}
					return &MCPRequest{Method: "hello"}, nil
				},
			},
			value:       `{"foo":"bar"}`,
			want:        &MCPRequest{Method: "hello"},
			expectError: false,
		},
		{
			name: "Parse returns error",
			mcpImpl: &mockMCPImpl{
				parseFn: func(b []byte) (*MCPRequest, error) {
					return nil, errors.New("parse err")
				},
			},
			value:       `anything`,
			want:        "parse err",
			expectError: true,
		},
		{
			name:        "cannot convert mcp native type",
			mcpImpl:     nil, // purposely pass wrong type
			value:       `{"foo":"bar"}`,
			want:        "type conversion error",
			expectError: true,
		},
		{
			name: "cannot convert value to string",
			mcpImpl: &mockMCPImpl{
				parseFn: func(b []byte) (*MCPRequest, error) { return &MCPRequest{}, nil },
			},
			value:       1234, // should not be convertible to string in utils.ConvertToNative[string]
			want:        "unsupported type conversion",
			expectError: true,
		},
	}

	impl := setupTestEnvWithMCP(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mcpVal any
			if tt.mcpImpl != nil {
				mcpVal = MCP{tt.mcpImpl}
			} else {
				mcpVal = struct{}{} // trigger conversion error
			}
			mcpRefVal := impl.NativeToValue(mcpVal)
			valueRefVal := impl.NativeToValue(tt.value)

			got := impl.mcp_parse(mcpRefVal, valueRefVal)

			if tt.expectError {
				if got.Type() != types.ErrType {
					t.Errorf("Expected error, got %v (type %v)", got, got.Type())
				}
				if tt.want != nil && tt.want != "" {
					if s, ok := got.Value().(error); ok {
						if !strings.Contains(s.Error(), tt.want.(string)) {
							t.Errorf("Expected error containing %q, got %v", tt.want, s)
						}
					}
				}
			} else {
				native, err := got.ConvertToNative(reflect.TypeFor[*MCPRequest]())
				if err != nil {
					t.Errorf("ConvertToNative err: %v", err)
				}
				gotReq, ok := native.(*MCPRequest)
				if !ok {
					t.Fatalf("Failed to convert output to *MCPRequest: %T", native)
				}
				if !reflect.DeepEqual(gotReq, tt.want) {
					t.Errorf("got %v, want %v", gotReq, tt.want)
				}
			}
		})
	}
}

func TestImpl_mcp_get_string(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"name": "test"},
				},
			},
			key:          "name",
			defaultValue: "default",
			want:         "test",
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "name",
			defaultValue: "default",
			want:         "default",
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "name",
			defaultValue: "default",
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"name": "test"},
				},
			},
			key:          1234, // not a string
			defaultValue: "default",
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"name": "test"},
				},
			},
			key:          "name",
			defaultValue: 1234, // not a string
			want:         "unsupported type conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[string](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_string(req, key, def)
		})
	}
}

func TestImpl_mcp_get_int(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get int",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"count": 42},
				},
			},
			key:          "count",
			defaultValue: 0,
			want:         42,
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "count",
			defaultValue: 10,
			want:         10,
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "count",
			defaultValue: 0,
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"count": 42},
				},
			},
			key:          1234, // not a string
			defaultValue: 0,
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to int",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"count": 42},
				},
			},
			key:          "count",
			defaultValue: "not an int",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[int](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_int(req, key, def)
		})
	}
}

func TestImpl_mcp_get_float(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get float",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"double": 3.14},
				},
			},
			key:          "double",
			defaultValue: 0.0,
			want:         3.14,
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "double",
			defaultValue: 1.5,
			want:         1.5,
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "double",
			defaultValue: 0.0,
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"double": 3.14},
				},
			},
			key:          1234, // not a string
			defaultValue: 0.0,
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to float64",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"double": 3.14},
				},
			},
			key:          "double",
			defaultValue: "not a float",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[float64](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_float(req, key, def)
		})
	}
}

func TestImpl_mcp_get_bool(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get bool",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"enabled": true},
				},
			},
			key:          "enabled",
			defaultValue: false,
			want:         true,
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "enabled",
			defaultValue: false,
			want:         false,
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "enabled",
			defaultValue: false,
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"enabled": true},
				},
			},
			key:          1234, // not a string
			defaultValue: false,
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to bool",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"enabled": true},
				},
			},
			key:          "enabled",
			defaultValue: "not a bool",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[bool](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_bool(req, key, def)
		})
	}
}

func TestImpl_mcp_get_string_slice(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get string slice",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"tags": []string{"a", "b", "c"}},
				},
			},
			key:          "tags",
			defaultValue: []string{},
			want:         []string{"a", "b", "c"},
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "tags",
			defaultValue: []string{"default"},
			want:         []string{"default"},
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "tags",
			defaultValue: []string{},
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"tags": []string{"a"}},
				},
			},
			key:          1234, // not a string
			defaultValue: []string{},
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to []string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"tags": []string{"a"}},
				},
			},
			key:          "tags",
			defaultValue: "not a slice",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[[]string](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_string_slice(req, key, def)
		})
	}
}

func TestImpl_mcp_get_int_slice(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get int slice",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"numbers": []int{1, 2, 3}},
				},
			},
			key:          "numbers",
			defaultValue: []int{},
			want:         []int{1, 2, 3},
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "numbers",
			defaultValue: []int{0},
			want:         []int{0},
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "numbers",
			defaultValue: []int{},
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"numbers": []int{1}},
				},
			},
			key:          1234, // not a string
			defaultValue: []int{},
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to []int",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"numbers": []int{1}},
				},
			},
			key:          "numbers",
			defaultValue: "not a slice",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[[]int](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_int_slice(req, key, def)
		})
	}
}

func TestImpl_mcp_get_float_slice(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get float slice",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"values": []float64{1.1, 2.2, 3.3}},
				},
			},
			key:          "values",
			defaultValue: []float64{},
			want:         []float64{1.1, 2.2, 3.3},
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "values",
			defaultValue: []float64{0.0},
			want:         []float64{0.0},
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "values",
			defaultValue: []float64{},
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"values": []float64{1.1}},
				},
			},
			key:          1234, // not a string
			defaultValue: []float64{},
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to []float64",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"values": []float64{1.1}},
				},
			},
			key:          "values",
			defaultValue: "not a slice",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[[]float64](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_float_slice(req, key, def)
		})
	}
}

func TestImpl_mcp_get_bool_slice(t *testing.T) {
	tests := []struct {
		name         string
		mcpRequest   *MCPRequest
		key          any
		defaultValue any
		want         any
		expectError  bool
	}{
		{
			name: "successful get bool slice",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"flags": []bool{true, false, true}},
				},
			},
			key:          "flags",
			defaultValue: []bool{},
			want:         []bool{true, false, true},
			expectError:  false,
		},
		{
			name: "key not found returns default",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"other": "value"},
				},
			},
			key:          "flags",
			defaultValue: []bool{false},
			want:         []bool{false},
			expectError:  false,
		},
		{
			name:         "cannot convert MCPRequest",
			mcpRequest:   nil,
			key:          "flags",
			defaultValue: []bool{},
			want:         "type conversion error",
			expectError:  true,
		},
		{
			name: "cannot convert key to string",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"flags": []bool{true}},
				},
			},
			key:          1234, // not a string
			defaultValue: []bool{},
			want:         "unsupported type conversion",
			expectError:  true,
		},
		{
			name: "cannot convert defaultValue to []bool",
			mcpRequest: &MCPRequest{
				ToolCall: &mcp.CallToolParams{
					Arguments: map[string]any{"flags": []bool{true}},
				},
			},
			key:          "flags",
			defaultValue: "not a slice",
			want:         "unsupported native conversion",
			expectError:  true,
		},
	}

	testImpl := setupTestEnv(t)

	for _, tt := range tests {
		wantErrSubstring := ""
		if tt.expectError && tt.want != nil {
			if s, ok := tt.want.(string); ok {
				wantErrSubstring = s
			}
		}
		runGetterTest[[]bool](t, testImpl, tt.name, tt.mcpRequest, tt.key, tt.defaultValue, tt.want, tt.expectError, wantErrSubstring, func(i *impl, req, key, def ref.Val) ref.Val {
			return i.mcp_get_bool_slice(req, key, def)
		})
	}
}
