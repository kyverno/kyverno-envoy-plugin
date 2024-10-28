package v1alpha1

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

type AuthorizationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizationPolicySpec `json:"spec"`
}

type AuthorizationPolicySpec struct {
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
	Variables []admissionregistrationv1.Variable `json:"variables,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=variables"`

	// +listType=atomic
	// +optional
	Rules []AuthorizationRule `json:"rules,omitempty"`
}

type AuthorizationRule struct {
	// +kubebuilder:default:={}
	// +optional
	When When `json:"when"`

	// +required
	Return string `json:"return"`
}

type When struct {
	// +optional
	Expression string `json:"expression,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthorizationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizationPolicy `json:"items"`
}
