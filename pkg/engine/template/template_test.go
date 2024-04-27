package template

import (
	"context"
	"reflect"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

func TestExecute_JWTDecode(t *testing.T) {
	// Define a valid JWT token and secret key
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"

	// Define an invalid JWT token
	invalidToken := "invalid.jwt.token"

	tests := []struct {
		name      string
		statement string
		value     any
		bindings  binding.Bindings
		opts      []Option
		expected  any
		shouldErr bool
	}{
		{
			name:      "Valid JWT token",
			statement: "jwt_decode(@ , 'secret')",
			value:     validToken,
			bindings:  nil,
			opts:      nil,
			expected: map[string]interface{}{
				"header": map[string]interface{}{
					"alg": "HS256",
					"typ": "JWT",
				},
				"payload": map[string]interface{}{
					"exp":  2.241081539e+09,
					"nbf":  1.514851139e+09,
					"role": "guest",
					"sub":  "YWxpY2U=",
				},
				"sig": "6a61316267764974343733393362615f576253426d33354e72556864784d346d4f56514e3869587a386c6b",
			},
			shouldErr: false,
		},
		{
			name:      "Invalid JWT token",
			statement: "jwt_decode(@, 'secret')",
			value:     invalidToken,
			bindings:  nil,
			opts:      nil,
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Execute(context.Background(), test.statement, test.value, test.bindings, test.opts...)
			if test.shouldErr && err == nil {
				t.Errorf("Expected error but got nil")
			} else if !test.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if !test.shouldErr && !equalMaps(result.(map[string]interface{}), test.expected.(map[string]interface{})) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestExecute_JWTDecodepayloadrole(t *testing.T) {
	// Define a valid JWT token and secret key
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"

	// Define an invalid JWT token
	invalidToken := "invalid.jwt.token"

	tests := []struct {
		name      string
		statement string
		value     any
		bindings  binding.Bindings
		opts      []Option
		expected  any
		shouldErr bool
	}{
		{
			name:      "Valid JWT token",
			statement: "jwt_decode(@ , 'secret').payload.role",
			value:     validToken,
			bindings:  nil,
			opts:      nil,
			expected:  "guest",
			shouldErr: false,
		},
		{
			name:      "Invalid JWT token",
			statement: "jwt_decode(@, 'secret').payload.role",
			value:     invalidToken,
			bindings:  nil,
			opts:      nil,
			expected:  nil,
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Execute(context.Background(), test.statement, test.value, test.bindings, test.opts...)
			if test.shouldErr && err == nil {
				t.Errorf("Expected error but got nil")
			} else if !test.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if !test.shouldErr && !reflect.DeepEqual(result.(string), test.expected.(string)) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func equalMaps(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || !reflect.DeepEqual(v1, v2) {
			return false
		}
	}

	return true
}
