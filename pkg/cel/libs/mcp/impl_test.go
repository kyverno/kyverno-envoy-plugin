package mcp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
)

// mockMCPImpl is a mock implementation of MCPImpl for tests
type mockMCPImpl struct {
	parseFn func([]byte) (*MCPRequest, error)
}

func (m *mockMCPImpl) Parse(b []byte) (*MCPRequest, error) {
	return m.parseFn(b)
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

	// Setup env with NativeTypes and types
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
	impl := &impl{adapter}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mcpVal any
			if tt.mcpImpl != nil {
				mcpVal = MCP{tt.mcpImpl}
			} else {
				mcpVal = struct{}{} // trigger conversion error
			}
			mcpRefVal := adapter.NativeToValue(mcpVal)
			valueRefVal := adapter.NativeToValue(tt.value)

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
