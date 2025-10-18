# kyverno-http-authorizer-control-plane

Kyverno policies based HTTP authorization server ❤️

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A server to enforce kyverno policies for HTTP authorization. This server allows applying Kyverno policies to HTTP requests.

## Overview

This authorization server provides HTTP-based policy enforcement using Kyverno policies. It can be integrated with various proxies and gateways that support external authorization services to make authorization decisions based on Kyverno policies.

## Installing the Chart

Add `kyverno-http-authorizer` Helm repository:

```shell
helm repo add kyverno-http-authorizer https://kyverno.github.io/kyverno-http-authorizer/
```

Install `kyverno-http-authorizer-control-plane` Helm chart:

```shell
helm install kyverno-http-authorizer-control-plane --namespace kyverno --create-namespace kyverno-http-authorizer/kyverno-http-authorizer-control-plane
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
| containers.server.image.registry | string | `"ghcr.io"` | Image registry |
| containers.server.image.repository | string | `"kyverno/kyverno-http-authorizer"` | Image repository |
| containers.server.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| containers.server.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| containers.server.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| containers.server.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| containers.server.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| containers.server.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| containers.server.ports | list | `[{"containerPort":9080,"name":"http","protocol":"TCP"},{"containerPort":9081,"name":"grpc","protocol":"TCP"}]` | Container ports. |
| containers.server.args | list | `["serve","authz-server","--probes-address=:9080","--grpc-address=:9081","--metrics-address=:9082"]` | Container args. |
| service.port | int | `9081` | Service port. |
| service.type | string | `"ClusterIP"` | Service type. |
| service.nodePort | string | `nil` | Service node port. Only used if `type` is `NodePort`. |
| service.annotations | object | `{}` | Service annotations. |
| service.appProtocol | string | `nil` | Service application protocol. Setting app protocol is only needed in specific cases like integration with certain gateways. ref: https://kubernetes.io/docs/concepts/services-networking/service/#application-protocol |
| webhook.annotations | object | `{}` | Webhook annotations |
| webhook.failurePolicy | string | `"Fail"` | Webhook failure policy |
| webhook.objectSelector | string | `nil` | Webhook object selector |
| webhook.namespaceSelector | object | `{"matchExpressions":[{"key":"kyverno-injection","operator":"In","values":["enabled"]}]}` | Webhook namespace selector |
| pdb | string | `nil` |  |
| externalPolicySources | list | `[]` | External policy sources |
| crds.install | bool | `true` |  |

## Source Code

* <https://github.com/kyverno/kyverno-http-authorizer>

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
