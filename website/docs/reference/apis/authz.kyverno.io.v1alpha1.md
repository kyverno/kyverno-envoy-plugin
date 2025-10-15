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
| `sources` | [`[]AuthorizationServerPolicySource`](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource) | :white_check_mark: |  | <p>AuthorizationServerPolicySource contains all the sources of policies for the authorization server.</p> |

## ExternalPolicySource     {#authz-kyverno-io-v1alpha1-ExternalPolicySource}

**Appears in:**
    
- [AuthorizationServerPolicySource](#authz-kyverno-io-v1alpha1-AuthorizationServerPolicySource)

<p>ExternalSource defines an external policy source.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `url` | `string` | :white_check_mark: |  | <p>URL is the URL of the external policy source Supported schemes are: file://, oci://, https://, etc</p> |

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

<p>KubernetesSource defines a Kubernetes-based policy source.
KubernetesPolicySource defines a reference to a Kubernetes policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `policyRef` | [`PolicyObjectReference`](#authz-kyverno-io-v1alpha1-PolicyObjectReference) |  |  | *No description provided.* |

## ObjectName     {#authz-kyverno-io-v1alpha1-ObjectName}

(Alias of `string`)

**Appears in:**
    
- [PolicyObjectReference](#authz-kyverno-io-v1alpha1-PolicyObjectReference)

<p>ObjectName refers to the name of a Kubernetes object.
Object names can have a variety of forms, including RFC 1123 subdomains,
RFC 1123 labels, or RFC 1035 labels.</p>


## PolicyObjectReference     {#authz-kyverno-io-v1alpha1-PolicyObjectReference}

**Appears in:**
    
- [KubernetesPolicySource](#authz-kyverno-io-v1alpha1-KubernetesPolicySource)

<p>PolicyObjectReference represents a reference to a policy resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `group` | [`Group`](#authz-kyverno-io-v1alpha1-Group) |  |  | *No description provided.* |
| `kind` | [`Kind`](#authz-kyverno-io-v1alpha1-Kind) |  |  | <p>Kind is the kind of the referent. For example, "ValidatingPolicy".</p> |
| `name` | [`ObjectName`](#authz-kyverno-io-v1alpha1-ObjectName) |  |  | <p>Name is the name of the referent. Mutually exclusive with Selector.</p> |
| `selector` | [`meta/v1.LabelSelector`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#labelselector-v1-meta) |  |  | <p>Selector is a label selector to select the Kubernetes policy resource. Mutually exclusive with Name.</p> |

  