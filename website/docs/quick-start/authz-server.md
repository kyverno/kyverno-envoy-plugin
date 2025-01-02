# Authz server

## Setup

In this quick start guide we will deploy the Kyverno Authz Server inside a cluster.

Then you will interface [Istio](https://istio.io/latest/), an open source service mesh with the Kyverno Authz Server to delegate the request authorisation based on policies installed in the cluster.

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install the Kyverno Authz Server
- [istioctl](https://istio.io/latest/docs/setup/getting-started/#download) to configure the mesh
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
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

Notice that in the configuration, we define an `extensionProviders` section that points to the Kyverno Authz Server we will install in the next step:

```yaml
[...]
    extensionProviders:
    - name: kyverno-authz-server.local
      envoyExtAuthzGrpc:
        service: kyverno-authz-server.kyverno.svc.cluster.local
        port: '9081'
[...]
```

### Deploy the Kyverno Authz Server

The first step is to deploy the Kyverno Authz Server.

```bash
# create the kyverno namespace
kubectl create ns kyverno

# label the namespace to inject the envoy proxy
kubectl label namespace kyverno istio-injection=enabled

# deploy the kyverno authz server
helm install kyverno-authz-server --namespace kyverno --wait  \
  --repo https://kyverno.github.io/kyverno-envoy-plugin       \
  kyverno-authz-server
```

### Deploy a sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# label the namespace to inject the envoy proxy
kubectl label namespace demo istio-injection=enabled

# deploy the httpbin application
kubectl apply -n demo -f \
  https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
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

### Deploy a Kyverno AuthorizationPolicy

A Kyverno `AuthorizationPolicy` defines the rules used by the Kyverno authz server to make a decision based on a given Envoy [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest).

It uses the [CEL language](https://github.com/google/cel-spec) to analyse the incoming `CheckRequest` and is expected to produce a [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) in return.

```bash
# deploy kyverno authorization policy
kubectl apply -f - <<EOF
apiVersion: envoy.kyverno.io/v1alpha1
kind: AuthorizationPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  variables:
  - name: force_authorized
    expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  deny:
  - match: >
      !variables.allowed
    response: >
      envoy.Denied(403).Response()
EOF
```

This simple policy will deny requests if they don't contain the header `x-force-authorized` with the value `enabled` or `true`.

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

Now we can send requests to the sample application and verify the result.

The following request will return `403` (denied by our policy):

```bash
curl -s -w "\nhttp_code=%{http_code}" httpbin:8000/get
```

The following request will return `200` (allowed by our policy):

```bash
curl -s -w "\nhttp_code=%{http_code}" httpbin:8000/get -H "x-force-authorized: true"
```

## Wrap Up

Congratulations on completing the quick start guide!

This tutorial demonstrated how to configure Istio’s EnvoyFilter to utilize the Kyverno Authz Server as an external authorization service.
