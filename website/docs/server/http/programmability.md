# Programmability

The Kyverno Authz server provides programmability through CEL (Common Expression Language) expressions that allow you to transform request and response attributes dynamically.

## Overview

Two flags enable request and response transformation:

- **`--input-expression`**: Transforms incoming requests before authorization
- **`--output-expression`**: Transforms outgoing responses after authorization

Both flags accept CEL expressions that take a specific input type and evaluate to the same type.

## Input Expression

### Purpose

The input expression transforms incoming requests before they are processed by the authorization engine.

### Input Type

`http.CheckRequest`

### Use Cases

- Modify request headers
- Extract information from custom headers
- Transform request attributes
- Normalize request data

### Example

```cel
http.CheckRequest{
  attributes: http.CheckRequestAttributes{
    method: object.attributes.Header("x-original-method")[0],
    header: object.attributes.header,
    host: url(object.attributes.Header("x-original-url")[0]).getHostname(),
    scheme: url(object.attributes.Header("x-original-url")[0]).getScheme(),
    path: url(object.attributes.Header("x-original-url")[0]).getEscapedPath(),
    query: url(object.attributes.Header("x-original-url")[0]).getQuery(),
    body: object.attributes.body,
    fragment: "todo",
  }
}
```

This example demonstrates how to use the `url` library to extract information from headers and reconstruct the request attributes.

## Output Expression

### Purpose

The output expression transforms responses before they are sent back to the client.

### Input Type

`httpserver.HttpResponse`

### Use Cases

- Modify response status codes
- Add or modify response headers
- Customize response body
- Add authentication metadata

### Example

```cel
httpserver.HttpResponse{
  status: 401,
  body: bytes(object.denied.reason),
  header: {"authenticated-by": ["kyverno-authz-server"]}
}
```

### Modifiable Fields

- **`status`**: HTTP status code
- **`body`**: Response body (as bytes)
- **`header`**: Response headers (map of string arrays)
