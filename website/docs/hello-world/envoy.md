---
title: Envoy Authz Server
description: Example demonstrating how to run and test the Kyverno Authorization Server in Envoy mode using JWT-based authentication.
---

# Envoy Authz Server

This example demonstrates how to run and test the **Kyverno Authorization Server** in **Envoy** mode with a JWT-based authentication policy.

It shows how the server evaluates gRPC requests and enforces access rules using token validation.

## Prerequisites

Before starting, ensure you have:

- [Docker](https://www.docker.com) installed and running.
- [grpcurl](https://github.com/fullstorydev/grpcurl) for sending gRPC requests.

!!! tip
    You **do not need a Kubernetes cluster** — the server runs locally and loads policies directly from files.

## Step 1: Create a Policy

Save the following policy as `quick-start.yaml` in your working directory:

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
  - expression: >
      variables.token == null || !variables.token.Valid
        ? envoy.Denied(401).Response()
        : null
  - expression: >
      variables.token.Claims.?role.orValue("") != "admin"
        ? envoy.Denied(403).Response()
        : null
  - expression: >
      envoy.Allowed().Response()
```

This policy:

- Denies requests without a valid token (`401`).
- Denies authenticated users without an admin role (`403`).
- Allows admin users (`200`).

## Step 2: Run the Envoy Authz Server

Start the Authz Server with Docker:

```bash
docker run --rm \
  -v ${HOME}/.kube/:/etc/kubeconfig/ \
  -v ${PWD}/quick-start.yaml/:/data/policies/quick-start.yaml \
  -e KUBECONFIG=/etc/kubeconfig/config \
  -p 9081:9081 \
  ghcr.io/kyverno/kyverno-envoy-plugin:a83ddce53efe0a35dfe239d3089bdefa19ca4f80 \
  serve envoy authz-server --kube-policy-source=false \
  --external-policy-source file://data/policies
```

Expected startup output:

```
2025-11-04T17:42:32Z    INFO    HTTP Server starting... {"address": ":9080", "cert": "", "key": ""}
2025-11-04T17:42:32Z    INFO    GRPC Server starting... {"address": "[::]:9081", "network": "tcp"}
```

The gRPC server listens on port **9081**.

## Step 3: Send Requests to the Server

### Example 1 — Missing Authorization Header

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

Response:

```json
{
  "status": { "code": 7 },
  "deniedResponse": { "status": { "code": "Unauthorized" } }
}
```

### Example 2 — Non-admin User Token

```bash
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"

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

Response:

```json
{
  "status": { "code": 7 },
  "deniedResponse": { "status": { "code": "Forbidden" } }
}
```

### Example 3 — Admin User Token

```bash
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

Response:

```json
{
  "status": {},
  "okResponse": {}
}
```

## Step 4: Summary

1. Define a JWT-based Envoy authorization policy.  
2. Run the Kyverno Envoy Authz Server via Docker.  
3. Test requests with grpcurl and validate access control decisions.
