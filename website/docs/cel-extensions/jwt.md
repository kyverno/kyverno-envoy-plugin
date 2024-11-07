# Jwt library

Policies have native functionality to decode and verify the contents of JWT tokens in order to enforce additional authorization logic on requests.

## Types

### `<Token>`

*CEL Type / Proto* `jwt.Token`

| Field | CEL Type / Proto | Docs |
|---|---|---|
| Valid | `bool` | |
| Header | `google.protobuf.Struct` | [Docs](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) |
| Claims | `google.protobuf.Struct` | [Docs](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) |

## Functions

### jwt.Decode

The `jwt.Decode` function decodes and validates a JWT token.
It accepts two arguments: the token and the secret to verify the signature.

#### Signature and overloads

```
jwt.Decode(<string> token, <string> key) -> <Token>
```

#### Example

```
jwt.Decode("eyJhbGciOiJIUzI1NiI....", "secret")
```
