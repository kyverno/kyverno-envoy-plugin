
# Meta

[meta]: #meta

- Name: HTTP authorization mode for ValidatingPolicy

- Start Date: 2025-09-01

- Author(s): aerosouund

  

# Table of Contents

[table-of-contents]: #table-of-contents

- [Meta](#meta)

- [Table of Contents](#table-of-contents)

- [Overview](#overview)

- [Motivation](#motivation)

- [Proposal](#proposal)

- [Implementation](#implementation)

- [Open questions](#open-questions)

- [CRD Changes (OPTIONAL)](#crd-changes-optional)

  

# Overview

[overview]: #overview

Administrators and operations engineers run a lot into the question of how to provide authorization to their HTTP endpoints. Using the power of CEL, and the flexibility of the ValidatingPolicy resource, it's possible to turn those policies into a spec for authorization into an HTTP endpoint. This KDP aims to propose a library that will be integrated into the envoy plugin, giving it the aforementioned capabilities
  

# Motivation

[motivation]: #motivation


- Why should we do this?

Expanding Kyverno capabilities to new territory like HTTP authorization is a mark of its flexibility and its ability to solve many challenges in the authorization space. Proving it capable to solve challenges for things away from Kubernetes.

- What use cases does it support?

Managing HTTP authorization using CEL.

- What is the expected outcome?

An evaluation system, capable of handling policies that use the HTTP library, and produce HTTP responses from them.

# Proposal

  
HTTP authorization policies are a flavor of `ValidatingPolicies`, with their evaluation mode set to `HTTP`. 
Rules are applied serially in the policy they appear in. Ordering and result consistency is not guaranteed across multiple policies.

The system takes as input a golang `http.Request` object (from the standard library), produces a `http.Response`, and wires it back to the requester. 

An `http.request` CEL object will be introduced in the environment. Which will have the following fields with their respective types

```golang
type Request struct {
	Method   string `cel:"method"`
	Headers  KV     `cel:"headers"`
	Path     string `cel:"path"`
	Host     string `cel:"host"`
	Scheme   string `cel:"scheme"`
	Query    KV     `cel:"queryParams"`
	Fragment string `cel:"fragment"`
	Size     int64  `cel:"size"`
	Protocol string `cel:"protocol"`
	Body     string `cel:"body"`
	RawBody  []byte `cel:"rawBody"`
}
```

The return type of the entire expression will be `http.response()`, structured as such:

``` golang
type Response struct {
	StatusCode  int    `cel:"statusCode"`
	Status      string `cel:"status"`
	Headers     KV     `cel:"headers"`
	Body        string `cel:"body"`
}
```

## Library spec


### Headers:

| Function                  | Arguments                  | Return |
|---------------------------|---------------------------|--------|
| `http.request.headers.get()`   | `key` (`string`) – get the first value for header with key | `string` |
| `http.request.headers.getAll()`   | `key` (`string`) get all values for a header that was passed multiple times | `[]string` |


### Query parameters:

| Function                  | Arguments                  | Return |
|---------------------------|---------------------------|--------|
| `http.request.queryParams.get()`   | `key` (`string`) – get the first value for a set parameter| `string` |
| `http.request.queryParams.getAll()`   | `key` (`string`) get all values for a parameter that was passed multiple times | `[]string` |


### Response type

| Function                  | Arguments                  | Return |
|---------------------------|---------------------------|--------|
| `http.response().withHeader()`   | `key` (`string`), `value` (`string`): header key and value | `response` |
| `http.response().status()`   | `status` (`int`): integer value representing the status | `response` |
| `http.response().withBody()`   | `body` (`string`): add http response body | `response` |


## Example policies:

- Allow requests that contain a certain header


```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: require-foo-header
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") != "" 
        ? http.response().status(200)
        : http.response().status(400).withBody("header 'foo' is required")
```

- Allow requests where a certain header is a certain value

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: validate-foo-header-value
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") == "bar" 
        ? http.response().status(200)
        : http.response().status(400).withBody("header 'foo' must have value 'bar'")
```

- Allow requests where a certain header is a certain value to a particular path

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: validate-users-endpoint-access
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") == "bar" && http.request.path == "/v1/users"
        ? http.response().status(200)
        : http.response().status(400).withBody("header 'foo' must have value 'bar' when calling /v1/users")
```

- Deny requests where a certain header is a certain value to a particular path regex

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: validate-post-users-regex
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") == "bar" && http.request.path.matches("/*/users") && http.request.method == "POST"
        ? http.response().status(400)
        : http.response().status(200)
```

- Allow requests where a certain header is a certain value to a particular path prefix

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: validate-post-users-prefix
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") == "bar" && http.request.path.startsWith("/users") && http.request.method == "POST"
        ? http.response().status(200)
        : http.response().status(400)
```

- Allow requests where header and query parameter match specific values

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: validate-header-and-query-params
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
      http.request.headers.get("foo") == "bar" && http.request.queryParams.get("something") == "someone"
        ? http.response().status(200)
        : http.response().status(400)
```

- Allow requests where a value exists in a header that was passed multiple times

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: allow-multi-header-value
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
       "bar" in http.request.headers.getAll("foo")
        ? http.response().status(200)
        : http.response().status(400)
```

- Deny/Allow rule chaining

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: deny-allow-rule-chaining
spec:
  evaluation:
    mode: HTTP
  validations:
  - expression: >
       "undesiredHeaderVal" in http.request.headers.getAll("foo")
        ? http.response().status(400)
        : null
  - expression: >
       "undesiredParamVal" in http.request.queryParams.getAll("foo")
        ? http.response().status(400)
        : null
  - expression: >
      http.request.headers.get("users") == "allowedUser"
        && http.response().status(200)
```

## Advanced Real-World Use Cases:


- Token-based authorization (using the JWT library)

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: jwt-token-authorization
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: jwks
    expression: "https://myidp.com/.well-known/jwks.json"
  - name: authorization
    expression: http.request.headers.get("authorization").split(" ")
  - name: token
    expression: >
      size(variables.authorization) == 2 &&
      variables.authorization[0].lowerAscii() == "bearer"
        ? jwt.Decode(variables.authorization[1], variables.jwks)
        : null
  validations:
  - expression: >
      variables.token.Claims["myidp:groups"] in ["devops", "backend"]
        ? http.response().status(200)
        : http.response().status(403).withBody("Insufficient permissions")
```

- Request body size validation using rawBody

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: request-body-size-limit
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: bodySize
    expression: size(http.request.rawBody)
  - name: maxSizeBytes
    expression: 1048576  # 1MB limit
  validations:
  - expression: >
      variables.bodySize > variables.maxSizeBytes
        ? http.response().status(413).withBody("Request body too large: " + string(variables.bodySize) + " bytes (max: " + string(variables.maxSizeBytes) + ")")
        : http.response().status(200)
```

- File upload restrictions with security headers

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: secure-file-upload
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: contentType
    expression: http.request.headers.get("content-type")
  - name: contentLength
    expression: int(http.request.headers.get("content-length"))
  - name: allowedTypes
    expression: ["image/jpeg", "image/png", "application/pdf"]
  validations:
  - expression: >
      http.request.method == "POST" && http.request.path.startsWith("/upload")
        && variables.contentLength > 10485760  # 10MB limit
        ? http.response().status(413).withBody("File too large")
        : null
  - expression: >
      http.request.method == "POST" && http.request.path.startsWith("/upload")
        && !(variables.contentType in variables.allowedTypes)
        ? http.response().status(415).withBody("Unsupported file type")
        : null
  - expression: >
      http.request.method == "POST" && http.request.path.startsWith("/upload")
        ? http.response().status(200)
            .withHeader("x-content-type-options", "nosniff")
            .withHeader("x-frame-options", "DENY")
            .withHeader("strict-transport-security", "max-age=31536000")
        : http.response().status(200)
```

- Block bots and particular IPs

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: geo-ip-restrictions
spec:
  evaluation:
    mode: HTTP
  variables:
  - name: clientIp
    expression: http.request.headers.get("x-forwarded-for").split(",")[0]
  - name: userAgent
    expression: http.request.headers.get("user-agent")
  - name: blockedIps
    expression: ["192.168.1.100", "10.0.0.50"]
  validations:
  - expression: >
      variables.clientIp in variables.blockedIps
        ? http.response().status(403).withBody("IP address blocked")
        : null
  - expression: >
      variables.userAgent.contains("bot") || variables.userAgent.contains("crawler")
        ? http.response().status(403).withBody("Automated requests not allowed")
        : null
  - expression: > # return 200 otherwise
      http.response().status(200)
```


# Implementation

- Add a new CEL HTTP library

- Expose an HTTP endpoint that will receive requests and evaluate them through the existing policies

- Declare `http.request` as a variable in the CEL environment

- Declare a new global overload: `http.response()`, to instantiate a response object

- Declare member overloads on the response type: `withBody`, `status`, `withHeader`


# Open questions

- What should the default code be if the user specifies no `status`?

- What's the most appropriate protocol for serving this? gRPC will make it quicker, and easier to develop. But REST will be more easy to adopt for users because they will be applying no transformations before sending requests

- Should we allow users to specify both the status (the status string in the response) and the status code (the number representing status)? or do we aim for minimalism by only exposing `status`?

- Should we integrate this into the envoy plugin? or fork it?

- If all rules return null, what final HTTP response do we send back to the user ?

- Should we give a different name to the package to avoid clashing with the existing CEL HTTP library? Granted that we aren't integrating this library into the environment

- Do we make clients include requester information? 

# Additional potential features
  
- Geo-IP lookup
- Caching
- Rate limiting

# CRD Changes

New evaluation mode `HTTP` for `ValidatingPolicy`