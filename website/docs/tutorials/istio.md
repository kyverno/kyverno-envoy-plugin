# Istio 

[Istio](https://istio.io/latest/) is an open source service mesh for managing the different microservices that make up a cloud-native application. Istio provides a mechanism to use a service as an external authorizer with the [AuthorizationPolicy API](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/).

This tutorial shows how Istio’s AuthorizationPolicy can be configured to delegate authorization decisions to Kyverno-envoy-plugin.

## Prerequisites

This tutorial requires Kubernetes 1.20 or later. To run the tutorial locally ensure you start a cluster with Kubernetes version 1.20+, we recommend using [minikube](https://kubernetes.io/docs/getting-started-guides/minikube) or [KIND](https://kind.sigs.k8s.io/).

The tutorial also requries istio v1.19.0 or later. To install istio, follow the instructions [here](https://istio.io/latest/docs/setup/getting-started/) or run the below script it will create a kind cluster and install istio

```sh
#!/bin/bash

KIND_IMAGE=kindest/node:v1.29.2
ISTIO_REPO=https://istio-release.storage.googleapis.com/charts
ISTIO_NS=istio-system

# Create Kind cluster
kind create cluster --image $KIND_IMAGE --wait 1m --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |-
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
        protocol: TCP
      - containerPort: 443
        hostPort: 443
        protocol: TCP
  - role: worker
EOF

# Install Istio components
helm upgrade --install istio-base       --namespace $ISTIO_NS           --create-namespace --wait --repo $ISTIO_REPO base
helm upgrade --install istiod           --namespace $ISTIO_NS           --create-namespace --wait --repo $ISTIO_REPO istiod

```
The tutorial requires admission controller in the `kyverno-envoy-sidecar-injector` namespace that automatically injects the kyverno-envoy-plugin sidecar into pods in namespaces labelled with `kyverno-envoy-sidecar/injection=enabled`. To install the sidecar-injector admission controller then checkout the [installation guide](https://github.com/kyverno/kyverno-envoy-plugin/tree/main/sidecar-injector).

### Creating a simple authorization policy

This tutorial assumes you have some basic knowledge of [validatingPolicy](https://kyverno.github.io/kyverno-json/latest/policies/policies/#policy-structure) and [assertion trees](https://kyverno.github.io/kyverno-json/latest/policies/asserts/). In summary the policy below does the following:

- Checks that the JWT token is valid
- Checks that the action is allowed based on the token payload `role` and the request path
- Guests have read-only access to the `/book` endpoint, admins can create users too as long as the name is not the same as the admin's name.

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
### Deploying the application

Create a namespace called demo and label it with `istio-injection=enabled` to enable sidecar injection:

```shell
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: demo
  labels:
    istio-injection: enabled
EOF
```

First we need to apply kyverno policy configmap this policy will be passed to kyverno-envoy-sidecar injector admission controller:

```shell
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: policy-files
  namespace: demo
data:
  policy.yaml: |
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
EOF                   
```

Deploy the sample application which provides information about books in a collection and exposes APIs to get, create and delete Book resources at `/book` endpoint and make it accessible in the cluster, and enable sidecar injection of the kyverno-envoy-plugin sidecar by adding the `kyverno-envoy-sidecar/injection: enabled` label to the deployment:

```shell
$ kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: testapp
  namespace: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: testapp
  template:
    metadata:
      labels:
        kyverno-envoy-sidecar/injection: enabled
        app: testapp
    spec:
      containers:
      - name: testapp
        image: sanskardevops/test-application:0.0.1
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: testapp
  namespace: demo
spec:
  type: ClusterIP
  selector:
    app: testapp
  ports:
  - port: 8080
    targetPort: 8080
EOF
```
Check that their should be three containers should be running in the pod.

```shell
$ kubectl -n demo  get all
NAME                        READY   STATUS    RESTARTS   AGE
pod/echo-55c77757f4-w6979   3/3     Running   0          3h59m

NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/echo   ClusterIP   10.96.110.173   <none>        8080/TCP   4h5m

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/echo   1/1     1            1           3h59m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/echo-55c77757f4   1         1         1       3h59m
```

### ServiceEntry

ServiceEntry to registor the kyverno-envoy-plugin sidecar as external authorizer and ServiceEntry to allow Istio to find the Kyverno-Envoy-Plugin sidecar.

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

AuthorizationPolicy to direct authorization checks to the Kyverno-Envoy-Plugin sidecar.

```shell
$ kubectl apply -f - <<EOF
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
    # Allowed all path except /healthz
    - operation:
        notPaths: ["/healthz"]
EOF 
```
This policy configures an external service for authorization. Note that the service is not specified directly in the policy but using a `provider.name` field.

### Verify the authorization 

For convenience, we’ll want to store Alice’s and Bob’s tokens in environment variables. Here bob is assigned the admin role and alice is assigned the guest role.

```bash 
export ALICE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
export BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6ImFkbWluIiwic3ViIjoiWVd4cFkyVT0ifQ.veMeVDYlulTdieeX-jxFZ_tCmqQ_K8rwx2OktUHv5Z0"
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

Check on kyverno-envoy-plugin container logs 
```bash
kubectl logs "$(kubectl get pod -l app=testapp -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
```

### Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure Istio’s EnvoyFilter to utilize the kyverno-envoy-plugin as an external authorization service.

Additionally, the tutorial provided an example policy using the kyverno-envoy-plugin that returns a boolean decision to determine whether a request should be permitted.

Further details about the tutorial can be found [here](https://github.com/kyverno/kyverno-envoy-plugin/tree/main/demo/istio).








