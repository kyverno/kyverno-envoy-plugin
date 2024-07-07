# Gloo Edge demo

This GLoo Edge demo is the prototype of Kyverno-Envoy-plugin. 

## Overview

[Gloo Edge](https://docs.solo.io/gloo-edge/latest/) is a Kubernetes-native API gateway built on top of Envoy Proxy. Gloo Edge is designed to be extensible and pluggable, and can be used to secure, control, and observe traffic to and from any application. `Gloo Edge` an API Gateway built on Envoy, offers a `Kubernetes` Custom Resource Definition (CRD) for managing Envoy configurations to handle traffic management and routing.

`Gloo Edge` supports the creation of a [Custom External Authorization Service](https://docs.solo.io/gloo-edge/latest/guides/security/auth/custom_auth/) that adheres to the Envoy specification for an External Authorization Server.

This tutorial demonstrates how to use the Kyverno-Envoy-Plugin with Gloo Edge to enforce security policies for upstream services.

## Demo instructions

### Required tools

1. [`minikube`](https://minikube.sigs.k8s.io/docs/start/)
2. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
3. [`helm`](https://helm.sh/docs/intro/install/)

### Create a local cluster and install Gloo Edge

Start a local cluster with `minikube`:

```bash
$ minikube start    
```

Setup and configure `Gloo Edge` using the below commands:

```bash
$ helm repo add gloo https://storage.googleapis.com/solo-public-helm
$ helm upgrade --install --namespace gloo-system --create-namespace gloo gloo/gloo
$ kubectl config set-context $(kubectl config current-context) --namespace=gloo-system
```
### Create a VirtualService and Upstream

[VirtualService](https://docs.solo.io/gloo-edge/latest/introduction/architecture/concepts/#virtual-services) define a set of route rules that live under a domain or set of domains. Route rules consist of matchers, which specify the kind of function calls to match (requests and events, are currently supported), and the name of the destination (or destinations) where to route them.

[Upstreams](https://docs.solo.io/gloo-edge/latest/introduction/architecture/concepts/#upstreams) define destinations for routes. Upstreams tell Gloo Edge what to route to and how to route to them. Gloo Edge determines how to handle routing for the Upstream based on its spec field. Upstreams have a type-specific spec field that provides routing information to Gloo Edge.

In this tutorial, we will create a VirtualService and Upstream that will route requests to the `httpbin.org` service.

```bash
$ kubectl apply -f - <<EOF
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: httpbin
spec:
  static:
    hosts:
      - addr: httpbin.org
        port: 80
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: httpbin
spec:
  virtualHost:
    domains:
      - '*'
    routes:
      - matchers:
         - prefix: /
        routeAction:
          single:
            upstream:
              name: httpbin
              namespace: gloo-system
        options:
          autoHostRewrite: true
EOF          
```

### Test Gloo

For simplification port-forwarding will be used. Open another terminal and execute.

```bash
$ kubectl port-forward deployment/gateway-proxy 8080:8080
```    

For simplicity, we will use the `curl` command to send requests to the `httpbin.org` service.
The `VirtualService` we created earlier will route all requests to the [https://httpbin.org/](https://httpbin.org/).
Let's test the Gloo works properly by running the following command in the first terminal:

```bash
$ curl -XGET -Is http://localhost:8080/get
HTTP/1.1 200 OK
date: Sat, 06 Jul 2024 22:00:56 GMT
content-type: application/json
content-length: 302
server: envoy
access-control-allow-origin: *
access-control-allow-credentials: true
x-envoy-upstream-service-time: 853
```
```bash
$ curl http -XGET -Is http://localhost:8080/post
HTTP/1.1 200 OK
date: Sat, 06 Jul 2024 22:00:56 GMT
content-type: application/json
content-length: 302
server: envoy
access-control-allow-origin: *
access-control-allow-credentials: true
x-envoy-upstream-service-time: 853
```

### Create a Kyverno-Json policy

The following policy will block all other method requests expect `GET` request to the `httpbin.org` service.

policy.yaml

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: checkrequest
spec:
  rules:
    - name: deny-post-request
      assert:
        any:
        - message: "POST method calls are not allowed only GET call are allowed"
          check:
            request:
                http:
                    method: GET
```

Store the policy in Kubernetes as a secret.

```bash
$ kubectl create secret generic kyverno-policy --from-file=policy.yaml
```

### Setup a Kyverno-Envoy-Plugin

Create a deployment for the Kyverno-Envoy-Plugin using the below commands:

```bash
$ kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kyverno
  labels:
    app: kyverno
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kyverno
  template:
    metadata:
      labels:
        app: kyverno
    spec:
      containers:
      - name: kyverno
        image: sanskardevops/plugin:0.0.34
        imagePullPolicy: IfNotPresent
        ports:
          - containerPort: 8181
          - containerPort: 9191
        volumeMounts:
          - readOnly: true
            mountPath: /policy
            name: kyverno-policy
        args:
          - "serve"
          - "--policy=/policy/policy.yaml"
          - "--address=:9191"
          - "--healthaddress=:8181"
      volumes:
      - name: kyverno-policy
        secret:
          secretName: kyverno-policy
EOF          
```

Next, define a Kubernetes Service for Kyverno-Envoy-Plugin. This is required to create a DNS record and thereby create a Gloo Upstream object.

> Note: Since the name of the service port is `grpc`, `Gloo` will understand that traffic should be routed using HTTP2 protocol.

```bash
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: kyverno
  labels:
    app: kyverno
spec:
  ports:
  - name: grpc
    port: 9191
    targetPort: 9191
    protocol: TCP
  selector:
    app: kyverno
EOF    
```

### Configure Gloo Edge to use Kyverno-Envoy-Plugin

To use the Kyverno-Envoy-Plugin as a custom auth server, we need to configure Gloo Edge to use it by adding the `extauth` attribute to the `GatewayProxy` resource.

gloo.yaml

```yaml
global:
  extensions:
    extAuth:
      extauthzServerRef:
        name: gloo-system-kyverno-9191
        namespace: gloo-system
```

To apply it, run the following command:

```bash
$ helm upgrade --install --namespace gloo-system --create-namespace -f gloo.yaml gloo gloo/gloo
```

Then, configure Gloo Edge routes to perform authorization via configured extauth before regular processing.

vs-patch.yaml

```yaml
spec:
  virtualHost:
    options:
      extauth:
        customAuth: {}
```

Then apply the patch to our `VirtualService` as shown below:

```bash
$ kubectl patch vs httpbin --type=merge --patch "$(cat vs-patch.yaml)"
```

### Exercise the Kyverno-Envoy-Plugin with Kyverno-Json policy

After the patch is applied, let's verify that OPA allows only allows GET requests.

```bash 
$ curl -XGET -Is localhost:8080/get
HTTP/1.1 200 OK
date: Sat, 06 Jul 2024 23:07:26 GMT
content-type: application/json
content-length: 302
server: envoy
access-control-allow-origin: *
access-control-allow-credentials: true
x-envoy-upstream-service-time: 577
```

```bash
$ curl http -XPOST -Is localhost:8080/post 
HTTP/1.1 403 Forbidden
content-length: 192
content-type: text/plain
date: Sat, 06 Jul 2024 23:26:10 GMT
server: envoy
```
Check the logs of the Kyverno-Envoy-Plugin to see the policy is applied.
 
```bash
kubectl logs deployments/kyverno 
2024/07/06 22:35:49 Starting GRPC server on port :9191
2024/07/06 22:35:49 Starting HTTP health checks on port :8181
2024/07/06 22:38:22 Request is initialized in kyvernojson engine .
2024/07/06 22:38:22 Request passed the deny-post-request policy rule.
2024/07/06 23:26:10 Request is initialized in kyvernojson engine .
2024/07/06 23:26:10 Request violation: -> POST method calls are not allowed only GET call are allowed
 -> any[0].check.request.http.method: Invalid value: "POST": Expected value: "GET"
```