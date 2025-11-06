# Programmability

The authorization server CRD can be used to change things about the HTTP reponse being returned from the authz server

```yaml
apiVersion: authz.kyverno.io/v1alpha1
kind: AuthorizationServer
metadata:
  name: http-server
  namespace: default
spec:
  type:
    http:
      modifiers:
        request:
        response:
```
