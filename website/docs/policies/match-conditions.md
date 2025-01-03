# Match conditions

You can define match conditions if you need fine-grained request filtering.

Match conditions are **CEL expressions**. All match conditions must evaluate to `true` for the request to be evaluated.

!!!info

    Match conditions have access to the same CEL variables as validation expressions.

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

- If an incoming request doesn't have the `x-force-deny` header, then the condition will return `false` and the policy won't apply
- If an incoming request has the `x-force-deny` header, then the condition will return `true` and the `deny` rule will deny the request with status code `403`

## Error handling

In the event of an error evaluating a match condition the policy is not evaluated. Whether to reject the request is determined as follows:

1. If any match condition evaluated to `false` (regardless of other errors), the policy is skipped.
1. Otherwise:
    - for `failurePolicy: Fail`, reject the request (without evaluating the policy).
    - for `failurePolicy: Ignore`, proceed with the request but skip the policy.
