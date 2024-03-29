# Istio Demo 

This Istio Demo is prototype of the kyverno envoy plugin.

## Overview 

The goal of the demo to show user how kyverno-envoy-plugin will work with istio and how it can be used to enforce policies to the traffic between services. The Kyverno-envoy-plugin allows configuring these Envoy proxies to query Kyverno-json for policy decisions on incoming requests.

## Demo instructions

### Required tools

1. [`kind`](https://kind.sigs.k8s.io/)
1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
1. [`helm`](https://helm.sh/docs/intro/install/)

### Create a local cluster and install Istio

The [bootstrap.sh](bootstrap.sh) script contains everything needed to create a local cluster and install Istio.

```console
# create a local cluster and install istio
./bootstrap.sh
```

### Sample application

Manifests for the sample application are available in [sample-application.yaml](manifests/sample-application.yaml).

```console
# deploy sample application
kubectl apply -f ./manifests/sample-application.yaml
```

### Calling the sample application

We are going to call the sample application using a pod in the cluster.

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - echo.demo.svc.cluster.local:8080/foo

{
  "path": "/foo",
  "headers": {
    "host": "echo.demo.svc.cluster.local:8080",
    "user-agent": "Wget",
    "x-forwarded-proto": "http",
    "x-request-id": "1badcd84-75eb-4911-9835-b3588e3c5eee",
    "x-b3-traceid": "904f847c3db71758fa4076e48440800a",
    "x-b3-spanid": "fa4076e48440800a",
    "x-b3-sampled": "0"
  },
  "method": "GET",
  "body": "",
  "fresh": false,
  "hostname": "echo.demo.svc.cluster.local",
  "ip": "::ffff:127.0.0.6",
  "ips": [],
  "protocol": "http",
  "query": {},
  "subdomains": [
    "svc",
    "demo",
    "echo"
  ],
  "xhr": false,
  "os": {
    "hostname": "echo-6847f9f85-wbgbx"
  },
  "connection": {}
}
```

### Authorization policy

Now we can deploy an istio `AuthorizationPolicy`:

```console
# deploy authorisation policy
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: ext-authz
  namespace: demo
spec:
  action: CUSTOM
  provider:
    name: kyverno-ext-authz-http
  rules:
  - to:
    - operation:
        paths: ["/foo"]
EOF
```

This policy configures an external service for authorization. Note that the service is not specified directly in the policy but using a `provider.name` field.

The provider will be registered later in the istio config map.

### Calling the sample application again

Calling the sample application again at the `/foo` path will return `403 Forbidden`.

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - echo.demo.svc.cluster.local:8080/foo

wget: server returned error: HTTP/1.1 403 Forbidden
```

Note that calling another path (like `/bar`) succeeds as it's not part of the policy.

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - echo.demo.svc.cluster.local:8080/bar

{
  "path": "/bar",
  "headers": {
    "host": "echo.demo.svc.cluster.local:8080",
    "user-agent": "Wget",
    "x-forwarded-proto": "http",
    "x-request-id": "ca22cf4c-fd28-4dff-94a1-bc0611d710a4",
    "x-b3-traceid": "202ef8abae854851c12c033ff52252e4",
    "x-b3-spanid": "c12c033ff52252e4",
    "x-b3-sampled": "0"
  },
  "method": "GET",
  "body": "",
  "fresh": false,
  "hostname": "echo.demo.svc.cluster.local",
  "ip": "::ffff:127.0.0.6",
  "ips": [],
  "protocol": "http",
  "query": {},
  "subdomains": [
    "svc",
    "demo",
    "echo"
  ],
  "xhr": false,
  "os": {
    "hostname": "echo-6847f9f85-wbgbx"
  },
  "connection": {}
}
```

### Register authorization provider

Edit the mesh configmap to register authorization provider with the following command: 

```console  
kubectl edit configmap istio -n istio-system 
```

In the editor, add the extension provider definitions to the mesh configmap.

```yaml
  data:
    mesh: |-   
      extensionProviders:
      - name: "kyverno-ext-authz-grpc"
        envoyExtAuthzGrpc:
          service: "ext-authz.demo.svc.cluster.local"
          port: "9000"
      - name: "kyverno-ext-authz-http"
        envoyExtAuthzHttp:
          service: "ext-authz.demo.svc.cluster.local"
          port: "8000"
          includeRequestHeadersInCheck: ["x-ext-authz"]
```

### Authorization service

The following command will deploy the sample external authorizer which allows requests with the header `x-ext-authz: allow`:

```console
kubectl apply -n demo -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/extauthz/ext-authz.yaml

```
Verify the sample external authorizer is up and running:

```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz

2024/03/12 11:46:42 Starting gRPC server at [::]:9000
2024/03/12 11:46:42 Starting HTTP server at [::]:8000

```


### Calling the sample application again

Calling the sample application again at the `/foo` path with with header `x-ext-authz: allow` will succeed. 

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="x-ext-authz: allow" --output-document - echo.demo.svc.cluster.local:8080/foo

{
  "path": "/foo",
  "headers": {
    "host": "echo.demo.svc.cluster.local:8080",
    "user-agent": "Wget",
    "x-ext-authz": "allow",
    "x-forwarded-proto": "http",
    "x-request-id": "2ef1a0ce-6948-413e-a9a9-91c5b9242b5c",
    "x-ext-authz-check-result": "allowed",
    "x-ext-authz-check-received": "source:{address:{socket_address:{address:\"10.244.1.7\" port_value:52396}}} destination:{address:{socket_address:{address:\"10.244.1.3\" port_value:8080}}} request:{time:{seconds:1710245883 nanos:556386000} http:{id:\"15150282829336904450\" method:\"GET\" headers:{key:\":authority\" value:\"echo.demo.svc.cluster.local:8080\"} headers:{key:\":method\" value:\"GET\"} headers:{key:\":path\" value:\"/foo\"} headers:{key:\":scheme\" value:\"http\"} headers:{key:\"user-agent\" value:\"Wget\"} headers:{key:\"x-ext-authz\" value:\"allow\"} headers:{key:\"x-forwarded-proto\" value:\"http\"} headers:{key:\"x-request-id\" value:\"2ef1a0ce-6948-413e-a9a9-91c5b9242b5c\"} path:\"/foo\" host:\"echo.demo.svc.cluster.local:8080\" scheme:\"http\" protocol:\"HTTP/1.1\"}} metadata_context:{}",
    "x-ext-authz-additional-header-override": "grpc-additional-header-override-value",
    "x-b3-traceid": "ddc174607e9d88bf1830b48578b53e79",
    "x-b3-spanid": "1830b48578b53e79",
    "x-b3-sampled": "0"
  },
  "method": "GET",
  "body": "",
  "fresh": false,
  "hostname": "echo.demo.svc.cluster.local",
  "ip": "::ffff:127.0.0.6",
  "ips": [],
  "protocol": "http",
  "query": {},
  "subdomains": [
    "svc",
    "demo",
    "echo"
  ],
  "xhr": false,
  "os": {
    "hostname": "echo-6847f9f85-fg9pd"
  },
  "connection": {}
}pod "test" deleted
```

Calling the sample application again at the `/foo` path with with header `x-ext-authz: deny` will be denied. 

```console

kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="x-ext-authz: deny" --output-document - ec
ho.demo.svc.cluster.local:8080/foo

wget: server returned error: HTTP/1.1 403 Forbidden
pod "test" deleted
pod default/test terminated (Error)

```
Check the log of the sample ext_authz server to confirm it was called twice. The first one was allowed and the second one was denied:

```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f


2024/03/12 11:55:26 Starting HTTP server at [::]:8000
2024/03/12 11:55:26 Starting gRPC server at [::]:9000
2024/03/12 12:18:03 [gRPCv3][allowed]: echo.demo.svc.cluster.local:8080/foo, attributes: source:{address:{socket_address:{address:"10.244.1.7" port_value:52396}}} destination:{address:{socket_address:{address:"10.244.1.3" port_value:8080}}} request:{time:{seconds:1710245883 nanos:556386000} http:{id:"15150282829336904450" method:"GET" headers:{key:":authority" value:"echo.demo.svc.cluster.local:8080"} headers:{key:":method" value:"GET"} headers:{key:":path" value:"/foo"} headers:{key:":scheme" value:"http"} headers:{key:"user-agent" value:"Wget"} headers:{key:"x-ext-authz" value:"allow"} headers:{key:"x-forwarded-proto" value:"http"} headers:{key:"x-request-id" value:"2ef1a0ce-6948-413e-a9a9-91c5b9242b5c"} path:"/foo" host:"echo.demo.svc.cluster.local:8080" scheme:"http" protocol:"HTTP/1.1"}} metadata_context:{}
2024/03/12 12:18:37 [gRPCv3][denied]: echo.demo.svc.cluster.local:8080/foo, attributes: source:{address:{socket_address:{address:"10.244.1.8" port_value:45762}}} destination:{address:{socket_address:{address:"10.244.1.3" port_value:8080}}} request:{time:{seconds:1710245917 nanos:57648000} http:{id:"2185755048778078711" method:"GET" headers:{key:":authority" value:"echo.demo.svc.cluster.local:8080"} headers:{key:":method" value:"GET"} headers:{key:":path" value:"/foo"} headers:{key:":scheme" value:"http"} headers:{key:"user-agent" value:"Wget"} headers:{key:"x-ext-authz" value:"deny"} headers:{key:"x-forwarded-proto" value:"http"} headers:{key:"x-request-id" value:"007781a3-519e-400f-8562-cabf75e989c1"} path:"/foo" host:"echo.demo.svc.cluster.local:8080" scheme:"http" protocol:"HTTP/1.1"}} metadata_context:{}

```

## Architecture

The below architecture illustrates a scenario where no service mesh or Envoy-like components have been pre-installed or already installed.

![Architecture](architecture1.png)

The below architecture illustrates a scenario where a service mesh or Envoy-like components have been pre-installed or already installed.
![Architecture](architecture2.png)

## Requirements

- Istio Authorizationpolicy manifest  to add "extension provider " concept in MeshConfig to specify Where/how to talk to envoy ext-authz service 
-
