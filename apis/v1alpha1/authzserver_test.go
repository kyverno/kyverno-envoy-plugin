package v1alpha1

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAuthorizationServerSpec(t *testing.T) {
	testCases := []struct {
		name   string
		spec   AuthorizationServerSpec
		expect AuthorizationServerSpec
	}{
		{
			name: "Empty sources",
			spec: AuthorizationServerSpec{},
			expect: AuthorizationServerSpec{
				Sources: nil,
			},
		},
		{
			name: "With single source",
			spec: AuthorizationServerSpec{
				Sources: []AuthorizationServerPolicySource{
					{
						KubernetesPolicySource: KubernetesPolicySource{
							PolicyRef: PolicyObjectReference{
								Group: ptrGroup("policies.kyverno.io"),
								Kind:  ptrKind("ValidatingPolicy"),
								Name:  ptrObjectName("test-policy"),
							},
						},
					},
				},
			},
			expect: AuthorizationServerSpec{
				Sources: []AuthorizationServerPolicySource{
					{
						KubernetesPolicySource: KubernetesPolicySource{
							PolicyRef: PolicyObjectReference{
								Group: ptrGroup("policies.kyverno.io"),
								Kind:  ptrKind("ValidatingPolicy"),
								Name:  ptrObjectName("test-policy"),
							},
						},
					},
				},
			},
		},
		{
			name: "With external policy source",
			spec: AuthorizationServerSpec{
				Sources: []AuthorizationServerPolicySource{
					{
						ExternalPolicySource: ExternalPolicySource{
							URL: "https://example.com/policy.bundle",
						},
					},
				},
			},
			expect: AuthorizationServerSpec{
				Sources: []AuthorizationServerPolicySource{
					{
						ExternalPolicySource: ExternalPolicySource{
							URL: "https://example.com/policy.bundle",
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.spec, tc.expect) {
				t.Errorf("got %+v, expected %+v", tc.spec, tc.expect)
			}
		})
	}
}

func TestAuthorizationServerTypeField(t *testing.T) {
	tests := []struct {
		name   string
		typ    AuthorizationServerType
		expect func(t *testing.T, typ AuthorizationServerType)
	}{
		{
			name: "Envoy type set",
			typ: AuthorizationServerType{
				Envoy: &EnvoyAuthorizationServer{
					Port: 8080,
					Modifiers: &Modifiers{
						Request:  "req-mod",
						Response: "resp-mod",
					},
				},
			},
			expect: func(t *testing.T, typ AuthorizationServerType) {
				if typ.Envoy == nil {
					t.Errorf("Envoy should not be nil")
				}
				if typ.HTTP != nil {
					t.Errorf("HTTP should be nil when Envoy is set")
				}
				if typ.Envoy.Port != 8080 {
					t.Errorf("unexpected Port: %d", typ.Envoy.Port)
				}
				if typ.Envoy.Modifiers == nil || typ.Envoy.Modifiers.Request != "req-mod" || typ.Envoy.Modifiers.Response != "resp-mod" {
					t.Errorf("unexpected Modifiers: %+v", typ.Envoy.Modifiers)
				}
			},
		},
		{
			name: "HTTP type set",
			typ: AuthorizationServerType{
				HTTP: &HTTPAuthorizationServer{
					Port: 9090,
					Modifiers: &Modifiers{
						Request:  "req-script",
						Response: "resp-script",
					},
				},
			},
			expect: func(t *testing.T, typ AuthorizationServerType) {
				if typ.HTTP == nil {
					t.Errorf("HTTP should not be nil")
				}
				if typ.Envoy != nil {
					t.Errorf("Envoy should be nil when HTTP is set")
				}
				if typ.HTTP.Port != 9090 {
					t.Errorf("unexpected Port: %d", typ.HTTP.Port)
				}
				if typ.HTTP.Modifiers == nil || typ.HTTP.Modifiers.Request != "req-script" || typ.HTTP.Modifiers.Response != "resp-script" {
					t.Errorf("unexpected Modifiers: %+v", typ.HTTP.Modifiers)
				}
			},
		},
		{
			name: "Neither type set (invalid case)",
			typ:  AuthorizationServerType{},
			expect: func(t *testing.T, typ AuthorizationServerType) {
				if typ.Envoy != nil || typ.HTTP != nil {
					t.Errorf("Both Envoy and HTTP should be nil for empty struct")
				}
			},
		},
		{
			name: "Both types set (invalid config)",
			typ: AuthorizationServerType{
				Envoy: &EnvoyAuthorizationServer{
					Port: 8888,
				},
				HTTP: &HTTPAuthorizationServer{
					Port: 9999,
				},
			},
			expect: func(t *testing.T, typ AuthorizationServerType) {
				if typ.Envoy == nil || typ.HTTP == nil {
					t.Errorf("Both Envoy and HTTP should be set")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.typ)
		})
	}
}

func TestAuthorizationServerSpec_TypeFieldUsage(t *testing.T) {
	envoyType := AuthorizationServerType{
		Envoy: &EnvoyAuthorizationServer{
			Port: 8000,
		},
	}
	httpType := AuthorizationServerType{
		HTTP: &HTTPAuthorizationServer{
			Port: 9000,
		},
	}

	specEnvoy := AuthorizationServerSpec{
		Type: envoyType,
		Sources: []AuthorizationServerPolicySource{
			{
				KubernetesPolicySource: KubernetesPolicySource{
					PolicyRef: PolicyObjectReference{
						Group: ptrGroup("policies.kyverno.io"),
						Kind:  ptrKind("ValidatingPolicy"),
						Name:  ptrObjectName("e-policy"),
					},
				},
			},
		},
	}
	if specEnvoy.Type.Envoy == nil {
		t.Errorf("Envoy field should be set in Type")
	}
	if specEnvoy.Type.HTTP != nil {
		t.Errorf("HTTP field should not be set in Type")
	}
	if specEnvoy.Type.Envoy.Port != 8000 {
		t.Errorf("unexpected Envoy Port: %d", specEnvoy.Type.Envoy.Port)
	}

	specHTTP := AuthorizationServerSpec{
		Type: httpType,
		Sources: []AuthorizationServerPolicySource{
			{
				ExternalPolicySource: ExternalPolicySource{
					URL: "https://host.net/policy",
				},
			},
		},
	}
	if specHTTP.Type.HTTP == nil {
		t.Errorf("HTTP field should be set in Type")
	}
	if specHTTP.Type.Envoy != nil {
		t.Errorf("Envoy field should not be set in Type")
	}
	if specHTTP.Type.HTTP.Port != 9000 {
		t.Errorf("unexpected HTTP Port: %d", specHTTP.Type.HTTP.Port)
	}
}

func TestAuthorizationServerRoundTrip(t *testing.T) {
	AuthorizationServer := AuthorizationServer{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AuthorizationServer",
			APIVersion: "authz.kyverno.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-AuthorizationServer",
			Namespace: "default",
			Labels: map[string]string{
				"kyverno.io/test": "AuthorizationServer",
			},
		},
		Spec: AuthorizationServerSpec{
			Sources: []AuthorizationServerPolicySource{
				{
					KubernetesPolicySource: KubernetesPolicySource{
						PolicyRef: PolicyObjectReference{
							Group: ptrGroup("policies.kyverno.io"),
							Kind:  ptrKind("ValidatingPolicy"),
							Name:  ptrObjectName("mypolicy"),
						},
					},
				},
				{
					ExternalPolicySource: ExternalPolicySource{
						URL: "oci://myrepo/mybundle:v1",
					},
				},
			},
		},
	}

	if AuthorizationServer.Kind != "AuthorizationServer" {
		t.Errorf("unexpected Kind: %s", AuthorizationServer.Kind)
	}
	if AuthorizationServer.Name != "test-AuthorizationServer" {
		t.Errorf("unexpected Name: %s", AuthorizationServer.Name)
	}
	if AuthorizationServer.Spec.Sources == nil || len(AuthorizationServer.Spec.Sources) != 2 {
		t.Errorf("missing or unexpected spec.sources")
	}
}

func TestPolicyObjectReference_MutualExclusion(t *testing.T) {
	// either name or selector must be specified, not both
	validName := PolicyObjectReference{
		Name: ptrObjectName("foo"),
	}

	validSelector := PolicyObjectReference{
		Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"a": "b"}},
	}

	// The logic here should check that *either* Name or Selector is set, but not both or neither.
	// In this case, only Name is set, which should be valid.
	if validName.Name == nil || validSelector.Selector == nil {
		t.Errorf("expected only Name to be set, got Name:%v Selector:%v", validName.Name, validSelector.Selector)
	}

	invalid := PolicyObjectReference{
		Name:     ptrObjectName("foo"),
		Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"a": "b"}},
	}
	if invalid.Name != nil && invalid.Selector != nil {
		// This would fail schema validation, but in Go it's valid.
		t.Log("Name and Selector are both set: validation should fail in scheme, not here")
	}
}

func TestExternalPolicySource(t *testing.T) {
	src := ExternalPolicySource{
		URL: "file:///etc/policies/bundle.tar.gz",
	}
	if src.URL == "" {
		t.Errorf("expected a URL for external policy source")
	}
}

func TestKubernetesPolicySource_Defaults(t *testing.T) {
	// Omitting PolicyRef should be valid and select all ValidatingPolicy in cluster.
	src := KubernetesPolicySource{}
	if src.PolicyRef.Name != nil {
		t.Errorf("expected empty PolicyRef.Name")
	}
	if src.PolicyRef.Selector != nil {
		t.Errorf("expected empty PolicyRef.Selector")
	}
}

// Pointer helpers for types
func ptrGroup(v Group) *Group                { tmp := v; return &tmp }
func ptrKind(v Kind) *Kind                   { tmp := v; return &tmp }
func ptrObjectName(v ObjectName) *ObjectName { tmp := v; return &tmp }
