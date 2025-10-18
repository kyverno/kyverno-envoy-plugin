package compiler_test

import (
	"context"
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/client-go/dynamic"
)

var pol = &vpol.ValidatingPolicy{
	Spec: vpol.ValidatingPolicySpec{
		EvaluationConfiguration: &vpol.EvaluationConfiguration{
			Mode: v1alpha1.EvaluationModeEnvoy,
		},
		Variables: []admissionregistrationv1.Variable{
			{
				Name:       "force_authorized",
				Expression: `object.attributes.request.http.headers[?"x-force-authorized"].orValue("") in ["enabled", "true"]`,
			},
			{
				Name:       "force_unauthenticated",
				Expression: `object.attributes.request.http.headers[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]`,
			},
			{
				Name:       "metadata",
				Expression: `{"my-new-metadata": "my-new-value"}`,
			},
		},
		Validations: []admissionregistrationv1.Validation{
			{
				Expression: `variables.force_unauthenticated ? envoy.Denied(401).WithBody("Authentication Failed").Response() : null`,
			},
			{
				Expression: `!variables.force_authorized ? envoy.Denied(403).WithBody("Unauthorized Request").Response() : null`,
			},
			{
				Expression: `envoy.Allowed().WithHeader("x-validated-by", "my-security-checkpoint").WithoutHeader("x-force-authorized").WithResponseHeader("x-add-custom-response-header", "added").Response().WithMetadata(variables.metadata)`,
			},
		},
	},
}

func TestCompiler(t *testing.T) {
	compiler := compiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()

	compiled, errList := compiler.Compile(pol)
	assert.NoError(t, errList.ToAggregate())

	type testCase struct {
		request           *authv3.CheckRequest
		responseType      any
		expectedHeaders   map[string]string
		unexpectedHeaders []string
		body              string
	}

	tests := []testCase{
		{
			request: &authv3.CheckRequest{
				Attributes: &authv3.AttributeContext{
					Request: &authv3.AttributeContext_Request{
						Http: &authv3.AttributeContext_HttpRequest{
							Headers: map[string]string{
								"x-force-authorized": "true",
							},
						},
					},
				},
			},
			responseType: &authv3.CheckResponse_OkResponse{},
			expectedHeaders: map[string]string{
				"x-validated-by": "my-security-checkpoint",
			},
			unexpectedHeaders: []string{"x-force-authorized"},
		},
		{
			request: &authv3.CheckRequest{
				Attributes: &authv3.AttributeContext{
					Request: &authv3.AttributeContext_Request{
						Http: &authv3.AttributeContext_HttpRequest{
							Headers: map[string]string{
								"x-force-unauthenticated": "enabled",
							},
						},
					},
				},
			},
			responseType: &authv3.CheckResponse_DeniedResponse{},
			body:         "Authentication Failed",
		},
		{
			request: &authv3.CheckRequest{
				Attributes: &authv3.AttributeContext{
					Request: &authv3.AttributeContext_Request{
						Http: &authv3.AttributeContext_HttpRequest{
							Headers: make(map[string]string),
						},
					},
				},
			},
			responseType: &authv3.CheckResponse_DeniedResponse{},
			body:         "Unauthorized Request",
		},
	}

	for _, test := range tests {
		resp, err := compiled.Evaluate(context.TODO(), nil, test.request)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		ok := assert.IsType(t, test.responseType, resp.HttpResponse)
		if !ok {
			return
		}

		var headers []*corev3.HeaderValueOption
		var body string

		switch r := resp.HttpResponse.(type) {
		case *authv3.CheckResponse_OkResponse:
			headers = r.OkResponse.Headers
		case *authv3.CheckResponse_DeniedResponse:
			headers = r.DeniedResponse.Headers
			body = r.DeniedResponse.Body
		}

		for k, v := range test.expectedHeaders {
			for _, header := range headers {
				if header.Header.Key == k {
					assert.Equal(t, header.Header.Value, v)
					break
				}

				assert.Failf(t, "missing '%s' header", k)
			}
		}

		for _, h := range test.expectedHeaders {
			for _, header := range headers {
				if header.Header.Key == h {
					assert.Failf(t, "unexpected '%s' header", h)
				}
			}
		}

		if test.body != "" {
			assert.Equal(t, test.body, body)
		}
	}
}
