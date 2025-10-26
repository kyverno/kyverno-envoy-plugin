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
| `external` | [`ExternalPolicySource`](#authz-kyverno-io-v1alpha1-ExternalPolicySource) | :white_check_mark: |  | *No description provided.* |

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
| `envoy` | [`EnvoyAuthorizationServer`](#authz-kyverno-io-v1alpha1-EnvoyAuthorizationServer) |  |  | <p>Envoy configures an Envoy-based authorization server.</p> |
| `http` | [`HTTPAuthorizationServer`](#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer) |  |  | <p>HTTP configures a custom HTTP authorization server.</p> |

## EnvoyAuthorizationServer     {#authz-kyverno-io-v1alpha1-EnvoyAuthorizationServer}

**Appears in:**
    
- [AuthorizationServerType](#authz-kyverno-io-v1alpha1-AuthorizationServerType)

<p>EnvoyAuthorizationServer defines the Envoy authorization server configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `port` | `int` | :white_check_mark: |  | <p>Port is the port the server listens on.</p> |
| `modifiers` | [`Modifiers`](#authz-kyverno-io-v1alpha1-Modifiers) |  |  | <p>Modifiers to apply to requests and responses.</p> |

## ExternalPolicySource     {#authz-kyverno-io-v1alpha1-ExternalPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>ExternalSource defines an external policy source.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `url` | `string` | :white_check_mark: |  | <p>URL is the URL of the external policy source Supported schemes are: file://, oci://, https://, etc</p> |

## HTTPAuthorizationServer     {#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer}

**Appears in:**
    
- [AuthorizationServerType](#authz-kyverno-io-v1alpha1-AuthorizationServerType)

<p>HTTPAuthorizationServer defines the HTTP authorization server configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `port` | `int` | :white_check_mark: |  | <p>Port is the port the server listens on.</p> |
| `modifiers` | [`Modifiers`](#authz-kyverno-io-v1alpha1-Modifiers) |  |  | <p>Modifiers to apply to requests and responses.</p> |

## KubernetesPolicySource     {#authz-kyverno-io-v1alpha1-KubernetesPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>KubernetesPolicySource defines a reference to a Kubernetes policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `policyRef` | [`PolicyObjectReference`](#authz-kyverno-io-v1alpha1-PolicyObjectReference) |  |  | *No description provided.* |

## Modifiers     {#authz-kyverno-io-v1alpha1-Modifiers}

**Appears in:**
    
- [EnvoyAuthorizationServer](#authz-kyverno-io-v1alpha1-EnvoyAuthorizationServer)
- [HTTPAuthorizationServer](#authz-kyverno-io-v1alpha1-HTTPAuthorizationServer)

<p>Modifiers defines the request/response modifiers for the authorization server.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `request` | `string` |  |  | <p>Request is a script or expression for modifying the incoming request.</p> |
| `response` | `string` |  |  | <p>Response is a script or expression for modifying the outgoing response.</p> |

## PolicyObjectReference     {#authz-kyverno-io-v1alpha1-PolicyObjectReference}

**Appears in:**
    
- [KubernetesPolicySource](#authz-kyverno-io-v1alpha1-KubernetesPolicySource)

<p>PolicyObjectReference represents a reference to a policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `group` | `string` |  |  | *No description provided.* |
| `kind` | `string` |  |  | <p>Kind is the kind of the referent. For example, "ValidatingPolicy".</p> |
| `name` | `string` |  |  | <p>Name is the name of the referent. Mutually exclusive with Selector.</p> |
| `selector` | [`meta/v1.LabelSelector`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#labelselector-v1-meta) |  |  | <p>Selector is a label selector to select the Kubernetes policy resource. Mutually exclusive with Name.</p> |

  