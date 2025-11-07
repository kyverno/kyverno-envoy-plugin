# HTTP library

The `http` library provides types and functions for working with HTTP requests and responses in CEL expressions. It enables policies to inspect incoming HTTP requests and construct authorization responses.

## Types

### `http.CheckRequest`

Represents the top-level HTTP check request object.

| Field | CEL Type | Description |
|---|---|---|
| `attributes` | `http.CheckRequestAttributes` | Request attributes containing all HTTP request details |

**Example:**
```cel
object.attributes.method == "POST"
```

### `http.CheckRequestAttributes`

Contains all the attributes of an HTTP request.

| Field | CEL Type | Description |
|---|---|---|
| `method` | `string` | HTTP method (GET, POST, etc.) |
| `header` | `map<string, list<string>>` | Request headers (multi-value map) |
| `host` | `string` | Host header value |
| `protocol` | `string` | HTTP protocol version (HTTP/1.1, HTTP/2, etc.) |
| `contentLength` | `int` | Content length in bytes |
| `body` | `bytes` | Request body as raw bytes |
| `scheme` | `string` | URL scheme (http, https) |
| `path` | `string` | URL path |
| `query` | `map<string, list<string>>` | Query parameters (multi-value map) |
| `fragment` | `string` | URL fragment |

**Example:**
```cel
object.attributes.method == "POST" && object.attributes.path.startsWith("/api")
```

### `http.CheckResponseOk`

Represents an allowed/approved response (empty struct).

### `http.CheckResponseDenied`

Represents a denied response with a reason.

| Field | CEL Type | Description |
|---|---|---|
| `reason` | `string` | Reason for denial |

### `http.CheckResponse`

The final response object that contains either an OK or Denied response.

| Field | CEL Type | Description |
|---|---|---|
| `ok` | `http.CheckResponseOk` | Set if request is allowed |
| `denied` | `http.CheckResponseDenied` | Set if request is denied |

## Functions

### http.Allowed()

Creates an allowed response (CheckResponseOk).

**Signature:**
```cel
http.Allowed() -> http.CheckResponseOk
```

**Example:**
```cel
http.Allowed()
```

### http.Denied()

Creates a denied response with a reason string.

**Signature:**
```cel
http.Denied(string) -> http.CheckResponseDenied
```

**Example:**
```cel
http.Denied("Access denied: insufficient permissions")
http.Denied("Invalid authentication token")
```

### Header()

Gets all values for a specific header from the request attributes. Returns a list of strings.

**Signature:**
```cel
http.CheckRequestAttributes.Header(string) -> list<string>
```

**Example:**
```cel
object.attributes.Header("authorization")
object.attributes.Header("content-type")
```

### QueryParam()

Gets all values for a specific query parameter from the request attributes. Returns a list of strings.

**Signature:**
```cel
http.CheckRequestAttributes.QueryParam(string) -> list<string>
```

**Example:**
```cel
object.attributes.QueryParam("token")
object.attributes.QueryParam("api_key")
```

### Response()

Converts a CheckResponseOk or CheckResponseDenied into a final CheckResponse.

**Signature:**
```cel
http.CheckResponseOk.Response() -> http.CheckResponse
http.CheckResponseDenied.Response() -> http.CheckResponse
```

**Example:**
```cel
http.Allowed().Response()
http.Denied("Forbidden").Response()
```

## Complete Examples

### Allow all requests

```cel
http.Allowed().Response()
```

### Deny request with reason

```cel
http.Denied("Access denied: insufficient permissions").Response()
```

### Check authorization header

```cel
size(object.attributes.Header("authorization")) > 0
  ? http.Allowed().Response()
  : http.Denied("Missing authorization header").Response()
```

### Validate HTTP method

```cel
object.attributes.method == "GET" || object.attributes.method == "POST"
  ? http.Allowed().Response()
  : http.Denied("Method not allowed").Response()
```

### Check path prefix

```cel
object.attributes.path.startsWith("/api/v1")
  ? http.Allowed().Response()
  : http.Denied("Invalid API path").Response()
```

### Validate query parameter

```cel
size(object.attributes.QueryParam("api_key")) > 0
  ? http.Allowed().Response()
  : http.Denied("Missing api_key parameter").Response()
```

### Check header value

```cel
"application/json" in object.attributes.Header("content-type")
  ? http.Allowed().Response()
  : http.Denied("Invalid content type").Response()
```

### Complex authorization logic

```cel
object.attributes.method == "POST" && 
object.attributes.path.startsWith("/api/admin") &&
size(object.attributes.Header("x-admin-token")) > 0
  ? http.Allowed().Response()
  : http.Denied("Admin access required").Response()
```
