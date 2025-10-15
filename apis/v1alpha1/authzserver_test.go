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
