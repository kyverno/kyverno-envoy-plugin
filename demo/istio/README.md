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
### Install kyverno-envoy-sidecar-injector admission controller 

First, we need to install the kyverno-envoy-sidecar-injector admission controller to inject kyverno-envoy-plugin sidecar into the sample-application pods(upstream pod which we need to authorization).

Flow this [README](./../../sidecar-injector/README.md) for installation of kyverno-envoy-sidecar-injector admission controller 

### Sample application

Manifests for the sample application are available in [test-application.yaml](manifests/test-application.yaml). The sample app provides information about books in a collection and exposes APIs to get, create and delete Book resources.

But first we need to apply [kyverno policy configmap](manifests/policy-config.yaml) this policy will be passed to kyverno-envoy-sidecar:

```console
kubectl apply -f ./manifests/namespace.yaml
```

```console
kubectl apply -f ./manifests/policy-config.yaml
```

```console
# deploy sample application
kubectl apply -f ./manifests/test-application.yaml      
```
Check that their should be three containers should be running in the pod.

```console
kubectl -n demo get all 
```
```bash
sanskar@sanskar-HP-Laptop-15s-du1xxx:~$ kubectl -n demo  get all
NAME                        READY   STATUS    RESTARTS   AGE
pod/echo-55c77757f4-w6979   3/3     Running   0          3h59m

NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/echo   ClusterIP   10.96.110.173   <none>        8080/TCP   4h5m

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/echo   1/1     1            1           3h59m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/echo-55c77757f4   1         1         1       3h59m

```

### Calling the sample application

We are going to call the sample application using a pod in the cluster.

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - testapp.demo.svc.cluster.local:8080/book
```
output 
```
[{"id":"1","bookname":"Harry Potter","author":"J.K. Rowling"},{"id":"2","bookname":"Animal Farm","author":"George Orwell"}]
pod "test" deleted
```

### ServiceEntry

ServiceEntry to registor the kyverno-envoy-plugin sidecar as external authorizer.

```console
kubectl apply -f ./manifests/service-entry.yaml
```
```yaml
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: kyverno-ext-authz-grpc-local
spec:
  hosts:
  - "kyverno-ext-authz-grpc.local"
  # The service name to be used in the extension provider in the mesh config.
  endpoints:
  - address: "127.0.0.1"
  ports:
  - name: grpc
    number: 9000
    # The port number to be used in the extension provider in the mesh config.
    protocol: GRPC
  resolution: STATIC
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
          service: "kyverno-ext-authz-grpc.local"
          port: "9000"
```

### Authorization policy

Now we can deploy an istio `AuthorizationPolicy`:
AuthorizationPolicy to tell Istio to use kyverno-envoy-plugin as the Authz Server

```console
kubectl apply -f ./manifests/authorizationpolicy.yaml
```

```yaml
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: kyverno-ext-authz-grpc
  namespace: demo
spec:
  action: CUSTOM
  provider:
    # The provider name must match the extension provider defined in the mesh config.
    name: kyverno-ext-authz-grpc
  rules:
  # The rules specify when to trigger the external authorizer.
  - to:
    - operation:
        notPaths: ["/healthz"]
    # Allowed all path except /healthz
```

This policy configures an external service for authorization. Note that the service is not specified directly in the policy but using a `provider.name` field.

### Verify the authorization 

For convenience, we’ll want to store Alice’s and Bob’s tokens in environment variables. Here bob is assigned the admin role and alice is assigned the guest role.

```bash 
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"
```

In this policy , we are using the Kyverno JSON API to define a Validating which we have already applied to the cluster. 
The policy checks the conditions of the incoming request and denies the request if the user is a guest and the request method is POST at the /book path.

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

Here is the Comman format of CheckRequest payload, Envoy transmits a [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) in Protobuf format to an external authorization service (which is our kyverno-envoy-plugin) for making access control decisions. This payload is then converted into a JSON format (inside kyverno-envoy-plugin) and evaluated against the defined policy within the Kyverno JSON engine.

```json
{
  "source": {
    "address": {
      "socketAddress": {
        "address": "10.244.1.10",
        "portValue": 59252
      }
    }
  },
  "destination": {
    "address": {
      "socketAddress": {
        "address": "10.244.1.4",
        "portValue": 8080
      }
    }
  },
  "request": {
    "time": "2024-04-09T07:42:29.634453Z",
    "http": {
      "id": "14694995155993896575",
      "method": "GET",
      "headers": {
        ":authority": "testapp.demo.svc.cluster.local:8080",
        ":method": "GET",
        ":path": "/book",
        ":scheme": "http",
        "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk",
        "user-agent": "Wget",
        "x-forwarded-proto": "http",
        "x-request-id": "27cd2724-e0f4-4a69-a1b1-9a94edfa31bb"
      },
      "path": "/book",
      "host": "echo.demo.svc.cluster.local:8080",
      "scheme": "http",
      "protocol": "HTTP/1.1"
    }
  },
  "metadataContext": {},
  "routeMetadataContext": {}
}
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
kubectl logs "$(kubectl get pod -l app=testapp -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f

```