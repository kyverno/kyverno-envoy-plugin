# HTTP Policy Breakdown

This guide provides a comprehensive breakdown of how to write `ValidatingPolicy` resources for plain HTTP authorization.

## Overview

When using the Kyverno Authz Server as an HTTP authorization service, policies analyze HTTP requests and return HTTP responses with authorization decisions. This mode allows you to protect any HTTP service without requiring Envoy.

## Policy Structure

A Kyverno `ValidatingPolicy` for HTTP consists of:

1. **Evaluation Mode**: Must be set to `HTTP`
2. **Failure Policy**: How to handle policy evaluation failures
3. **Match Conditions** (optional): Fine-grained request filtering
4. **Variables** (optional): Reusable expressions
5. **Validation Rules**: Authorization logic

## Evaluation Mode

For HTTP authorization, the evaluation mode **must** be set to `HTTP`:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: http-policy
spec:
  evaluation:
    mode: HTTP  # Required for HTTP authorization
  validations:
  - expression: ...
```

## Failure Policy

The `failurePolicy` defines how to handle failures during policy evaluation (parse errors, type check errors, runtime errors).

Allowed values:
- `Fail` (default): Deny the request if policy evaluation fails
- `Ignore`: Allow the request if policy evaluation fails

### Example: Fail Policy

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail  # Deny on failure
  evaluation:
    mode: HTTP
  variables:
  - name: force_authorized
    expression: object.headers.get("x-force-authorized")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
```

### Example: Ignore Policy

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  failurePolicy: Ignore  # Allow on failure
  evaluation:
    mode: HTTP
  variables:
  - name: force_authorized
    expression: object.headers.get("x-force-authorized")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
```

## Match Conditions

Match conditions provide fine-grained request filtering using CEL expressions. All match conditions must evaluate to `true` for the policy to apply.

!!!info
    Variables are NOT available in match conditions because they are evaluated before the rest of the policy.

### Example: Path-Based Matching

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: api-protection
spec:
  failurePolicy: Fail
  evaluation:
    mode: HTTP
  matchConditions:
  - name: is-api-path
    expression: object.path.startsWith("/api/")
  validations:
  - expression: |
      http.response().status(403).withBody("API access denied")
```

### Example: Header-Based Matching

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: admin-only
spec:
  evaluation:
    mode: HTTP
  matchConditions:
  - name: has-admin-header
    expression: object.headers.get("x-user-role") == "admin"
  validations:
  - expression: |
      http.response().status(200).withBody("Admin access granted")
```

### Error Handling

If a match condition evaluation fails:
1. If any match condition evaluated to `false`, the policy is skipped
2. Otherwise:
   - For `failurePolicy: Fail`: Reject the request
   - For `failurePolicy: Ignore`: Skip the policy and allow the request

## Variables

Variables are named CEL expressions that can be reused throughout the policy. They are available under the `variables` identifier.

!!!info
    The incoming HTTP request is available under the `object` identifier.

### Example: Using Variables

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
    # Extract header value
  - name: force_authorized
    expression: object.headers.get("x-force-authorized")
    # Compute authorization status
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
  - expression: |
      http.response().status(200)
```

**Important**: Variables must be sorted by order of first appearance. A variable can reference earlier variables but not later ones.

## Validation Rules

Validation rules contain the authorization logic. Each rule is a CEL expression that returns either an `http.Response` object or `null`.

### Evaluation Order

1. Rules are evaluated sequentially in the order they appear
2. If a rule returns an `http.Response` (non-null), that response is returned immediately
3. If a rule returns `null`, evaluation continues to the next rule
4. If all rules return `null`, the request is denied by default

!!!warning
    When multiple policies match a request, a random policy will be selected. Use strict match conditions to avoid conflicts.

### Basic Example

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
  # Deny if not allowed
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
  # Allow the request
  - expression: |
      http.response().status(200)
```

### Advanced Example

This example demonstrates multiple validation checks with custom headers and status codes:

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
  # Check 1: Return 401 if unauthenticated
  - expression: |
      variables.force_unauthenticated
        ? http.response().status(401).withBody("Authentication Failed")
        : null
  # Check 2: Return 403 if not authorized
  - expression: |
      !variables.force_authorized
        ? http.response().status(403).withBody("Unauthorized Request")
        : null
  # Check 3: Allow with custom headers
  - expression: |
      http.response()
        .status(200)
        .withHeader("x-validated-by", "kyverno")
        .withHeader("x-custom-header", "custom-value")
```

This policy demonstrates:
- **Sequential evaluation**: Each validation is checked in order
- **Conditional responses**: Using ternary operators to return responses or null
- **Multiple status codes**: Different responses for authentication (401) vs authorization (403) failures
- **Custom headers**: Adding headers to successful responses
- **Custom body**: Setting response body content

## CEL HTTP Extension Library

The CEL engine includes helper functions for creating HTTP responses:

### Key Functions

- **`http.response()`**: Creates a new HTTP response builder
- **`.status(code)`**: Sets the HTTP status code
- **`.withBody(content)`**: Sets the response body
- **`.withHeader(key, value)`**: Adds a response header

### Response Builder Pattern

The HTTP library uses a builder pattern for constructing responses:

```yaml
validations:
- expression: |
    http.response()
      .status(200)
      .withHeader("Content-Type", "application/json")
      .withHeader("X-Custom-Header", "value")
      .withBody('{"status": "ok"}')
```

## Using External Data

You can fetch data from external sources to make authorization decisions:

### Example: External API Call

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

### Example: Multiple External Calls

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: multi-source-policy
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: userInfo
    expression: |
      http.Get("http://user-service:8080/users/" + object.headers.get("x-user-id"))
  - name: permissions
    expression: |
      http.Get("http://auth-service:8080/permissions/" + variables.userInfo.role)
  validations:
  - expression: |
      variables.permissions.canAccess
        ? http.response().status(200).withHeader("x-user-role", variables.userInfo.role)
        : http.response().status(403).withBody("Insufficient permissions")
```

## Complete Example

Here's a complete policy that combines all concepts:

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: complete-http-policy
spec:
  failurePolicy: Fail
  evaluation:
    mode: HTTP
  matchConditions:
  - name: is-api-path
    expression: object.path.startsWith("/api/")
  - name: is-post-or-put
    expression: object.method in ["POST", "PUT"]
  variables:
  - name: auth_header
    expression: object.headers.get("authorization")
  - name: is_authenticated
    expression: variables.auth_header.startsWith("Bearer ")
  - name: user_role
    expression: object.headers.get("x-user-role")
  - name: is_admin
    expression: variables.user_role == "admin"
  - name: request_path
    expression: object.path
  validations:
  # Deny if not authenticated
  - expression: |
      !variables.is_authenticated
        ? http.response()
            .status(401)
            .withHeader("WWW-Authenticate", "Bearer")
            .withBody("Authentication required")
        : null
  # Deny if not admin
  - expression: |
      !variables.is_admin
        ? http.response()
            .status(403)
            .withBody("Admin access required for " + variables.request_path)
        : null
  # Allow with tracking headers
  - expression: |
      http.response()
        .status(200)
        .withHeader("x-auth-validated", "true")
        .withHeader("x-policy-applied", "complete-http-policy")
        .withHeader("x-user-role", variables.user_role)
```

## Request Object Structure

The HTTP request object (`object`) has the following structure:

```javascript
{
  method: "GET",              // HTTP method
  path: "/api/users",         // Request path
  headers: {                  // Request headers (map)
    "authorization": "Bearer token",
    "content-type": "application/json"
  },
  body: "...",               // Request body (if present)
  query: {                   // Query parameters (map)
    "page": "1",
    "limit": "10"
  }
}
```

### Accessing Request Data

```yaml
variables:
# Access method
- name: method
  expression: object.method

# Access path
- name: path
  expression: object.path

# Access headers
- name: auth_header
  expression: object.headers.get("authorization")

# Access query parameters
- name: page
  expression: object.query.get("page")

# Check if header exists
- name: has_auth
  expression: object.headers.has("authorization")
```

## Common Patterns

### Pattern 1: Token Validation

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: token-validation
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: token
    expression: object.headers.get("authorization").replace("Bearer ", "")
  - name: is_valid
    expression: variables.token.size() > 20
  validations:
  - expression: |
      !variables.is_valid
        ? http.response().status(401).withBody("Invalid token")
        : http.response().status(200)
```

### Pattern 2: Role-Based Access Control

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: rbac-policy
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: user_role
    expression: object.headers.get("x-user-role")
  - name: allowed_roles
    expression: ["admin", "editor"]
  validations:
  - expression: |
      !(variables.user_role in variables.allowed_roles)
        ? http.response().status(403).withBody("Role not authorized")
        : http.response().status(200)
```

### Pattern 3: Rate Limiting Check

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: rate-limit-check
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: user_id
    expression: object.headers.get("x-user-id")
  - name: rate_limit_status
    expression: |
      http.Get("http://rate-limiter:8080/check/" + variables.user_id)
  validations:
  - expression: |
      variables.rate_limit_status.exceeded
        ? http.response()
            .status(429)
            .withHeader("Retry-After", "60")
            .withBody("Rate limit exceeded")
        : http.response().status(200)
```

### Pattern 4: Path-Based Authorization

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: path-authorization
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: is_admin_path
    expression: object.path.startsWith("/admin/")
  - name: is_admin_user
    expression: object.headers.get("x-user-role") == "admin"
  validations:
  - expression: |
      variables.is_admin_path && !variables.is_admin_user
        ? http.response().status(403).withBody("Admin access required")
        : null
  - expression: |
      http.response().status(200)
```

## Best Practices

1. **Use specific match conditions** to avoid policy conflicts when multiple policies exist
2. **Order validations carefully** - put most common deny conditions first
3. **Use variables** to avoid repeating complex expressions
4. **Return `null`** from validations that don't make a decision
5. **Set appropriate failure policies** based on your security requirements
6. **Use descriptive policy names** for easier troubleshooting
7. **Add custom headers** for observability and debugging
8. **Cache external data** when possible to improve performance
9. **Use meaningful HTTP status codes** (401 for authentication, 403 for authorization)
10. **Provide clear error messages** in response bodies

## Debugging Tips

1. **Add debug headers** to responses to track which policy and validation matched:
   ```yaml
   .withHeader("x-policy-name", "my-policy")
   .withHeader("x-validation-step", "step-2")
   ```

2. **Use variables to break down complex logic** for easier debugging

3. **Test match conditions separately** to ensure they work as expected

4. **Check failure policy behavior** in development before deploying to production

## Additional Resources

- [CEL HTTP Extension Library](../cel-extensions/http.md)
- [CEL Language Specification](https://github.com/google/cel-spec)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
