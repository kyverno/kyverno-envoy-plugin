---
title: authz.kyverno.io (v1alpha1)
content_type: tool-reference
package: authz.kyverno.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [AuthorizationServer](#authz-kyverno-io-v1alpha1-AuthorizationServer)
  
## AuthorizationServer     {#authz-kyverno-io-v1alpha1-AuthorizationServer}

<p>AuthorizationServer is a resource that represents a new kyverno authorization server.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `authz.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `AuthorizationServer` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`AuthorizationServerSpec`](#authz-kyverno-io-v1alpha1-AuthorizationServerSpec) | :white_check_mark: |  | *No description provided.* |

## AuthorizationServerPolicySource     {#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource}

**Appears in:**
    
- [AuthorizationServerSpec](#authz-kyverno-io-v1alpha1-AuthorizationServerSpec)

<p>AuthorizationServerPolicySource represents where the authorization server will get its policies from.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `kubernetes` | [`KubernetesPolicySource`](#authz-kyverno-io-v1alpha1-KubernetesPolicySource) | :white_check_mark: |  | *No description provided.* |
| `fs` | [`FsPolicySource`](#authz-kyverno-io-v1alpha1-FsPolicySource) | :white_check_mark: |  | *No description provided.* |
| `git` | [`GitPolicySource`](#authz-kyverno-io-v1alpha1-GitPolicySource) | :white_check_mark: |  | *No description provided.* |
| `oci` | [`OciPolicySource`](#authz-kyverno-io-v1alpha1-OciPolicySource) | :white_check_mark: |  | *No description provided.* |

## AuthorizationServerSpec     {#authz-kyverno-io-v1alpha1-AuthorizationServerSpec}

**Appears in:**
    
- [AuthorizationServer](#authz-kyverno-io-v1alpha1-AuthorizationServer)

<p>AuthorizationServerSpec defines the spec of a authorization server.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `type` | [`AuthorizationServerType`](#authz-kyverno-io-v1alpha1-AuthorizationServerType) | :white_check_mark: |  | <p>Type defines the type of authorization server.</p> |
| `sources` | [`[]AuthorizationServerPolicySource`](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource) | :white_check_mark: |  | <p>AuthorizationServerPolicySource contains all the sources of policies for the authorization server.</p> |

## AuthorizationServerType     {#authz-kyverno-io-v1alpha1-AuthorizationServerType}

**Appears in:**
    
- [AuthorizationServerSpec](#authz-kyverno-io-v1alpha1-AuthorizationServerSpec)

<p>AuthorizationServerType defines the type of authorization server.
Only one of the fields should be set at a time (mutually exclusive).</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `envoy` | [`EnvoyAuthorizationServer`](#authz-kyverno-io-v1alpha1-EnvoyAuthorizationServer) | :white_check_mark: |  | <p>Envoy configures an Envoy-based authorization server.</p> |
| `http` | [`HTTPAuthorizationServer`](#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer) | :white_check_mark: |  | <p>HTTP configures a custom HTTP authorization server.</p> |

## EnvoyAuthorizationServer     {#authz-kyverno-io-v1alpha1-EnvoyAuthorizationServer}

**Appears in:**
    
- [AuthorizationServerType](#authz-kyverno-io-v1alpha1-AuthorizationServerType)

<p>EnvoyAuthorizationServer defines the Envoy authorization server configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `network` | `string` |  |  | <p>Network is the network the server listens on.</p> |
| `address` | `string` | :white_check_mark: |  | <p>Address is the network address the server listens on.</p> |

## FsPolicySource     {#authz-kyverno-io-v1alpha1-FsPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>FsPolicySource defines the configuration for loading a policy
from a local or mounted filesystem path.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `path` | `string` | :white_check_mark: |  | <p>Path specifies the filesystem location where the policy files are stored.</p> |

## GitPolicySource     {#authz-kyverno-io-v1alpha1-GitPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>GitPolicySource defines the configuration for retrieving a policy
from a Git repository.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `url` | `string` | :white_check_mark: |  | <p>URL specifies the Git repository location that contains the policy files or definitions. Supported formats typically include HTTPS or SSH Git URLs.</p> |

## Group     {#authz-kyverno-io-v1alpha1-Group}

(Alias of `string`)

**Appears in:**
    
- [PolicyObjectReference](#authz-kyverno-io-v1alpha1-PolicyObjectReference)

<p>Group refers to a Kubernetes Group. It must either be an empty string or a
RFC 1123 subdomain.</p>
<p>This validation is based off of the corresponding Kubernetes validation:
https://github.com/kubernetes/apimachinery/blob/02cfb53916346d085a6c6c7c66f882e3c6b0eca6/pkg/util/validation/validation.go#L208</p>
<p>Valid values include:</p>
<ul>
<li>&quot;&quot; - empty string implies core Kubernetes API group</li>
<li>&quot;authz.kyverno.io&quot;</li>
<li>&quot;policies.kyverno.io&quot;</li>
</ul>
<p>Invalid values include:</p>
<ul>
<li>&quot;example.com/bar&quot; - &quot;/&quot; is an invalid character</li>
</ul>


## HTTPAuthorizationServer     {#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer}

**Appears in:**
    
- [AuthorizationServerType](#authz-kyverno-io-v1alpha1-AuthorizationServerType)

<p>HTTPAuthorizationServer defines the HTTP authorization server configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `address` | `string` | :white_check_mark: |  | <p>Address is the network address the server listens on.</p> |
| `modifiers` | [`Modifiers`](#authz-kyverno-io-v1alpha1-Modifiers) | :white_check_mark: |  | <p>Modifiers to apply to requests and responses.</p> |

## Kind     {#authz-kyverno-io-v1alpha1-Kind}

(Alias of `string`)

**Appears in:**
    
- [PolicyObjectReference](#authz-kyverno-io-v1alpha1-PolicyObjectReference)

<p>Kind refers to a Kubernetes Kind.</p>
<p>Valid values include:</p>
<ul>
<li>&quot;Service&quot;</li>
<li>&quot;HTTPRoute&quot;</li>
</ul>
<p>Invalid values include:</p>
<ul>
<li>&quot;invalid/kind&quot; - &quot;/&quot; is an invalid character</li>
</ul>


## KubernetesPolicySource     {#authz-kyverno-io-v1alpha1-KubernetesPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>KubernetesPolicySource defines a reference to a Kubernetes policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `policyRef` | [`PolicyObjectReference`](#authz-kyverno-io-v1alpha1-PolicyObjectReference) | :white_check_mark: |  | <p>PolicyRef is a reference to Kubernetes policy resources. When omitted, all ValidatingPolicy resources in the cluster are selected. When present, filters policies by name or selector.</p> |

## Modifiers     {#authz-kyverno-io-v1alpha1-Modifiers}

**Appears in:**
    
- [HTTPAuthorizationServer](#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer)

<p>Modifiers defines the request/response modifiers for the authorization server.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `request` | `string` | :white_check_mark: |  | <p>Request is a script or expression for modifying the incoming request.</p> |
| `response` | `string` | :white_check_mark: |  | <p>Response is a script or expression for modifying the outgoing response.</p> |

## ObjectName     {#authz-kyverno-io-v1alpha1-ObjectName}

(Alias of `string`)

**Appears in:**
    
- [PolicyObjectReference](#authz-kyverno-io-v1alpha1-PolicyObjectReference)

<p>ObjectName refers to the name of a Kubernetes object.
Object names can have a variety of forms, including RFC 1123 subdomains,
RFC 1123 labels, or RFC 1035 labels.</p>


## OciPolicySource     {#authz-kyverno-io-v1alpha1-OciPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>OciPolicySource defines the configuration for fetching policies
from an OCI (Open Container Initiative) registry.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `url` | `string` | :white_check_mark: |  | <p>URL specifies the location of the OCI registry or image that contains the policy definitions.</p> |
| `allowInsecureRegistry` | `bool` | :white_check_mark: |  | <p>AllowInsecureRegistry indicates whether connections to an insecure (HTTP or self-signed HTTPS) registry are permitted. This should generally be false in production environments to ensure secure communication.</p> |
| `imagePullSecrets` | `[]string` | :white_check_mark: |  | <p>ImagePullSecrets lists the names of Kubernetes secrets that contain credentials needed to authenticate with the OCI registry. These are typically referenced in Kubernetes to pull images from private registries.</p> |

## PolicyObjectReference     {#authz-kyverno-io-v1alpha1-PolicyObjectReference}

**Appears in:**
    
- [KubernetesPolicySource](#authz-kyverno-io-v1alpha1-KubernetesPolicySource)

<p>PolicyObjectReference represents a reference to a policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `group` | [`Group`](#authz-kyverno-io-v1alpha1-Group) | :white_check_mark: |  | <p>Group is the group of the referent. For example, "policies.kyverno.io". When unspecified or empty string, core API group is inferred.</p> |
| `kind` | [`Kind`](#authz-kyverno-io-v1alpha1-Kind) | :white_check_mark: |  | <p>Kind is the kind of the referent. For example, "ValidatingPolicy".</p> |
| `name` | [`ObjectName`](#authz-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: |  | <p>Name is the name of the referent. Mutually exclusive with Selector.</p> |
| `selector` | [`meta/v1.LabelSelector`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#labelselector-v1-meta) | :white_check_mark: |  | <p>Selector is a label selector to select the Kubernetes policy resource. Mutually exclusive with Name.</p> |

  