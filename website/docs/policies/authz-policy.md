# Policy Reference 

This page provides guidance on writing policies for request content processed by the kyverno-json validating policy, utilizing Envoy’s [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html).

### Writing Policies 

Let start with an example policy that restricts access to an endpoint based on user's role and permissions.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
    name: checkrequest
spec:
    rules:
    - name: deny-guest-request-at-post
        assert:
        any:
        - message: "POST method calls at path /book are not allowed to guests users"
            check:
            request:
                http:
                    method: POST
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): admin
                    path: /book                             
        - message: "GET method call is allowed to both guest and admin users"
            check:
            request:
                http:
                    method: GET
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): admin
                    path: /book 
        - message: "GET method call is allowed to both guest and admin users"
            check:
            request:
                http:
                    method: GET
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): guest
                    path: /book  
```

The above policy uses the `jwt_decode` builtin function to parse and verify the JWT containing information about the user making the request. it uses other builtins like `split`, `base64_decode`, `campare`, `contains` etc kyverno has many different [function](https://kyverno.github.io/kyverno-json/latest/jp/functions/) which can be used in policy.

Sample input recevied by kyverno-json validating policy is shown below:

```json
{
  "source": {
    "address": {
      "socketAddress": {
        "address": "10.244.1.10",
        "portValue": 59252
      }
    }
  },
  "destination": {
    "address": {
      "socketAddress": {
        "address": "10.244.1.4",
        "portValue": 8080
      }
    }
  },
  "request": {
    "time": "2024-04-09T07:42:29.634453Z",
    "http": {
      "id": "14694995155993896575",
      "method": "GET",
      "headers": {
        ":authority": "testapp.demo.svc.cluster.local:8080",
        ":method": "GET",
        ":path": "/book",
        ":scheme": "http",
        "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk",
        "user-agent": "Wget",
        "x-forwarded-proto": "http",
        "x-request-id": "27cd2724-e0f4-4a69-a1b1-9a94edfa31bb"
      },
      "path": "/book",
      "host": "echo.demo.svc.cluster.local:8080",
      "scheme": "http",
      "protocol": "HTTP/1.1"
    }
  },
  "metadataContext": {},
  "routeMetadataContext": {}
}
```

With the help of assertion tree, we can write policies that can be used to validate the request content. 

An `assert` declaration contains an `any` or `all` list in which each entry contains a `check` and a `message`. The `check` contains a JMESPath expression that is evaluated against the request content. The `message` is a string that is returned when the check fails.
A check can contain one or more JMESPath expressions. Expressions represent projections of seleted data in the JSON payload and the result of this projection is passed to descendants for futher analysis. All comparisons happen in the leaves of the assertion tree.

For more detail checkout [Policy Structure](https://kyverno.github.io/kyverno-json/latest/policies/policies/) and [Assertion trees](https://kyverno.github.io/kyverno-json/latest/policies/asserts/#assertion-trees).

- HTTP method `request.http.method`
- Request path `request.http.path`
- Authorization header `request.http.headers.authorization`

when we decode this above mentioned JWT token in the request payload we get payload.role `guest`:

```json
{
  "exp": 2241081539,
  "nbf": 1514851139,
  "role": "guest",
  "sub": "YWxpY2U="
}
```
With the input value above, the answer is:
```
true
```
