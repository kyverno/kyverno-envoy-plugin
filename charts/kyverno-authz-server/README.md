# kyverno-authz-server

Kyverno policies based authorization plugin for Envoy ❤️

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A plugin to enforce kyverno policies with Envoy. This plugin allows applying Kyverno policies to APIs managed by Envoy.

## Overview

[Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/what_is_envoy) is a L7 proxy and communication bus designed for large modern service oriented architectures . Envoy (v1.7.0+) supports an External Authorization filter which calls an authorization service to check if the incoming request is authorized or not. [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) feature will help us to make a decision based on Kyverno policies .

## Installing the Chart

Add `kyverno-envoy-plugin` Helm repository:

```shell
helm repo add kyverno-json https://kyverno.github.io/kyverno-envoy-plugin/
```

Install `kyverno-authz-server` Helm chart:

```shell
helm install kyverno-authz-server --namespace kyverno --create-namespace kyverno-envoy-plugin/kyverno-authz-server
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `nil` | Override the name of the chart |
| fullnameOverride | string | `nil` | Override the expanded name of the chart |
| rbac.create | bool | `true` | Create RBAC resources |
| rbac.serviceAccount.name | string | `nil` | The ServiceAccount name |
| rbac.serviceAccount.annotations | object | `{}` | Annotations for the ServiceAccount |
| config.type | string | `"envoy"` | Authz server type (`envoy` or `http`) |
| config.grpc.network | string | `"tcp"` | GRPC network type (tcp, unix, etc.) |
| config.grpc.address | string | `":9081"` | GRPC address |
| config.http.address | string | `":9081"` | HTTP address |
| config.http.nestedRequest | bool | `true` | Expect the requests to validate to be in the body of the original request |
| config.sources.kube | bool | `true` | Enable in-cluster kubernetes policy source |
| config.sources.external | list | `[]` | External policy sources |
| config.sources.controlPlane.address | string | `""` | Control plane address (leave empty for standalone mode) |
| config.sources.controlPlane.reconnectWait | string | `"3s"` | Duration to wait before retrying connecting to the control plane |
| config.sources.controlPlane.maxDialInterval | string | `"8s"` | Duration to wait before stopping attempts of sending a policy to a client |
| config.sources.controlPlane.healthCheckInterval | string | `"30s"` | Interval for sending health checks |
| config.allowInsecureRegistry | bool | `false` | Allow insecure registry for pulling policy images |
| config.imagePullSecrets | list | `[]` | Image pull secrets for fetching policies from OCI registries |
| authzServer.deployment.replicas | int | `nil` | Desired number of pods |
| authzServer.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| authzServer.deployment.annotations | object | `{}` | Deployment annotations. |
| authzServer.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| authzServer.pod.labels | object | `{}` | Additional labels to add to each pod |
| authzServer.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| authzServer.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| authzServer.pod.securityContext | object | `{}` | Security context |
| authzServer.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| authzServer.pod.tolerations | list | `[]` | List of node taints to tolerate |
| authzServer.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| authzServer.pod.priorityClassName | string | `""` | Optional priority class |
| authzServer.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| authzServer.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| authzServer.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| authzServer.pod.affinity | object | `{}` | Pod affinity constraints. |
| authzServer.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| authzServer.container.image.registry | string | `"ghcr.io"` | Image registry |
| authzServer.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| authzServer.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| authzServer.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| authzServer.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| authzServer.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| authzServer.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| authzServer.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9081,"name":"authz","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"}]` | Container ports. |
| authzServer.service.annotations | object | `{}` | Service annotations. |
| authzServer.service.type | string | `"ClusterIP"` | Service type. |
| authzServer.service.port | int | `9081` | Service port. |
| authzServer.service.appProtocol | string | `nil` |  |
| authzServer.service.nodePort | string | `nil` |  |
| validatingWebhookConfiguration.annotations | object | `{}` | Webhook annotations |
| validatingWebhookConfiguration.certificates.static | object | `{}` | Static data to set in certificate secret |
| validatingWebhookConfiguration.certificates.certManager | object | `{}` | Infos for creating certificate with cert manager |
| validatingWebhookConfiguration.failurePolicy | string | `"Fail"` | Webhook failure policy |
| validatingWebhookConfiguration.deployment.replicas | int | `nil` | Desired number of pods |
| validatingWebhookConfiguration.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| validatingWebhookConfiguration.deployment.annotations | object | `{}` | Deployment annotations. |
| validatingWebhookConfiguration.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| validatingWebhookConfiguration.pod.labels | object | `{}` | Additional labels to add to each pod |
| validatingWebhookConfiguration.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| validatingWebhookConfiguration.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| validatingWebhookConfiguration.pod.securityContext | object | `{}` | Security context |
| validatingWebhookConfiguration.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| validatingWebhookConfiguration.pod.tolerations | list | `[]` | List of node taints to tolerate |
| validatingWebhookConfiguration.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| validatingWebhookConfiguration.pod.priorityClassName | string | `""` | Optional priority class |
| validatingWebhookConfiguration.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| validatingWebhookConfiguration.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| validatingWebhookConfiguration.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| validatingWebhookConfiguration.pod.affinity | object | `{}` | Pod affinity constraints. |
| validatingWebhookConfiguration.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| validatingWebhookConfiguration.container.image.registry | string | `"ghcr.io"` | Image registry |
| validatingWebhookConfiguration.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| validatingWebhookConfiguration.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| validatingWebhookConfiguration.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| validatingWebhookConfiguration.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| validatingWebhookConfiguration.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| validatingWebhookConfiguration.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| validatingWebhookConfiguration.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"},{"containerPort":9443,"name":"webhook","protocol":"TCP"}]` | Container ports. |
| crds.install | bool | `false` |  |

## Source Code

* <https://github.com/kyverno/kyverno-envoy-plugin>

## Requirements

Kubernetes: `>=1.25.0-0`

| Repository | Name | Version |
|------------|------|---------|
| file://../kyverno-lib | kyverno-lib | 0.0.0 |

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Kyverno |  | <https://kyverno.io/> |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
