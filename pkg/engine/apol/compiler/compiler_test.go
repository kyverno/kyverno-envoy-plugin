package compiler_test

import (
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/compiler"
	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

var pol = &v1alpha1.AuthorizationPolicy{
	Spec: v1alpha1.AuthorizationPolicySpec{
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
		Deny: []v1alpha1.Authorization{
			{
				Match:    "variables.force_unauthenticated",
				Response: `envoy.Denied(401).WithBody("Authentication Failed").Response()`,
			},
			{
				Match:    "!variables.force_authorized",
				Response: `envoy.Denied(403).WithBody("Unauthorized Request").Response()`,
			},
		},
		Allow: []v1alpha1.Authorization{
			{
				Response: `envoy.Allowed().WithHeader("x-validated-by", "my-security-checkpoint").WithoutHeader("x-force-authorized").WithResponseHeader("x-add-custom-response-header", "added").Response().WithMetadata(variables.metadata)`,
			},
		},
	},
}

func TestCompiler(t *testing.T) {
	compiler := compiler.NewCompiler()

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
		allow, deny := compiled.For(test.request)
		resp, err := deny()
		assert.NoError(t, err)
		if resp == nil {
			resp, err = allow()
			assert.NoError(t, err)
		}
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
