# Policies

A Kyverno `AuthorizationPolicy` is a custom [Kubernetes resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and can be easily managed via Kubernetes APIs, GitOps workflows, and other existing tools.

## Resource Scope

A Kyverno `AuthorizationPolicy` is a cluster-wide resource.

## API Group and Kind

An `AuthorizationPolicy` belongs to the `envoy.kyverno.io/v1alpha1` group and can only be of kind `AuthorizationPolicy`.

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  # if something fails the request will be denied
  failurePolicy: Fail
  variables:
    # `force_authorized` references the 'x-force-authorized' header
    # from the envoy check request (or '' if not present)
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
    # `allowed` will be `true` if `variables.force_authorized` has the
    # value 'enabled' or 'true'
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  deny:
    # make an authorisation decision based on the value of `variables.allowed`
  - match: >
      !variables.allowed
    response: >
      envoy.Denied(403).Response()
```

## Envoy External Authorization

The Kyverno Authz Server implements the [Envoy External Authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter) API.

A Kyverno `AuthorizationPolicy` analyses an Envoy [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) and can make a decision by returning an Envoy [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) (or nothing if no decision is made).

## CEL language

An `AuthorizationPolicy` uses the [CEL language](https://github.com/google/cel-spec) to process the `CheckRequest` sent by Envoy.

CEL is an expression language that’s fast, portable, and safe to execute in performance-critical applications.

## Policy structure

A Kyverno `AuthorizationPolicy` is made of:

- A [failure policy](./failure-policy.md)
- Eventually some [variables](./variables.md)
- The [authorization rules](./authorization-rules.md)
