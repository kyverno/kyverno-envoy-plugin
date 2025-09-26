# ValidatingPolicy

A Kyverno `ValidatingPolicy` is a custom [Kubernetes resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and can be easily managed via Kubernetes APIs, GitOps workflows, and other existing tools.

## Resource Scope

A Kyverno `ValidatingPolicy` is a cluster-wide resource.

## API Group and Kind

A `ValidatingPolicy` belongs to the `policies.kyverno.io/v1alpha1` group and can only be of kind `ValidatingPolicy`.

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  # if something fails the request will be denied
  failurePolicy: Fail
  evaluation:
    mode: Envoy
  variables:
    # `force_authorized` references the 'x-force-authorized' header
    # from the envoy check request (or '' if not present)
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
    # `allowed` will be `true` if `variables.force_authorized` has the
    # value 'enabled' or 'true'
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
    # make an authorisation decision based on the value of `variables.allowed`
  - expression: >
      !variables.allowed
        ? envoy.Denied(403).Response()
        : envoy.Allowed().Response()
```

## Envoy External Authorization

The Kyverno Authz Server implements the [Envoy External Authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter) API.

A Kyverno `ValidatingPolicy` analyses an Envoy [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) and can make a decision by returning a [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse).

## CEL language

A `ValidatingPolicy` uses the [CEL language](https://github.com/google/cel-spec) to process the `CheckRequest` sent by Envoy.

CEL is an expression language that’s fast, portable, and safe to execute in performance-critical applications.

## Policy structure

A Kyverno `ValidatingPolicy` is made of:

- An [evaluation mode](./evaluation-mode.md)
- A [failure policy](../failure-policy.md)
- [Match conditions](../match-conditions.md) if needed
- Eventually some [variables](../variables.md)
- The [validation rules](./validation-rules.md)
