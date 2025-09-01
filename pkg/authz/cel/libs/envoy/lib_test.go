package envoy_test

import (
	"reflect"
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/interpreter"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/envoy"
	"github.com/stretchr/testify/assert"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestOkResponse(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   *authv3.CheckResponse
	}{{
		name: "fluent",
		source: `
		envoy
			.Allowed()
			.WithHeader(envoy.Header("foo", "bar").KeepEmptyValue())
			.Response()
			.WithMetadata({"my-new-metadata": "my-new-value"})
		`,
		want: &authv3.CheckResponse{
			Status: &status.Status{
				Code: 0,
			},
			HttpResponse: &authv3.CheckResponse_OkResponse{
				OkResponse: &authv3.OkHttpResponse{
					Headers: []*corev3.HeaderValueOption{{
						Header: &corev3.HeaderValue{
							Key:   "foo",
							Value: "bar",
						},
						KeepEmptyValue: true,
					}},
				},
			},
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"my-new-metadata": structpb.NewStringValue("my-new-value"),
				},
			},
		},
	}, {
		name: "empty",
		want: &authv3.CheckResponse{},
		source: `
		envoy.service.auth.v3.CheckResponse{}
		`,
	}, {
		name: "with status",
		want: &authv3.CheckResponse{
			Status: &status.Status{
				Code: 0,
			},
		},
		source: `
		envoy.service.auth.v3.CheckResponse{
			status: google.rpc.Status{
				code: 0
			}
		}
		`,
	}, {
		name: "with metadata",
		want: &authv3.CheckResponse{
			DynamicMetadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"foo": structpb.NewStringValue("bar"),
				},
			},
		},
		source: `
		envoy.service.auth.v3.CheckResponse{
			dynamic_metadata: {
				"foo": "bar"
			}
		}
		`,
	}, {
		name: "with response",
		want: &authv3.CheckResponse{
			HttpResponse: &authv3.CheckResponse_OkResponse{OkResponse: &authv3.OkHttpResponse{}},
		},
		source: `
		envoy.service.auth.v3.CheckResponse{
			ok_response: envoy.service.auth.v3.OkHttpResponse{}
		}
		`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(envoy.Lib())
			assert.NoError(t, err)
			ast, issues := env.Compile(tt.source)
			assert.Nil(t, issues)
			prog, err := env.Program(ast)
			assert.NoError(t, err)
			assert.NotNil(t, prog)
			out, _, err := prog.Eval(interpreter.EmptyActivation())
			assert.NoError(t, err)
			assert.NotNil(t, out)
			got, err := out.ConvertToNative(reflect.TypeFor[*authv3.CheckResponse]())
			assert.NoError(t, err)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}
