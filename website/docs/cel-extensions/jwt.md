# Jwt library

Policies have native functionality to decode and verify the contents of JWT tokens in order to enforce additional authorization logic on requests.

## Functions

### jwt.Decode

The `jwt.Decode` function decodes and validates a JWT token. It accepts two arguments: the token and the secret to verify the signature.

#### Signature and overloads

```
jwt.Decode(<string> token, <string> key) -> <token>
```

#### Example

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Ignore
  variables:
  - name: token
    expression: >
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
  - name: secret
    expression: >
      "secret"
  authorizations:
  - expression: >
      jwt.Decode(variables.token, variables.secret)....
```
