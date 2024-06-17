# Standalone Envoy

The tutorial shows how Envoy's [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) can be used with Kyverno as an authorization service to enforce security policies over API requests received by Envoy.

## Overview

In this tutorial we'll see how to use Kyverno-envoy-plugin as an External Authorization service for the Envoy proxy. The goal of the demo to show user how kyverno-envoy-plugin will work with standalone envoy and how it can be used to enforce policies to the traffic between services. The Kyverno-envoy-plugin allows configuring these Envoy proxies to query Kyverno-json for policy decisions on incoming requests. The kyverno-envoy-plugin is cofigured as a static binary and can be run as a sidecar container in the same pod as the application.

We'll do this by:

- Running a local Kubernetes cluster
- Creating a simple authorization policy in [ValidatingPolicy](https://kyverno.github.io/kyverno-json/latest/policies/policies/#api-group-and-kind) 
- Deploying a sample application with Envoy and kyverno-envoy-plugin sidecars
- Run some sample requests to see the policy in action

Note that other than the HTTP client and bundle server, all components are co-located in the same pod.

## Demo instructions

### Required tools

1. [`kind`](https://kind.sigs.k8s.io/)
2. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

{{< info >}}
If you haven't used `kind` before, you can find installation instructions
in the [project documentation](https://kind.sigs.k8s.io/#installation-and-usage).
{{</ info >}}

### Running a local Kubernetes cluster 

To start a local kubernetes cluster to run our demo, we'll be using [kind](https://kind.sigs.k8s.io/). In order to use the kind command, you’ll need to have Docker installed on your machine. 

Create a cluster with the following command:

```shell
$ kind create cluster --name kyverno-tutorial --image kindest/node:v1.29.2
Creating cluster "kyverno-tutorial" ...
 ✓ Ensuring node image (kindest/node:v1.29.2) 🖼
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-kyverno-tutorial"
You can now use your cluster with:

kubectl cluster-info --context kind-kyverno-tutorial

Thanks for using kind! 😊
```

Listing the cluster nodes, should show something like this:

```shell
$ kubectl get nodes
NAME                             STATUS   ROLES           AGE   VERSION
kyverno-tutorial-control-plane   Ready    control-plane   79s   v1.29.2
```

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

Create a file called policy.yaml with the above content and store it in a configMap:

```shell
$ kubectl create configmap policy --from-file=policy.yaml
```

### Deploying an application with Envoy and Kyverno-Envoy-Plugin sidecars

In this tutorial, we are manually configuring the Envoy proxy sidecar to intermediate HTTP traffic from clients and our application. Envoy will consult Kyverno-Envoy-Plugin to make authorization decisions for each request by sending `CheckRequest` gRPC messages over a gRPC connection.

We will use the following Envoy configuration to achieve this. In summary, this configures Envoy to:

- Listen on Port `7000` for HTTP traffic
- Consult Kyverno-Envoy-Plugin at `127.0.0.1:9000` for authorization decisions and deny failing requests
- Forward request to the application at `127.0.0.1:8080` if ok.

```yaml
    static_resources:
      listeners:
      - address:
          socket_address:
            address: 0.0.0.0
            port_value: 7000
        filter_chains:
        - filters:
          - name: envoy.filters.network.http_connection_manager
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
              codec_type: auto
              stat_prefix: ingress_http
              route_config:
                name: local_route
                virtual_hosts:
                - name: backend
                  domains:
                  - "*"
                  routes:
                  - match:
                      prefix: "/"
                    route:
                      cluster: service
              http_filters:
              - name: envoy.ext_authz
                typed_config:
                  "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
                  transport_api_version: V3
                  with_request_body:
                    max_request_bytes: 8192
                    allow_partial_message: true
                  failure_mode_allow: false
                  grpc_service:
                    google_grpc:
                      target_uri: 127.0.0.1:9000
                      stat_prefix: ext_authz
                    timeout: 0.5s
              - name: envoy.filters.http.router
                typed_config:
                  "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
      clusters:
      - name: service
        connect_timeout: 0.25s
        type: strict_dns
        lb_policy: round_robin
        load_assignment:
          cluster_name: service
          endpoints:
          - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: 127.0.0.1
                    port_value: 8080
    admin:
      access_log_path: "/dev/null"
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 8001
    layered_runtime:
      layers:
        - name: static_layer_0
          static_layer:
            envoy:
              resource_limits:
                listener:
                  example_listener_name:
                    connection_limit: 10000
            overload:
              global_downstream_max_connections: 50000
```

Create a `ConfigMap` containing the above configuration by running:

```shell
$ kubectl create configmap proxy-config --from-file envoy.yaml 
```
Our application will be configured using a `Deployment` and `Service`. There are few things to note:

- The pods have an `initContainer` that configures the `iptables` rules to redirect traffic to the Envoy Proxy sidecar.
- The `test-application` container is simple go application stores book information in-memory state.
- The `envoy` container is configured to use `proxy-config` `ConfigMap` as the Envoy configuration we created earlier
- The `kyverno-envoy-plugin` container is configured to use `policy` `ConfigMap` as the Kyverno policy we created earlier

```yaml
# test-application.yaml
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
        app: testapp
    spec:
      initContainers:
        - name: proxy-init
          image: sanskardevops/proxyinit:latest
          # Configure the iptables bootstrap script to redirect traffic to the
          # Envoy proxy on port 8000, specify that Envoy will be running as user
          # 1111, and that we want to exclude port 8181 from the proxy for the Kyverno health checks.
          # These values must match up with the configuration
          # defined below for the "envoy" and "kyverno-envoy-plugin" containers.
          args: ["-p", "7000", "-u", "1111", -w, "8181"]
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
            runAsNonRoot: false
            runAsUser: 0
      containers:
        - name: test-application
          image: sanskardevops/test-application:0.0.1
          ports:
            - containerPort: 8080
        - name: envoy
          image: envoyproxy/envoy:v1.30-latest
          securityContext:
            runAsUser: 1111
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - readOnly: true
              mountPath: /config
              name: proxy-config
          args:
            - "envoy"
            - "--config-path"
            - "/config/envoy.yaml"
        - name: kyverno-envoy-plugin
          image: sanskardevops/plugin:0.0.34
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8181
            - containerPort: 9000
          volumeMounts:
            - readOnly: true
              mountPath: /policies
              name: policy-files
          args:
            - "serve"
            - "--policy=/policies/policy.yaml"
            - "--address=:9000"
            - "--healthaddress=:8181"
          livenessProbe:
            httpGet:
              path: /health
              scheme: HTTP
              port: 8181
            initialDelaySeconds: 5
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: /health
              scheme: HTTP
              port: 8181
            initialDelaySeconds: 5
            periodSeconds: 5  
      volumes:
        - name: proxy-config
          configMap:
            name: proxy-config
        - name: policy-files
          configMap:
            name: policy-files
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
```

Deploy the application and Kubernetes Service to the cluster with:

```shell
$ kubectl apply -f test-application.yaml
```
Check that everything is working by listing the pod and make sure all three pods are running ok.

```shell
$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
testapp-74b4bc88-5d4wh       3/3     Running   0          1m
```
### Policy in action 

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

### Cleanup

Delete the cluster by running:
```shell
$ kind delete cluster --name kyverno-tutorial
```
## Wrap Up
Congratulations on completing the tutorial!

In this tutorial, you learned how to utilize the kyverno-envoy-plugin as an external authorization service to enforce custom policies through Envoy’s external authorization filter.

The tutorial also included an example policy using kyverno-envoy-plugin that returns a boolean decision indicating whether a request should be permitted.

Moreover, Envoy’s external authorization filter supports the inclusion of optional response headers and body content that can be sent to either the downstream client or upstream server. An example of a rule that not only determines request authorization but also provides optional response headers, body content, and HTTP status is available here.