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
	metav1.ObjectMeta `json:"metadata"`
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

	// MatchConditions is a list of conditions that must be met for a request to be validated.
	// An empty list of matchConditions matches all requests.
	//
	// The exact matching logic is (in order):
	//   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.
	//   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.
	//   3. If any matchCondition evaluates to an error (but none are FALSE):
	//      - If failurePolicy=Fail, reject the request
	//      - If failurePolicy=Ignore, the policy is skipped
	//
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	MatchConditions []admissionregistrationv1.MatchCondition `json:"matchConditions,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

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

	// Deny contain CEL expressions which is used to deny a request.
	// +listType=atomic
	// +optional
	Deny []Authorization `json:"deny,omitempty"`

	// Allow contain CEL expressions which is used to allow a request.
	// +listType=atomic
	// +optional
	Allow []Authorization `json:"allow,omitempty"`
}

func (s *AuthorizationPolicySpec) GetFailurePolicy() admissionregistrationv1.FailurePolicyType {
	if s.FailurePolicy == nil {
		return admissionregistrationv1.Fail
	}
	return *s.FailurePolicy
}

// Authorization defines an authorization policy rule
type Authorization struct {
	// Match represents the match condition which will be evaluated by CEL. Must evaluate to bool.
	// +optional
	Match string `json:"match,omitempty"`
	// Response represents the response expression which will be evaluated by CEL.
	// ref: https://github.com/google/cel-spec
	// CEL expressions have access to CEL variables as well as some other useful variables:
	//
	// - 'object' - The object from the incoming request. (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest)
	//
	// CEL expressions are expected to return an envoy CheckResponse (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse).
	// +required
	Response string `json:"response"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizationPolicyList defines a list of authorization policies
type AuthorizationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizationPolicy `json:"items"`
}
