---
title: "kyverno-envoy-plugin serve sidecar-injector"
slug: "kyverno-envoy-plugin_serve_sidecar-injector"
description: "CLI reference for kyverno-envoy-plugin serve sidecar-injector"
---

## kyverno-envoy-plugin serve sidecar-injector

Start the Kubernetes mutating webhook injecting Kyverno Authorizer sidecars into pod containers

```
kyverno-envoy-plugin serve sidecar-injector [flags]
```

### Options

```
      --address string                           Address to listen on (default ":9443")
      --cert-file string                         File containing tls certificate
      --control-plane-address string             The control plane address to inject into the sidecars
      --control-plane-max-dial-interval string   Duration to wait before stopping attempts of sending a policy to a client (default "8s")
      --control-plane-reconnect-wait string      Duration to wait before retrying connecting to the control plane (default "3s")
      --health-check-interval string             Interval for sending health checks (default "30s")
  -h, --help                                     help for sidecar-injector
      --key-file string                          File containing tls private key
      --sidecar-image string                     Image to use in sidecar
```

### SEE ALSO

* [kyverno-envoy-plugin serve](kyverno-envoy-plugin_serve.md)	 - Run Kyverno Envoy Plugin servers

