package authz

import (
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/stretchr/testify/assert"
)

func Test_convert(t *testing.T) {
	tests := []struct {
		name    string
		attrs   *authv3.AttributeContext
		want    map[string]any
		wantErr bool
	}{
		{
			name: "empty AttributeContext",
			attrs: &authv3.AttributeContext{
				Source:            nil,
				Destination:       nil,
				Request:           nil,
				ContextExtensions: nil,
			},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "non-empty AttributeContext",
			attrs: &authv3.AttributeContext{
				Source: &authv3.AttributeContext_Peer{
					Principal: "test-principal",
				},
				Destination: &authv3.AttributeContext_Peer{
					Principal: "test-destination",
				},
				Request: &authv3.AttributeContext_Request{
					Time: nil,
					Http: &authv3.AttributeContext_HttpRequest{
						Method: "GET",
						Path:   "/test",
					},
				},
				ContextExtensions: map[string]string{
					"test-key": "test-value",
				},
			},
			want: map[string]any{
				"source": map[string]any{
					"principal": "test-principal",
				},
				"destination": map[string]any{
					"principal": "test-destination",
				},
				"request": map[string]any{
					"http": map[string]any{
						"method": "GET",
						"path":   "/test",
					},
				},
				"contextExtensions": map[string]any{
					"test-key": "test-value",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(tt.attrs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
