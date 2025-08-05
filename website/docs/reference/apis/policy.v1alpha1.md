---
title: policy (v1alpha1)
content_type: tool-reference
package: envoy.kyverno.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [AuthorizationPolicy](#envoy-kyverno-io-v1alpha1-AuthorizationPolicy)
- [ValidatingPolicy](#envoy-kyverno-io-v1alpha1-ValidatingPolicy)
  
## AuthorizationPolicy     {#envoy-kyverno-io-v1alpha1-AuthorizationPolicy}

<p>AuthorizationPolicy defines an authorization policy resource</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `envoy.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `AuthorizationPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`AuthorizationPolicySpec`](#envoy-kyverno-io-v1alpha1-AuthorizationPolicySpec) | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicy     {#envoy-kyverno-io-v1alpha1-ValidatingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `envoy.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidatingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ValidatingPolicySpec`](#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`ValidatingPolicyStatus`](#envoy-kyverno-io-v1alpha1-ValidatingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## AdmissionConfiguration     {#envoy-kyverno-io-v1alpha1-AdmissionConfiguration}

**Appears in:**
    
- [EvaluationConfiguration](#envoy-kyverno-io-v1alpha1-EvaluationConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls if rules are applied during admission. Optional. Default value is "true".</p> |

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

## BackgroundConfiguration     {#envoy-kyverno-io-v1alpha1-BackgroundConfiguration}

**Appears in:**
    
- [EvaluationConfiguration](#envoy-kyverno-io-v1alpha1-EvaluationConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls if rules are applied to existing resources during a background scan. Optional. Default value is "true". The value must be set to "false" if the policy rule uses variables that are only available in the admission review request (e.g. user name).</p> |

## ConditionStatus     {#envoy-kyverno-io-v1alpha1-ConditionStatus}

**Appears in:**
    
- [ValidatingPolicyStatus](#envoy-kyverno-io-v1alpha1-ValidatingPolicyStatus)

<p>ConditionStatus is the shared status across all policy types</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ready` | `bool` |  |  | <p>The ready of a policy is a high-level summary of where the policy is in its lifecycle. The conditions array, the reason and message fields contain more detail about the policy's status.</p> |
| `conditions` | [`[]meta/v1.Condition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#condition-v1-meta) |  |  | *No description provided.* |
| `message` | `string` |  |  | <p>Message is a human readable message indicating details about the generation of ValidatingAdmissionPolicy/MutatingAdmissionPolicy It is an empty string when ValidatingAdmissionPolicy/MutatingAdmissionPolicy is successfully generated.</p> |

## EvaluationConfiguration     {#envoy-kyverno-io-v1alpha1-EvaluationConfiguration}

**Appears in:**
    
- [ValidatingPolicySpec](#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `mode` | [`EvaluationMode`](#envoy-kyverno-io-v1alpha1-EvaluationMode) |  |  | <p>Mode is the mode of policy evaluation. Allowed values are "Kubernetes" or "JSON". Optional. Default value is "Kubernetes".</p> |
| `admission` | [`AdmissionConfiguration`](#envoy-kyverno-io-v1alpha1-AdmissionConfiguration) |  |  | <p>Admission controls policy evaluation during admission.</p> |
| `background` | [`BackgroundConfiguration`](#envoy-kyverno-io-v1alpha1-BackgroundConfiguration) |  |  | <p>Background  controls policy evaluation during background scan.</p> |

## EvaluationMode     {#envoy-kyverno-io-v1alpha1-EvaluationMode}

(Alias of `string`)

**Appears in:**
    
- [EvaluationConfiguration](#envoy-kyverno-io-v1alpha1-EvaluationConfiguration)

## PodControllersGenerationConfiguration     {#envoy-kyverno-io-v1alpha1-PodControllersGenerationConfiguration}

**Appears in:**
    
- [ValidatingPolicyAutogenConfiguration](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `controllers` | `[]string` | :white_check_mark: |  | *No description provided.* |

## Target     {#envoy-kyverno-io-v1alpha1-Target}

**Appears in:**
    
- [ValidatingPolicyAutogen](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogen)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `group` | `string` | :white_check_mark: |  | *No description provided.* |
| `version` | `string` | :white_check_mark: |  | *No description provided.* |
| `resource` | `string` | :white_check_mark: |  | *No description provided.* |
| `kind` | `string` | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicyAutogen     {#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogen}

**Appears in:**
    
- [ValidatingPolicyAutogenStatus](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `targets` | [`[]Target`](#envoy-kyverno-io-v1alpha1-Target) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ValidatingPolicySpec`](#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicyAutogenConfiguration     {#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration}

**Appears in:**
    
- [ValidatingPolicySpec](#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `podControllers` | [`PodControllersGenerationConfiguration`](#envoy-kyverno-io-v1alpha1-PodControllersGenerationConfiguration) | :white_check_mark: |  | <p>PodControllers specifies whether to generate a pod controllers rules.</p> |
| `validatingAdmissionPolicy` | [`VapGenerationConfiguration`](#envoy-kyverno-io-v1alpha1-VapGenerationConfiguration) | :white_check_mark: |  | <p>ValidatingAdmissionPolicy specifies whether to generate a Kubernetes ValidatingAdmissionPolicy.</p> |

## ValidatingPolicyAutogenStatus     {#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus}

**Appears in:**
    
- [ValidatingPolicyStatus](#envoy-kyverno-io-v1alpha1-ValidatingPolicyStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `configs` | [`map[string]ValidatingPolicyAutogen`](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogen) | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicySpec     {#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec}

**Appears in:**
    
- [ValidatingPolicy](#envoy-kyverno-io-v1alpha1-ValidatingPolicy)
- [ValidatingPolicyAutogen](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogen)

<p>ValidatingPolicySpec is the specification of the desired behavior of the ValidatingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1-admissionregistration) | :white_check_mark: |  | <p>MatchConstraints specifies what resources this policy is designed to validate. The AdmissionPolicy cares about a request if it matches _all_ Constraints. Required.</p> |
| `validations` | [`[]admissionregistration/v1.Validation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validation-v1-admissionregistration) |  |  | <p>Validations contain CEL expressions which is used to apply the validation. Validations and AuditAnnotations may not both be empty; a minimum of one Validations or AuditAnnotations is required.</p> |
| `failurePolicy` | [`admissionregistration/v1.FailurePolicyType`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#failurepolicytype-v1-admissionregistration) |  |  | <p>failurePolicy defines how to handle failures for the admission policy. Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions or bindings. failurePolicy does not define how validations that evaluate to false are handled. When failurePolicy is set to Fail, the validationActions field define how failures are enforced. Allowed values are Ignore or Fail. Defaults to Fail.</p> |
| `auditAnnotations` | [`[]admissionregistration/v1.AuditAnnotation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#auditannotation-v1-admissionregistration) |  |  | <p>auditAnnotations contains CEL expressions which are used to produce audit annotations for the audit event of the API request. validations and auditAnnotations may not both be empty; a least one of validations or auditAnnotations is required.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. If a parameter object is provided, it can be accessed via the `params` handle in the same manner as validation expressions. The exact matching logic is (in order):   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.   3. If any matchCondition evaluates to an error (but none are FALSE):      - If failurePolicy=Fail, reject the request      - If failurePolicy=Ignore, the policy is skipped</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `autogen` | [`ValidatingPolicyAutogenConfiguration`](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration) |  |  | <p>AutogenConfiguration defines the configuration for the generation controller.</p> |
| `validationActions` | [`[]admissionregistration/v1.ValidationAction`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validationaction-v1-admissionregistration) | :white_check_mark: |  | <p>ValidationAction specifies the action to be taken when the matched resource violates the policy. Required.</p> |
| `webhookConfiguration` | [`WebhookConfiguration`](#envoy-kyverno-io-v1alpha1-WebhookConfiguration) |  |  | <p>WebhookConfiguration defines the configuration for the webhook.</p> |
| `evaluation` | [`EvaluationConfiguration`](#envoy-kyverno-io-v1alpha1-EvaluationConfiguration) |  |  | <p>EvaluationConfiguration defines the configuration for the policy evaluation.</p> |

## ValidatingPolicyStatus     {#envoy-kyverno-io-v1alpha1-ValidatingPolicyStatus}

**Appears in:**
    
- [ValidatingPolicy](#envoy-kyverno-io-v1alpha1-ValidatingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#envoy-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |
| `autogen` | [`ValidatingPolicyAutogenStatus`](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus) |  |  | *No description provided.* |
| `generated` | `bool` |  |  | <p>Generated indicates whether a ValidatingAdmissionPolicy/MutatingAdmissionPolicy is generated from the policy or not</p> |

## VapGenerationConfiguration     {#envoy-kyverno-io-v1alpha1-VapGenerationConfiguration}

**Appears in:**
    
- [ValidatingPolicyAutogenConfiguration](#envoy-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` | :white_check_mark: |  | <p>Enabled specifies whether to generate a Kubernetes ValidatingAdmissionPolicy. Optional. Defaults to "false" if not specified.</p> |

## WebhookConfiguration     {#envoy-kyverno-io-v1alpha1-WebhookConfiguration}

**Appears in:**
    
- [ValidatingPolicySpec](#envoy-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeoutSeconds` | `int32` | :white_check_mark: |  | <p>TimeoutSeconds specifies the maximum time in seconds allowed to apply this policy. After the configured time expires, the admission request may fail, or may simply ignore the policy results, based on the failure policy. The default timeout is 10s, the value must be between 1 and 30 seconds.</p> |

  