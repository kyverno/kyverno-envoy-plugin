# Evualuation mode

A `ValidatingPolicy` is a generic policy definition resource and can be consumed by different tools.

To be considered by the Kyverno Authz Server, a `ValidatingPolicy` **must** have its `spec.evaluation.mode` set to `Envoy`.

## Example

Below is a policy with `spec.evaluation.mode` set to `Envoy`. This policy will apply to the Kyverno Authz Server:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  evaluation:
    # this policy will apply to authz server
    mode: Envoy
  validations:
  - expression: ...
```

This policy doesn't apply to the Kyverno Authz Server, its `spec.evaluation.mode` field is set to `Kubernetes`:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  evaluation:
    # this policy doesn't apply to authz server
    mode: Kubernetes
  validations:
  - expression: ...
```
