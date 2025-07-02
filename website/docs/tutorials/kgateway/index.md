# KGateway

[kgateway](https://https://kgateway.dev/) is the most mature and widely deployed Envoy-based gateway in the market today. Built on open source and open standards, kgateway implements the Kubernetes Gateway API with a control plane that scales from lightweight microgateway deployments between services, to massively parallel centralized gateways handling billions of API calls, to advanced AI gateway use cases for safety, security, and governance when serving models or integrating applications with third-party LLMs. kgateway brings omni-directional API connectivity to any cloud and any environment.

This tutorial shows how kgateway can be configured to delegate authorization decisions to the Kyverno Authz Server.

## Setup

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install kgateway the Kyverno Authz Server
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Install KGateway

First we need to install KGateway in the cluster.

```bash
# install gateway API CDRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# install kgateway CDRDs
helm upgrade -i --create-namespace --namespace kgateway-system --version v2.0.3 --wait kgateway-crds oci://cr.kgateway.dev/kgateway-dev/charts/kgateway-crds

# install kgateway
helm upgrade -i --namespace kgateway-system --version v2.0.3 --wait kgateway oci://cr.kgateway.dev/kgateway-dev/charts/kgateway
```

### Deploy a sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# deploy the httpbin application
kubectl apply \
  -n demo \
  -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
```

### Set up an API gateway

Create an API gateway with an HTTP listener by using the Kubernetes Gateway API.

```bash
kubectl apply -f - <<EOF
kind: Gateway
apiVersion: gateway.networking.k8s.io/v1
metadata:
  name: http
  namespace: kgateway-system
spec:
  gatewayClassName: kgateway
  listeners:
  - name: http
    protocol: HTTP
    port: 80
    allowedRoutes:
      namespaces:
        from: All
EOF
```

!!!info

    Using Kind and getting a CrashLoopBackOff error with a Failed to create temporary file message in the logs? You might have a multi-arch platform issue on macOS. In your Docker Desktop settings, uncheck **Use Rosetta**, restart Docker, re-create your Kind cluster, and try again.

### Expose the app on the gateway

Now that you have an app and a gateway proxy, you can create a route to access the app.

```bash
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: httpbin
  namespace: demo
  labels:
    example: httpbin-route
spec:
  parentRefs:
  - name: http
    namespace: kgateway-system
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

TBC








## Testing

At this we have deployed and configured KGateway, the Kyverno Authz Server, a sample application, and the authorization and security policies.

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
curl -s -w "\nhttp_code=%{http_code}" http.kgateway-system/get
```

Calling with Alice’s JWT token will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" http.kgateway-system/get -H "authorization: Bearer $ALICE_TOKEN"
```

Calling with Bob’s JWT token will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" http.kgateway-system/get -H "authorization: Bearer $BOB_TOKEN"
```

## Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure kgateway to utilize the Kyverno Authz Server as an external authorization service.

Additionally, the tutorial provided an example policy to decode a JWT token and make a decision based on it.
