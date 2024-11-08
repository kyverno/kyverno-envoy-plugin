# Istio mTLS demo

This Istio mTLS Demo is a prototype of kyverno-envoy-plugin .

## Overview And Architecture 

Istio is an open source service mesh for managing the different microservices that make up a cloud-native application. Istio provides a mechanism to use a service as an external authorizer with the [AuthorizationPolicy API](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/). 

The kyverno-envoy-plugin is a custom Envoy filter that is used to intercept the incoming request to the service and validate the request using the kyverno engine. 

In this tutorial we will create a two simple microservices which are going to make external authorization to a single kyverno-envoy-plugin service in the mesh. With this tutorial we are going to understand how to use multiple microservices to make authorization decisions to a single ext-authz server. 

![arch-istio-mtls](arch-istio-mtls.png)

To handle multiple different requests effectively, we leverage the `match/exclude` declarations to route the specific authz-request to the appropriate validating policy within the Kyverno engine. This approach allows us to execute the right validating policy for each request, enabling efficient and targeted request processing.

### Example Policy

The following policies will be executed by the kyverno-envoy-plugin to validate incoming requests made specifically to the `testapp-1` service. By leveraging the match declarations, we ensure that these policies are executed only when the incoming request is destined for the `testapp-1` service. This targeted approach allows us to apply the appropriate validation rules and policies based on the specific service being accessed.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test-policy
spec:
  rules:
    - name: deny-external-calls-testapp-1
      match:
        any:
        - request:
            http:
                host: 'testapp-1.demo.svc.cluster.local:8080'
      assert:
        all:
        - message: "The GET method is restricted to the /book path."
          check:
            request:
                http:
                    method: 'GET'
                    path: '/book'
```
To execute the policy when the incoming request is made to `testapp-2` service we need to use the `match` declarations.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test-policy
spec:
  rules:
    - name: deny-external-calls-testapp-2
      match:
        any:
        - request:
            http:
                host: 'testapp-2.demo.svc.cluster.local:8080'
      assert:
        all:
        - message: "The GET method is restricted to the /movies path."
          check:
            request:
                http:
                    method: 'GET'
                    path: '/movie'   
```
The example json request for above payload will be like below.

```json
{
  "source": {
    "address": {
      "socketAddress": {
        "address": "10.244.0.71",
        "portValue": 33880
      }
    }
  },
  "destination": {
    "address": {
      "socketAddress": {
        "address": "10.244.0.65",
        "portValue": 8080
      }
    }
  },
  "request": {
    "time": "2024-05-20T07:52:01.566887Z",
    "http": {
      "id": "5415544797791892902",
      "method": "GET",
      "headers": {
        ":authority": "testapp-2.demo.svc.cluster.local:8080",
        ":method": "GET",
        ":path": "/movie",
        ":scheme": "http",
        "user-agent": "Wget",
        "x-forwarded-proto": "http",
        "x-request-id": "a3ad9f03-c9cd-4eab-97d1-83e90e0cee1b"
      },
      "path": "/movie",
      "host": "testapp-2.demo.svc.cluster.local:8080",
      "scheme": "http",
      "protocol": "HTTP/1.1"
    }
  },
  "metadataContext": {},
  "routeMetadataContext": {}
}
```

To enhance security, we can implement Mutual TLS (mTLS) for peer authentication between test services and kyverno-envoy-plugin. Since we are currently using JSON request data to validate incoming requests, there is a potential risk of this data being tampered with during transit. Implementing mTLS would ensure that communication between services is encrypted and authenticated, mitigating the risk of unauthorized data modification.

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: mtls-demo
  namespace: demo
spec:
  mtls:
    mode: STRICT
---
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: mtls-testapp-1
  namespace: demo
spec:
  selector:
    matchLabels:
      app: testapp-1
  mtls:
    mode: STRICT
  portLevelMtls:
    8080:
      mode: PERMISSIVE
---
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: mtls-testapp-2
  namespace: demo
spec:
  selector:
    matchLabels:
      app: testapp-2
  mtls:
    mode: STRICT
  portLevelMtls:
    8080:
      mode: PERMISSIVE
```
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
### Sample applications

Manifests for the sample applications are available in [test-application-1.yaml](manifests/test-application-1.yaml) and [test-application-2.yaml](manifests/test-application-2.yaml). The sample app `testapp-1` provides information about books in a collection and exposes APIs to get, create and delete Book resources. The sample app `testapp-2` provides information about movies in a collection and exposes APIs to get, create and delete Movie resources.

```console
# Create a namespace `demo`
kubectl apply -f ./manifests/namespace.yaml
```

```console
# deploy sample application testapp-1 and testapp-2
kubectl apply -f ./manifests/test-application-1.yaml     
kubectl apply -f ./manifests/test-application-1.yaml
```
### Calling the sample applications

We are going to call the sample applications using a pod in the cluster.

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - testapp-1.demo.svc.cluster.local:8080/book

[{"id":"1","bookname":"Harry Potter","author":"J.K. Rowling"},{"id":"2","bookname":"Animal Farm","author":"George Orwell"}]
pod "test" deleted

```
```console 
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - testapp-2.demo.svc.cluster.local:8080/movie

[{"id":"1","Moviename":"Inception","Actor":"Leonardo DiCaprio"},{"id":"2","Moviename":"Batman","Actor":"Jack Nicholson"}]
pod "test" deleted

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
        paths: ["/book","/movie"]
```

This policy configures an external service for authorization. Note that the service is not specified directly in the policy but using a provider.name field. The rules specify that requests to paths `/book` and `/movies`.

### Authorization service deployment 

The deployment manifest of the authorization service is available in [ext-auth-server.yaml](manifests/ext-auth-server.yaml). This deployment require policy through configmap .

Apply the policy configmap with the following command.

```console 
kubectl apply -f ./manifests/policy-configmap.yaml
```
```console
#Deploy the kyverno external authorizer server
kubectl apply -f ./manifests/ext-auth-server.yaml
```
Verify the sample external authorizer is up and running:
```console
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
Starting GRPC server on Port 9000
Starting HTTP server on Port 8000

```
### Apply PeerAuthentication Policy

Apply the PeerAuthentication policy to enable mTLS for the sample applications and external authorizer.

```console
kubectl apply -f ./manifests/peerAuthentication.yaml
```

### Test the sample applications

Check on the logs of the sample applications to see that the requests are accepted and rejected

Check on `GET` request on `testapp-1` which is allowed according to policy `deny-external-calls-testapp-1`

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - testapp-1.demo.svc.cluster.local:8080/book

[{"id":"1","bookname":"Harry Potter","author":"J.K. Rowling"},{"id":"2","bookname":"Animal Farm","author":"George Orwell"}]
pod "test" deleted
```

Check on `GET` request on `testapp-2` which is allowed according to policy `deny-external-calls-testapp-2`

```console
kubectl run test -it --rm --restart=Never --image=busybox -- wget -q --output-document - testapp-2.demo.svc.cluster.local:8080/movie

[{"id":"1","Moviename":"Inception","Actor":"Leonardo DiCaprio"},{"id":"2","Moviename":"Batman","Actor":"Jack Nicholson"}]
pod "test" deleted

```

Check logs of external authorizer to see that the requests are which policy was executed for a perticular request .

```console 
kubectl logs "$(kubectl get pod -l app=ext-authz -n demo -o jsonpath={.items..metadata.name})" -n demo -c ext-authz -f
Starting GRPC server on Port 9000
Starting HTTP server on Port 8000
2024/05/21 07:41:33 Request is initialized in kyvernojson engine .
2024/05/21 07:41:33 Request passed the deny-external-calls-testapp-1 policy rule.
2024/05/21 07:42:22 Request is initialized in kyvernojson engine .
2024/05/21 07:42:22 Request passed the deny-external-calls-testapp-2 policy rule.
```
First request was directed to testapp-1 which was allowed by the policy `deny-external-calls-testapp-1` and the second request was directed to testapp-2 which was allowed by the policy `deny-external-calls-testapp-2`.
  
