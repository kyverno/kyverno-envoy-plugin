# Authorization rules

An `AuthorizationPolicy` main element is the authorization rules defined in `authorizations`.

Every authorization rule must contain a [CEL](https://github.com/google/cel-spec) `expression`. It is expected to return an Envoy `CheckResponse` describing the decision made by the rule (or nothing if no decision is made).

Creating the Envoy [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) can be a tedious task, you need to remember the different types names and format.

The CEL engine used to evaluate the authorization rules has been extended with a library to make the creation of `CheckResponse` easier. (TODO)

## Authorization rules

The policy below will allow requests if they contain the header `x-force-authorized` with the value `enabled` or `true`.
If the header is not present or has a different value, the request will be denied.

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  authorizations:
    # make an authorisation decision based on the value of `variables.allowed`
    # - allow the request if it is `true`
    # - deny the request with 403 status code if it is `false`
  - expression: >
      variables.allowed
        ? envoy.Allowed().Response()
        : envoy.Denied(403).Response()
```

In this simple rule:

- `envoy.Allowed().Response()`

    Creates a `CheckResponse` to allow the request

- `envoy.Denied(403).Response()`

    Creates a `CheckResponse` to deny the request with status code `403`

However, we can do a lot more with Envoy's `CheckResponse`.
Envoy can add or remove headers, query parameters, register dynamic metadata passed along the filters chain, and even change the response body. (TODO)

![dynamic metadata](../schemas/dynamic-metadata.png)

### Multiple rules

In the example above, we combined allow and denied response handling in a single expression.
However it is possible to use multiple expressions, the first one returning a non null response will be used by the Kyverno Authz Server:

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  authorizations:
    # allow the request if `variables.allowed` is `true`
    # or delegate the decision to the next rule
  - expression: >
      variables.allowed
        ? envoy.Allowed().Response()
        : null
    # deny the request with 403 status code
  - expression: >
      envoy.Denied(403).Response()
```

### The hard way

Below is the same policy, creating the `CheckResponses` manually:

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  authorizations:
  - expression: >
      variables.allowed
        ? envoy.service.auth.v3.CheckResponse{
            status: google.rpc.Status{
              code: 0
            },
            ok_response: envoy.service.auth.v3.OkHttpResponse{}
          }
        : envoy.service.auth.v3.CheckResponse{
            status: google.rpc.Status{
              code: 7
            },
            denied_response: envoy.service.auth.v3.DeniedHttpResponse{
              status: envoy.type.v3.HttpStatus{
                code: 403
              }
            }
          }
```

### Advanced example

This second policy showcases a more advanced example.

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.attributes.request.http.headers[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]
  authorizations:
  - expression: >
      variables.force_authorized && !variables.force_unauthenticated
      ? envoy
          .Allowed()
          .WithHeader("x-validated-by", "my-security-checkpoint")
          .WithoutHeader("x-force-authorized")
          .WithResponseHeader("x-add-custom-response-header", "added")
          .Response()
          .WithMetadata({"my-new-metadata": "my-new-value"})
      : envoy
          .Denied(variables.force_unauthenticated ? 401 : 403)
          .WithBody(variables.force_unauthenticated ? "Authentication Failed" : "Unauthorized Request")
          .Response()
EOF
```

Notice this policy uses helper functions:

- `WithHeader`

    To add a request header

- `WithoutHeader`

    To remove a request header

- `WithResponseHeader`

    To add a response header

- `WithBody`

    To modify the response body

- `WithMetadata`

    To add dynamic metadata in the envoy filter chain (this is useful when you want to pass data to another filter in the chain or you want to print it in the application logs)
