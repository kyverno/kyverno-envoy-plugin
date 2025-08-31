
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

- [Definitions](#definitions)

- [Motivation](#motivation)

- [Proposal](#proposal)

- [Implementation](#implementation)

- [Migration (OPTIONAL)](#migration-optional)

- [Drawbacks](#drawbacks)

- [Alternatives](#alternatives)

- [Prior Art](#prior-art)

- [Unresolved Questions](#unresolved-questions)

- [CRD Changes (OPTIONAL)](#crd-changes-optional)

  

# Overview

[overview]: #overview

Administrators and operations engineers run alot into the question of how to provide authentication to their HTTP endpoints. Using the power of CEL, and the flexibility of the ValidatingPolicy resource, its possible to turn those policies into a spec for authorization into a HTTP endpoint. This KDP aims to propose a library that will be intergrated into the envoy plugin, giving it the aformentioned capabilities
  

# Motivation

[motivation]: #motivation

  

- Why should we do this?

- What use cases does it support?

- What is the expected outcome?

  

# Proposal

  
HTTP authorization policies are a flavor of `ValidatingPolicies`, with their evaluation mode set to `HTTP`. 
An `http.request` CEL object will be introduced in the environment. Which will have the following fields with their respective types

```golang
type request struct {
    // todo: unexport fields
	method   string              `cel:"method"`
	headers  map[string][]string `cel:"headers"`
	path     string              `cel:"path"`
	host     string              `cel:"host"`
	scheme   string              `cel:"scheme"`
	query    map[string][]string `cel:"queryParams"`
	fragment string              `cel:"fragment"`
	size     int64               `cel:"size"`
	protocol string              `cel:"protocol"`
	body     string              `cel:"body"`
	rawBody  []byte              `cel:"rawBody"`
}
```

The return type of the entire expression will be `http.response()`, structured as such:

``` golang
type response struct {
    // todo: unexport fields
    statusCode int                 `cel:"statusCode"`
    headers    map[string][]string `cel:"headers"`
    body       []byte              `cel:"body"`   
    protocol   string              `cel:"protocol"`
    reason     string              `cel:"reason"`
}
```

## Library spec


### Headers:

| Function                  | Arguments                  | Return |
|---------------------------|---------------------------|--------|
| `http.request.headers.get()`   | `key` (`string`) – get the first value for header with key | `[]string` |
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
        ? http.response().status(400).withBody("header 'foo' is required")
        : http.response().status(200)
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
        ? http.response().status(400).withBody("header 'foo' must have value 'bar'")
        : http.response().status(200)
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
        ? http.response().status(400).withBody("header 'foo' must have value 'bar' when calling /v1/users")
        : http.response().status(200)
```

- Allow requests where a certain header is a certain value to a particular path regex

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
        ? http.response().status(400)
        : http.response().status(200)
```

- Allow requests where a certain header is a certain value to a particular path prefix

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
        ? http.response().status(400)
        : http.response().status(200)
```

- Token based authorization (using the JWT library)

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
        ? http.response().status(403).withHeader("foo", "bar")
        : http.response().status(200)
```

  

# Implementation

- Add a new CEL HTTP library

- Expose a HTTP endpoint that will receive requests and evalute them through the existing policies

- Declare `http.request` as a variable in the CEL environment
  

# CRD Changes (OPTIONAL)

  
New evaluation mode `HTTP` for `ValidatingPolicy`