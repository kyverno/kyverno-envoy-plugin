package v1alpha1

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// AuthorizationPolicy defines an authorization policy resource
type AuthorizationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizationPolicySpec `json:"spec"`
}

// AuthorizationPolicySpec defines the spec of an authorization policy
type AuthorizationPolicySpec struct {
	// FailurePolicy defines how to handle failures for the policy. Failures can
	// occur from CEL expression parse errors, type check errors, runtime errors and invalid
	// or mis-configured policy definitions.
	//
	// FailurePolicy does not define how validations that evaluate to false are handled.
	//
	// Allowed values are Ignore or Fail. Defaults to Fail.
	// +optional
	FailurePolicy *admissionregistrationv1.FailurePolicyType `json:"failurePolicy,omitempty"`

	// Variables contain definitions of variables that can be used in composition of other expressions.
	// Each variable is defined as a named CEL expression.
	// The variables defined here will be available under `variables` in other expressions of the policy
	// except MatchConditions because MatchConditions are evaluated before the rest of the policy.
	//
	// The expression of a variable can refer to other variables defined earlier in the list but not those after.
	// Thus, Variables must be sorted by the order of first appearance and acyclic.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	Variables []admissionregistrationv1.Variable `json:"variables,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Authorizations contain CEL expressions which is used to apply the authorization.
	// +listType=atomic
	// +optional
	Authorizations []Authorization `json:"authorizations,omitempty"`
}

func (s *AuthorizationPolicySpec) GetFailurePolicy() admissionregistrationv1.FailurePolicyType {
	if s.FailurePolicy == nil {
		return admissionregistrationv1.Fail
	}
	return *s.FailurePolicy
}

// Authorization defines an authorization policy rule
type Authorization struct {
	// Expression represents the expression which will be evaluated by CEL.
	// ref: https://github.com/google/cel-spec
	// CEL expressions have access to CEL variables as well as some other useful variables:
	//
	// - 'object' - The object from the incoming request. (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest)
	//
	// CEL expressions are expected to return an envoy CheckResponse (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse).
	// +required
	Expression string `json:"expression"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizationPolicyList defines a list of authorization policies
type AuthorizationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizationPolicy `json:"items"`
}
