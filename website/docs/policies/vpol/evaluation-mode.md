# Evualuation mode

A `ValidatingPolicy` is a generic policy definition resource and can be consumed by different tools.

To be considered by the Kyverno Authz Server, a `ValidatingPolicy` **must** have its `spec.evaluation.mode` set to either `Envoy` or `HTTP`.

## Examples

### Envoy Mode

Below is a policy with `spec.evaluation.mode` set to `Envoy`. This policy will apply to the Kyverno Authz Server for Envoy proxy authorization:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo-envoy
spec:
  evaluation:
    # this policy will apply to authz server for Envoy
    mode: Envoy
  validations:
  - expression: ...
```

### HTTP Mode

Below is a policy with `spec.evaluation.mode` set to `HTTP`. This policy will apply to the Kyverno Authz Server for plain HTTP authorization:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo-http
spec:
  evaluation:
    # this policy will apply to authz server for HTTP
    mode: HTTP
  validations:
  - expression: ...
```

### Non-applicable Mode

This policy doesn't apply to the Kyverno Authz Server because its `spec.evaluation.mode` field is set to `Kubernetes`:

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
