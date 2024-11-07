# Istio 

[Istio](https://istio.io/latest/) is an open source service mesh for managing the different microservices that make up a cloud-native application. Istio provides a mechanism to use a service as an external authorizer with the [AuthorizationPolicy API](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/).

This tutorial shows how Istio’s AuthorizationPolicy can be configured to delegate authorization decisions to the Kyverno Authz Server.

## Setup

### Prerequisites

- A Kubernetes cluster with Istio installed
- [Helm](https://helm.sh/) to install the Kyverno Authz Server
- [istioctl](https://istio.io/latest/docs/setup/getting-started/#download) to configure the mesh
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) and istall Istio with Helm.

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m

# install istio
helm install istio-base --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts base
helm install istiod --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts istiod
```

### Deploy the Kyverno Authz Server

The first step is to deploy the Kyverno Authz Server.

```bash
# create the kyverno namespace
kubectl create ns kyverno

# label the namespace to inject the envoy proxy
kubectl label namespace kyverno istio-injection=enabled

# deploy the kyverno authz server
helm install kyverno-authz-server --namespace kyverno --wait --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server
```

### Configure the mesh

We need to register the Kyverno Authz Server with Istio.

```bash
# configure the mesh
istioctl install -y -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
    accessLogFile: /dev/stdout
    extensionProviders:
    - name: kyverno-authz-server.local
      envoyExtAuthzGrpc:
        service: kyverno-authz-server.kyverno.svc.cluster.local
        port: '9081'
EOF
```

Notice that in the configuration, we define an `extensionProviders` section that points to the Kyverno Authz Server installation:

```yaml
[...]
    extensionProviders:
    - name: kyverno-authz-server.local
      envoyExtAuthzGrpc:
        service: kyverno-authz-server.kyverno.svc.cluster.local
        port: '9081'
[...]
```

### Deploy a sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# label the namespace to inject the envoy proxy
kubectl label namespace demo istio-injection=enabled

# deploy the httpbin application
kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml -n demo
```

### Deploy an Istio AuthorizationPolicy

An `AuthorizationPolicy` is the custom Istio resource that defines the services that will be protected by the Kyverno Authz Server.

```bash
# deploy istio authorization policy
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: kyverno-authz-server
  namespace: demo
spec:
  action: CUSTOM
  provider:
    name: kyverno-authz-server.local
  rules:
  - {} # empty rules, it will apply to all requests
EOF
```

Notice that in this resource, we define the Kyverno Authz Server `extensionProvider` you set in the Istio configuration:

```yaml
[...]
  provider:
    name: kyverno-authz-server.local
[...]
```

## Creating a Kyverno AuthorizationPolicy

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
  authorizations:
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
EOF
```

## Testing

At this we have deployed and configured Istio, the Kyverno Authz Server, a sample application, and the authorization policies.

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
curl -s -w "\nhttp_code=%{http_code}" httpbin:8000/get
```

Calling with Alice’s JWT token will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" httpbin:8000/get -H "authorization: Bearer $ALICE_TOKEN"
```

Calling with Bob’s JWT token will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" httpbin:8000/get -H "authorization: Bearer $BOB_TOKEN"
```

## Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure Istio’s EnvoyFilter to utilize the Kyverno Authz Server as an external authorization service.

Additionally, the tutorial provided an example policy to decode a JWT token and make a decision based on it.
