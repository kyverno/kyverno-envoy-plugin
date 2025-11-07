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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  validations:
  - expression: |
      !variables.force_authorized
        ? http.Denied("forbidden").Response()
        : http.Allowed().Response()
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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  validations:
  - expression: |
      !variables.force_authorized
        ? http.Denied("Forbidden").Response()
        : http.Allowed().Response()
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
    expression: object.attributes.path.startsWith("/api/")
  validations:
  - expression: |
      http.Denied("API access denied").Response()
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
    expression: object.attributes.header[?"x-user-role"].orValue("") == "admin"
  validations:
  - expression: |
      http.Allowed().Response()
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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  validations:
  - expression: |
      !variables.force_authorized
        ? http.Denied("Forbidden").Response()
        : null
  - expression: |
      http.Allowed().Response()
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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  validations:
  # Deny if not allowed
  - expression: |
      !variables.force_authorized
        ? http.Denied("Forbidden").Response()
        : null
  # Allow the request
  - expression: |
      http.Allowed().Response()
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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.attributes.header[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]
  validations:
  # Check 1: Deny if unauthenticated
  - expression: |
      variables.force_unauthenticated
        ? http.Denied("Authentication Failed").Response()
        : null
  # Check 2: Deny if not authorized
  - expression: |
      !variables.force_authorized
        ? http.Denied("Unauthorized Request").Response()
        : null
  # Check 3: Allow the request
  - expression: |
      http.Allowed().Response()
```

This policy demonstrates:
- **Sequential evaluation**: Each validation is checked in order
- **Conditional responses**: Using ternary operators to return responses or null
- **Different denial reasons**: Providing specific reasons for authentication vs authorization failures

## CEL HTTP Extension Library

The CEL engine includes helper functions for creating HTTP responses:

### Key Functions

- **`http.Allowed()`**: Creates an allowed response
- **`http.Denied(reason)`**: Creates a denied response with a reason string
- **`.Response()`**: Converts the response to the final CheckResponse type

### Response Pattern

The HTTP library provides simple functions for authorization decisions:

```yaml
validations:
# Allow the request
- expression: |
    http.Allowed().Response()

# Deny with a reason
- expression: |
    http.Denied("Access denied: insufficient permissions").Response()
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
      object.attributes.header[?"secret-header"].orValue("") == variables.secretWord
        ? http.Allowed().Response()
        : http.Denied("Invalid secret").Response()
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
  - name: userId
    expression: object.attributes.header[?"x-user-id"].orValue("")
  - name: userInfo
    expression: |
      http.Get("http://user-service:8080/users/" + variables.userId)
  - name: permissions
    expression: |
      http.Get("http://auth-service:8080/permissions/" + variables.userInfo.role)
  validations:
  - expression: |
      variables.permissions.canAccess
        ? http.Allowed().Response()
        : http.Denied("Insufficient permissions").Response()
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
    expression: object.attributes.path.startsWith("/api/")
  - name: is-post-or-put
    expression: object.attributes.method in ["POST", "PUT"]
  variables:
  - name: auth_header
    expression: object.attributes.header[?"authorization"].orValue("")
  - name: is_authenticated
    expression: variables.auth_header.startsWith("Bearer ")
  - name: user_role
    expression: object.attributes.header[?"x-user-role"].orValue("")
  - name: is_admin
    expression: variables.user_role == "admin"
  - name: request_path
    expression: object.attributes.path
  validations:
  # Deny if not authenticated
  - expression: |
      !variables.is_authenticated
        ? http.Denied("Authentication required").Response()
        : null
  # Deny if not admin
  - expression: |
      !variables.is_admin
        ? http.Denied("Admin access required for " + variables.request_path).Response()
        : null
  # Allow the request
  - expression: |
      http.Allowed().Response()
```

## Request Object Structure

The HTTP request object (`object`) has the following structure:

```javascript
{
  attributes: {
    method: "GET",                           // HTTP method
    path: "/api/users",                      // Request path
    header: {                                // Request headers (multi-value map)
      "authorization": ["Bearer token"],
      "content-type": ["application/json"]
    },
    body: bytes,                             // Request body as bytes
    query: {                                 // Query parameters (multi-value map)
      "page": ["1"],
      "limit": ["10"]
    },
    host: "example.com",                     // Host header
    protocol: "HTTP/1.1",                    // Protocol version
    contentLength: 123,                      // Content length
    scheme: "https",                         // URL scheme
    fragment: ""                             // URL fragment
  }
}
```

### Accessing Request Data

```yaml
variables:
# Access method
- name: method
  expression: object.attributes.method

# Access path
- name: path
  expression: object.attributes.path

# Access headers using map syntax with optional chaining
- name: auth_header
  expression: object.attributes.header[?"authorization"].orValue("")

# Access query parameters using map syntax
- name: page
  expression: object.attributes.query[?"page"].orValue("")

# Check if header exists
- name: has_auth
  expression: object.attributes.header[?"authorization"].hasValue()

# Get header value with default
- name: content_type
  expression: object.attributes.header[?"content-type"].orValue("application/octet-stream")
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
  - name: auth_header
    expression: object.attributes.header[?"authorization"].orValue("")
  - name: token
    expression: variables.auth_header.startsWith("Bearer ") ? variables.auth_header.replace("Bearer ", "") : ""
  - name: is_valid
    expression: variables.token.size() > 20
  validations:
  - expression: |
      !variables.is_valid
        ? http.Denied("Invalid token").Response()
        : http.Allowed().Response()
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
    expression: object.attributes.header[?"x-user-role"].orValue("")
  - name: allowed_roles
    expression: ["admin", "editor"]
  validations:
  - expression: |
      !(variables.user_role in variables.allowed_roles)
        ? http.Denied("Role not authorized").Response()
        : http.Allowed().Response()
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
    expression: object.attributes.header[?"x-user-id"].orValue("")
  - name: rate_limit_status
    expression: |
      http.Get("http://rate-limiter:8080/check/" + variables.user_id)
  validations:
  - expression: |
      variables.rate_limit_status.exceeded
        ? http.Denied("Rate limit exceeded").Response()
        : http.Allowed().Response()
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
    expression: object.attributes.path.startsWith("/admin/")
  - name: is_admin_user
    expression: object.attributes.header[?"x-user-role"].orValue("") == "admin"
  validations:
  - expression: |
      variables.is_admin_path && !variables.is_admin_user
        ? http.Denied("Admin access required").Response()
        : null
  - expression: |
      http.Allowed().Response()
```

## Best Practices

1. **Use specific match conditions** to avoid policy conflicts when multiple policies exist
2. **Order validations carefully** - put most common deny conditions first
3. **Use variables** to avoid repeating complex expressions
4. **Return `null`** from validations that don't make a decision
5. **Set appropriate failure policies** based on your security requirements
6. **Use descriptive policy names** for easier troubleshooting
7. **Cache external data** when possible to improve performance
8. **Provide clear denial reasons** in `http.Denied()` calls for better debugging
9. **Use optional chaining** with `.orValue("")` to safely access headers and query parameters
10. **Use `.hasValue()`** to check if a header or query parameter exists

## Debugging Tips

1. **Include detailed denial reasons** to track which validation failed:
   ```yaml
   http.Denied("Policy: my-policy, Step: authentication-check").Response()
   ```

2. **Use variables to break down complex logic** for easier debugging

3. **Test match conditions separately** to ensure they work as expected

4. **Check failure policy behavior** in development before deploying to production

5. **Use optional chaining syntax** - `object.attributes.header[?"key"].orValue("")` for safe access

## Additional Resources

- [CEL HTTP Extension Library](../cel-extensions/http.md)
- [CEL Language Specification](https://github.com/google/cel-spec)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
