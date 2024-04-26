# Standalone Envoy Demo 

This Standalone Envoy Demo is prototype of the kyverno envoy plugin.

## Overview 

The goal of the demo to show user how kyverno-envoy-plugin will work with standalone envoy and how it can be used to enforce policies to the traffic between services. The Kyverno-envoy-plugin allows configuring these Envoy proxies to query Kyverno-json for policy decisions on incoming requests. The kyverno-envoy-plugin si cofigured as a static binary and can be run as a sidecar container in the same pod as the application.

## Demo instructions

### Required tools

1. [`kind`](https://kind.sigs.k8s.io/)
1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)


### Create a local cluster

The [bootstrap.sh](bootstrap.sh) script contains everything needed to create a local cluster.

```console
./bootstrap.sh
```
### Architecture Overview 

### Install kyverno-envoy sidecar with application 

First, we need to create a namespace `demo` where we will deploy the application with envoy and kyverno-envoy-plugin as a sidecar containers.
 
```console
kubectl create namespace demo
```

Install application with envoy and kyverno-envoy-plugin as a sidecar container.

```console
kubectl apply -f ./manifests/application.yaml
```
The `applicaition.yaml` manifest defines the following resource:

- The Deployment includes an example Go application that provides information of books in the library books collection and exposes APIs to `get`, `create` and `delete` books collection. Check this out for more information about the [Go test application](https://github.com/Sanskarzz/kyverno-envoy-demos/tree/main/test-application) . 

- The Deployment also includes a kyverno-envoy-plugin sidecar container in addition to the Envoy sidecar container. When Envoy recevies API request destined for the Go test applicaiton, it will check with kyverno-envoy-plugin to decide if the request should be allowed and the kyverno-envoy-plugin sidecar container is configured to query Kyverno-json engine for policy decisions on incoming requests.

- A ConfigMap `policy-config` is used to pass the policy to kyverno-envoy-plugin sidecar in the namespace `demo` where the application is deployed .

- A ConfigMap `envoy-config` is used to pass an Envoy configuration with an External Authorization Filter to direct authorization checks to the kyverno-envoy-plugin sidecar. 

- The Deployment also includes an init container that install iptables rules to redirect all container traffic to the Envoy proxy sidecar container , more about init container can be found [here](./envoy_iptables)

### Make Test application accessible in the cluster .

```console 
kubectl apply -f ./manifests/service.yaml
```

### Calling the sample test application and verify the authorization 

For convenience, we’ll want to store Alice’s and Bob’s tokens in environment variables. Here bob is assigned the admin role and alice is assigned the guest role.

```bash 
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"
```

The policy we passed to kyverno-envoy-plugin sidecar in the ConfigMap `policy-config` is configured to check the conditions of the incoming request and denies the request if the user is a guest and the request method is `POST` at the `/book` path.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
    name: checkrequest
spec:
    rules:
    - name: deny-guest-request-at-post
        assert:
        any:
        - message: "POST method calls at path /book are not allowed to guests users"
            check:
            request:
                http:
                    method: POST
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): admin
                    path: /book                             
        - message: "GET method call is allowed to both guest and admin users"
            check:
            request:
                http:
                    method: GET
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): admin
                    path: /book 
        - message: "GET method call is allowed to both guest and admin users"
            check:
            request:
                http:
                    method: GET
                    headers:
                        authorization:
                            (split(@, ' ')[1]):
                                (jwt_decode(@ , 'secret').payload.role): guest
                    path: /book               
```

Check for `Alice` which can get book but cannot create book.

```bash
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Bearer "$ALICE_TOKEN"" --output-document - testapp.demo.svc.cluster.local:8080/book
```
```bash
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Bearer "$ALICE_TOKEN"" --post-data='{"bookname":"Harry Potter", "author":"J.K. Rowling"}' --output-document - testapp.demo.svc.cluster.local:8080/book
```
Check the `Bob` which can get book also create the book 

```bash
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Bearer "$BOB_TOKEN"" --output-document - testapp.demo.svc.cluster.local:8080/book
```

```bash
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Bearer "$BOB_TOKEN"" --post-data='{"bookname":"Harry Potter", "author":"J.K. Rowling"}' --output-document - testapp.demo.svc.cluster.local:8080/book
```

Check on logs 
```bash
kubectl logs "$(kubectl get pod -l app=testapp -n demo -o jsonpath={.items..metadata.name})" -n demo -c kyverno-envoy-plugin -f
```
First , third and last request is passed but second request is failed.

```console 
sanskar@sanskar-HP-Laptop-15s-du1xxx:~$ kubectl logs "$(kubectl get pod -l app=testapp -n demo -o jsonpath={.items..metadata.name})" -n demo -c kyverno-envoy-plugin -f
Starting HTTP server on Port 8000
Starting GRPC server on Port 9000
Request is initialized in kyvernojson engine .
2024/04/26 17:11:42 Request passed the deny-guest-request-at-post policy rule.
Request is initialized in kyvernojson engine .
2024/04/26 17:22:11 Request violation: -> POST method calls at path /book are not allowed to guests users
 -> any[0].check.request.http.headers.authorization.(split(@, ' ')[1]).(jwt_decode(@ , 'secret').payload.role): Invalid value: "guest": Expected value: "admin"
-> GET method call is allowed to both guest and admin users
 -> any[1].check.request.http.headers.authorization.(split(@, ' ')[1]).(jwt_decode(@ , 'secret').payload.role): Invalid value: "guest": Expected value: "admin"
 -> any[1].check.request.http.method: Invalid value: "POST": Expected value: "GET"
-> GET method call is allowed to both guest and admin users
 -> any[2].check.request.http.method: Invalid value: "POST": Expected value: "GET"
Request is initialized in kyvernojson engine .
2024/04/26 17:23:13 Request passed the deny-guest-request-at-post policy rule.
Request is initialized in kyvernojson engine .
2024/04/26 17:23:55 Request passed the deny-guest-request-at-post policy rule.
```