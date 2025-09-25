---
title: "kyverno-envoy-plugin serve sidecar-injector"
slug: "kyverno-envoy-plugin_serve_sidecar-injector"
description: "CLI reference for kyverno-envoy-plugin serve sidecar-injector"
---

## kyverno-envoy-plugin serve sidecar-injector

Start the Kubernetes mutating webhook injecting Kyverno Authz Server sidecars into pod containers

```
kyverno-envoy-plugin serve sidecar-injector [flags]
```

### Options

```
      --address string                       Address to listen on (default ":9443")
      --cert-file string                     File containing tls certificate
      --external-policy-source stringArray   External policy sources
  -h, --help                                 help for sidecar-injector
      --key-file string                      File containing tls private key
      --sidecar-image string                 Image to use in sidecar
```

### SEE ALSO

* [kyverno-envoy-plugin serve](kyverno-envoy-plugin_serve.md)	 - Run Kyverno Envoy Plugin servers

