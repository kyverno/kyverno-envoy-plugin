# JSON schemas

JSON schemas for the Kyverno Envoy Plugin are available:

- [AuthorizationPolicy (v1alpha1)](https://github.com/kyverno/kyverno-envoy-plugin/blob/main/.schemas/json/authorizationpolicy-envoy-v1alpha1.json)

They can be used to enable validation and autocompletion in your IDE.

## VS code

In VS code, simply add a comment on top of your YAML resources.

### AuthorizationPolicy

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/main/.schemas/json/authorizationpolicy-envoy-v1alpha1.json
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo-policy.example.com
spec:
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.attributes.request.http.headers[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]
  authorizations:
  - expression: >
      variables.force_authorized && !variables.force_unauthenticated
      ? envoy
          .Allowed()
          .WithHeader("x-validated-by", "my-security-checkpoint")
          .WithoutHeader("x-force-authorized")
          .WithResponseHeader("x-add-custom-response-header", "added")
          .Response()
          .WithMetadata({"my-new-metadata": "my-new-value"})
      : envoy
          .Denied(variables.force_unauthenticated ? 401 : 403)
          .WithBody(variables.force_unauthenticated ? "Authentication Failed" : "Unauthorized Request")
          .Response()
```
