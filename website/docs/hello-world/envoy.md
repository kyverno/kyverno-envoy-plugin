# Envoy Authz Server

## Create a policy

```yaml
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: quick-start
spec:
  evaluation:
    mode: Envoy 
  failurePolicy: Fail 
  variables: 
  - name: authorization
    expression: object.attributes.request.http.headers[?"authorization"].orValue("").split(" ")
  - name: token
    expression: >
      size(variables.authorization) == 2 && variables.authorization[0].lowerAscii() == "bearer"
        ? jwt.Decode(variables.authorization[1], "secret")
        : null
  validations: 
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

## Run the Envoy Authz Server

We will use Docker to run the Kyverno Authz Server image. From a terminal run the following command:

```bash
docker run --rm                                     \
    -v ${HOME}/.kube/:/etc/kubeconfig/              \
    -v ${PWD}/quick-start.yaml/:/data/policies/quick-start.yaml              \
    -e KUBECONFIG=/etc/kubeconfig/config            \
    -p 9081:9081 \
    ghcr.io/kyverno/kyverno-envoy-plugin:a83ddce53efe0a35dfe239d3089bdefa19ca4f80  \
    serve envoy authz-server --kube-policy-source=false \
    --external-policy-source file://data/policies
```

The server will start on port `9081`:

```
2025-11-04T17:42:32Z    INFO    HTTP Server starting... {"address": ":9080", "cert": "", "key": ""}
2025-11-04T17:42:32Z    INFO    GRPC Server starting... {"address": "[::]:9081", "network": "tcp"}
```

## Send requests to the server

```bash
grpcurl -plaintext -d @ localhost:9081 envoy.service.auth.v3.Authorization/Check <<EOF
{
  "attributes": {
    "request": {
      "http": {
        "headers": {
          "authorization": "empty"
        }
      }
    }
  }
}
EOF
```

```json
{
  "status": {
    "code": 7
  },
  "deniedResponse": {
    "status": {
      "code": "Unauthorized"
    }
  }
}
```

```bash
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"

grpcurl -plaintext -d @ localhost:9081 envoy.service.auth.v3.Authorization/Check <<EOF
{
  "attributes": {
    "request": {
      "http": {
        "headers": {
          "authorization": "bearer $ALICE_TOKEN"
        }
      }
    }
  }
}
EOF
```

```json
{
  "status": {
    "code": 7
  },
  "deniedResponse": {
    "status": {
      "code": "Forbidden"
    }
  }
}
```

```bash
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"

grpcurl -plaintext -d @ localhost:9081 envoy.service.auth.v3.Authorization/Check <<EOF
{
  "attributes": {
    "request": {
      "http": {
        "headers": {
          "authorization": "bearer $BOB_TOKEN"
        }
      }
    }
  }
}
EOF
```

```json
{
  "status": {},
  "okResponse": {}
}
```