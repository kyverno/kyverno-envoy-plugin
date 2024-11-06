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
| crds.install | bool | `true` | Whether to have Helm install the CRDs, if the CRDs are not installed by Helm, they must be added before policies can be created |
| crds.annotations | object | `{}` | Additional CRDs annotations |
| crds.labels | object | `{}` | Additional CRDs labels |
| rbac.create | bool | `true` | Create RBAC resources |
| rbac.serviceAccount.name | string | `nil` | The ServiceAccount name |
| rbac.serviceAccount.annotations | object | `{}` | Annotations for the ServiceAccount |
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
| containers.server.image.registry | string | `"ghcr.io"` | Image registry |
| containers.server.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| containers.server.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| containers.server.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| containers.server.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| containers.server.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| containers.server.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| containers.server.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.ports | list | `[{"containerPort":9080,"name":"http","protocol":"TCP"},{"containerPort":9081,"name":"grpc","protocol":"TCP"}]` | Container ports. |
| containers.server.args | list | `["serve","--http-address=:9080","--grpc-address=:9081"]` | Container args. |
| service.port | int | `9081` | Service port. |
| service.type | string | `"ClusterIP"` | Service type. |
| service.nodePort | string | `nil` | Service node port. Only used if `type` is `NodePort`. |
| service.annotations | object | `{}` | Service annotations. |
| pdb | string | `nil` |  |

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
