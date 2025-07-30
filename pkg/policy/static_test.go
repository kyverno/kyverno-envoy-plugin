package policy_test

import (
	"context"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"testing"
	"testing/fstest"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/policy"
	"github.com/stretchr/testify/assert"
)

func TestNewFsProvider(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		expectedError  string
		expectedCount  int
		mockCompileErr bool
	}{
		{
			name: "valid single policy",
			files: map[string]string{
				"policy1.yaml": `
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: test-policy
spec:
  allow:
    - response: 'envoy.Allowed().Response()'`,
			},
			expectedCount: 1,
		},
		{
			name: "multiple valid policies in single file",
			files: map[string]string{
				"policies.yaml": `
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: policy1
spec:
  allow:
    - response: 'envoy.Allowed().Response()'
---
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: policy2
spec:
  allow:
    - response: 'envoy.Allowed().Response()'`,
			},
			expectedCount: 2,
		},
		{
			name: "ignore non-yaml files",
			files: map[string]string{
				"policy.yaml": `
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: test-policy
spec:
  allow:
    - response: 'envoy.Allowed().Response()'`,
				"readme.txt": "This should be ignored",
			},
			expectedCount: 1,
		},
		{
			name: "invalid yaml content",
			files: map[string]string{
				"invalid.yaml": `invalid: yaml:`,
			},
			expectedCount: 0,
		},
		{
			name:          "directory access",
			files:         nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer un système de fichiers virtuel pour les tests
			fsys := fstest.MapFS{}
			for name, content := range tt.files {
				fsys[name] = &fstest.MapFile{
					Data: []byte(content),
					Mode: 0644,
				}
			}

			// Mock compiler
			mockCompiler := &MockCompiler{
				shouldError: tt.mockCompileErr,
			}

			// Create provider
			provider := policy.NewFsProvider(mockCompiler, fsys)
			policies, err := provider.CompiledPolicies(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, policies, tt.expectedCount)
		})
	}
}

// MockCompiler est un mock du Compiler pour les tests
type MockCompiler struct {
	shouldError bool
}

func (m *MockCompiler) Compile(policy *v1alpha1.AuthorizationPolicy) (policy.CompiledPolicy, field.ErrorList) {
	if m.shouldError {
		return nil, field.ErrorList{
			&field.Error{
				Type:   field.ErrorTypeInternal,
				Field:  "",
				Detail: "Internal compiler error",
			},
		}
	}
	return &MockCompiledPolicy{}, nil
}

// MockCompiledPolicy est un mock de CompiledPolicy pour les tests
type MockCompiledPolicy struct{}

func (m *MockCompiledPolicy) For(r *authv3.CheckRequest) (policy.AllowFunc, policy.DenyFunc) {
	return func() (*authv3.CheckResponse, error) { return nil, nil },
		func() (*authv3.CheckResponse, error) { return nil, nil }
}
