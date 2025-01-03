# Envoy library

The `envoy` library adds some types and function to simplify the creation of [OkResponse](#okresponse) and [DeniedResponse](#deniedresponse) objects.

## Types

### `<CheckRequest>`

*CEL Type / Proto:* [`envoy.service.auth.v3.CheckRequest`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest)

### `<OkResponse>`

*CEL Type / Proto:* `envoy.OkResponse`

| Field | CEL Type / Proto | Docs |
|---|---|---|
| status | `google.rpc.Status` | [Docs](#status) |
| http_response | `envoy.service.auth.v3.OkHttpResponse` | [Docs](#okhttpresponse) |
| dynamic_metadata | `google.protobuf.Struct` | [Docs](#metadata) |

### `<DeniedResponse>`

*CEL Type / Proto:* `envoy.DeniedResponse`

| Field | CEL Type / Proto | Docs |
|---|---|---|
| status | `google.rpc.Status` | [Docs](#status) |
| http_response | `envoy.service.auth.v3.DeniedHttpResponse` | [Docs](#deniedhttpresponse) |
| dynamic_metadata | `google.protobuf.Struct` | [Docs](#metadata) |

### `<OkHttpResponse>`

*CEL Type / Proto:* [`envoy.service.auth.v3.OkHttpResponse`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-okhttpresponse)

### `<DeniedHttpResponse>`

*CEL Type / Proto:* [`envoy.service.auth.v3.DeniedHttpResponse`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-deniedhttpresponse)

### `<Metadata>`

*CEL Type / Proto:* [`google.protobuf.Struct`](https://protobuf.dev/reference/protobuf/google.protobuf/#struct)

### `<HeaderValueOption>`

*CEL Type / Proto:* [`envoy.config.core.v3.HeaderValueOption`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-headervalueoption)

### `<QueryParameter>`

*CEL Type / Proto:* [`envoy.config.core.v3.QueryParameter`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-queryparameter)

### `<Status>`

*CEL Type / Proto:* [`google.rpc.Status`](https://cloud.google.com/natural-language/docs/reference/rpc/google.rpc#status)

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

### envoy.QueryParam

This function creates a `<QueryParameter>` object.

#### Signature and overloads

```
envoy.QueryParam(<string> key, <string> value) -> <QueryParameter>
```

#### Example

```
envoy.QueryParam("foo", "bar")
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
```
<OkHttpResponse>.WithQueryParam(<string> key, <string> value) -> <OkHttpResponse>
```

#### Example

```
envoy.Allowed().WithQueryParam(envoy.QueryParam("foo", "bar"))
```
```
envoy.Allowed().WithQueryParam("foo", "bar")
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

This function creates a `<OkResponse>` / `DeniedResponse` object from an `<OkHttpResponse>` / `<DeniedHttpResponse>`.

#### Signature and overloads

```
<OkHttpResponse>.Response() -> <OkResponse>
```
```
<DeniedHttpResponse>.Response() -> <DeniedResponse>
```

#### Example

```
envoy.Allowed().Response()
```
```
envoy.Denied(401).Response()
```

### WithMessage

This function sets the `status.message` field of an `<OkResponse>` / `DeniedResponse` object.

#### Signature and overloads

```
<OkResponse>.WithMessage(<string> message) -> <OkResponse>
```
```
<DeniedResponse>.WithMessage(<string> message) -> <DeniedResponse>
```

#### Example

```
envoy.Allowed().Response().WithMessage("hello world!")
```
```
envoy.Denied(401).Response().WithMessage("hello world!")
```

### WithMetadata

This function sets the `dynamic_metadata` field of an `<OkResponse>` / `DeniedResponse` object.

#### Signature and overloads

```
<OkResponse>.WithMetadata(<Metadata> metadata) -> <OkResponse>
```
```
<DeniedResponse>.WithMetadata(<Metadata> metadata) -> <DeniedResponse>
```

#### Example

```
envoy.Allowed().Response().WithMetadata({ "foo": "bar" })
```
```
envoy.Denied(401).Response().WithMetadata({ "foo": "bar" })
```
