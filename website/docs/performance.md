# Performance

This page offers guidance and best practices for benchmarking the performance of the kyverno-envoy-plugin, helping users understand the associated overhead. It outlines an example setup for conducting benchmarks, various benchmarking scenarios, and key metrics to capture for assessing the impact of the kyverno-envoy-plugin.

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

The third component is the `kyverno-envoy-plugin` itself, which is configured to load and enforce Kyverno policies on incoming requests. 

```yaml
containers:
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
```

## Benchmark Scenarios

The following scenarios should be tested to compare the performance of the `kyverno-envoy-plugin` under different conditions:

1. **App Only**: Requests are sent directly to the application, without Envoy or the `kyverno-envoy-plugin`.
2. **App and Envoy**: Envoy is included in the request path, but the `kyverno-envoy-plugin` is not (i.e., Envoy External Authorization API is disabled).
3. **App, Envoy, and Kyverno (RBAC policy)**: Envoy External Authorization API is enabled, and a sample real-world RBAC policy is loaded into the `kyverno-envoy-plugin`.

## Load Testing with k6

To perform load testing, we'll use the k6 tool. Follow these steps:

1. **Install k6**: Install k6 on your machine by following the instructions on the official website: https://k6.io/docs/getting-started/installation/

2. **Write the k6 script**:  Below is the example k6 script. 

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

/*
Replace ip for every scenerio
export SERVICE_PORT=$(kubectl -n demo get service testapp -o jsonpath='{.spec.ports[?(@.port==8080)].nodePort}')
export SERVICE_HOST=$(minikube ip)
export SERVICE_URL=$SERVICE_HOST:$SERVICE_PORT
echo $SERVICE_URL

http://192.168.49.2:31541

*/
const BASE_URL = 'http://192.168.49.2:31541'; 

export default function () {
  group('GET /book with guest token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  sleep(1); // Sleep for 1 second between iterations
}
```

3. **Run the k6 test**: Run the load test with the following command:

```shell
$ k6 run -f - <<EOF
import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // Ramp-up to 100 virtual users over 30 seconds
    { duration: '1m', target: 100 }, // Stay at 100 virtual users for 1 minute
    { duration: '30s', target: 0 }, // Ramp-down to 0 virtual users over 30 seconds
  ],
};


const BASE_URL = 'http://192.168.49.2:31700'; // Replace with your application URL 

export default function () {
  group('GET /book with guest token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  sleep(1); // Sleep for 1 second between iterations
}
EOF
```
4. **Analyze the results**: Generate an json report with detailed insight by running:

```console
k6 run --out json=report.json k6-script.js
```
5. ***Repeat for different scenarios**: 

- # App only 
    In this case , request are sent directly to the sample application ie no Envoy and Kyverno-plugin in the request path .
    For this run this command to apply the sample applicaition and then test with k6

    ```shell
    $ kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/main/tests/performance-test/manifest/app.yaml
    ```
    Results of the k6 when only application is applied
    ```bash

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: k6-script.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 2m30s max duration (incl. graceful stop):
              * default: Up to 100 looping VUs for 2m0s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


     █ GET /book with guest token

       ✓ is status 200

     checks.........................: 100.00% ✓ 9048      ✗ 0    
     data_received..................: 2.1 MB  18 kB/s
     data_sent......................: 2.6 MB  21 kB/s
     group_duration.................: avg=1.01ms   min=166.46µs med=775.01µs max=36ms    p(90)=1.72ms   p(95)=2.31ms  
     http_req_blocked...............: avg=15.08µs  min=1.55µs   med=6.54µs   max=4.09ms  p(90)=12.07µs  p(95)=15.25µs 
     http_req_connecting............: avg=4.58µs   min=0s       med=0s       max=1.57ms  p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=745.73µs min=103.06µs med=549.17µs max=35.88ms p(90)=1.26ms   p(95)=1.75ms  
       { expected_response:true }...: avg=745.73µs min=103.06µs med=549.17µs max=35.88ms p(90)=1.26ms   p(95)=1.75ms  
     http_req_failed................: 0.00%   ✓ 0         ✗ 9048 
     http_req_receiving.............: avg=119.69µs min=11.33µs  med=77.78µs  max=10.97ms p(90)=193.73µs p(95)=285.58µs
     http_req_sending...............: avg=41µs     min=6.96µs   med=31.12µs  max=2.39ms  p(90)=61.88µs  p(95)=78.15µs 
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=585.04µs min=75.52µs  med=407.87µs max=35.84ms p(90)=965.49µs p(95)=1.33ms  
     http_reqs......................: 9048    75.050438/s
     iteration_duration.............: avg=1s       min=1s       med=1s       max=1.06s   p(90)=1s       p(95)=1s      
     iterations.....................: 9048    75.050438/s
     vus............................: 2       min=2       max=100
     vus_max........................: 100     min=100     max=100


running (2m00.6s), 000/100 VUs, 9048 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  2m0s
    ``` 

- # App and Envoy
    In this case, Kyverno-envoy-plugin is not included in the path but Envoy is but Envoy External Authorization API disabled 
    For this run this command to apply the sample application with envoy.

    ```shell
    $ kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/main/tests/performance-test/manifest/app-envoy.yaml
    ```

    Results of k6 after applying sample-application with envoy.
    ```bash

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: k6-script.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 2m30s max duration (incl. graceful stop):
              * default: Up to 100 looping VUs for 2m0s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


     █ GET /book with guest token

       ✓ is status 200

     checks.........................: 100.00% ✓ 9031      ✗ 0    
     data_received..................: 2.5 MB  21 kB/s
     data_sent......................: 2.6 MB  21 kB/s
     group_duration.................: avg=2.66ms  min=457.22µs med=1.8ms   max=65.53ms p(90)=4.85ms   p(95)=6.58ms  
     http_req_blocked...............: avg=12.81µs min=1.52µs   med=5.98µs  max=2.41ms  p(90)=11.84µs  p(95)=13.9µs  
     http_req_connecting............: avg=3.82µs  min=0s       med=0s      max=2.34ms  p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=2.38ms  min=383.7µs  med=1.58ms  max=65.22ms p(90)=4.36ms   p(95)=5.92ms  
       { expected_response:true }...: avg=2.38ms  min=383.7µs  med=1.58ms  max=65.22ms p(90)=4.36ms   p(95)=5.92ms  
     http_req_failed................: 0.00%   ✓ 0         ✗ 9031 
     http_req_receiving.............: avg=136.3µs min=12.53µs  med=76.74µs max=12.75ms p(90)=183.23µs p(95)=272.91µs
     http_req_sending...............: avg=41.54µs min=6.58µs   med=28.1µs  max=4.15ms  p(90)=59.62µs  p(95)=74.85µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=2.2ms   min=349.23µs med=1.43ms  max=65.08ms p(90)=4.05ms   p(95)=5.52ms  
     http_reqs......................: 9031    74.825497/s
     iteration_duration.............: avg=1s      min=1s       med=1s      max=1.06s   p(90)=1s       p(95)=1s      
     iterations.....................: 9031    74.825497/s
     vus............................: 3       min=3       max=100
     vus_max........................: 100     min=100     max=100


running (2m00.7s), 000/100 VUs, 9031 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  2m0s
    ```

- # App, Envoy and Kyverno-envoy-plugin 
    In this case, performance measurements are observed with Envoy External Authorization API enabled and a sample real-world RBAC policy loaded into kyverno-envoy-plugin .
    For this apply this command to apply sample-application, envoy and kyverno-envoy-plugin

    ```shell
    $ kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/main/tests/performance-test/manifest/app-envoy-plugin.yaml
    ```

    Results of k6 after applying sample-application, Envoy and kyverno-envoy-plugin . 
    ```console

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: k6-script.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 2m30s max duration (incl. graceful stop):
              * default: Up to 100 looping VUs for 2m0s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


     █ GET /book with guest token

       ✓ is status 200

     checks.........................: 100.00% ✓ 8655      ✗ 0    
     data_received..................: 2.4 MB  20 kB/s
     data_sent......................: 2.4 MB  20 kB/s
     group_duration.................: avg=46.54ms min=4.59ms  med=29.69ms max=337.79ms p(90)=109.35ms p(95)=140.51ms
     http_req_blocked...............: avg=11.88µs min=1.21µs  med=4.15µs  max=2.83ms   p(90)=9.87µs   p(95)=11.4µs  
     http_req_connecting............: avg=4.98µs  min=0s      med=0s      max=2.18ms   p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=46.37ms min=4.49ms  med=29.49ms max=337.69ms p(90)=109.26ms p(95)=140.28ms
       { expected_response:true }...: avg=46.37ms min=4.49ms  med=29.49ms max=337.69ms p(90)=109.26ms p(95)=140.28ms
     http_req_failed................: 0.00%   ✓ 0         ✗ 8655 
     http_req_receiving.............: avg=65.19µs min=11.14µs med=56.47µs max=5.58ms   p(90)=102.86µs p(95)=145.19µs
     http_req_sending...............: avg=30.35µs min=5.43µs  med=18.48µs max=5.29ms   p(90)=46.63µs  p(95)=58µs    
     http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=46.27ms min=4.43ms  med=29.42ms max=337.65ms p(90)=109.22ms p(95)=140.24ms
     http_reqs......................: 8655    71.999297/s
     iteration_duration.............: avg=1.04s   min=1s      med=1.03s   max=1.33s    p(90)=1.11s    p(95)=1.14s   
     iterations.....................: 8655    71.999297/s
     vus............................: 2       min=2       max=100
     vus_max........................: 100     min=100     max=100


running (2m00.2s), 000/100 VUs, 8655 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  2m0s
    ```
## Measuring Performance

The following metrics should be measured to evaluate the performance impact of the `kyverno-envoy-plugin`:

- **End-to-end latency**
  The end-to-end latency represents the time taken for a request to complete, from the client sending the request to receiving the response. Based on the k6 results, the average end-to-end latency for the different scenarios is as follows:

  - App Only: `avg=1.01ms` (from `group_duration` or `http_req_duration`)
  - App and Envoy: `avg=2.38ms` (from `http_req_duration`)
  - App, Envoy, and Kyverno-envoy-plugin: `avg=46.37ms` (from `http_req_duration`)

- **Kyverno evaluation latency**
  The Kyverno evaluation latency represents the time taken by the kyverno-envoy-plugin to evaluate the request against the configured policies. While the k6 results do not directly provide this metric, an estimate can be inferred by analyzing the differences in latency between the "App and Envoy" scenario and the "App, Envoy, and Kyverno-envoy-plugin" scenario.

  The difference in average latency between these two scenarios is:
  `46.37ms` - `2.38ms` = `43.99ms`

  This difference can be attributed to the Kyverno evaluation latency and the gRPC server handler latency combined. Assuming the gRPC server handler latency is relatively small compared to the Kyverno evaluation latency, the estimated range for the Kyverno evaluation latency is around 40ms to 45ms.

- **Resource utilization**
  Refers to CPU and memory usage of the Kyverno-Envoy-Plugin container , `kubectl top` utility can be laveraged to measure the resource utilization.

  Get the resource utilization of the kyverno-envoy-plugin container using the following command:

  ```shell
  $ kubectl top pod -n demo --containers
  ```

  To monitor resource utilization overtime use the following command:

  ```shell
  $ watch -n 1 "kubectl top pod -n demo --containers"
  ```

  Now run the k6 script in different terminal window and observe the resource utilization of the kyverno-envoy-plugin container.

  Initial resource utilization of the kyverno-envoy-plugin container:

  ```console
  POD                        NAME                   CPU(cores)   MEMORY(bytes)
  testapp-5955cd6f8b-dbvgd   envoy                  4m           70Mi
  testapp-5955cd6f8b-dbvgd   kyverno-envoy-plugin   1m           51Mi
  testapp-5955cd6f8b-dbvgd   test-application       1m           11Mi
  ```

  Resource utilization of the kyverno-envoy-plugin container after 100 requests:

  ```console
  POD                        NAME                   CPU(cores)   MEMORY(bytes)
  testapp-5955cd6f8b-dbvgd   envoy                  110m         70Mi
  testapp-5955cd6f8b-dbvgd   kyverno-envoy-plugin   895m         60Mi
  testapp-5955cd6f8b-dbvgd   test-application       17m          15Mi

  ```

  Observations:

  - The CPU utilization of the kyverno-envoy-plugin container increased significantly from 1m to 895m after receiving 100   requests during the load test.
  - The memory utilization also increased, but to a lesser extent, from 51Mi to 60Mi.

  Resource utilization of the kyverno-envoy-plugin container after load completion:

  ```console
  POD                        NAME                   CPU(cores)   MEMORY(bytes)
  testapp-5955cd6f8b-dbvgd   envoy                  4m           70Mi
  testapp-5955cd6f8b-dbvgd   kyverno-envoy-plugin   1m           51Mi
  testapp-5955cd6f8b-dbvgd   test-application       1m           11Mi
  ```

  Observations:
  - After the load test completed and the request volume returned to normal levels, the CPU and memory utilization of the kyverno-envoy-plugin container returned to their initial values. This indicates that the kyverno-envoy-plugin can efficiently handle the increased load during the test and release the additional resources when the load subsides.

  Correlation with k6 results:
  - The k6 script simulated a load test scenario with 100 virtual users, ramping up over 30 seconds, staying at 100 users for 1 minute, and then ramping down over 30 seconds.
  - During the load test, when the request volume was at its peak (100 virtual users), the kyverno-envoy-plugin container experienced a significant increase in CPU utilization, reaching 895m.
  - This CPU utilization spike aligns with the increased processing demand on the kyverno-envoy-plugin to evaluate the incoming requests against the configured Kyverno policies.
  - The memory utilization increase during the load test was relatively modest, suggesting that the policy evaluation did not significantly impact the memory requirements of the kyverno-envoy-plugin.



 

