# Match conditions

You can define match conditions if you need fine-grained request filtering.

Match conditions are **CEL expressions**. All match conditions must evaluate to `true` for the request to be evaluated.

## Example

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  matchConditions:
  - name: has-header
    expression: has(object.attributes.request.http.headers["x-force-deny"])
  deny:
  - response: >
      envoy.Denied(403).Response()
```

In the policy above, the `matchConditions` will be used to deny all requests having the `x-force-deny` header.

- If an incoming request doesn't have the `x-force-deny` header the condition will return `false` and the policy won't apply
- If an incoming request has the `x-force-deny` header the condition will return `true` and the `deny` rule will deny the request with status code `403`