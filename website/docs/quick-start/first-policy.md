# First policy

The Kyverno Authz Server uses `ValidatingPolicy` resources to define authorization rules.

Policies for Envoy and HTTP authorization follow the same structure and logic bu have subtle differences.

## Key Concepts

- **Evaluation Mode**: Set to `Envoy` or `HTTP` to determine the request type
- **Failure Policy**: Controls behavior when policy evaluation fails (`Fail` or `Ignore`)
- **Match Conditions**: Optional CEL expressions for fine-grained request filtering
- **Variables**: Reusable named expressions available throughout the policy
- **Validation Rules**: CEL expressions that return authorization decisions (or `null` to continue processing with the next rule)

## Example policy

The policy below does the following:

1.    Parse the incoming request `authorization` header to decode a bearer token.
2.    Decoding uses the `jwt` CEL lib.
3.    Based on the presence, validity of the token, and roles contained in the claims, the policy will make a decision to allow or deny the request.

=== "Envoy"

    ```yaml
    apiVersion: policies.kyverno.io/v1alpha1
    kind: ValidatingPolicy
    metadata:
      name: quick-start
    spec:
      evaluation:
        mode: Envoy # (1)!
      failurePolicy: Fail # (2)!
      variables: # (3)!
      - name: authorization
        expression: object.attributes.request.http.headers[?"authorization"].orValue("").split(" ")
      - name: token
        expression: >
          size(variables.authorization) == 2 && variables.authorization[0].lowerAscii() == "bearer"
            ? jwt.Decode(variables.authorization[1], "secret")
            : null
      validations: # (4)!
        # request not authenticated -> 401
      - expression: >
          variables.token == null || !variables.token.Valid
            ? envoy.Denied(401).Response()
            : null
        # request authenticated but not admin role -> 403
      - expression: >
          variables.token.Claims.?role.orValue("") != "admin"
            ? envoy.Denied(403).Response()
            : null
        # request authenticated and admin role -> 200
      - expression: >
          envoy.Allowed().Response()
    ```
    
    1.    **Evaluation Mode**: Set to `Envoy` or `HTTP` to determine the request type
    2.    **Failure Policy**: Controls behavior when policy evaluation fails (`Fail` or `Ignore`)
    3.    **Variables**: Reusable named expressions available throughout the policy
    4.    **Validation Rules**: CEL expressions that return authorization decisions (or `null` to continue processing with the next rule)

=== "HTTP"


    ```yaml
    apiVersion: policies.kyverno.io/v1alpha1
    kind: ValidatingPolicy
    metadata:
      name: quick-start
    spec:
      evaluation:
        mode: HTTP # (1)!
      failurePolicy: Fail # (2)!
      variables: # (3)!
      - name: authorization
        expression: object.attributes.header[?"authorization"].orValue("").split(" ")
      - name: token
        expression: >
          size(variables.authorization) == 2 && variables.authorization[0].lowerAscii() == "bearer"
            ? jwt.Decode(variables.authorization[1], "secret")
            : null
      validations: # (4)!
        # request not authenticated -> allowed
      - expression: >
          variables.token == null || !variables.token.Valid
            ? http.Allowed().Response()
            : null
        # request authenticated but not admin role -> denied
      - expression: >
          variables.token.Claims.?role.orValue("") != "admin"
            ? http.Denied("authenticated but not an admin").Response()
            : null
        # request authenticated and admin role -> allowed
      - expression: >
          http.Allowed().Response()
    ```
    
    1.    **Evaluation Mode**: Set to `Envoy` or `HTTP` to determine the request type
    2.    **Failure Policy**: Controls behavior when policy evaluation fails (`Fail` or `Ignore`)
    3.    **Variables**: Reusable named expressions available throughout the policy
    4.    **Validation Rules**: CEL expressions that return authorization decisions (or `null` to continue processing with the next rule)
