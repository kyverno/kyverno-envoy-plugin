# Envoy library

The `envoy` library adds some types and function to simplify the creation of Envoy [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) objects.

## Types

| CEL Type | Proto | Docs |
|---|---|---|
| `<CheckRequest>` | `envoy.service.auth.v3.CheckRequest` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) |
| `<CheckResponse>` | `envoy.service.auth.v3.CheckResponse` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) |
| `<OkHttpResponse>` | `envoy.service.auth.v3.OkHttpResponse` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-okhttpresponse) |
| `<DeniedHttpResponse>` | `envoy.service.auth.v3.DeniedHttpResponse` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-deniedhttpresponse) |
| `<Metadata>` | `google.protobuf.Struct` | [Docs](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) |
| `<HeaderValueOption>` | `envoy.config.core.v3.HeaderValueOption` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-headervalueoption) |
| `<QueryParameter>` | `envoy.config.core.v3.QueryParameter` | [Docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-queryparameter) |

## Functions

### envoy.Allowed

This function creates an `<OkHttpResponse>` object.

#### Signature and overloads

```
envoy.Allowed() -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed()
```

### envoy.Denied

This function creates a `<DeniedHttpResponse>` object.

#### Signature and overloads

```
envoy.Denied(<int> code) -> <DeniedHttpResponse>
```

#### Example

```
envoy.Denied(401)
```

### envoy.Response

This function creates a `<CheckResponse>` object.

#### Signature and overloads

```
envoy.Response(<int> code) -> <CheckResponse>
```
```
envoy.Response(<OkHttpResponse> ok) -> <CheckResponse>
```
```
envoy.Response(<DeniedHttpResponse> denied) -> <CheckResponse>
```

#### Example

```
// ok
envoy.Response(0)

// permission denied
envoy.Response(7)
```
```
envoy.Response(envoy.Allowed())
```
```
envoy.Response(envoy.Denied(401))
```

### envoy.Null

This function creates a null `<CheckResponse>` (useful when an authorisation rule doesn't make a decision).

#### Signature and overloads

```
envoy.Null() -> <CheckResponse>
```

#### Example

```
envoy.Null()
```

### envoy.Header

This function creates an `<HeaderValueOption>` object.

#### Signature and overloads

```
envoy.Header(<string> key, <string> value) -> <HeaderValueOption>
```

#### Example

```
envoy.Header("foo", "bar")
```

### WithBody

This function sets the body of a `<DeniedHttpResponse>` object.

#### Signature and overloads

```
<DeniedHttpResponse>.WithBody(<string> body) -> <DeniedHttpResponse>
```

#### Example

```
envoy.Denied(401).WithBody("Unauthorized Request")
```

### WithHeader

This function adds a `<HeaderValueOption>`:

- When the request is sent upstream by Envoy, in the case of an `<OkHttpResponse>`.
- When the response is sent downstream by Envoy, in the case of a `<DeniedHttpResponse>`.

#### Signature and overloads

```
<OkHttpResponse>.WithHeader(<HeaderValueOption> header) -> <OkHttpResponse>
```
```
<OkHttpResponse>.WithHeader(<string> key, <string> value) -> <OkHttpResponse>
```
```
<DeniedHttpResponse>.WithHeader(<HeaderValueOption> header) -> <DeniedHttpResponse>
```
```
<DeniedHttpResponse>.WithHeader(<string> key, <string> value) -> <DeniedHttpResponse>
```

#### Example

```
envoy.Allowed().WithHeader(envoy.Header("foo", "bar"))
```
```
envoy.Allowed().WithHeader("foo", "bar")
```
```
envoy.Denied(401).WithHeader(envoy.Header("foo", "bar"))
```
```
envoy.Denied(401).WithHeader("foo", "bar")
```

### WithoutHeader

This function marks a header to be removed when the request is sent upstream by Envoy.

#### Signature and overloads

```
<OkHttpResponse>.WithoutHeader(<string> header) -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed().WithoutHeader("foo")
```

### WithResponseHeader

This function adds a `<HeaderValueOption>` when the response is sent downstream by Envoy.

#### Signature and overloads

```
<OkHttpResponse>.WithResponseHeader(<HeaderValueOption> header) -> <OkHttpResponse>
```
```
<OkHttpResponse>.WithResponseHeader(<string> key, <string> value) -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed().WithResponseHeader(envoy.Header("foo", "bar"))
```
```
envoy.Allowed().WithResponseHeader("foo", "bar")
```

### WithQueryParam

This function adds a `<QueryParameter>` to be added when the request is sent upstream by Envoy.

#### Signature and overloads

```
<OkHttpResponse>.WithQueryParam(<QueryParameter> param) -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed().WithQueryParam(envoy.config.core.v3.QueryParameter{
    key: "foo",
    bar: "bar"
})
```

### WithoutQueryParam

This function marks a query parameter to be removed when the request is sent upstream by Envoy.

#### Signature and overloads

```
<OkHttpResponse>.WithoutQueryParam(<string> param) -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed().WithoutQueryParam("foo")
```

### KeepEmptyValue

This function sets the `keep_empty_value` field of an `<HeaderValueOption>` object.

#### Signature and overloads

```
<HeaderValueOption>.KeepEmptyValue() -> <HeaderValueOption>
```
```
<HeaderValueOption>.KeepEmptyValue(<bool> keep) -> <HeaderValueOption>
```

#### Example

```
envoy.Header("foo", "bar").KeepEmptyValue()
```
```
envoy.Header("foo", "bar").KeepEmptyValue(true)
```

### Response

This function creates a `<CheckResponse>` object from an `<OkHttpResponse>` or `<DeniedHttpResponse>`.

#### Signature and overloads

```
<OkHttpResponse>.Response() -> <CheckResponse>
```
```
<DeniedHttpResponse>.Response() -> <CheckResponse>
```

#### Example

```
envoy.Allowed().Response()
```
```
envoy.Denied(401).Response()
```

### WithMessage

This function sets the `status.message` field of a `<CheckResponse>` object.

#### Signature and overloads

```
<CheckResponse>.WithMessage(<string> message) -> <CheckResponse>
```

#### Example

```
envoy.Allowed().Response().WithMessage("hello world!")
```
```
envoy.Denied(401).Response().WithMessage("hello world!")
```

### WithMetadata

This function sets the `dynamic_metadata` field of a `<CheckResponse>` object.

#### Signature and overloads

```
<CheckResponse>.WithMetadata(<Metadata> metadata) -> <CheckResponse>
```

#### Example

```
envoy.Allowed().Response().WithMetadata({ "foo": "bar" })
```
```
envoy.Denied(401).Response().WithMetadata({ "foo": "bar" })
```
