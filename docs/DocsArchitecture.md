# Overview 

This document provides a high-level overview of the architecture and design for the Kyverno Envoy plugin. 

## Architecture 

Building the authorization plugin requires three components:

- A sample upstream service: This service simulates real-world applications that the plugin will authorize it.
- An Envoy proxy service: This service acts as a mediator, routing traffic and enforcing authorization decisions made by the Kyverno authorization server.
- The Kyverno authorization service: This is the core of the plugin, responsible for making authorization decisions based on Kyverno policies.

Basically the Architecture of the plugin will be different in different scenario which means when cluster already uses application like Istio or gloo-edge where envoy is already deployed in the cluster in that case we will integrate our kyverno authz server with existing envoy proxy setup. 

In this explanation, we'll assume there are no other Envoy proxy solutions like Istio or Gloo Edge already deployed in the cluster. I've explained how to achieve authorization in a service mesh with Kyverno as the authorization server below.

![Architecture](demo/istio/architecture1.png)

In this design the application/service pod will include envoy sidecar and kyverno authz server sidecar 

### init container 

To redirect all container traffic throught the Envoy proxy sidecare, In this init container the [Istio Proxy init script](https://github.com/open-policy-agent/contrib/tree/main/envoy_iptables) will be executed .

### Envoy sidecar container 

Envoy (v1.7.0+) supports an External Authorization filter which calls an authorization service to check if the incoming request is authorized or not . [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) feature will help us to make a decision based on Kyverno policies . 

This sidecar will include Envoy cofiguration to configure kyverno authz server to make decision on the incoming request 

```yml
admin:
  access_log_path: /dev/stdout
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

static_resources:
  listeners:
  - name: listener1
    address:
      socket_address: { address: 0.0.0.0, port_value: 51051 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          stat_prefix: testsrv
          codec_type: AUTO
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route: { cluster: upstream }
          http_filters:
          - name: envoy.ext_authz
            typed_config:
              '@type': type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
              transport_api_version: V3
              failure_mode_allow: false
              grpc_service:
                envoy_grpc:
                  cluster_name: kyverno-ext-authz
              with_request_body:
                allow_partial_message: true
                max_request_bytes: 1024
                pack_as_bytes: true
          - name: envoy.filters.http.router
            typed_config:
              '@type': type.googleapis.com/envoy.extensions.filters.http.router.v3.Router

  clusters:
  - name: upstream-service
    type: STRICT_DNS
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: upstream-service
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: 8080

  - name: kyverno-ext-authz
    connect_timeout: 1.25s
    type: LOGICAL_DNS
    lb_policy: ROUND_ROBIN
    http2_protocol_options: {}
    load_assignment:
      cluster_name: kyverno-ext-authz
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: kyverno-ext-authz
                port_value: 9000
```

The external authorization filter calls an external gRPC service to check whether an incoming HTTP request is authorized or not. If the request is deemed unauthorized, then the request will be denied normally with 403 (Forbidden) response. 

The content of the requests that are passed to kyverno ext authz service is specified by [CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) and Envoy receives response as [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse).

### Kyverno External Authz server 

The plugin server extends kyverno-JSON with a GRPC server that implements the Envoy [External Authorization API](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter#config-http-filters-ext-authz) 

This Kyverno authz server will accept Arguments like multiple kyverno-json validation policy file or path to policy file , address and path of the grpc server as mentioned in envoy config 
 
This authz server will recive [service.auth.v3.CheckRequest](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkrequest) JSON payload 

The request info will look like this and all authentication info send by user also be present here in this payload only. 

```json
{
  "source": {...},
  "destination": {...},
  "request": {...},
  "context_extensions": {...},
  "metadata_context": {...},
  "route_metadata_context": {...},
  "tls_session": {...}
}
```

This CheckRequest payload has all information like request information , authentication information , path , method etc 

For example 
```json
{
  "source": {
    "address": {
      "socket_address": {
        "address": "10.244.0.10",
        "port_value": 43466
      }
    }
  },
  "destination": {
    "address": {
      "socket_address": {
        "address": "10.244.0.7",
        "port_value": 8080
      }
    }
  },
  "request": {
    "time": {
      "seconds": 1710826031,
      "nanos": 497170000
    },
    "http": {
      "id": "5368491147255081293",
      "method": "GET",
      "headers": {
        ":authority": "echo.demo.svc.cluster.local:8080",
        ":method": "GET",
        ":path": "/foo",
        ":scheme": "http",
        "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiZ3Vlc3QiLCJzdWIiOiJZV3hwWTJVPSIsIm5iZiI6MTUxNDg1MTEzOSwiZXhwIjoxOTQxMDgxNTM5fQ.rN_hxMsoQzCjg6lav6mfzDlovKM9azaAjuwhjq3n9r8",
        "user-agent": "Wget",
        "x-forwarded-proto": "http",
        "x-request-id": "ee4c2a8e-1c4b-46b9-b54a-df5e042fe652"
      },
      "path": "/foo",
      "host": "echo.demo.svc.cluster.local:8080",
      "scheme": "http",
      "protocol": "HTTP/1.1"
    }
  },
  "metadata_context": {},
  "route_metadata_context": {}
}
``` 

For evaluation we will use Kyverno-Json Go Api which is a way to embed the kyverno-JSON engine in Go program that validate JSON payloads using Kyverno policies.
This CheckRequest json payload will executed/scan against accepted Argument validation policy.yaml in the kyverno-json engine and response of the engine will be converted into [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) and CheckResponse return back to envoy if request is deemed unauthorized then request will be denied normally with 403(Forbidden) reponse if the request is authorized then the request will be accepted with status OK response then envoy will redirect request to upstream service .


### Kyverno Policy 


How we will access and restrict request information in the policy we will use builtin and kyverno function to make it work . 
To decode the session token and jwt token the following function will be used for finding users and roles . 

Function like
- split(string, string) 
      To split the strings like to split authorization token "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiZ3Vlc3QiLCJzdWIiOiJZV3hwWTJVPSIsIm5iZiI6MTUxNDg1MTEzOSwiZXhwIjoxOTQxMDgxNTM5fQ.rN_hxMsoQzCjg6lav6mfzDlovKM9azaAjuwhjq3n9r8" 

      
      split(request.http.authorization, " ") it will split the Bearer and token , space is used as a delimiter

- base64_decode(string) 
      To decode the encoded session based authentication tokens to produce user role   

- We need to support the jwt verification functions to decode and verify signature jwt token like opa has built in [Token verification functions](https://www.openpolicyagent.org/docs/latest/policy-reference/#tokens)    



For example 

Here the user with there role

`guest` role and can perform a `GET` request  .
`admin` role and can perform a  `GET` and `POST` request. 

The users will request to service with jwt token within the headers specifying their role 

```json
{
    "attributes": {
        "request": {
            "http": {
                "method": "GET",
                "headers": {
                    "authorization": "Bearer eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJleHAiOiAyMjQxMDgxNTM5LCAibmJmIjogMTUxNDg1MTEzOSwgInJvbGUiOiAiZ3Vlc3QiLCAic3ViIjogIllXeHBZMlU9In0.Uk5hgUqMuUfDLvBLnlXMD0-X53aM_Hlziqg3vhOsCc8"
                }
            }
        }
    }
} 

```

The policy should have features like spliting and decoding jwt token or session token 

```yml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: authz-policy
spec:
  rules:
  - assert:
      any:
      - check: 
          (attributes.request.http.headers.authorization) | split(@, ' ')[1] | jwtdecode(@,secret)[2].role == 'guest' 
          && (attributes.request.http.method) != 'GET'
      - check: 
          (attributes.request.http.headers.authorization) | split(@, ' ')[1] | jwtdecode(@,secret)[2].role == 'admin' 
          && (attributes.request.http.method) != 'POST'  
      - check: 
          (attributes.request.http.headers.authorization) | split(@, ' ')[1] | jwtdecode(@,secret)[2].role == 'admin' 
          && (attributes.request.http.method) != 'GET'       
          
```

When we decode the jwt token the payload data is 

```json
{
  "exp": 2241081539,
  "nbf": 1514851139,
  "role": "guest",
  "sub": "YWxpY2U="
}

```

To make it more user friendly experiance we need more builtin and custom function for matching and decoding which should be easy to use.

The policy evaluation result, either 'pass' or 'fail', will be embed into the [CheckResponse](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#service-auth-v3-checkresponse) message and sent back to Envoy. If policy result is `pass` then CheckResponse status is set to be `OK` or `200` then the incoming request will be allowed and if policy result `fail` then CheckResponse status is set to be `403` then the incoming request will be denied .


### Deployment of Upstream sevice with the sidecar containers 

Upstream App deployment with kyverno-envoy and Envoy as a sidecar 

Example The deployment be like 

```yml
apiversion: apps/v1
kind: Deployment
metadata:
  name: example-app
  labels:
    app: example-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: example-app
  template:
    metadata:
      labels:
        app: example-app
    spec:
      initContainers:
        - name: proxy-init
          image: istio/proxyinit:latest
          # Configure the iptables bootstrap script to redirect traffic to the
          # Envoy proxy on port 8000, specify that Envoy will be running as user
          # 1111, and that we want to exclude port 8282 from the proxy for the
          # kyverno health checks. These values must match up with the configuration
          # defined below for the "envoy" and "kyvermo" containers.
          args: ["-p", "8000", "-u", "1111", "-w", "8282"]
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
            runAsNonRoot: false
            runAsUser: 0
      containers:
        - name: app
          image: sanskardevops/testservice:latest # Upstream service
          ports:
            - containerPort: 8080
        - name: envoy
          image: envoyproxy/envoy:v1.20.0
          env:
            - name: ENVOY_UID
              value: "1111"
          volumeMounts:
            - readOnly: true
              mountPath: /config
              name: proxy-config
            - readOnly: false
              mountPath: /run/sockets
              name: emptyDir
          args:
            - "envoy"
            - "--log-level"
            - "debug"
            - "--config-path"
            - "/config/envoy.yaml"
        - name: kyverno-envoy
          image: sanskardevops/kyverno-envoy:0.0.1  #authorization service
          securityContext:
            runAsUser: 1111
          volumeMounts:
            - readOnly: true
              mountPath: /policy
              name: kyverno-policy
            - readOnly: false
              mountPath: /run/sockets
              name: emptyDir
          containerPort: 9002    
          args:
            - "serve"
            - "--policy=/policy/kyverno-policy.yaml"
            - "--address=localhost:9002"         
      volumes:
        - name: proxy-config
          configMap:
            name: proxy-config
        - name: kyverno-policy
          configMap:
            name: kyverno-policy
        - name: emptyDir
          emptyDir: {}
```

The configuration such as policy.yaml, and Envoy configuration will be provided through volume mounted `ConfigMaps` within the deployment.

```yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kyverno-policy
data:
  policy.rego: |
    apiVersion: json.kyverno.io/v1alpha1
    kind: ValidatingPolicy
    metadata:
      name: check-external-auth
    spec:
      rules:
        - name: 
          assert:
            all:
            - message: "GET calls are not allowed"
              check:
                (request.http.method == 'GET'): false
```

## Kyverno Authorization Server with Istio Service Mesh 

Istio is an open source service mesh for managing the different microservices that make up a cloud-native application. Istio provides a mechanism to use a service as an external authorizer with the [AuthorizationPolicy API](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/).

The Istio service mesh already uses the envoy proxy so we will integrate our kyverno authorization server with istio envoy proxy, for this external authorization server implementation istio provides feature in the Istio authorization policy using action field value set to be `CUSTOM` to delegate the access control to an external authorization system which will be our kyverno authorization server 

#### Deployment of external authorization server

Istio provides three type of deployment of external authorization server 
- Deploy External authorizer in a standalone pod in the mesh 
- Deploy External authorizer outside of the mesh 
- Deploy External authorizer as a separate container as a sidecar container in the same pod of the application which needs authorization 

For last two deployment of external authorizer we need to create a service entry resource to register the service to the mesh and make sure it is accessible to the proxy  

Here are the pros and cons of each Deployment types 

- Standalone Pod in the Mesh:
    Example of Deployment standalone Pod is show [Istio Demo](./demo/istio/README.md)
  - Pros:
    - Simpler deployment within the service and no need to register the service by service entry resource .
    - Easier to manage alongside istio components
  - Cons:
    - Adds another service to the mesh, increasing complexity.
    - Failure of the authorizer can impact overall mesh functionality. 
    - Scalability might be limited compared to external deployment. 

- External Deployment (Outside Mesh):
  - Pros:
    - Isolates authorization server from istio, improving resilience. 
    - Easier scaling of the authoricer independently.      
  - Cons: 
    - Increased network communication overhead for authorization checks.
    - Latency will be higher compared to deployemnt with sidecare container.
    - Potential management overhead for a separate service.

- Sidecar Container in the same Pod:
  - Pros:
    - Tight integration with the application needing authorization.
    - Minimizes network overhead for authorization checks and fastest on authorization checks as compared to others options.
    - Efficient resource utilization by sharing a pod.
  - Cons: 
    - Very complex installation of the sidecar. 
    - Failure of the application can impact authorization and vice versa.
  

Explaining Deployement of kyverno external authorization as sidecar container in same Pod 

- Consider tight coupling if the application and authorization logic are highly dependent and considering minimum network overhead or lowest latency for authorization checks. If we automate or improve the installation then this method of deployment of external authorization will be best way to deploy kyverno-envoy server as separate container in same pod or as sidecar container.

![Architecture](demo/istio/architecture2.png)

To automate or improve the installation 
  -  we can add a Mutate webhook admission controller to add/inject our sidecar container with the pod, if the pod configuration has annotation like `kyverno-envoy-injection=enabled` then the admission controller automatically inject the kyverno-envoy sidecar container into pods and opa also uses this type of admission controller which injects the sidecar 

  -  To build the mutating webhook admission controller for injecting sidecar container we can take reference from an open-source project [tumblr/k8s-sidecar-injector](https://github.com/tumblr/k8s-sidecar-injector) 

  -  This Configuration will be injected by admission controller
    ```yml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: test-injectionconfig1
      namespace: default
      labels
        app: k8s-sidecar-injector
    data:
      sidecar-v1: |
        name: kyverno-envoy
        containers: 
        - name: kyverno-envoy
          image: sanskardevops/kyverno-envoy:0.0.1  #authorization service
          securityContext:
            runAsUser: 1111
          volumeMounts:
            - readOnly: true
              mountPath: /policy
              name: kyverno-policy
          containerPort: 9002    
          args:
            - "serve"
            - "--policy=/policy/kyverno-policy.yaml"
            - "--address=localhost:9002"
        volumes:
        - name: kyverno-policy
          configMap:
            name: kyverno-policy
    ``` 
     
we need to define external authorizer that is allowed to be used in the mesh, so we need to define extension provider in the mesh config 

```
kubectl edit configmap istio -n istio-system
```
The following content will register external provider as kyverno external authorization server 

```yml
data:
  mesh: |-
    # Add the following content to define the external authorizers.
    extensionProviders:
    - name: "kyverno-ext-authz-grpc"
      envoyExtAuthzGrpc:
        service: "ext-authz.foo.svc.cluster.local"
        port: "9000"
    - name: "kyverno-ext-authz-http"
      envoyExtAuthzHttp:
        service: "ext-authz.foo.svc.cluster.local"
        port: "8000"
```
And we are deploying external authorizer as a separate container in the same pod of the application so we also need to create a service entry resource to register the service to the mesh and make it accessible to the proxy 

```yml
#Define the service entry for the local ext-authz service on port 9000.
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: kyverno-ext-authz-grpc
spec:
  hosts:
  - "ext-authz-grpc.local"
  endpoints:
  - address: "127.0.0.1"
  ports:
  - name: grpc
    number: 9000
    protocol: GRPC
  resolution: STATIC
```
Then we have to apply authorization policy with the `CUSTOM` action value. 
```yml

apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: ext-authz
  namespace: demo
spec:
  action: CUSTOM
  provider:
    # The provider name must match the extension provider defined in the mesh config.
    # You can also replace this with sample-ext-authz-http to test the other external authorizer definition.
    name: kyverno-ext-authz-grpc
  rules:
  # The rules specify when to trigger the external authorizer.
  - to:
    - operation:
        paths: ["/headers"]

```




