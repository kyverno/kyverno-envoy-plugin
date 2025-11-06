# Programmability

There are two flags that you can pass to the Authz server to transform attributes of requests and responses coming into and out of it. The `--input-expression` (for requests) and the `--output-expression` (for responses).
Those flags take a string represting a CEL expression take an input of a `http.CheckRequest` (for input flag) or `httpserver.HttpResponse` (for output flag) and evaluate to the same type. They can be used to change things like add a header, modify a header.. etc. But the most valuable use case is changing status code of a response coming out of the authorizer. an example for this is:
```
httpserver.HttpResponse{ status: 401, body: bytes(object.denied.reason), header: {"authenticated-by": ["kyverno-authz-server"]} }
```
