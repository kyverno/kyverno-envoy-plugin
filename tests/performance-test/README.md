# Performance Testing for kyverno-envoy-plugin

This document outlines the approach for performance testing the `kyverno-envoy-plugin`, which is a plugin for Envoy that enforces Kyverno policies on API requests.

## Benchmark Setup

The benchmark setup consists of the following components:

### Sample Application

The first component is a simple Go application that provides information of books in the library books collection and exposes APIs to `get`, `create` and `delete` books collection. Check this out for more information about the [Go test application](https://github.com/Sanskarzz/kyverno-envoy-demos/tree/main/test-application) . 

### Envoy

The second component is the Envoy proxy, which runs alongside the example application. The Envoy configuration defines an external authorization filter `envoy.ext_authz` for a gRPC authorization server. The config uses Envoy's in-built gRPC client to make external gRPC calls.

```
static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 8000
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
                  target_uri: 127.0.0.1:9191
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

### Kyverno-envoy-plugin

The third component is the `kyverno-envoy-plugin` itself, which is configured to load and enforce Kyverno policies on incoming requests. checkout for [kyverno-envoy-plugin](application.yaml)

## Benchmark Scenarios

The following scenarios should be tested to compare the performance of the `kyverno-envoy-plugin` under different conditions:

1. **App Only**: Requests are sent directly to the application, without Envoy or the `kyverno-envoy-plugin`.
2. **App and Envoy**: Envoy is included in the request path, but the `kyverno-envoy-plugin` is not (i.e., Envoy External Authorization API is disabled).
3. **App, Envoy, and Kyverno (RBAC policy)**: Envoy External Authorization API is enabled, and a sample real-world RBAC policy is loaded into the `kyverno-envoy-plugin`.

## Load Testing with k6

To perform load testing, we'll use the k6 tool. Follow these steps:

1. **Install k6**: Install k6 on your machine by following the instructions on the official website: https://k6.io/docs/getting-started/installation/

2. **Write the k6 script**:  An example script is provided in the repository [k6-script.js](k6-script.js)

```js
import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // Ramp-up to 100 virtual users over 30 seconds
    { duration: '1m', target: 100 }, // Stay at 100 virtual users for 1 minute
    { duration: '30s', target: 0 }, // Ramp-down to 0 virtual users over 30 seconds
  ],
};

const BASE_URL = 'minikube ip with sample application'; // Replace with your application URL

export default function () {
  group('GET /book with admin token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_admin_token' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  group('GET /book with guest token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_guest_token' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  group('POST /book with guest token', () => {
    const res = http.post(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_guest_token' },
    });
    check(res, {
      'is status 403': (r) => r.status === 403,
    });
  });

  sleep(1); // Sleep for 1 second between iterations
}
```

3. **Run the k6 test**: Run the load test with the following command:

```console
k6 run k6-script.yaml
```
4. **Analyze the results**: Generate an HTML report with detailed insight by running:

```console
k6 run --out html=report.html k6-script.js
```
5. ***Repeat for different scenarios**: 

- # App only
    Results 
    ```html

    ``` 

- # App and Envoy
    Results
    ```html

    ```

- # App, Envoy and Kyverno-envoy-plugin 
    Results
    ```html

    ```
## Measuring Performance

The following metrics should be measured to evaluate the performance impact of the `kyverno-envoy-plugin`:

- End-to-end latency
- Kyverno evaluation latency
- gRPC server handler latency
- Resource utilization (CPU, memory)


