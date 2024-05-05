package request

import (
	"reflect"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

func TestConvert(t *testing.T) {
	type args struct {
		attrs *authv3.AttributeContext
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty AttributeContext",
			args: args{
				attrs: &authv3.AttributeContext{
					Source:            nil,
					Destination:       nil,
					Request:           nil,
					ContextExtensions: nil,
				},
			},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "non-empty AttributeContext",
			args: args{
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
			},
			want: map[string]interface{}{
				"source": map[string]interface{}{
					"principal": "test-principal",
				},
				"destination": map[string]interface{}{
					"principal": "test-destination",
				},
				"request": map[string]interface{}{
					"http": map[string]interface{}{
						"method": "GET",
						"path":   "/test",
					},
				},
				"contextExtensions": map[string]interface{}{
					"test-key": "test-value",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
