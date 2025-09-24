# JSON schemas

JSON schemas for the Kyverno Envoy Plugin are available:

- [ValidatingPolicy (v1alpha1)](https://raw.githubusercontent.com/kyverno/playground/main/schemas/json/v3/validatingpolicy-policies.kyverno.io-v1alpha1.json)

They can be used to enable validation and autocompletion in your IDE.

## VS code

In VS code, simply add a comment on top of your YAML resources.

### ValidatingPolicy

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/playground/main/schemas/json/v3/validatingpolicy-policies.kyverno.io-v1alpha1.json
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  evaluation:
    mode: Envoy
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.attributes.request.http.headers[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]
  - name: metadata
    expression: '{"my-new-metadata": "my-new-value"}'
  validations:
    # if force_unauthenticated -> 401
  - expression: >
      variables.force_unauthenticated
        ? envoy
            .Denied(401)
            .WithBody("Authentication Failed")
            .Response()
            .WithMetadata(variables.metadata)
        : null
    # if not force_authorized -> 403
  - expression: >
      !variables.force_authorized
        ? envoy
            .Denied(403)
            .WithBody("Unauthorized Request")
            .Response()
        : null
    # else -> 200
  - expression: >
      envoy
        .Allowed()
        .WithHeader("x-validated-by", "my-security-checkpoint")
        .WithoutHeader("x-force-authorized")
        .WithResponseHeader("x-add-custom-response-header", "added")
        .Response()
        .WithMetadata(variables.metadata)
```
