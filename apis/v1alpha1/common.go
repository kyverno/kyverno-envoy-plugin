package v1alpha1

// Group refers to a Kubernetes Group. It must either be an empty string or a
// RFC 1123 subdomain.
//
// This validation is based off of the corresponding Kubernetes validation:
// https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/util/validation/validation.go#L208
//
// Valid values include:
//
// * "" - empty string implies core Kubernetes API group
// * "authz.kyverno.io"
// * "policies.kyverno.io"
//
// Invalid values include:
//
// * "example.com/bar" - "/" is an invalid character
//
// +kubebuilder:validation:MaxLength=253
// +kubebuilder:validation:Pattern=`^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
type Group = string

// Kind refers to a Kubernetes Kind.
//
// Valid values include:
//
// * "Service"
// * "HTTPRoute"
//
// Invalid values include:
//
// * "invalid/kind" - "/" is an invalid character
//
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
type Kind = string

// ObjectName refers to the name of a Kubernetes object.
// Object names can have a variety of forms, including RFC 1123 subdomains,
// RFC 1123 labels, or RFC 1035 labels.
//
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=253
type ObjectName = string

// Namespace refers to a Kubernetes namespace. It must be a RFC 1123 label.
//
// This validation is based off of the corresponding Kubernetes validation:
// https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/util/validation/validation.go#L187
//
// This is used for Namespace name validation here:
// https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/api/validation/generic.go#L63
//
// Valid values include:
//
// * "example"
//
// Invalid values include:
//
// * "example.com" - "." is an invalid character
//
// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
type Namespace = string
