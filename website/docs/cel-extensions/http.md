# HTTP library

The `http` library provides types and functions for working with HTTP requests and responses in CEL expressions. It enables policies to inspect incoming HTTP requests and construct custom HTTP responses.

## Types

### `http.Request`

Represents an HTTP request with all its attributes.

| Field | CEL Type | Description |
|---|---|---|
| `method` | `string` | HTTP method (GET, POST, etc.) |
| `headers` | `http.KV` | Request headers |
| `path` | `string` | URL path |
| `host` | `string` | Host header value |
| `scheme` | `string` | URL scheme (http, https) |
| `queryParams` | `http.KV` | Query parameters |
| `fragment` | `string` | URL fragment |
| `size` | `int` | Request body size in bytes |
| `protocol` | `string` | HTTP protocol version (HTTP/1.1, HTTP/2) |
| `body` | `string` | Request body as string |
| `rawBody` | `bytes` | Request body as raw bytes |

**Example:**
```cel
object.method == "POST" && object.path.startsWith("/api")
```

### `http.KV`

Represents a key-value map for headers and query parameters. Supports multiple values per key.

**Methods:**
- `get(string) -> string`: Get the first value for a header/parameter
- `getAll(string) -> list<string>`: Get all values for a header/parameter

**Example:**
```cel
object.headers.get("content-type") == "application/json"
```

### `http.Response`

Represents an HTTP response that can be returned from a policy.

| Field | CEL Type | Description |
|---|---|---|
| `status` | `int` | HTTP status code |
| `headers` | `http.KV` | Response headers |
| `body` | `string` | Response body |

**Methods:**
- `status(int) -> http.Response`: Set the HTTP status code
- `withHeader(string, string) -> http.Response`: Add a response header
- `withBody(string) -> http.Response`: Set the response body

**Example:**
```cel
http.response().status(403).withBody("Access denied")
```

## Functions

### http.response()

Creates a new `http.Response` object that can be customized with status, headers, and body.

**Signature:**
```cel
http.response() -> http.Response
```

**Example:**
```cel
http.response().status(200).withBody("Success")
```

### get()

Gets the first value of a header or query parameter from an `http.KV` object. Returns an empty string if the key doesn't exist.

**Signature:**
```cel
http.KV.get(string) -> string
```

**Example:**
```cel
object.headers.get("authorization")
object.queryParams.get("token")
```

### getAll()

Gets all values of a header or query parameter from an `http.KV` object. Returns an empty list if the key doesn't exist.

**Signature:**
```cel
http.KV.getAll(string) -> list<string>
```

**Example:**
```cel
object.headers.getAll("accept")
```

### status()

Sets the HTTP status code for an `http.Response` object.

**Signature:**
```cel
http.Response.status(int) -> http.Response
```

**Example:**
```cel
http.response().status(403)
http.response().status(200)
http.response().status(401)
```

### withHeader()

Adds a header to an `http.Response` object. Can be called multiple times to add multiple headers.

**Signature:**
```cel
http.Response.withHeader(string, string) -> http.Response
```

**Example:**
```cel
http.response().status(200).withHeader("x-custom-header", "value")
http.response().status(403).withHeader("www-authenticate", "Bearer")
```

### withBody()

Sets the response body for an `http.Response` object.

**Signature:**
```cel
http.Response.withBody(string) -> http.Response
```

**Example:**
```cel
http.response().status(403).withBody("Access denied")
http.response().status(200).withBody("Request approved")
```

## Complete Examples

### Allow request with custom header

```cel
http.response().status(200).withHeader("x-validated-by", "kyverno")
```

### Deny request with custom status and body

```cel
http.response().status(403).withBody("Insufficient permissions")
```

### Check authorization header

```cel
object.headers.get("authorization").startsWith("Bearer ")
  ? http.response().status(200)
  : http.response().status(401).withBody("Missing authorization header")
```

### Validate content type

```cel
object.headers.get("content-type") == "application/json"
  ? http.response().status(200)
  : http.response().status(415).withBody("Unsupported media type")
```
