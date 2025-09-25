# Authorization rules

An `AuthorizationPolicy` main concern is to define authorization rules to `deny` or `allow` requests.

Every authorization rule is made of an optional `match` statement and a required `response` statement. Both statements are written in [CEL](https://github.com/google/cel-spec).

If the `match` statement is present and evaluates to `true`, the `response` statement is used to create the response payload returned to the envoy proxy.
Depending on the rule type, the response is expected to be an envoy.OkResponse or envoy.DeniedResponse.

Creating an [OkResponse](../../cel-extensions/envoy.md#okresponse) or [DeniedResponse](../../cel-extensions/envoy.md#deniedresponse) can be a tedious task, you need to remember the different types names and format.

The CEL engine used to evaluate the authorization rules has been extended with a library to make the creation of responses easier. Browse the [available libraries documentation](../../cel-extensions/index.md) for details.

## Evaluation order

1. All `deny` rules are evaluated first, the first matching rule is used to send the deny response to the envoy proxy.
1. If no `deny` rule matched, `allow` rules are evaluated and the first matching rule is used to send the response to the envoy proxy.
1. If no rule matched, the request is allowed by default.

!!!info

    When multiple policies are present, `deny` and `allow` rules are concatenated together in policy name alphabetical order.

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
  # make an authorisation decision based on the value of `variables.allowed`
  # - deny the request with 403 status code if it is `false`
  # - else allow the request
  deny:
  - match: >
      !variables.allowed
    response: >
      envoy.Denied(403).Response()
  allow:
  - response: >
      envoy.Allowed().Response()
```

In this simple rule:

- `envoy.Allowed().Response()`

    Creates an `OkResponse` to allow the request

- `envoy.Denied(403).Response()`

    Creates a `DeniedResponse` to deny the request with status code `403`

However, we can do a lot more.
Envoy can add or remove headers, query parameters, register dynamic metadata passed along the filters chain, and even change the response body.

![dynamic metadata](../../schemas/dynamic-metadata.png)

### The hard way

Below is the same policy, creating the `envoy.OkResponse` and `envoy.DeniedResponse` manually.

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
  deny:
  - match: >
      !variables.allowed
    response: >
      envoy.DeniedResponse{
        status: google.rpc.Status{
          code: 7
        },
        http_response: envoy.service.auth.v3.DeniedHttpResponse{
          status: envoy.type.v3.HttpStatus{
            code: 403
          }
        }
      }
  allow:
  - response: >
      envoy.OkResponse{
        status: google.rpc.Status{
          code: 0
        },
        http_response: envoy.service.auth.v3.OkHttpResponse{}
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
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  - name: force_unauthenticated
    expression: object.attributes.request.http.headers[?"x-force-unauthenticated"].orValue("") in ["enabled", "true"]
  - name: metadata
    expression: '{"my-new-metadata": "my-new-value"}'
  deny:
    # if force_unauthenticated -> 401
  - match: >
      variables.force_unauthenticated
    response: >
      envoy
        .Denied(401)
        .WithBody("Authentication Failed")
        .Response()
    # if not force_authorized -> 403
  - match: >
      !variables.force_authorized
    response: >
      envoy
        .Denied(403)
        .WithBody("Unauthorized Request")
        .Response()
  allow:
    # else -> 200
  - response: >
      envoy
        .Allowed()
        .WithHeader("x-validated-by", "my-security-checkpoint")
        .WithoutHeader("x-force-authorized")
        .WithResponseHeader("x-add-custom-response-header", "added")
        .Response()
        .WithMetadata(variables.metadata)
```

Notice this policy uses helper functions:

- [envoy.Allowed](../../cel-extensions/envoy.md#envoyallowed)

    To create an OK http response

- [envoy.Denied](../../cel-extensions/envoy.md#envoydenied)

    To create a DENIED http response

- [Response](../../cel-extensions/envoy.md#response)

    To create a check response from an http response

- [WithHeader](../../cel-extensions/envoy.md#withheader)

    To add a request header

- [WithoutHeader](../../cel-extensions/envoy.md#withoutheader)

    To remove a request header

- [WithResponseHeader](../../cel-extensions/envoy.md#withresponseheader)

    To add a response header

- [WithBody](../../cel-extensions/envoy.md#withbody)

    To modify the response body

- [WithMetadata](../../cel-extensions/envoy.md#withmetadata)

    To add dynamic metadata in the envoy filter chain (this is useful when you want to pass data to another filter in the chain or you want to print it in the application logs)

!!!info

    The full documentation of the CEL Envoy library is available [here](../../cel-extensions/envoy.md).