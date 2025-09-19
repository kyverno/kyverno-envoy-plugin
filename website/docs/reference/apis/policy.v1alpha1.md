---
title: policy (v1alpha1)
content_type: tool-reference
package: envoy.kyverno.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [AuthorizationPolicy](#envoy-kyverno-io-v1alpha1-AuthorizationPolicy)
  
## AuthorizationPolicy     {#envoy-kyverno-io-v1alpha1-AuthorizationPolicy}

<p>AuthorizationPolicy defines an authorization policy resource</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `envoy.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `AuthorizationPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`AuthorizationPolicySpec`](#envoy-kyverno-io-v1alpha1-AuthorizationPolicySpec) | :white_check_mark: |  | *No description provided.* |

## Authorization     {#envoy-kyverno-io-v1alpha1-Authorization}

**Appears in:**
    
- [AuthorizationPolicySpec](#envoy-kyverno-io-v1alpha1-AuthorizationPolicySpec)

<p>Authorization defines an authorization policy rule</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `match` | `string` |  |  | <p>Match represents the match condition which will be evaluated by CEL. Must evaluate to bool.</p> |
| `response` | `string` | :white_check_mark: |  | <p>Response represents the response expression which will be evaluated by CEL. ref: https://github.com/google/cel-spec CEL expressions have access to CEL variables as well as some other useful variables: - 'object' - The object from the incoming request. (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) CEL expressions are expected to return an envoy CheckResponse (https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse).</p> |

## AuthorizationPolicySpec     {#envoy-kyverno-io-v1alpha1-AuthorizationPolicySpec}

**Appears in:**
    
- [AuthorizationPolicy](#envoy-kyverno-io-v1alpha1-AuthorizationPolicy)

<p>AuthorizationPolicySpec defines the spec of an authorization policy</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `failurePolicy` | [`admissionregistration/v1.FailurePolicyType`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#failurepolicytype-v1-admissionregistration) |  |  | <p>FailurePolicy defines how to handle failures for the policy. Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions. FailurePolicy does not define how validations that evaluate to false are handled. Allowed values are Ignore or Fail. Defaults to Fail.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. An empty list of matchConditions matches all requests. The exact matching logic is (in order):   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.   3. If any matchCondition evaluates to an error (but none are FALSE):      - If failurePolicy=Fail, reject the request      - If failurePolicy=Ignore, the policy is skipped</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `deny` | [`[]Authorization`](#envoy-kyverno-io-v1alpha1-Authorization) |  |  | <p>Deny contain CEL expressions which is used to deny a request.</p> |
| `allow` | [`[]Authorization`](#envoy-kyverno-io-v1alpha1-Authorization) |  |  | <p>Allow contain CEL expressions which is used to allow a request.</p> |

  