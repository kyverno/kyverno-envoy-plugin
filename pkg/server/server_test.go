package server

import (
	"context"
	"reflect"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func Test_extAuthzServerV3_Check(t *testing.T) {
	type fields struct {
		policies []string
	}
	type args struct {
		ctx context.Context
		req *authv3.CheckRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *authv3.CheckResponse
		wantErr bool
	}{
		{
			name: "Allow admin GET method request at path /book",
			fields: fields{
				policies: []string{"./../../tests/policies/policy.yaml"},
			},
			args: args{
				ctx: context.Background(),
				req: &authv3.CheckRequest{
					Attributes: &authv3.AttributeContext{
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
								Path:   "/book",
								Headers: map[string]string{
									"Content-Type":  "application/json",
									"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0",
								},
							},
						},
					},
				},
			},
			want: &authv3.CheckResponse{
				Status: &status.Status{Code: int32(codes.OK)},
			},
			wantErr: false,
		},
		{
			name: "Deny guest GET method request at path /book",
			fields: fields{
				policies: []string{"./../../tests/policies/policy.yaml"},
			},
			args: args{
				ctx: context.Background(),
				req: &authv3.CheckRequest{
					Attributes: &authv3.AttributeContext{
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
								Path:   "/book",
								Headers: map[string]string{
									"Content-Type":  "application/json",
									"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk",
								},
							},
						},
					},
				},
			},
			want: &authv3.CheckResponse{
				Status: &status.Status{
					Code: int32(codes.PermissionDenied),
				},
				HttpResponse: &authv3.CheckResponse_DeniedResponse{
					DeniedResponse: &authv3.DeniedHttpResponse{
						Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden},
						Body:   "Request denied by Kyverno JSON engine. Reason: -> GET method calls at path /book are not allowed to guests users\n -> any[0].check.request.http.headers.authorization.(split(@, ' ')[1]).(jwt_decode(@ , 'secret').payload.role): Invalid value: \"guest\": Expected value: \"admin\"",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &extAuthzServerV3{
				policies: tt.fields.policies,
			}
			got, err := s.Check(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("extAuthzServerV3.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extAuthzServerV3.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
