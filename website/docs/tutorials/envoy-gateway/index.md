# Envoy Gateway

[Envoy Gateway](https://gateway.envoyproxy.io/) is an open source project for managing [Envoy Proxy](https://www.envoyproxy.io/) as a standalone or Kubernetes-based application
gateway. [Gateway API](https://gateway-api.sigs.k8s.io/) resources are used to dynamically provision and configure the managed Envoy Proxies.

This tutorial shows how Envoy Gateway can be configured to delegate authorization decisions to the Kyverno Authz Server.

## Setup

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install Envoy Gateway the Kyverno Authz Server
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Install Envoy Gateway

First we need to install Envoy Gateway in the cluster.

```bash
# install envoy gateway
helm install envoy-gateway -n envoy-gateway-system --create-namespace --wait --version v1.2.2 oci://docker.io/envoyproxy/gateway-helm
```

### Deploy a sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# deploy the httpbin application
kubectl apply -n demo -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
```

### Create a GatewayClass and a Gateway

With Envoy Gateway installed we can now create a `Gateway`. To do so we will also create a dedicated `GatewayClass`.

Depending on your setup you will potentially need to create an `EnvoyProxy` resource to customize the way Envoy Gateway will create the underlying `Service`. The script below creates one to set the name and type of the service because the kind cluster created in the first step doesn't come with load balancer support.

```yaml
# create a gateway
kubectl apply -n demo -f - <<EOF
apiVersion: gateway.envoyproxy.io/v1alpha1
kind: EnvoyProxy
metadata:
  name: demo
spec:
  provider:
    type: Kubernetes
    kubernetes:
      envoyService:
        name: internet   # use a known name for the created service
        type: ClusterIP  # because a kind cluster has no support for LB
---
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: demo
spec:
  controllerName: gateway.envoyproxy.io/gatewayclass-controller
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: demo
spec:
  gatewayClassName: demo
  infrastructure:
    parametersRef:
      group: gateway.envoyproxy.io
      kind: EnvoyProxy
      name: demo
  listeners:
  - name: http
    protocol: HTTP
    port: 80
EOF
```

### Create an HTTPRoute to the sample application

Next, we need to link the Gateway to our sample applicate with an `HTTPRoute`.

```yaml
# create an http route to the sample app
kubectl apply -n demo -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: demo
spec:
  parentRefs:
  - name: demo
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - group: ''
      kind: Service
      name: httpbin
      port: 8000
      weight: 1
EOF
```

### Deploy the Kyverno Authz Server

Now deploy the Kyverno Authz Server.

```bash
# deploy the kyverno authz server
helm install kyverno-authz-server --namespace kyverno --create-namespace --wait --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server
```

## Create a Kyverno AuthorizationPolicy

In summary the policy below does the following:

- Checks that the JWT token is valid
- Checks that the action is allowed based on the token payload `role` and the request path

```yaml
# deploy kyverno authorization policy
kubectl apply -f - <<EOF
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: authorization
    expression: object.attributes.request.http.headers[?"authorization"].orValue("").split(" ")
  - name: token
    expression: >
      size(variables.authorization) == 2 && variables.authorization[0].lowerAscii() == "bearer"
        ? jwt.Decode(variables.authorization[1], "secret")
        : null
  deny:
    # request not authenticated -> 401
  - match: >
      variables.token == null || !variables.token.Valid
    response: >
      envoy.Denied(401).Response()
    # request authenticated but not admin role -> 403
  - match: >
      variables.token.Claims.?role.orValue("") != "admin"
    response: >
      envoy.Denied(403).Response()
  allow:
    # request authenticated and admin role -> 200
  - response: >
      envoy
        .Allowed()
        .WithHeader("x-validated-by", "my-security-checkpoint")
        .WithoutHeader("x-force-authorized")
        .WithResponseHeader("x-add-custom-response-header", "added")
        .Response()
EOF
```

### Deploy an Envoy Gateway SecurityPolicy

A `SecurityPolicy` is the custom Envoy Gateway resource to configure underlying Envoy Proxy to use an external auth server (the Kyverno Authz Server we installed in a prior step).

```yaml
# deploy envoy gateway security policy
kubectl apply -n demo -f - <<EOF
apiVersion: gateway.envoyproxy.io/v1alpha1
kind: SecurityPolicy
metadata:
  name: demo
spec:
  targetRefs:
  - group: gateway.networking.k8s.io
    kind: HTTPRoute
    name: demo
  extAuth:
    grpc:
      backendRef:
        group: ''
        kind: Service
        name: kyverno-authz-server
        namespace: kyverno
        port: 9081
EOF
```

Notice that in this resource, we define the Kyverno Authz Server service as the GRPC backend:

```yaml
[...]
  extAuth:
    grpc:
      backendRef:
        group: ''
        kind: Service
        name: kyverno-authz-server
        namespace: kyverno
        port: 9081
[...]
```

Also notice that the security policy applies to the `demo` HTTPRoute:

```yaml
[...]
  targetRefs:
    - group: gateway.networking.k8s.io
      kind: HTTPRoute
      name: demo
[...]
```

### Grant access to the Kyverno Authz Server service

Last thing we need to configure is to grant access to the Kyverno Authz Server service for our SecurityPolicy to take effect.

```yaml
# grant access
kubectl apply -n kyverno -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: demo
spec:
  from:
  - group: gateway.envoyproxy.io
    kind: SecurityPolicy
    namespace: demo
  to:
  - group: ''
    kind: Service
EOF
```

## Testing

At this we have deployed and configured Envoy Gateway, the Kyverno Authz Server, a sample application, and the authorization and security policies.

### Start an in-cluster shell

Let's start a pod in the cluster with a shell to call into the sample application.

```bash
# run an in-cluster shell
kubectl run -i -t busybox --image=alpine --restart=Never -n demo
```

### Install curl

We will use curl to call into the sample application but it's not installed in our shell, let's install it in the pod.

```bash
# install curl
apk add curl
```

### Call into the sample application

Now we can send request to the sample application and verify the result.

For convenience, we will store Alice’s and Bob’s tokens in environment variables.

Here Bob is assigned the admin role and Alice is assigned the guest role.

```bash 
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"
```

Calling without a JWT token will return `401`:

```bash
curl -s -w "\nhttp_code=%{http_code}" internet.envoy-gateway-system/get
```

Calling with Alice’s JWT token will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" internet.envoy-gateway-system/get -H "authorization: Bearer $ALICE_TOKEN"
```

Calling with Bob’s JWT token will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" internet.envoy-gateway-system/get -H "authorization: Bearer $BOB_TOKEN"
```

## Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure Envoy Gateway to utilize the Kyverno Authz Server as an external authorization service.

Additionally, the tutorial provided an example policy to decode a JWT token and make a decision based on it.
