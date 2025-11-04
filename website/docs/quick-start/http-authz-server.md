# HTTP Authz Server

## Create a policy

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: policy
spec:
  evaluation:
    mode: HTTP
  matchConditions:
  - name: match-host
    expression: |
      object.attributes.host == "http-srv.app"
  validations:
  - expression: >
      http.Denied("not allowed").Response()
```

## Run the HTTP Authz Server

We will use Docker to run the Kyverno Authz Server image. From a terminal run the following command:

```bash
docker run --rm                                     \
    -v ${HOME}/.kube/:/etc/kubeconfig/              \
    -v ${PWD}/quick-start.yaml/:/data/policies/quick-start.yaml              \
    -e KUBECONFIG=/etc/kubeconfig/config            \
    -p 9083:9083 \
    ghcr.io/kyverno/kyverno-envoy-plugin:a83ddce53efe0a35dfe239d3089bdefa19ca4f80  \
    serve http authz-server --kube-policy-source=false \
    --external-policy-source file://data/policies
```

The server will start on port `9083`:

```
2025-11-04T10:23:08Z    INFO    HTTP Server starting... {"address": ":9080", "cert": "", "key": ""}
2025-11-04T10:23:08Z    INFO    HTTP Server starting... {"address": ":9083", "cert": "", "key": ""}
```

## Send requests to the server

```bash
curl -s -I -X POST http://127.0.0.1:9083
```

```bash
curl -s -I -X POST http://127.0.0.1:9083 -H "Host: http-srv.app"
```
