# CEL extensions

The CEL engine used to evaluate variables and authorization rules has been extended with libraries to help processing the input `CheckRequest` and forge the corresponding `OkResponse` and/or `DeniedResponse`.

## Envoy plugin libraries

- [Envoy](./envoy.md)
- [Jwk](./jwk.md)
- [Jwt](./jwt.md)
- [Json](./json.md)

## Common libraries

The libraries below are common CEL extensions enabled in the Kyverno Authz Server CEL engine:

- [Optional types](https://pkg.go.dev/github.com/google/cel-go/cel#OptionalTypes)
- [Cross type numeric comparisons](https://pkg.go.dev/github.com/google/cel-go/cel#CrossTypeNumericComparisons)
- [Bindings](https://pkg.go.dev/github.com/google/cel-go/ext#readme-bindings)
- [Encoders](https://pkg.go.dev/github.com/google/cel-go/ext#readme-encoders)
- [Lists](https://pkg.go.dev/github.com/google/cel-go/ext#readme-lists)
- [Math](https://pkg.go.dev/github.com/google/cel-go/ext#readme-math)
- [Protos](https://pkg.go.dev/github.com/google/cel-go/ext#readme-protos)
- [Sets](https://pkg.go.dev/github.com/google/cel-go/ext#readme-sets)
- [Strings](https://pkg.go.dev/github.com/google/cel-go/ext#readme-strings)

## Kubernetes libraries

The libraries below are imported from Kubernetes:

- [Lists](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-list-library)
- [Regex](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-regex-library)
- [URL](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-url-library)
- [IP](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-ip-address-library)
- [CIDR](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-cidr-library)
- [Format](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-format-library)
- [Quantity](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-quantity-library)
- [Semver](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-semver-library)

## Kyverno libraries

The libraries below are imported from Kyverno:

- [HTTP](https://kyverno.io/docs/policy-types/cel-libraries/#http-library)
- [Image](https://kyverno.io/docs/policy-types/cel-libraries/#image-library)
- [ImageData](https://kyverno.io/docs/policy-types/cel-libraries/#imagedata-library)
