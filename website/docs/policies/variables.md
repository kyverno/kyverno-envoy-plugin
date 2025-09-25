# Variables

A Kyverno `ValidatingPolicy` can define `variables` that will be made available to all authorization rules.

Variables can be used in composition of other expressions.
Each variable is defined as a named [CEL](https://github.com/google/cel-spec) expression.
The will be available under `variables` in other expressions of the policy.

The expression of a variable can refer to other variables defined earlier in the list but not those after. Thus, variables must be sorted by the order of first appearance and acyclic.

!!!info

    The incoming `CheckRequest` from Envoy is made available to the policy under the `object` identifier.

## Variables

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
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
        : null
```
