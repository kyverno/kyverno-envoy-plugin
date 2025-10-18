# kyverno-authz-server-control-plane

Kyverno control plane for managing sidecar authorization servers

## Description

This Helm chart deploys the Kyverno control plane, which manages and distributes policies to sidecar authorization servers running in your cluster. The control plane watches for ValidatingPolicy resources and streams them to connected sidecar instances.

## Installation

```bash
helm install kyverno-authz-server-control-plane ./charts/kyverno-authz-server-control-plane
```

## Configuration

The following table lists the configurable parameters specific to the control plane:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `config.grpcNetwork` | GRPC network type (tcp, unix, etc.) | `tcp` |
| `config.initialSendWait` | Duration to wait before retrying a send to a client | `5s` |
| `config.maxSendInterval` | Duration to wait before stopping attempts of sending a policy to a client | `10s` |
| `config.clientFlushInterval` | Interval for how often to remove dead client connections | `180s` |
| `config.maxClientInactiveDuration` | Duration to wait before declaring a client as inactive | `240s` |

## Architecture

The control plane:
- Watches ValidatingPolicy resources in the cluster
- Maintains gRPC connections with sidecar authorization servers
- Streams policy updates to connected sidecars
- Monitors client health and removes inactive connections

## Requirements

- Kubernetes 1.25+
- Helm 3+
- ValidatingPolicy CRD (installed via kyverno-lib dependency)
