package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:namespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,categories=kyverno

// AuthorizationServer is a resource that represents a new kyverno authorization server.
type AuthorizationServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizationServerSpec `json:"spec,omitempty"`
}

// AuthorizationServerSpec defines the spec of a authorization server.
type AuthorizationServerSpec struct {
	// AuthorizationServerPolicySource contains all the sources of policies for the authorization server.
	Sources []AuthorizationServerPolicySource `json:"sources,omitempty"`
}

// AuthorizationServerPolicySource represents where the authorization server will get its policies from.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:MaxProperties=1
type AuthorizationServerPolicySource struct {
	KubernetesPolicySource `json:"kubernetes,omitempty"`
	ExternalPolicySource   `json:"external,omitempty"`
}

// PolicyObjectReference represents a reference to a policy resource.
// +kubebuilder:validation:XValidation:rule="has(self.name) || has(self.selector)",message="either name or selector must be specified"
// +kubebuilder:validation:XValidation:rule="!(has(self.name) && has(self.selector))",message="name and selector are mutually exclusive"
type PolicyObjectReference struct {
	// Group is the group of the referent. For example, "policies.kyverno.io".
	// When unspecified or empty string, core API group is inferred.

	// +optional
	// +kubebuilder:default=policies.kyverno.io
	Group *Group `json:"group,omitempty"`

	// Kind is the kind of the referent. For example, "ValidatingPolicy".
	// +optional
	// +kubebuilder:default=ValidatingPolicy
	Kind *Kind `json:"kind"`

	// Name is the name of the referent.
	// Mutually exclusive with Selector.
	// +optional
	Name *ObjectName `json:"name,omitempty"`

	// Note: Namespace is omitted because policies are global, not namespaced.

	// Selector is a label selector to select the Kubernetes policy resource.
	// Mutually exclusive with Name.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// KubernetesSource defines a Kubernetes-based policy source.
// +kubebuilder:validation:Enum=Ignore;Fail
// KubernetesPolicySource defines a reference to a Kubernetes policy resource.
type KubernetesPolicySource struct {
	// PolicyRef is a reference to Kubernetes policy resources.
	// When omitted, all ValidatingPolicy resources in the cluster are selected.
	// When present, filters policies by name or selector.

	// +optional
	PolicyRef PolicyObjectReference `json:"policyRef"`
}

// ExternalSource defines an external policy source.
type ExternalPolicySource struct {
	// URL is the URL of the external policy source
	// Supported schemes are: file://, oci://, https://, etc
	URL string `json:"url"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AuthorizationServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthorizationServer `json:"items,omitempty"`
}
