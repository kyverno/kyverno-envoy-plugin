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
	// Type defines the type of authorization server.
	// +required
	Type AuthorizationServerType `json:"type"`
	// AuthorizationServerPolicySource contains all the sources of policies for the authorization server.
	Sources []AuthorizationServerPolicySource `json:"sources,omitempty"`
}

// AuthorizationServerType defines the type of authorization server.
// Only one of the fields should be set at a time (mutually exclusive).
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:MaxProperties=1
type AuthorizationServerType struct {
	// Envoy configures an Envoy-based authorization server.
	Envoy *EnvoyAuthorizationServer `json:"envoy,omitempty"`
	// HTTP configures a custom HTTP authorization server.
	HTTP *HTTPAuthorizationServer `json:"http,omitempty"`
}

// EnvoyAuthorizationServer defines the Envoy authorization server configuration.
type EnvoyAuthorizationServer struct {
	// Network is the network the server listens on.
	// +kubebuilder:default=tcp
	// +optional
	Network string `json:"network"`
	// Address is the network address the server listens on.
	Address string `json:"address"`
}

// HTTPAuthorizationServer defines the HTTP authorization server configuration.
type HTTPAuthorizationServer struct {
	// Port is the port the server listens on.
	Port int `json:"port,omitempty"`
	// Modifiers to apply to requests and responses.
	Modifiers *Modifiers `json:"modifiers,omitempty"`
}

// Modifiers defines the request/response modifiers for the authorization server.
type Modifiers struct {
	// Request is a script or expression for modifying the incoming request.
	Request string `json:"request,omitempty"`
	// Response is a script or expression for modifying the outgoing response.
	Response string `json:"response,omitempty"`
}

// AuthorizationServerPolicySource represents where the authorization server will get its policies from.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:MaxProperties=1
type AuthorizationServerPolicySource struct {
	Kubernetes *KubernetesPolicySource `json:"kubernetes,omitempty"`
	External   *ExternalPolicySource   `json:"external,omitempty"`
	Oci        *OciPolicySource        `json:"oci,omitempty"`
	Fs         *FsPolicySource         `json:"fs,omitempty"`
}

// PolicyObjectReference represents a reference to a policy resource.
// +kubebuilder:validation:XValidation:rule="has(self.name) || has(self.selector)",message="either name or selector must be specified"
// +kubebuilder:validation:XValidation:rule="!(has(self.name) && has(self.selector))",message="name and selector are mutually exclusive"
type PolicyObjectReference struct {
	// Group is the group of the referent. For example, "policies.kyverno.io".
	// When unspecified or empty string, core API group is inferred.
	// +kubebuilder:default=policies.kyverno.io
	Group Group `json:"group,omitempty"`

	// Kind is the kind of the referent. For example, "ValidatingPolicy".
	// +kubebuilder:default=ValidatingPolicy
	Kind Kind `json:"kind,omitempty"`

	// Name is the name of the referent.
	// Mutually exclusive with Selector.
	Name ObjectName `json:"name,omitempty"`

	// Selector is a label selector to select the Kubernetes policy resource.
	// Mutually exclusive with Name.
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// KubernetesPolicySource defines a reference to a Kubernetes policy resource.
type KubernetesPolicySource struct {
	// PolicyRef is a reference to Kubernetes policy resources.
	// When omitted, all ValidatingPolicy resources in the cluster are selected.
	// When present, filters policies by name or selector.
	PolicyRef *PolicyObjectReference `json:"policyRef,omitempty"`
}

// ExternalSource defines an external policy source.
type ExternalPolicySource struct {
	// URL is the URL of the external policy source
	// Supported schemes are: file://, oci://, https://, etc
	// +required
	URL string `json:"url"`
}

// OciPolicySource defines the configuration for fetching policies
// from an OCI (Open Container Initiative) registry.
type OciPolicySource struct {
	// URL specifies the location of the OCI registry or image
	// that contains the policy definitions.
	URL string `json:"url"`

	// AllowInsecureRegistry indicates whether connections to an
	// insecure (HTTP or self-signed HTTPS) registry are permitted.
	// This should generally be false in production environments
	// to ensure secure communication.
	AllowInsecureRegistry bool `json:"allowInsecureRegistry,omitempty"`

	// ImagePullSecrets lists the names of Kubernetes secrets that
	// contain credentials needed to authenticate with the OCI registry.
	// These are typically referenced in Kubernetes to pull images
	// from private registries.
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
}

// FsPolicySource defines the configuration for loading a policy
// from a local or mounted filesystem path.
type FsPolicySource struct {
	// Path specifies the filesystem location where the policy
	// files are stored.
	Path string `json:"path"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AuthorizationServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizationServer `json:"items,omitempty"`
}
