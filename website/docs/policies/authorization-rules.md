# Validating Policies

A `ValidatingPolicy` defines validation rules to authorize HTTP requests.

Each validation rule is a CEL expression that evaluates to either an `http.Response` object or `null`. All expressions are written in [CEL](https://github.com/google/cel-spec).

When a validation expression returns an `http.Response`, that response is immediately returned to the caller. If the expression returns `null`, evaluation continues to the next validation in the array.

The CEL engine has been extended with libraries to make working with HTTP requests and responses easier. Browse the [available libraries documentation](../cel-extensions/index.md) for details.

## Evaluation order

Validations are evaluated in the order they appear in the `validations` array:

1. Each validation expression is evaluated sequentially
2. If a validation returns an `http.Response` (non-null), that response is returned immediately and evaluation stops
3. If a validation returns `null`, evaluation continues to the next validation
4. If all validations return `null`, the request is denied by default

!!!warning

    When multiple policies match a request (based on `matchConditions`), a random policy from the matches will be selected and its response will be used. It is best practice to define strict match conditions for each policy to avoid conflicts.

## Validation rules

The policy below will allow requests if they contain the header `x-force-authorized` with the value `enabled` or `true`.
If the header is not present or has a different value, the request will be denied.

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  evaluation:
    mode: HTTP
  variables:
  - name: force_authorized
    expression: object.headers.get("x-force-authorized")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
  # deny the request with 403 if not allowed
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
  # allow the request
  - expression: |
      http.response().status(200)
```

In this policy:

- The first validation checks if the request is not allowed and returns a `403` response if true, or `null` to continue
- The second validation returns a `200` response to allow the request
- If the first validation returns `null`, the second validation is evaluated

You can customize responses with status codes, headers, and body content using the [http library](../cel-extensions/http.md).

### Advanced example

This policy showcases a more advanced example with multiple validation checks.

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  evaluation:
    mode: HTTP
  variables:
  - name: force_authorized
    expression: object.headers.get("x-force-authorized") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.headers.get("x-force-unauthenticated") in ["enabled", "true"]
  validations:
  # if force_unauthenticated -> 401
  - expression: |
      variables.force_unauthenticated
        ? http.response().status(401).withBody("Authentication Failed")
        : null
  # if not force_authorized -> 403
  - expression: |
      !variables.force_authorized
        ? http.response().status(403).withBody("Unauthorized Request")
        : null
  # else -> 200 with custom headers
  - expression: |
      http.response()
        .status(200)
        .withHeader("x-validated-by", "kyverno")
        .withHeader("x-custom-header", "custom-value")
```

This policy demonstrates:

- **Sequential evaluation**: Validations are checked in order
- **Conditional responses**: Using ternary operators to return responses or null
- **Custom headers**: Adding headers to successful responses
- **Multiple status codes**: Different responses for authentication (401) vs authorization (403) failures

### Using external data

You can fetch data from external sources to make authorization decisions:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: external-data-policy
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: secretWord
    expression: |
      http.Get("http://my-server:3000").secretWord
  validations:
  - expression: |
      object.headers.get("secret-header") == variables.secretWord
        ? http.response().status(200).withBody("Valid secret")
        : http.response().status(403).withBody("Invalid secret")
```

!!!info

    The full documentation of the CEL HTTP library is available [here](../cel-extensions/http.md).