# Quick Start 

This section presumes testing is conducted with Envoy version 1.10.0 or newer.

### Required tools

1. [`minikube`](https://minikube.sigs.k8s.io/docs/)
1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Create a local cluster

Start minikube cluster with the following command:

```console
minikube start
```
### Install kyverno-envoy sidecar with application 

Install application with envoy and kyverno-envoy-plugin as a sidecar container.

```console
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/main/quick_start.yaml 
```
The `applicaition.yaml` manifest defines the following resource:

- The Deployment includes an example Go application that provides information of books in the library books collection and exposes APIs to `get`, `create` and `delete` books collection. Check this out for more information about the [Go test application](https://github.com/Sanskarzz/kyverno-envoy-demos/tree/main/test-application) . 

- The Deployment also includes a kyverno-envoy-plugin sidecar container in addition to the Envoy sidecar container. When Envoy recevies API request destined for the Go test applicaiton, it will check with kyverno-envoy-plugin to decide if the request should be allowed and the kyverno-envoy-plugin sidecar container is configured to query Kyverno-json engine for policy decisions on incoming requests.

- A ConfigMap `policy-config` is used to pass the policy to kyverno-envoy-plugin sidecar in the namespace `default` where the application is deployed .

- A ConfigMap `envoy-config` is used to pass an Envoy configuration with an External Authorization Filter to direct authorization checks to the kyverno-envoy-plugin sidecar. 

- The Deployment also includes an init container that install iptables rules to redirect all container traffic to the Envoy proxy sidecar container , more about init container can be found [here](./envoy_iptables)

### Make Test application accessible in the cluster .

```console 
kubectl expose deployment testapp --type=NodePort --name=testapp-service --port=8080
```
 
### Set the `SERVICE_URL` environment variable to the service's IP/port.

minikube:

```sh
export SERVICE_PORT=$(kubectl get service testapp -o jsonpath='{.spec.ports[?(@.port==8080)].nodePort}')
export SERVICE_HOST=$(minikube ip)
export SERVICE_URL=$SERVICE_HOST:$SERVICE_PORT
echo $SERVICE_URL
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
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" http://$SERVICE_URL/book
```
```bash
curl -i -H "Authorization: Bearer "$ALICE_TOKEN"" -d '{"bookname":"Harry Potter", "author":"J.K. Rowling"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/book
```
Check the `Bob` which can get book also create the book 

```bash
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" http://$SERVICE_URL/book
```

```bash
curl -i -H "Authorization: Bearer "$BOB_TOKEN"" -d '{"bookname":"Harry Potter", "author":"J.K. Rowling"}' -H "Content-Type: application/json" -X POST http://$SERVICE_URL/book
```

Check on logs 
```bash
kubectl logs "$(kubectl get pod -l app=testapp -o jsonpath={.items..metadata.name})" -c kyverno-envoy-plugin -f
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

### Configuration 

To deploy Kyverno-Envoy include the following container in your kubernetes Deployments:

```yaml
- name: kyverno-envoy-plugin
  image: sanskardevops/plugin:0.0.34
  imagePullPolicy: IfNotPresent
  ports:
    - containerPort: 8181
    - containerPort: 9000
  volumeMounts:
    - readOnly: true
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
```