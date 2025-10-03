# Jwk library

The JWK lib helps working with jwk keys and key sets as described in [rfc7517](https://tools.ietf.org/html/rfc7517).

## Types

### `<Set>`

*CEL Type / Proto* `jwk.Set`

This is an opaque type with no available fields. Its purpose is to be used with [jwt.Decode](jwt.md#jwtdecode) to very a token issuer.

## Functions

### jwks.Fetch

The `jwks.Fetch` function fetches and parses a JWK resource specified by a URL.

#### Signature and overloads

```
jwks.Fetch(<string> url) -> <Set>
```

#### Example

```
jwks.Fetch("https://.../.well-known/jwks.json")
```
