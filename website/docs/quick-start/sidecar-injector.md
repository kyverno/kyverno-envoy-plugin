# Sidecar injector

## Setup

In this quick start guide we will deploy the Kyverno Authz Server as a sidecar using a mutating webhook.

Then you will interface [Istio](https://istio.io/latest/), an open source service mesh with the Kyverno Authz Server to delegate the request authorisation based on policies installed in the cluster.

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install the Kyverno Authz Server
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.34.0

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Configure the mesh

We need to register the Kyverno Authz Server with Istio.

```bash
# install itio base chart
helm install istio-base \
  --namespace istio-system --create-namespace \
  --wait \
  --repo https://istio-release.storage.googleapis.com/charts base

# install istiod chart
helm install istiod \
  --namespace istio-system --create-namespace \
  --wait \
  --repo https://istio-release.storage.googleapis.com/charts istiod \
  --values - <<EOF
meshConfig:
  accessLogFile: /dev/stdout
  extensionProviders:
  - name: kyverno-authz-server
    envoyExtAuthzGrpc:
      service: kyverno-authz-server.local
      port: 9081
EOF
```

Notice that in the configuration, we define an `extensionProviders` section that points to `kyverno-authz-server.local`, this the service entry we will use to let Envoy talk to our sidecar:

```yaml
[...]
meshConfig:
  extensionProviders:
  - name: kyverno-authz-server.local
    envoyExtAuthzGrpc:
      service: kyverno-authz-server.local
      port: '9081'
[...]
```

### Register the authz server sidecar ServiceEntry

We need to tell istio about the sidecar we injected and how to reach it.

```bash
# register authz server sidecar in the mesh
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1
kind: ServiceEntry
metadata:
  name: kyverno-authz-server
spec:
  hosts:
  - kyverno-authz-server.local
  endpoints:
  - address: 127.0.0.1
  ports:
  - name: grpc
    number: 9081
    protocol: GRPC
  resolution: STATIC
EOF
```

### Deploy cert-manager

The Kyverno Authz Server comes with a validation webhook and needs a certificate to let the api server call into it.

Let's deploy `cert-manager` to manage the certificate we need.

```bash
# install cert-manager
helm install cert-manager \
  --namespace cert-manager --create-namespace \
  --wait \
  --repo https://charts.jetstack.io cert-manager \
  --set crds.enabled=true

# create a self-signed cluster issuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF
```

For more certificate management options, refer to [Certificates management](../install/certificates.md).

### Install Kyverno ValidatingPolicy CRD

Before deploying the Kyverno Authz Server, we need to install the Kyverno ValidatingPolicy CRD.

```bash
kubectl apply \
  -f https://raw.githubusercontent.com/kyverno/kyverno/refs/heads/main/config/crds/policies.kyverno.io/policies.kyverno.io_validatingpolicies.yaml
```

### Create the demo Namespace

```bash
# create the demo namespace
kubectl create ns demo

# label the namespace to inject the envoy proxy
kubectl label namespace demo istio-injection=enabled

# label the namespace to inject the authz server sidecar
kubectl label namespace demo kyverno-injection=enabled
```

### Deploy a Kyverno ValidatingPolicy

A Kyverno `ValidatingPolicy` defines the rules used by the Kyverno authz server to make a decision based on a given Envoy [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest).

It uses the [CEL language](https://github.com/google/cel-spec) to analyse the incoming [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) and is expected to produce a [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) in return.

!!!note "Sidecar can't talk with API Server"
    Because the sidecar usually doesn't have the permissions to fetch policies from the API server, we need to provide the policies using an external source.
    In this example, we use a config map.

```bash
# deploy kyverno validating policy
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: kyverno-authz-server
  namespace: demo
data:
  policy.yaml: |
    apiVersion: policies.kyverno.io/v1alpha1
    kind: ValidatingPolicy
    metadata:
      name: demo
    spec:
      failurePolicy: Fail
      evaluation:
        mode: Envoy
      variables:
      - name: force_authorized
        expression: object.attributes.request.http.headers[?"x-force-authorized"].orValue("")
      - name: allowed
        expression: variables.force_authorized in ["enabled", "true"]
      validations:
      - expression: >-
          !variables.allowed ? envoy.Denied(403).Response() : null
EOF
```

This simple policy will deny requests if they don't contain the header `x-force-authorized` with the value `enabled` or `true`.

### Deploy the Sidecar injector

Now we can deploy the Kyverno Authz Server.

```bash
# create the kyverno namespace
kubectl create ns kyverno

# deploy the kyverno sidecar injector
helm install kyverno-authz-server \
  --namespace kyverno \
  --wait  \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-sidecar-injector \
  --set certificates.certManager.issuerRef.group=cert-manager.io \
  --set certificates.certManager.issuerRef.kind=ClusterIssuer \
  --set certificates.certManager.issuerRef.name=selfsigned-issuer
```

### Deploy the sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# deploy the httpbin application
kubectl apply \
  -n demo \
  -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
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
