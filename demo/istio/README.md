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
```

### Authorization policy

Now we can deploy an istio `AuthorizationPolicy`:

```
kubectl apply -f ./manifests/authorizationPolicy.yaml
```

```yaml
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: ext-authz
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
        paths: ["/foo"]
    - operation:    
        paths: ["/bar"]
```

This policy configures an external service for authorization. Note that the service is not specified directly in the policy but using a `provider.name` field. 
The `rules` specify that requests to paths `/foo` and `/bar` .

### Authorization service

Here is the deployment manifest of the ext-authz server , it require policy through configMap 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ext-authz
  namespace: demo 
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ext-authz 
  template:
    metadata:
      labels:
        app: ext-authz
    spec:
      containers:
      - image: sanskardevops/plugin:0.0.25
        imagePullPolicy: IfNotPresent
        name: ext-authz
        ports:
        - containerPort: 8000
        - containerPort: 9000
        args:
        - "serve"
        - "--policy=/policies/policy1.yaml"
        volumeMounts:
        - name: policy-files
          mountPath: /policies
      volumes:
      - name: policy-files
        configMap:
          name: policy-files
```

Apply the configMap which the following command:

```console
kubectl apply -f ./manifests/configmap.yaml
```
The policy allows `GET` method for path `/foo` only for alice not for bob.
Here is the ValidatingPolicy which is passed through configMap

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: check-Request
spec:
  rules:
    - name: deny-guest-request
      assert:
        all:
        - message: "GET method calls at path /foo are not allowed to guest"
          check:
            request:
                http:
                    method: GET
                    headers:
                        authorization:
                            (base64_decode(split(@, ' ')[1])): 
                                (split(@, ':')[0]): alice
                    path: /foo                              
```

The following command will deploy the kyverno external authorizer server:

```console
kubectl apply -f ./manifests/ext-authz.yaml
```
Verify the sample external authorizer is up and running:

```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
Starting HTTP server on Port 8000
Starting GRPC server on Port 9000

```

### Calling the sample application again

Calling the sample application again at the `/foo` path with the base64 encode authorization token of `alice` as a header ,  will succeed . 

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Basic YWxpY2U6cGFzc3dvcmQ=" --output-document - echo.demo.svc.cluster.local:8080/foo


{
  "path": "/foo",
  "headers": {
    "host": "echo.demo.svc.cluster.local:8080",
    "user-agent": "Wget",
    "x-forwarded-proto": "http",
    "x-request-id": "978ffdf4-28a3-4c97-83a9-059b110d625c",
    "x-b3-traceid": "db0f88897f788c540f752b8fc027b2b9",
    "x-b3-spanid": "0f752b8fc027b2b9",
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
    "hostname": "echo-6847f9f85-hdjzz"
  },
  "connection": {}
}pod "test" deleted
```
Check the log of the sample ext_authz server to confirm it was called .

```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
Starting HTTP server on Port 8000
Starting GRPC server on Port 9000
Request is initialized in kyvernojson engine .
2024/04/18 13:35:19 Request passed the deny-guest-request policy rule.

```

Calling the sample application again at the `/bar` path with the base64 encode authorization token of `bob` as a header  will be denied. 

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --header="authorization: Basic Ym9iOnBhc3N3b3Jk" --output-document - echo.demo.svc.cluster.local:8080/bar


wget: server returned error: HTTP/1.1 403 Forbidden
pod "test" deleted
pod default/test terminated (Error)
```

Check the log of the sample ext_authz server to confirm it was called twice. The first one was allowed and the second one was denied:

```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f

sanskar@sanskar-HP-Laptop-15s-du1xxx:~$ kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
Starting HTTP server on Port 8000
Starting GRPC server on Port 9000
Request is initialized in kyvernojson engine .
2024/04/18 13:35:19 Request passed the deny-guest-request policy rule.
Request is initialized in kyvernojson engine .
2024/04/18 13:36:48 Request violation: -> GET method calls at path /foo are not allowed to guest
 -> all[0].check.request.http.headers.authorization.(base64_decode(split(@, ' ')[1])).(split(@, ':')[0]): Invalid value: "bob": Expected value: "alice"
 -> all[0].check.request.http.path: Invalid value: "/bar": Expected value: "/foo"


```
