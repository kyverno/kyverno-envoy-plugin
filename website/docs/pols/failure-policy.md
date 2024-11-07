# Failure policy

FailurePolicy defines how to handle failures for the policy.

Failures can occur from CEL expression parse errors, type check errors, runtime errors and invalid or mis-configured policy definitions.

Allowed values are:

- `Ignore`
- `Fail`

If not set, the failure policy defaults to `Fail`.

!!!info

    FailurePolicy does not define how validations that evaluate to `false` are handled.

## Fail

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  # if something fails the request will be denied
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  authorizations:
  - expression: >
      variables.allowed
        ? envoy.Allowed().Response()
        : envoy.Denied(403).Response()
```

## Ignore

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  # if something fails the failure will be ignored and the request will be allowed
  failurePolicy: Ignore
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  authorizations:
  - expression: >
      variables.allowed
        ? envoy.Allowed().Response()
        : envoy.Denied(403).Response()
```
