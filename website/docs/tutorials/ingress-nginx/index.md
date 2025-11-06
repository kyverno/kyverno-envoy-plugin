# Ingress Nginx

[Ingress NGINX](https://kubernetes.github.io/ingress-nginx/) is an open source Ingress controller for Kubernetes that uses [NGINX](https://www.nginx.com/) as a reverse proxy and load balancer. It provides a flexible and powerful way to manage external access to services in a Kubernetes cluster.

This tutorial shows how Ingress NGINX can be configured to delegate authorization decisions to the Kyverno Authz Server using the external authentication feature.

## Setup

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install Ingress NGINX and the Kyverno Authz Server
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Install Ingress NGINX

First we need to install Ingress NGINX in the cluster.

```bash
# install ingress-nginx
helm install ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  --wait \
  --repo https://kubernetes.github.io/ingress-nginx ingress-nginx \
  --set controller.service.type=ClusterIP
```

The `controller.service.type=ClusterIP` setting is used because the kind cluster created in the previous step doesn't come with load balancer support. For production environments or cloud providers with load balancer support, you can omit this setting or use `LoadBalancer`.

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

For more certificate management options, refer to [Certificates management](../../install/certificates.md).


### Deploy the Authz server

```bash
# deploy the kyverno authz server
helm install kyverno-authz-server                                             \
  --namespace kyverno --create-namespace                                      \
  --wait                                                                      \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server  \
  --set config.type=http                                                      \
  --set certManager.issuerRef.name=selfsigned-issuer \
  --set certManager.issuerRef.kind=ClusterIssuer \
  --set certManager.issuerRef.group=cert-manager.io
```

## Create a Kyverno ValidatingPolicy

In summary the policy below does the following:

- Is triggered only when the host is `myapp.com` and the path starts with `/api/v1`
- Fetches a secret word from an external service
- Allows GET requests with a matching secret header
- Allows POST requests with `application/json` content type
- Denies all other requests with `403`

```yaml
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: example-api
spec:
  evaluation:
    mode: HTTP
  matchConditions:
  - expression: |
      object.attributes.host == "myapp.com"
    name: host
  - expression: |
      object.attributes.path.startsWith("/api/v1")
    name: v1-api
  variables:
  - name: secretWord
    expression: |
      http.Get("http://my-server:3000").secretWord
  - name: secretHeader
    expression: |
      size(object.attributes.Header("secret-header")) > 0 ? object.attributes.Header("secret-header")[0] : ""
  - name: contentType
    expression: |
      size(object.attributes.Header("content-type")) > 0 ? object.attributes.Header("content-type")[0] : ""
  validations:
  - expression: |
      variables.secretHeader == variables.secretWord && object.attributes.method == "GET"
        ? http.Allowed().Response()
        : null
  - expression: |
      variables.contentType == "application/json" && object.attributes.method == "POST"
        ? http.Allowed().Response()
        : null
  - expression: |
      http.Denied("validations didnt pass").Response()
```

### Deploy the external service

The policy will fetch a secret word from an external service. Let's deploy it first.

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Println("got a request")
		resp := map[string]string{"secretWord": "my-secret-word"}
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Server listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
```

Deploy this service to your cluster:

```bash
# create a deployment and service for the external service
kubectl apply -n demo -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-server
  template:
    metadata:
      labels:
        app: my-server
    spec:
      containers:
      - name: server
        image: your-registry/my-server:latest
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: my-server
spec:
  selector:
    app: my-server
  ports:
  - port: 3000
    targetPort: 3000
EOF
```

### Deploy a sample application

Httpbin is a well-known application that can be used to test HTTP requests and helps to show quickly how we can play with the request and response attributes.

```bash
# create the demo namespace
kubectl create ns demo

# deploy the httpbin application
kubectl apply \
  -n demo \
  -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
```

### Create an Ingress with External Authentication

Now create a separate Ingress resource for `myapp.com` with external authentication enabled.

```yaml
# create ingress with external auth for myapp.com
kubectl apply -n demo -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp
  annotations:
    nginx.ingress.kubernetes.io/auth-method: POST
    nginx.ingress.kubernetes.io/auth-url: "http://kyverno-authz-server.kyverno.svc.cluster.local:9083/validate"
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.com
    http:
      paths:
      - path: /api/v1
        pathType: Prefix
        backend:
          service:
            name: httpbin
            port:
              number: 8000
EOF
```

The `nginx.ingress.kubernetes.io/auth-url` annotation points to `localhost:9083` because the Kyverno Authz Server sidecar is injected into the Ingress NGINX controller pod and runs locally on port 9083 (HTTP). The Ingress is configured for host `myapp.com` and path `/api/v1/*` to match the ValidatingPolicy conditions.

## Testing

At this point we have deployed and configured Ingress NGINX, the Kyverno Authz Server, a sample application, and the authorization policies.

### Port-forward to the Ingress controller

To access the Ingress without setting up DNS, port-forward to the Ingress NGINX controller:

```bash
kubectl port-forward -n ingress-nginx service/ingress-nginx-controller 8080:80
```

### Call into the sample application

Now we can send requests to the sample application and verify the result.

The policy requires requests to `myapp.com` with path `/api/v1/*`. Let's test different scenarios:

GET request with the correct secret header will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -H "Host: myapp.com" \
  -H "secret-header: my-secret-word" \
  localhost:8080/api/v1/get
```

GET request with wrong secret header will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -H "Host: myapp.com" \
  -H "secret-header: wrong-word" \
  localhost:8080/api/v1/get
```

POST request with JSON content type will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -X POST \
  -H "Host: myapp.com" \
  -H "Content-Type: application/json" \
  -d '{"data":"test"}' \
  localhost:8080/api/v1/post
```

POST request without JSON content type will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -X POST \
  -H "Host: myapp.com" \
  localhost:8080/api/v1/post
```

Request to wrong host will not trigger the policy:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -H "Host: wronghost.com" \
  localhost:8080/api/v1/get
```

## Alternative: Using Kubernetes Resources

The previous policy fetched data from an external HTTP service using the `http.Get()` function. You can also fetch data from Kubernetes resources like ConfigMaps using the `resource.Get()` function. These functions are part of the [Kyverno CEL libraries](https://kyverno.io/docs/policy-types/cel-libraries/).

### Create a ConfigMap with the secret word

```bash
# create configmap with secret word
kubectl apply -n demo -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: secret-word
data:
  secret-word: "my-k8s-secret"
EOF
```

### Create a policy that reads from ConfigMap

```yaml
# create policy that reads from configmap
kubectl apply -f - <<EOF
apiVersion: envoy.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: acme-api
spec:
  evaluation:
    mode: HTTP
  matchConditions:
  - expression: |
      object.attributes.host == "acme.corp"
    name: host
  - expression: |
      object.attributes.path.startsWith("/api/v1")
    name: v1-api
  variables:
  - name: secretWord
    expression: |
      resource.Get("v1", "configmaps", "demo", "secret-word").data["secret-word"]
  - name: secretHeader
    expression: |
      size(object.attributes.Header("secret-header")) > 0 ? object.attributes.Header("secret-header")[0] : ""
  - name: contentType
    expression: |
      size(object.attributes.Header("content-type")) > 0 ? object.attributes.Header("content-type")[0] : ""
  validations:
  - expression: |
      variables.secretHeader == variables.secretWord && object.attributes.method == "GET"
        ? http.Allowed().Response()
        : null
  - expression: |
      variables.contentType == "application/json" && object.attributes.method == "POST"
        ? http.Allowed().Response()
        : null
  - expression: |
      http.Denied("validations didnt pass").Response()
EOF
```

This policy is similar to the previous one, but fetches the secret word from a ConfigMap in the `demo` namespace instead of an external HTTP service.

### Create an Ingress for acme.corp

```yaml
# create ingress for acme.corp
kubectl apply -n demo -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: acme
  annotations:
    nginx.ingress.kubernetes.io/auth-method: POST
    nginx.ingress.kubernetes.io/auth-url: "http://kyverno-authz-server.kyverno.svc.cluster.local:9083/validate"
spec:
  ingressClassName: nginx
  rules:
  - host: acme.corp
    http:
      paths:
      - path: /api/v1
        pathType: Prefix
        backend:
          service:
            name: httpbin
            port:
              number: 8000
EOF
```

### Test the acme.corp policy

GET request with the correct secret header from ConfigMap will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -H "Host: acme.corp" \
  -H "secret-header: my-k8s-secret" \
  localhost:8080/api/v1/get
```

GET request with wrong secret header will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -H "Host: acme.corp" \
  -H "secret-header: wrong-secret" \
  localhost:8080/api/v1/get
```

POST request with JSON content type will return `200`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -X POST \
  -H "Host: acme.corp" \
  -H "Content-Type: application/json" \
  -d '{"data":"test"}' \
  localhost:8080/api/v1/post
```

POST request without JSON content type will return `403`:

```bash
curl -s -w "\nhttp_code=%{http_code}" \
  -X POST \
  -H "Host: acme.corp" \
  localhost:8080/api/v1/post
```

## Wrap Up

Congratulations on completing the tutorial!

This tutorial demonstrated how to configure Ingress NGINX to utilize the Kyverno Authz Server as an external authorization service.

Additionally, the tutorial provided an example policy that fetches data from an external service and validates requests based on headers and HTTP methods.
