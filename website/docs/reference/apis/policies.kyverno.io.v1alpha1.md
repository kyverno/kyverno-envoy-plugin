---
title: policies.kyverno.io (v1alpha1)
content_type: tool-reference
package: policies.kyverno.io/v1alpha1
auto_generated: true
---


## Resource Types 


- [DeletingPolicy](#policies-kyverno-io-v1alpha1-DeletingPolicy)
- [GeneratingPolicy](#policies-kyverno-io-v1alpha1-GeneratingPolicy)
- [ImageValidatingPolicy](#policies-kyverno-io-v1alpha1-ImageValidatingPolicy)
- [MutatingPolicy](#policies-kyverno-io-v1alpha1-MutatingPolicy)
- [PolicyException](#policies-kyverno-io-v1alpha1-PolicyException)
- [ValidatingPolicy](#policies-kyverno-io-v1alpha1-ValidatingPolicy)
  
## DeletingPolicy     {#policies-kyverno-io-v1alpha1-DeletingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `DeletingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`DeletingPolicySpec`](#policies-kyverno-io-v1alpha1-DeletingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`DeletingPolicyStatus`](#policies-kyverno-io-v1alpha1-DeletingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## GeneratingPolicy     {#policies-kyverno-io-v1alpha1-GeneratingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `GeneratingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`GeneratingPolicySpec`](#policies-kyverno-io-v1alpha1-GeneratingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`GeneratingPolicyStatus`](#policies-kyverno-io-v1alpha1-GeneratingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## ImageValidatingPolicy     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ImageValidatingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ImageValidatingPolicySpec`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`ImageValidatingPolicyStatus`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## MutatingPolicy     {#policies-kyverno-io-v1alpha1-MutatingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `MutatingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`MutatingPolicySpec`](#policies-kyverno-io-v1alpha1-MutatingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`MutatingPolicyStatus`](#policies-kyverno-io-v1alpha1-MutatingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## PolicyException     {#policies-kyverno-io-v1alpha1-PolicyException}

<p>PolicyException declares resources to be excluded from specified policies.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `PolicyException` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`PolicyExceptionSpec`](#policies-kyverno-io-v1alpha1-PolicyExceptionSpec) | :white_check_mark: |  | <p>Spec declares policy exception behaviors.</p> |

## ValidatingPolicy     {#policies-kyverno-io-v1alpha1-ValidatingPolicy}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `policies.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidatingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ValidatingPolicySpec`](#policies-kyverno-io-v1alpha1-ValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |
| `status` | [`ValidatingPolicyStatus`](#policies-kyverno-io-v1alpha1-ValidatingPolicyStatus) |  |  | <p>Status contains policy runtime data.</p> |

## AdmissionConfiguration     {#policies-kyverno-io-v1alpha1-AdmissionConfiguration}

**Appears in:**
    
- [EvaluationConfiguration](#policies-kyverno-io-v1alpha1-EvaluationConfiguration)
- [GeneratingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration)
- [MutatingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-MutatingPolicyEvaluationConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls if rules are applied during admission. Optional. Default value is "true".</p> |

## Attestation     {#policies-kyverno-io-v1alpha1-Attestation}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

<p>Attestation defines the identification details of the  metadata that has to be verified</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the name for this attestation. It is used to refer to the attestation in verification</p> |
| `intoto` | [`InToto`](#policies-kyverno-io-v1alpha1-InToto) |  |  | <p>InToto defines the details of attestation attached using intoto format</p> |
| `referrer` | [`Referrer`](#policies-kyverno-io-v1alpha1-Referrer) |  |  | <p>Referrer defines the details of attestation attached using OCI 1.1 format</p> |

## Attestor     {#policies-kyverno-io-v1alpha1-Attestor}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

<p>Attestor is an identity that confirms or verifies the authenticity of an image or an attestation</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the name for this attestor. It is used to refer to the attestor in verification</p> |
| `cosign` | [`Cosign`](#policies-kyverno-io-v1alpha1-Cosign) |  |  | <p>Cosign defines attestor configuration for Cosign based signatures</p> |
| `notary` | [`Notary`](#policies-kyverno-io-v1alpha1-Notary) |  |  | <p>Notary defines attestor configuration for Notary based signatures</p> |

## BackgroundConfiguration     {#policies-kyverno-io-v1alpha1-BackgroundConfiguration}

**Appears in:**
    
- [EvaluationConfiguration](#policies-kyverno-io-v1alpha1-EvaluationConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls if rules are applied to existing resources during a background scan. Optional. Default value is "true". The value must be set to "false" if the policy rule uses variables that are only available in the admission review request (e.g. user name).</p> |

## CTLog     {#policies-kyverno-io-v1alpha1-CTLog}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>CTLog sets the configuration to verify the authority against a Rekor instance.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `url` | `string` |  |  | <p>URL sets the url to the rekor instance (by default the public rekor.sigstore.dev)</p> |
| `rekorPubKey` | `string` |  |  | <p>RekorPubKey is an optional PEM-encoded public key to use for a custom Rekor. If set, this will be used to validate transparency log signatures from a custom Rekor.</p> |
| `ctLogPubKey` | `string` |  |  | <p>CTLogPubKey, if set, is used to validate SCTs against a custom source.</p> |
| `tsaCertChain` | `string` |  |  | <p>TSACertChain, if set, is the PEM-encoded certificate chain file for the RFC3161 timestamp authority. Must contain the root CA certificate. Optionally may contain intermediate CA certificates, and may contain the leaf TSA certificate if not present in the timestamurce.</p> |
| `insecureIgnoreTlog` | `bool` |  |  | <p>InsecureIgnoreTlog skips transparency log verification.</p> |
| `insecureIgnoreSCT` | `bool` |  |  | <p>IgnoreSCT defines whether to use the Signed Certificate Timestamp (SCT) log to check for a certificate timestamp. Default is false. Set to true if this was opted out during signing.</p> |

## Certificate     {#policies-kyverno-io-v1alpha1-Certificate}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>Certificate defines the configuration for local signature verification</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `cert` | [`StringOrExpression`](#policies-kyverno-io-v1alpha1-StringOrExpression) |  |  | <p>Certificate is the to the public certificate for local signature verification.</p> |
| `certChain` | [`StringOrExpression`](#policies-kyverno-io-v1alpha1-StringOrExpression) |  |  | <p>CertificateChain is the list of CA certificates in PEM format which will be needed when building the certificate chain for the signing certificate. Must start with the parent intermediate CA certificate of the signing certificate and end with the root certificate</p> |

## ConditionStatus     {#policies-kyverno-io-v1alpha1-ConditionStatus}

**Appears in:**
    
- [DeletingPolicyStatus](#policies-kyverno-io-v1alpha1-DeletingPolicyStatus)
- [GeneratingPolicyStatus](#policies-kyverno-io-v1alpha1-GeneratingPolicyStatus)
- [ImageValidatingPolicyStatus](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyStatus)
- [MutatingPolicyStatus](#policies-kyverno-io-v1alpha1-MutatingPolicyStatus)
- [ValidatingPolicyStatus](#policies-kyverno-io-v1alpha1-ValidatingPolicyStatus)

<p>ConditionStatus is the shared status across all policy types</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ready` | `bool` |  |  | <p>The ready of a policy is a high-level summary of where the policy is in its lifecycle. The conditions array, the reason and message fields contain more detail about the policy's status.</p> |
| `conditions` | [`[]meta/v1.Condition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#condition-v1-meta) |  |  | *No description provided.* |
| `message` | `string` |  |  | <p>Message is a human readable message indicating details about the generation of ValidatingAdmissionPolicy/MutatingAdmissionPolicy It is an empty string when ValidatingAdmissionPolicy/MutatingAdmissionPolicy is successfully generated.</p> |

## Cosign     {#policies-kyverno-io-v1alpha1-Cosign}

**Appears in:**
    
- [Attestor](#policies-kyverno-io-v1alpha1-Attestor)

<p>Cosign defines attestor configuration for Cosign based signatures</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `key` | [`Key`](#policies-kyverno-io-v1alpha1-Key) |  |  | <p>Key defines the type of key to validate the image.</p> |
| `keyless` | [`Keyless`](#policies-kyverno-io-v1alpha1-Keyless) |  |  | <p>Keyless sets the configuration to verify the authority against a Fulcio instance.</p> |
| `certificate` | [`Certificate`](#policies-kyverno-io-v1alpha1-Certificate) |  |  | <p>Certificate defines the configuration for local signature verification</p> |
| `source` | [`Source`](#policies-kyverno-io-v1alpha1-Source) |  |  | <p>Sources sets the configuration to specify the sources from where to consume the signature and attestations.</p> |
| `ctlog` | [`CTLog`](#policies-kyverno-io-v1alpha1-CTLog) |  |  | <p>CTLog sets the configuration to verify the authority against a Rekor instance.</p> |
| `tuf` | [`TUF`](#policies-kyverno-io-v1alpha1-TUF) |  |  | <p>TUF defines the configuration to fetch sigstore root</p> |
| `annotations` | `map[string]string` |  |  | <p>Annotations are used for image verification. Every specified key-value pair must exist and match in the verified payload. The payload may contain other key-value pairs.</p> |

## Credentials     {#policies-kyverno-io-v1alpha1-Credentials}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `allowInsecureRegistry` | `bool` |  |  | <p>AllowInsecureRegistry allows insecure access to a registry.</p> |
| `providers` | [`[]CredentialsProvidersType`](#policies-kyverno-io-v1alpha1-CredentialsProvidersType) |  |  | <p>Providers specifies a list of OCI Registry names, whose authentication providers are provided. It can be of one of these values: default,google,azure,amazon,github.</p> |
| `secrets` | `[]string` |  |  | <p>Secrets specifies a list of secrets that are provided for credentials. Secrets must live in the Kyverno namespace.</p> |

## CredentialsProvidersType     {#policies-kyverno-io-v1alpha1-CredentialsProvidersType}

(Alias of `string`)

**Appears in:**
    
- [Credentials](#policies-kyverno-io-v1alpha1-Credentials)

<p>CredentialsProvidersType provides the list of credential providers required.</p>


## DeletingPolicySpec     {#policies-kyverno-io-v1alpha1-DeletingPolicySpec}

**Appears in:**
    
- [DeletingPolicy](#policies-kyverno-io-v1alpha1-DeletingPolicy)

<p>DeletingPolicySpec is the specification of the desired behavior of the DeletingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1-admissionregistration) | :white_check_mark: |  | <p>MatchConstraints specifies what resources this policy is designed to validate. The AdmissionPolicy cares about a request if it matches _all_ Constraints. Required.</p> |
| `conditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>Conditions is a list of conditions that must be met for a resource to be deleted. Conditions filter resources that have already been matched by the match constraints, namespaceSelector, and objectSelector. An empty list of conditions matches all resources. There are a maximum of 64 conditions allowed. The exact matching logic is (in order):   1. If ANY condition evaluates to FALSE, the policy is skipped.   2. If ALL conditions evaluate to TRUE, the policy is executed.</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `schedule` | `string` | :white_check_mark: |  | <p>The schedule in Cron format Required.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy defines how resources will be deleted (Foreground, Background, Orphan).</p> |

## DeletingPolicyStatus     {#policies-kyverno-io-v1alpha1-DeletingPolicyStatus}

**Appears in:**
    
- [DeletingPolicy](#policies-kyverno-io-v1alpha1-DeletingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#policies-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |
| `lastExecutionTime` | [`meta/v1.Time`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta) | :white_check_mark: |  | *No description provided.* |

## EvaluationConfiguration     {#policies-kyverno-io-v1alpha1-EvaluationConfiguration}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)
- [ValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `mode` | [`EvaluationMode`](#policies-kyverno-io-v1alpha1-EvaluationMode) |  |  | <p>Mode is the mode of policy evaluation. Allowed values are "Kubernetes" or "JSON". Optional. Default value is "Kubernetes".</p> |
| `admission` | [`AdmissionConfiguration`](#policies-kyverno-io-v1alpha1-AdmissionConfiguration) |  |  | <p>Admission controls policy evaluation during admission.</p> |
| `background` | [`BackgroundConfiguration`](#policies-kyverno-io-v1alpha1-BackgroundConfiguration) |  |  | <p>Background  controls policy evaluation during background scan.</p> |

## EvaluationMode     {#policies-kyverno-io-v1alpha1-EvaluationMode}

(Alias of `string`)

**Appears in:**
    
- [EvaluationConfiguration](#policies-kyverno-io-v1alpha1-EvaluationConfiguration)

## GenerateExistingConfiguration     {#policies-kyverno-io-v1alpha1-GenerateExistingConfiguration}

**Appears in:**
    
- [GeneratingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration)

<p>GenerateExistingConfiguration defines the configuration for generating resources for existing triggers.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls whether to trigger the policy for existing resources If is set to "true" the policy will be triggered and applied to existing matched resources. Optional. Defaults to "false" if not specified.</p> |

## GeneratingPolicyEvaluationConfiguration     {#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration}

**Appears in:**
    
- [GeneratingPolicySpec](#policies-kyverno-io-v1alpha1-GeneratingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `admission` | [`AdmissionConfiguration`](#policies-kyverno-io-v1alpha1-AdmissionConfiguration) |  |  | <p>Admission controls policy evaluation during admission.</p> |
| `generateExisting` | [`GenerateExistingConfiguration`](#policies-kyverno-io-v1alpha1-GenerateExistingConfiguration) |  |  | <p>GenerateExisting defines the configuration for generating resources for existing triggeres.</p> |
| `synchronize` | [`SynchronizationConfiguration`](#policies-kyverno-io-v1alpha1-SynchronizationConfiguration) |  |  | <p>Synchronization defines the configuration for the synchronization of generated resources.</p> |
| `orphanDownstreamOnPolicyDelete` | [`OrphanDownstreamOnPolicyDeleteConfiguration`](#policies-kyverno-io-v1alpha1-OrphanDownstreamOnPolicyDeleteConfiguration) | :white_check_mark: |  | <p>OrphanDownstreamOnPolicyDelete defines the configuration for orphaning downstream resources on policy delete.</p> |

## GeneratingPolicySpec     {#policies-kyverno-io-v1alpha1-GeneratingPolicySpec}

**Appears in:**
    
- [GeneratingPolicy](#policies-kyverno-io-v1alpha1-GeneratingPolicy)

<p>GeneratingPolicySpec is the specification of the desired behavior of the GeneratingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1-admissionregistration) | :white_check_mark: |  | <p>MatchConstraints specifies what resources will trigger this policy. The AdmissionPolicy cares about a request if it matches _all_ Constraints. Required.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. If a parameter object is provided, it can be accessed via the `params` handle in the same manner as validation expressions. The exact matching logic is (in order):   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.   3. If any matchCondition evaluates to an error (but none are FALSE):      - If failurePolicy=Fail, reject the request      - If failurePolicy=Ignore, the policy is skipped</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `evaluation` | [`GeneratingPolicyEvaluationConfiguration`](#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration) |  |  | <p>EvaluationConfiguration defines the configuration for the policy evaluation.</p> |
| `webhookConfiguration` | [`WebhookConfiguration`](#policies-kyverno-io-v1alpha1-WebhookConfiguration) |  |  | <p>WebhookConfiguration defines the configuration for the webhook.</p> |
| `generate` | [`[]Generation`](#policies-kyverno-io-v1alpha1-Generation) | :white_check_mark: |  | <p>Generation defines a set of CEL expressions that will be evaluated to generate resources. Required.</p> |

## GeneratingPolicyStatus     {#policies-kyverno-io-v1alpha1-GeneratingPolicyStatus}

**Appears in:**
    
- [GeneratingPolicy](#policies-kyverno-io-v1alpha1-GeneratingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#policies-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |

## Generation     {#policies-kyverno-io-v1alpha1-Generation}

**Appears in:**
    
- [GeneratingPolicySpec](#policies-kyverno-io-v1alpha1-GeneratingPolicySpec)

<p>Generation defines the configuration for the generation of resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `expression` | `string` | :white_check_mark: |  | <p>Expression is a CEL expression that takes a list of resources to be generated.</p> |

## Identity     {#policies-kyverno-io-v1alpha1-Identity}

**Appears in:**
    
- [Keyless](#policies-kyverno-io-v1alpha1-Keyless)

<p>Identity may contain the issuer and/or the subject found in the transparency
log.
Issuer/Subject uses a strict match, while IssuerRegExp and SubjectRegExp
apply a regexp for matching.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `issuer` | `string` |  |  | <p>Issuer defines the issuer for this identity.</p> |
| `subject` | `string` |  |  | <p>Subject defines the subject for this identity.</p> |
| `issuerRegExp` | `string` |  |  | <p>IssuerRegExp specifies a regular expression to match the issuer for this identity.</p> |
| `subjectRegExp` | `string` |  |  | <p>SubjectRegExp specifies a regular expression to match the subject for this identity.</p> |

## ImageExtractor     {#policies-kyverno-io-v1alpha1-ImageExtractor}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the name for this imageList. It is used to refer to the images in verification block as images.<name></p> |
| `expression` | `string` | :white_check_mark: |  | <p>Expression defines CEL expression to extract images from the resource.</p> |

## ImageValidatingPolicyAutogen     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogen}

**Appears in:**
    
- [ImageValidatingPolicyAutogenStatus](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `targets` | [`[]Target`](#policies-kyverno-io-v1alpha1-Target) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ImageValidatingPolicySpec`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |

## ImageValidatingPolicyAutogenConfiguration     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenConfiguration}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `podControllers` | [`PodControllersGenerationConfiguration`](#policies-kyverno-io-v1alpha1-PodControllersGenerationConfiguration) | :white_check_mark: |  | <p>PodControllers specifies whether to generate a pod controllers rules.</p> |

## ImageValidatingPolicyAutogenStatus     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenStatus}

**Appears in:**
    
- [ImageValidatingPolicyStatus](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `configs` | [`map[string]ImageValidatingPolicyAutogen`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogen) | :white_check_mark: |  | *No description provided.* |

## ImageValidatingPolicySpec     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec}

**Appears in:**
    
- [ImageValidatingPolicy](#policies-kyverno-io-v1alpha1-ImageValidatingPolicy)
- [ImageValidatingPolicyAutogen](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogen)

<p>ImageValidatingPolicySpec is the specification of the desired behavior of the ImageValidatingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1-admissionregistration) |  |  | <p>MatchConstraints specifies what resources this policy is designed to validate.</p> |
| `failurePolicy` | [`admissionregistration/v1.FailurePolicyType`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#failurepolicytype-v1-admissionregistration) |  |  | <p>FailurePolicy defines how to handle failures for the admission policy. Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions or bindings.</p> |
| `auditAnnotations` | [`[]admissionregistration/v1.AuditAnnotation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#auditannotation-v1-admissionregistration) |  |  | <p>auditAnnotations contains CEL expressions which are used to produce audit annotations for the audit event of the API request. validations and auditAnnotations may not both be empty; a least one of validations or auditAnnotations is required.</p> |
| `validationActions` | [`[]admissionregistration/v1.ValidationAction`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validationaction-v1-admissionregistration) | :white_check_mark: |  | <p>ValidationAction specifies the action to be taken when the matched resource violates the policy. Required.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed.</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression.</p> |
| `validationConfigurations` | [`ValidationConfiguration`](#policies-kyverno-io-v1alpha1-ValidationConfiguration) |  |  | <p>ValidationConfigurations defines settings for mutating and verifying image digests, and enforcing image verification through signatures.</p> |
| `matchImageReferences` | [`[]MatchImageReference`](#policies-kyverno-io-v1alpha1-MatchImageReference) |  |  | <p>MatchImageReferences is a list of Glob and CELExpressions to match images. Any image that matches one of the rules is considered for validation Any image that does not match a rule is skipped, even when they are passed as arguments to image verification functions</p> |
| `credentials` | [`Credentials`](#policies-kyverno-io-v1alpha1-Credentials) | :white_check_mark: |  | <p>Credentials provides credentials that will be used for authentication with registry.</p> |
| `images` | [`[]ImageExtractor`](#policies-kyverno-io-v1alpha1-ImageExtractor) |  |  | <p>ImageExtractors is a list of CEL expression to extract images from the resource</p> |
| `attestors` | [`[]Attestor`](#policies-kyverno-io-v1alpha1-Attestor) | :white_check_mark: |  | <p>Attestors provides a list of trusted authorities.</p> |
| `attestations` | [`[]Attestation`](#policies-kyverno-io-v1alpha1-Attestation) |  |  | <p>Attestations provides a list of image metadata to verify</p> |
| `validations` | [`[]admissionregistration/v1.Validation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validation-v1-admissionregistration) | :white_check_mark: |  | <p>Validations contain CEL expressions which is used to apply the image validation checks.</p> |
| `webhookConfiguration` | [`WebhookConfiguration`](#policies-kyverno-io-v1alpha1-WebhookConfiguration) |  |  | <p>WebhookConfiguration defines the configuration for the webhook.</p> |
| `evaluation` | [`EvaluationConfiguration`](#policies-kyverno-io-v1alpha1-EvaluationConfiguration) |  |  | <p>EvaluationConfiguration defines the configuration for the policy evaluation.</p> |
| `autogen` | [`ImageValidatingPolicyAutogenConfiguration`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenConfiguration) |  |  | <p>AutogenConfiguration defines the configuration for the generation controller.</p> |

## ImageValidatingPolicyStatus     {#policies-kyverno-io-v1alpha1-ImageValidatingPolicyStatus}

**Appears in:**
    
- [ImageValidatingPolicy](#policies-kyverno-io-v1alpha1-ImageValidatingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#policies-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |
| `autogen` | [`ImageValidatingPolicyAutogenStatus`](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenStatus) |  |  | *No description provided.* |

## InToto     {#policies-kyverno-io-v1alpha1-InToto}

**Appears in:**
    
- [Attestation](#policies-kyverno-io-v1alpha1-Attestation)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `type` | `string` | :white_check_mark: |  | <p>Type defines the type of attestation contained within the statement.</p> |

## Key     {#policies-kyverno-io-v1alpha1-Key}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>A Key must specify only one of CEL, Data or KMS</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `data` | `string` |  |  | <p>Data contains the inline public key</p> |
| `kms` | `string` |  |  | <p>KMS contains the KMS url of the public key Supported formats differ based on the KMS system used.</p> |
| `hashAlgorithm` | `string` |  |  | <p>HashAlgorithm specifues signature algorithm for public keys. Supported values are sha224, sha256, sha384 and sha512. Defaults to sha256.</p> |
| `expression` | `string` |  |  | <p>Expression is a Expression expression that returns the public key.</p> |

## Keyless     {#policies-kyverno-io-v1alpha1-Keyless}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>Keyless contains location of the validating certificate and the identities
against which to verify.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `identities` | [`[]Identity`](#policies-kyverno-io-v1alpha1-Identity) | :white_check_mark: |  | <p>Identities sets a list of identities.</p> |
| `roots` | `string` | :white_check_mark: |  | <p>Roots is an optional set of PEM encoded trusted root certificates. If not provided, the system roots are used.</p> |

## MAPGenerationConfiguration     {#policies-kyverno-io-v1alpha1-MAPGenerationConfiguration}

**Appears in:**
    
- [MutatingPolicyAutogenConfiguration](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` | :white_check_mark: |  | <p>Enabled specifies whether to generate a Kubernetes MutatingAdmissionPolicy. Optional. Defaults to "false" if not specified.</p> |

## MatchImageReference     {#policies-kyverno-io-v1alpha1-MatchImageReference}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

<p>MatchImageReference defines a Glob or a CEL expression for matching images</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `glob` | `string` |  |  | <p>Glob defines a globbing pattern for matching images</p> |
| `expression` | `string` |  |  | <p>Expression defines CEL Expressions for matching images</p> |

## MutateExistingConfiguration     {#policies-kyverno-io-v1alpha1-MutateExistingConfiguration}

**Appears in:**
    
- [MutatingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-MutatingPolicyEvaluationConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled enables mutation of existing resources. Default is false. When spec.targetMatchConstraints is not defined, Kyverno mutates existing resources matched in spec.matchConstraints.</p> |

## MutatingPolicyAutogen     {#policies-kyverno-io-v1alpha1-MutatingPolicyAutogen}

**Appears in:**
    
- [MutatingPolicyAutogenStatus](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `targets` | [`[]Target`](#policies-kyverno-io-v1alpha1-Target) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`MutatingPolicySpec`](#policies-kyverno-io-v1alpha1-MutatingPolicySpec) | :white_check_mark: |  | *No description provided.* |

## MutatingPolicyAutogenConfiguration     {#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenConfiguration}

**Appears in:**
    
- [MutatingPolicySpec](#policies-kyverno-io-v1alpha1-MutatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `podControllers` | [`PodControllersGenerationConfiguration`](#policies-kyverno-io-v1alpha1-PodControllersGenerationConfiguration) | :white_check_mark: |  | <p>PodControllers specifies whether to generate a pod controllers rules.</p> |
| `mutatingAdmissionPolicy` | [`MAPGenerationConfiguration`](#policies-kyverno-io-v1alpha1-MAPGenerationConfiguration) | :white_check_mark: |  | <p>MutatingAdmissionPolicy specifies whether to generate a Kubernetes MutatingAdmissionPolicy.</p> |

## MutatingPolicyAutogenStatus     {#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenStatus}

**Appears in:**
    
- [MutatingPolicyStatus](#policies-kyverno-io-v1alpha1-MutatingPolicyStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `configs` | [`map[string]MutatingPolicyAutogen`](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogen) | :white_check_mark: |  | *No description provided.* |

## MutatingPolicyEvaluationConfiguration     {#policies-kyverno-io-v1alpha1-MutatingPolicyEvaluationConfiguration}

**Appears in:**
    
- [MutatingPolicySpec](#policies-kyverno-io-v1alpha1-MutatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `admission` | [`AdmissionConfiguration`](#policies-kyverno-io-v1alpha1-AdmissionConfiguration) |  |  | <p>Admission controls policy evaluation during admission.</p> |
| `mutateExisting` | [`MutateExistingConfiguration`](#policies-kyverno-io-v1alpha1-MutateExistingConfiguration) |  |  | <p>MutateExisting controls whether existing resources are mutated.</p> |

## MutatingPolicySpec     {#policies-kyverno-io-v1alpha1-MutatingPolicySpec}

**Appears in:**
    
- [MutatingPolicy](#policies-kyverno-io-v1alpha1-MutatingPolicy)
- [MutatingPolicyAutogen](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogen)

<p>MutatingPolicySpec is the specification of the desired behavior of the MutatingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1alpha1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1alpha1-admissionregistration) | :white_check_mark: |  | <p>MatchConstraints specifies what resources this policy is designed to evaluate. The AdmissionPolicy cares about a request if it matches _all_ Constraints. Required.</p> |
| `failurePolicy` | [`admissionregistration/v1alpha1.FailurePolicyType`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#failurepolicytype-v1alpha1-admissionregistration) |  |  | <p>failurePolicy defines how to handle failures for the admission policy. Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions or bindings. failurePolicy does not define how validations that evaluate to false are handled. When failurePolicy is set to Fail, the validationActions field define how failures are enforced. Allowed values are Ignore or Fail. Defaults to Fail.</p> |
| `matchConditions` | [`[]admissionregistration/v1alpha1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1alpha1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. If a parameter object is provided, it can be accessed via the `params` handle in the same manner as validation expressions. The exact matching logic is (in order):   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.   3. If any matchCondition evaluates to an error (but none are FALSE):      - If failurePolicy=Fail, reject the request      - If failurePolicy=Ignore, the policy is skipped</p> |
| `variables` | [`[]admissionregistration/v1alpha1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1alpha1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `autogen` | [`MutatingPolicyAutogenConfiguration`](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenConfiguration) |  |  | <p>AutogenConfiguration defines the configuration for the generation controller.</p> |
| `targetMatchConstraints` | [`admissionregistration/v1alpha1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1alpha1-admissionregistration) |  |  | <p>TargetMatchConstraints specifies what target mutation resources this policy is designed to evaluate.</p> |
| `mutations` | [`[]admissionregistration/v1alpha1.Mutation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#mutation-v1alpha1-admissionregistration) |  |  | <p>mutations contain operations to perform on matching objects. mutations may not be empty; a minimum of one mutation is required. mutations are evaluated in order, and are reinvoked according to the reinvocationPolicy. The mutations of a policy are invoked for each binding of this policy and reinvocation of mutations occurs on a per binding basis.</p> |
| `webhookConfiguration` | [`WebhookConfiguration`](#policies-kyverno-io-v1alpha1-WebhookConfiguration) |  |  | <p>WebhookConfiguration defines the configuration for the webhook.</p> |
| `evaluation` | [`MutatingPolicyEvaluationConfiguration`](#policies-kyverno-io-v1alpha1-MutatingPolicyEvaluationConfiguration) |  |  | <p>EvaluationConfiguration defines the configuration for mutating policy evaluation.</p> |
| `reinvocationPolicy` | `admissionregistration/v1alpha1.ReinvocationPolicyType` | :white_check_mark: |  | <p>reinvocationPolicy indicates whether mutations may be called multiple times per MutatingAdmissionPolicyBinding as part of a single admission evaluation. Allowed values are "Never" and "IfNeeded". Never: These mutations will not be called more than once per binding in a single admission evaluation. IfNeeded: These mutations may be invoked more than once per binding for a single admission request and there is no guarantee of order with respect to other admission plugins, admission webhooks, bindings of this policy and admission policies.  Mutations are only reinvoked when mutations change the object after this mutation is invoked. Required.</p> |

## MutatingPolicyStatus     {#policies-kyverno-io-v1alpha1-MutatingPolicyStatus}

**Appears in:**
    
- [MutatingPolicy](#policies-kyverno-io-v1alpha1-MutatingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#policies-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |
| `autogen` | [`MutatingPolicyAutogenStatus`](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenStatus) |  |  | *No description provided.* |
| `generated` | `bool` |  |  | <p>Generated indicates whether a MutatingAdmissionPolicy is generated from the policy or not</p> |

## Notary     {#policies-kyverno-io-v1alpha1-Notary}

**Appears in:**
    
- [Attestor](#policies-kyverno-io-v1alpha1-Attestor)

<p>Notary defines attestor configuration for Notary based signatures</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `certs` | [`StringOrExpression`](#policies-kyverno-io-v1alpha1-StringOrExpression) |  |  | <p>Certs define the cert chain for Notary signature verification</p> |
| `tsaCerts` | [`StringOrExpression`](#policies-kyverno-io-v1alpha1-StringOrExpression) |  |  | <p>TSACerts define the cert chain for verifying timestamps of notary signature</p> |

## OrphanDownstreamOnPolicyDeleteConfiguration     {#policies-kyverno-io-v1alpha1-OrphanDownstreamOnPolicyDeleteConfiguration}

**Appears in:**
    
- [GeneratingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration)

<p>OrphanDownstreamOnPolicyDeleteConfiguration defines the configuration for orphaning downstream resources on policy delete.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls whether generated resources should be deleted when the policy that generated them is deleted with synchronization enabled. This option is only applicable to generate rules of the data type. Optional. Defaults to "false" if not specified.</p> |

## PodControllersGenerationConfiguration     {#policies-kyverno-io-v1alpha1-PodControllersGenerationConfiguration}

**Appears in:**
    
- [ImageValidatingPolicyAutogenConfiguration](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogenConfiguration)
- [MutatingPolicyAutogenConfiguration](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogenConfiguration)
- [ValidatingPolicyAutogenConfiguration](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `controllers` | `[]string` | :white_check_mark: |  | *No description provided.* |

## PolicyExceptionSpec     {#policies-kyverno-io-v1alpha1-PolicyExceptionSpec}

**Appears in:**
    
- [PolicyException](#policies-kyverno-io-v1alpha1-PolicyException)

<p>PolicyExceptionSpec stores policy exception spec</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `policyRefs` | [`[]PolicyRef`](#policies-kyverno-io-v1alpha1-PolicyRef) | :white_check_mark: |  | <p>PolicyRefs identifies the policies to which the exception is applied.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of CEL expressions that must be met for a resource to be excluded.</p> |
| `images` | `[]string` |  |  | <p>Images specifies container images to be excluded from policy evaluation. These excluded images can be referenced in CEL expressions via `exceptions.allowedImages`.</p> |
| `allowedValues` | `[]string` |  |  | <p>AllowedValues specifies values that can be used in CEL expressions to bypass policy checks. These values can be referenced in CEL expressions via `exceptions.allowedValues`.</p> |

## PolicyRef     {#policies-kyverno-io-v1alpha1-PolicyRef}

**Appears in:**
    
- [PolicyExceptionSpec](#policies-kyverno-io-v1alpha1-PolicyExceptionSpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the name of the policy</p> |
| `kind` | `string` | :white_check_mark: |  | <p>Kind is the kind of the policy</p> |

## Referrer     {#policies-kyverno-io-v1alpha1-Referrer}

**Appears in:**
    
- [Attestation](#policies-kyverno-io-v1alpha1-Attestation)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `type` | `string` | :white_check_mark: |  | <p>Type defines the type of attestation attached to the image.</p> |

## Source     {#policies-kyverno-io-v1alpha1-Source}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>Source specifies the location of the signature / attestations.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `repository` | `string` |  |  | <p>Repository defines the location from where to pull the signature / attestations.</p> |
| `PullSecrets` | [`[]core/v1.LocalObjectReference`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#localobjectreference-v1-core) |  |  | <p>SignaturePullSecrets is an optional list of references to secrets in the same namespace as the deploying resource for pulling any of the signatures used by this Source.</p> |
| `tagPrefix` | `string` |  |  | <p>TagPrefix is an optional prefix that signature and attestations have. This is the 'tag based discovery' and in the future once references are fully supported that should likely be the preferred way to handle these.</p> |

## StringOrExpression     {#policies-kyverno-io-v1alpha1-StringOrExpression}

**Appears in:**
    
- [Certificate](#policies-kyverno-io-v1alpha1-Certificate)
- [Notary](#policies-kyverno-io-v1alpha1-Notary)

<p>StringOrExpression contains either a raw string input or a CEL expression</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `value` | `string` |  |  | <p>Value defines the raw string input.</p> |
| `expression` | `string` |  |  | <p>Expression defines the a CEL expression input.</p> |

## SynchronizationConfiguration     {#policies-kyverno-io-v1alpha1-SynchronizationConfiguration}

**Appears in:**
    
- [GeneratingPolicyEvaluationConfiguration](#policies-kyverno-io-v1alpha1-GeneratingPolicyEvaluationConfiguration)

<p>SynchronizationConfiguration defines the configuration for the synchronization of generated resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled controls if generated resources should be kept in-sync with their source resource. If Synchronize is set to "true" changes to generated resources will be overwritten with resource data from Data or the resource specified in the Clone declaration. Optional. Defaults to "false" if not specified.</p> |

## TUF     {#policies-kyverno-io-v1alpha1-TUF}

**Appears in:**
    
- [Cosign](#policies-kyverno-io-v1alpha1-Cosign)

<p>TUF defines the configuration to fetch sigstore root</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `root` | [`TUFRoot`](#policies-kyverno-io-v1alpha1-TUFRoot) |  |  | <p>Root defines the path or data of the trusted root</p> |
| `mirror` | `string` |  |  | <p>Mirror is the base URL of Sigstore TUF repository</p> |

## TUFRoot     {#policies-kyverno-io-v1alpha1-TUFRoot}

**Appears in:**
    
- [TUF](#policies-kyverno-io-v1alpha1-TUF)

<p>TUFRoot defines the path or data of the trusted root</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `path` | `string` |  |  | <p>Path is the URL or File location of the TUF root</p> |
| `data` | `string` |  |  | <p>Data is the base64 encoded TUF root</p> |

## Target     {#policies-kyverno-io-v1alpha1-Target}

**Appears in:**
    
- [ImageValidatingPolicyAutogen](#policies-kyverno-io-v1alpha1-ImageValidatingPolicyAutogen)
- [MutatingPolicyAutogen](#policies-kyverno-io-v1alpha1-MutatingPolicyAutogen)
- [ValidatingPolicyAutogen](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogen)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `group` | `string` | :white_check_mark: |  | *No description provided.* |
| `version` | `string` | :white_check_mark: |  | *No description provided.* |
| `resource` | `string` | :white_check_mark: |  | *No description provided.* |
| `kind` | `string` | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicyAutogen     {#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogen}

**Appears in:**
    
- [ValidatingPolicyAutogenStatus](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `targets` | [`[]Target`](#policies-kyverno-io-v1alpha1-Target) | :white_check_mark: |  | *No description provided.* |
| `spec` | [`ValidatingPolicySpec`](#policies-kyverno-io-v1alpha1-ValidatingPolicySpec) | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicyAutogenConfiguration     {#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration}

**Appears in:**
    
- [ValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `podControllers` | [`PodControllersGenerationConfiguration`](#policies-kyverno-io-v1alpha1-PodControllersGenerationConfiguration) | :white_check_mark: |  | <p>PodControllers specifies whether to generate a pod controllers rules.</p> |
| `validatingAdmissionPolicy` | [`VapGenerationConfiguration`](#policies-kyverno-io-v1alpha1-VapGenerationConfiguration) | :white_check_mark: |  | <p>ValidatingAdmissionPolicy specifies whether to generate a Kubernetes ValidatingAdmissionPolicy.</p> |

## ValidatingPolicyAutogenStatus     {#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus}

**Appears in:**
    
- [ValidatingPolicyStatus](#policies-kyverno-io-v1alpha1-ValidatingPolicyStatus)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `configs` | [`map[string]ValidatingPolicyAutogen`](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogen) | :white_check_mark: |  | *No description provided.* |

## ValidatingPolicySpec     {#policies-kyverno-io-v1alpha1-ValidatingPolicySpec}

**Appears in:**
    
- [ValidatingPolicy](#policies-kyverno-io-v1alpha1-ValidatingPolicy)
- [ValidatingPolicyAutogen](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogen)

<p>ValidatingPolicySpec is the specification of the desired behavior of the ValidatingPolicy.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `matchConstraints` | [`admissionregistration/v1.MatchResources`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchresources-v1-admissionregistration) | :white_check_mark: |  | <p>MatchConstraints specifies what resources this policy is designed to validate. The AdmissionPolicy cares about a request if it matches _all_ Constraints. Required.</p> |
| `validations` | [`[]admissionregistration/v1.Validation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validation-v1-admissionregistration) |  |  | <p>Validations contain CEL expressions which is used to apply the validation. Validations and AuditAnnotations may not both be empty; a minimum of one Validations or AuditAnnotations is required.</p> |
| `failurePolicy` | [`admissionregistration/v1.FailurePolicyType`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#failurepolicytype-v1-admissionregistration) |  |  | <p>failurePolicy defines how to handle failures for the admission policy. Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions or bindings. failurePolicy does not define how validations that evaluate to false are handled. When failurePolicy is set to Fail, the validationActions field define how failures are enforced. Allowed values are Ignore or Fail. Defaults to Fail.</p> |
| `auditAnnotations` | [`[]admissionregistration/v1.AuditAnnotation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#auditannotation-v1-admissionregistration) |  |  | <p>auditAnnotations contains CEL expressions which are used to produce audit annotations for the audit event of the API request. validations and auditAnnotations may not both be empty; a least one of validations or auditAnnotations is required.</p> |
| `matchConditions` | [`[]admissionregistration/v1.MatchCondition`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#matchcondition-v1-admissionregistration) |  |  | <p>MatchConditions is a list of conditions that must be met for a request to be validated. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. If a parameter object is provided, it can be accessed via the `params` handle in the same manner as validation expressions. The exact matching logic is (in order):   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.   3. If any matchCondition evaluates to an error (but none are FALSE):      - If failurePolicy=Fail, reject the request      - If failurePolicy=Ignore, the policy is skipped</p> |
| `variables` | [`[]admissionregistration/v1.Variable`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#variable-v1-admissionregistration) |  |  | <p>Variables contain definitions of variables that can be used in composition of other expressions. Each variable is defined as a named CEL expression. The variables defined here will be available under `variables` in other expressions of the policy except MatchConditions because MatchConditions are evaluated before the rest of the policy. The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, Variables must be sorted by the order of first appearance and acyclic.</p> |
| `autogen` | [`ValidatingPolicyAutogenConfiguration`](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration) |  |  | <p>AutogenConfiguration defines the configuration for the generation controller.</p> |
| `validationActions` | [`[]admissionregistration/v1.ValidationAction`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#validationaction-v1-admissionregistration) | :white_check_mark: |  | <p>ValidationAction specifies the action to be taken when the matched resource violates the policy. Required.</p> |
| `webhookConfiguration` | [`WebhookConfiguration`](#policies-kyverno-io-v1alpha1-WebhookConfiguration) |  |  | <p>WebhookConfiguration defines the configuration for the webhook.</p> |
| `evaluation` | [`EvaluationConfiguration`](#policies-kyverno-io-v1alpha1-EvaluationConfiguration) |  |  | <p>EvaluationConfiguration defines the configuration for the policy evaluation.</p> |

## ValidatingPolicyStatus     {#policies-kyverno-io-v1alpha1-ValidatingPolicyStatus}

**Appears in:**
    
- [ValidatingPolicy](#policies-kyverno-io-v1alpha1-ValidatingPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `conditionStatus` | [`ConditionStatus`](#policies-kyverno-io-v1alpha1-ConditionStatus) |  |  | *No description provided.* |
| `autogen` | [`ValidatingPolicyAutogenStatus`](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenStatus) |  |  | *No description provided.* |
| `generated` | `bool` |  |  | <p>Generated indicates whether a ValidatingAdmissionPolicy/MutatingAdmissionPolicy is generated from the policy or not</p> |

## ValidationConfiguration     {#policies-kyverno-io-v1alpha1-ValidationConfiguration}

**Appears in:**
    
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `mutateDigest` | `bool` |  |  | <p>MutateDigest enables replacement of image tags with digests. Defaults to true.</p> |
| `verifyDigest` | `bool` |  |  | <p>VerifyDigest validates that images have a digest.</p> |
| `required` | `bool` |  |  | <p>Required validates that images are verified, i.e., have passed a signature or attestation check.</p> |

## VapGenerationConfiguration     {#policies-kyverno-io-v1alpha1-VapGenerationConfiguration}

**Appears in:**
    
- [ValidatingPolicyAutogenConfiguration](#policies-kyverno-io-v1alpha1-ValidatingPolicyAutogenConfiguration)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` | :white_check_mark: |  | <p>Enabled specifies whether to generate a Kubernetes ValidatingAdmissionPolicy. Optional. Defaults to "false" if not specified.</p> |

## WebhookConfiguration     {#policies-kyverno-io-v1alpha1-WebhookConfiguration}

**Appears in:**
    
- [GeneratingPolicySpec](#policies-kyverno-io-v1alpha1-GeneratingPolicySpec)
- [ImageValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ImageValidatingPolicySpec)
- [MutatingPolicySpec](#policies-kyverno-io-v1alpha1-MutatingPolicySpec)
- [ValidatingPolicySpec](#policies-kyverno-io-v1alpha1-ValidatingPolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeoutSeconds` | `int32` | :white_check_mark: |  | <p>TimeoutSeconds specifies the maximum time in seconds allowed to apply this policy. After the configured time expires, the admission request may fail, or may simply ignore the policy results, based on the failure policy. The default timeout is 10s, the value must be between 1 and 30 seconds.</p> |

  