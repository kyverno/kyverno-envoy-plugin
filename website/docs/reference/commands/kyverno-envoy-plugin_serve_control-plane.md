---
title: "kyverno-envoy-plugin serve control-plane"
slug: "kyverno-envoy-plugin_serve_control-plane"
description: "CLI reference for kyverno-envoy-plugin serve control-plane"
---

## kyverno-envoy-plugin serve control-plane

Start the Kyverno authorizer control plane

```
kyverno-envoy-plugin serve control-plane [flags]
```

### Options

```
      --client-flush-interval duration          Interval for how often to remove dead client connections (default 3m0s)
      --grpc-address string                     Address to listen on (default ":9081")
      --grpc-network string                     Network to listen on (default "tcp")
  -h, --help                                    help for control-plane
      --initial-send-wait duration              Duration to wait before retrying a send to a client (default 5s)
      --kube-as string                          Username to impersonate for the operation
      --kube-as-group stringArray               Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --kube-as-uid string                      UID to impersonate for the operation
      --kube-certificate-authority string       Path to a cert file for the certificate authority
      --kube-client-certificate string          Path to a client certificate file for TLS
      --kube-client-key string                  Path to a client key file for TLS
      --kube-cluster string                     The name of the kubeconfig cluster to use
      --kube-context string                     The name of the kubeconfig context to use
      --kube-disable-compression                If true, opt-out of response compression for all requests to the server
      --kube-insecure-skip-tls-verify           If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
  -n, --kube-namespace string                   If present, the namespace scope for this CLI request
      --kube-password string                    Password for basic authentication to the API server
      --kube-proxy-url string                   If provided, this URL will be used to connect via proxy
      --kube-request-timeout string             The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
      --kube-server string                      The address and port of the Kubernetes API server
      --kube-tls-server-name string             If provided, this name will be used to validate server certificate. If this is not provided, hostname used to contact the server is used.
      --kube-token string                       Bearer token for authentication to the API server
      --kube-user string                        The name of the kubeconfig user to use
      --kube-username string                    Username for basic authentication to the API server
      --leader-election                         Enable leader election
      --leader-election-id string               Leader election ID
      --max-client-inactive-duration duration   Duration to wait before declaring a client as inactive (default 4m0s)
      --max-send-interval duration              Duration to wait before stopping attempts of sending a policy to a client (default 10s)
      --metrics-address string                  Address to listen on for metrics (default ":9082")
      --probes-address string                   Address to listen on for health checks (default ":9080")
```

### SEE ALSO

* [kyverno-envoy-plugin serve](kyverno-envoy-plugin_serve.md)	 - Run Kyverno Envoy Plugin servers

