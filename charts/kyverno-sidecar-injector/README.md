# kyverno-sidecar-injector

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

Install `kyverno-sidecar-injector` Helm chart:

```shell
helm install kyverno-sidecar-injector --namespace kyverno --create-namespace kyverno-envoy-plugin/kyverno-sidecar-injector
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `nil` | Override the name of the chart |
| fullnameOverride | string | `nil` | Override the expanded name of the chart |
| rbac.create | bool | `true` | Create RBAC resources |
| rbac.serviceAccount.name | string | `nil` | The ServiceAccount name |
| rbac.serviceAccount.annotations | object | `{}` | Annotations for the ServiceAccount |
| certificates.static | object | `{}` | Static data to set in certificate secret |
| certificates.certManager | object | `{}` | Infos for creating certificate with cert manager |
| deployment.replicas | int | `nil` | Desired number of pods |
| deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| deployment.annotations | object | `{}` | Deployment annotations. |
| deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| pod.labels | object | `{}` | Additional labels to add to each pod |
| pod.annotations | object | `{}` | Additional annotations to add to each pod |
| pod.imagePullSecrets | list | `[]` | Image pull secrets |
| pod.securityContext | object | `{}` | Security context |
| pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| pod.tolerations | list | `[]` | List of node taints to tolerate |
| pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| pod.priorityClassName | string | `""` | Optional priority class |
| pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| pod.affinity | object | `{}` | Pod affinity constraints. |
| pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| containers.injector.image.registry | string | `"ghcr.io"` | Image registry |
| containers.injector.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| containers.injector.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| containers.injector.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| containers.injector.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| containers.injector.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| containers.injector.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| containers.injector.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.injector.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.injector.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.injector.ports | list | `[{"containerPort":9443,"name":"https","protocol":"TCP"}]` | Container ports. |
| containers.injector.args | list | `["serve","sidecar-injector","--address=:9443","--cert-file=/opt/kubernetes-sidecar-injector/certs/tls.crt","--key-file=/opt/kubernetes-sidecar-injector/certs/tls.key","--config-file=/opt/kubernetes-sidecar-injector/config/sidecar.yaml"]` | Container args. |
| service.port | int | `443` | Service port. |
| service.type | string | `"ClusterIP"` | Service type. |
| service.nodePort | string | `nil` | Service node port. Only used if `type` is `NodePort`. |
| service.annotations | object | `{}` | Service annotations. |
| webhook.annotations | object | `{}` | Webhook annotations |
| webhook.failurePolicy | string | `"Fail"` | Webhook failure policy |
| webhook.objectSelector | string | `nil` | Webhook object selector |
| webhook.namespaceSelector | object | `{"matchExpressions":[{"key":"kyverno-injection","operator":"In","values":["enabled"]}]}` | Webhook namespace selector |
| sidecar.name | string | `"kyverno-authz-server"` | Sidecar container name |
| sidecar.image.registry | string | `"ghcr.io"` | Image registry |
| sidecar.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| sidecar.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| sidecar.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| sidecar.externalPolicySources | list | `[]` | External policy sources |
| sidecar.volumes | list | `[]` | Additional volumes |
| sidecar.volumeMounts | list | `[]` | Additional sidecar container volume mounts |
| sidecar.imagePullSecrets | list | `[]` | Additional image pull secrets |
| sidecar.config.grpcNetwork | string | `"tcp"` | GRPC network type (tcp, unix, etc.) |
| sidecar.config.httpAuthServerAddress | string | `":9083"` | HTTP authorization server address |
| sidecar.config.allowInsecureRegistry | bool | `false` | Allow insecure registry for pulling policy images |
| sidecar.config.nestedRequest | bool | `false` | Expect the requests to validate to be in the body of the original request |
| sidecar.config.imagePullSecrets | list | `[]` | Image pull secrets for fetching policies from OCI registries |
| sidecar.config.controlPlane.address | string | `""` | Control plane address (required for sidecar mode) |
| sidecar.config.controlPlane.reconnectWait | string | `"3s"` | Duration to wait before retrying connecting to the control plane |
| sidecar.config.controlPlane.maxDialInterval | string | `"8s"` | Duration to wait before stopping attempts of sending a policy to a client |
| sidecar.config.controlPlane.healthCheckInterval | string | `"30s"` | Interval for sending health checks |
| crds.install | bool | `true` |  |

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
