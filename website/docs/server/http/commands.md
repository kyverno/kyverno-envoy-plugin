# Commands

## Authz server

Run the authz server to handle HTTP policies and authenticate HTTP requests

```bash
/kyverno-authz-server serve http authz-server <FLAGS>
```

### Flags:

- `--probes-address` (string, default: `:9080`) - Address to listen on for health checks
- `--metrics-address` (string, default: `:9082`) - Address to listen on for metrics
- `--external-policy-source` (stringArray) - External policy sources
- `--image-pull-secret` (stringArray) - Image pull secrets
- `--allow-insecure-registry` (bool, default: `false`) - Allow insecure registry
- `--kube-policy-source` (bool, default: `true`) - Enable in-cluster kubernetes policy source
- `--leader-election` (bool, default: `false`) - Enable leader election
- `--leader-election-id` (string) - Leader election ID
- `--server-address` (string, default: `:9083`) - Address to serve the http authorization server on
- `--nested-request` (bool, default: `false`) - Expect the requests to validate to be in the body of the original request
- `--control-plane-reconnect-wait` (duration, default: `3s`) - Duration to wait before retrying connecting to the control plane
- `--control-plane-max-dial-interval` (duration, default: `8s`) - Duration to wait before stopping attempts of sending a policy to a client
- `--health-check-interval` (duration, default: `30s`) - Interval for sending health checks
- `--control-plane-address` (string) - Control plane address
- `--kube-*` - Kubernetes configuration override flags (see kubectl documentation)

## Validation webhook

Serve a validating admission webhook to validate if the incoming policies contain valid CEL expressions

```bash
/kyverno-authz-server serve http validation-webhook <FLAGS>
```

### Flags:

- `--probes-address` (string, default: `:9080`) - Address to listen on for health checks
- `--metrics-address` (string, default: `:9082`) - Address to listen on for metrics
