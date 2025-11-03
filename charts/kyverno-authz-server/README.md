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
| authzServer.envoy.deployment.replicas | int | `nil` | Desired number of pods |
| authzServer.envoy.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| authzServer.envoy.deployment.annotations | object | `{}` | Deployment annotations. |
| authzServer.envoy.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| authzServer.envoy.pod.labels | object | `{}` | Additional labels to add to each pod |
| authzServer.envoy.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| authzServer.envoy.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| authzServer.envoy.pod.securityContext | object | `{}` | Security context |
| authzServer.envoy.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| authzServer.envoy.pod.tolerations | list | `[]` | List of node taints to tolerate |
| authzServer.envoy.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| authzServer.envoy.pod.priorityClassName | string | `""` | Optional priority class |
| authzServer.envoy.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| authzServer.envoy.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| authzServer.envoy.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| authzServer.envoy.pod.affinity | object | `{}` | Pod affinity constraints. |
| authzServer.envoy.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| authzServer.envoy.container.image.registry | string | `"ghcr.io"` | Image registry |
| authzServer.envoy.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| authzServer.envoy.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| authzServer.envoy.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| authzServer.envoy.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| authzServer.envoy.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| authzServer.envoy.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| authzServer.envoy.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.envoy.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.envoy.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.envoy.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9081,"name":"grpc","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"}]` | Container ports. |
| authzServer.envoy.container.args | list | `["serve","envoy","authz-server","--probes-address=:9080","--grpc-address=:9081","--grpc-network={{ .Values.authzServer.envoy.config.grpcNetwork }}","--metrics-address=:9082","--leader-election=true","--leader-election-id={{ template \"kyverno-authz-server.name\" . }}","--kube-policy-source={{ .Values.authzServer.envoy.config.kubePolicySource }}","--allow-insecure-registry={{ .Values.authzServer.envoy.config.allowInsecureRegistry }}"]` | Container args. |
| authzServer.envoy.service.annotations | object | `{}` | Service annotations. |
| authzServer.envoy.service.type | string | `"ClusterIP"` | Service type. |
| authzServer.envoy.service.ports | list | `[{"appProtocol":null,"name":"grpc","nodePort":null,"port":9081,"protocol":"TCP","targetPort":"grpc"}]` | Service ports. |
| authzServer.envoy.config.grpcNetwork | string | `"tcp"` | GRPC network type (tcp, unix, etc.) |
| authzServer.envoy.config.kubePolicySource | bool | `true` | Enable in-cluster kubernetes policy source |
| authzServer.envoy.config.externalPolicySources | list | `[]` | External policy sources |
| authzServer.envoy.config.allowInsecureRegistry | bool | `false` | Allow insecure registry for pulling policy images |
| authzServer.envoy.config.imagePullSecrets | list | `[]` | Image pull secrets for fetching policies from OCI registries |
| authzServer.http.deployment.replicas | int | `nil` | Desired number of pods |
| authzServer.http.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| authzServer.http.deployment.annotations | object | `{}` | Deployment annotations. |
| authzServer.http.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| authzServer.http.pod.labels | object | `{}` | Additional labels to add to each pod |
| authzServer.http.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| authzServer.http.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| authzServer.http.pod.securityContext | object | `{}` | Security context |
| authzServer.http.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| authzServer.http.pod.tolerations | list | `[]` | List of node taints to tolerate |
| authzServer.http.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| authzServer.http.pod.priorityClassName | string | `""` | Optional priority class |
| authzServer.http.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| authzServer.http.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| authzServer.http.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| authzServer.http.pod.affinity | object | `{}` | Pod affinity constraints. |
| authzServer.http.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| authzServer.http.container.image.registry | string | `"ghcr.io"` | Image registry |
| authzServer.http.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| authzServer.http.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| authzServer.http.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| authzServer.http.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| authzServer.http.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| authzServer.http.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| authzServer.http.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.http.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.http.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| authzServer.http.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"},{"containerPort":9083,"name":"http","protocol":"TCP"}]` | Container ports. |
| authzServer.http.container.args | list | `["serve","http","authz-server","--probes-address=:9080","--metrics-address=:9082","--server-address=:9083","--leader-election=true","--leader-election-id={{ template \"kyverno-authz-server.name\" . }}","--kube-policy-source={{ .Values.authzServer.http.config.kubePolicySource }}","--nested-request={{ .Values.authzServer.http.config.nestedRequest }}","--allow-insecure-registry={{ .Values.authzServer.http.config.allowInsecureRegistry }}"]` | Container args. |
| authzServer.http.service.annotations | object | `{}` | Service annotations. |
| authzServer.http.service.type | string | `"ClusterIP"` | Service type. |
| authzServer.http.service.ports | list | `[{"appProtocol":null,"name":"http","nodePort":null,"port":9083,"protocol":"TCP","targetPort":"http"}]` | Service ports. |
| authzServer.http.config.kubePolicySource | bool | `true` | Enable in-cluster kubernetes policy source |
| authzServer.http.config.externalPolicySources | list | `[]` | External policy sources |
| authzServer.http.config.allowInsecureRegistry | bool | `false` | Allow insecure registry for pulling policy images |
| authzServer.http.config.imagePullSecrets | list | `[]` | Image pull secrets for fetching policies from OCI registries |
| authzServer.http.config.nestedRequest | bool | `true` | Expect the requests to validate to be in the body of the original request |
| authzServer.http.config.controlPlane.address | string | `""` | Control plane address (leave empty for standalone mode) |
| authzServer.http.config.controlPlane.reconnectWait | string | `"3s"` | Duration to wait before retrying connecting to the control plane |
| authzServer.http.config.controlPlane.maxDialInterval | string | `"8s"` | Duration to wait before stopping attempts of sending a policy to a client |
| authzServer.http.config.controlPlane.healthCheckInterval | string | `"30s"` | Interval for sending health checks |
| validatingWebhookConfiguration.annotations | object | `{}` | Webhook annotations |
| validatingWebhookConfiguration.certificates.static | object | `{}` | Static data to set in certificate secret |
| validatingWebhookConfiguration.certificates.certManager | object | `{}` | Infos for creating certificate with cert manager |
| validatingWebhookConfiguration.webhooks.envoy.failurePolicy | string | `"Fail"` | Webhook failure policy |
| validatingWebhookConfiguration.webhooks.envoy.deployment.replicas | int | `nil` | Desired number of pods |
| validatingWebhookConfiguration.webhooks.envoy.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| validatingWebhookConfiguration.webhooks.envoy.deployment.annotations | object | `{}` | Deployment annotations. |
| validatingWebhookConfiguration.webhooks.envoy.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| validatingWebhookConfiguration.webhooks.envoy.pod.labels | object | `{}` | Additional labels to add to each pod |
| validatingWebhookConfiguration.webhooks.envoy.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| validatingWebhookConfiguration.webhooks.envoy.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| validatingWebhookConfiguration.webhooks.envoy.pod.securityContext | object | `{}` | Security context |
| validatingWebhookConfiguration.webhooks.envoy.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| validatingWebhookConfiguration.webhooks.envoy.pod.tolerations | list | `[]` | List of node taints to tolerate |
| validatingWebhookConfiguration.webhooks.envoy.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| validatingWebhookConfiguration.webhooks.envoy.pod.priorityClassName | string | `""` | Optional priority class |
| validatingWebhookConfiguration.webhooks.envoy.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| validatingWebhookConfiguration.webhooks.envoy.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| validatingWebhookConfiguration.webhooks.envoy.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| validatingWebhookConfiguration.webhooks.envoy.pod.affinity | object | `{}` | Pod affinity constraints. |
| validatingWebhookConfiguration.webhooks.envoy.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| validatingWebhookConfiguration.webhooks.envoy.container.image.registry | string | `"ghcr.io"` | Image registry |
| validatingWebhookConfiguration.webhooks.envoy.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| validatingWebhookConfiguration.webhooks.envoy.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| validatingWebhookConfiguration.webhooks.envoy.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| validatingWebhookConfiguration.webhooks.envoy.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| validatingWebhookConfiguration.webhooks.envoy.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| validatingWebhookConfiguration.webhooks.envoy.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| validatingWebhookConfiguration.webhooks.envoy.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.envoy.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.envoy.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.envoy.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9081,"name":"grpc","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"}]` | Container ports. |
| validatingWebhookConfiguration.webhooks.envoy.container.args | list | `["serve","envoy","validation-webhook","--probes-address=:9080","--metrics-address=:9082"]` | Container args. |
| validatingWebhookConfiguration.webhooks.http.failurePolicy | string | `"Fail"` | Webhook failure policy |
| validatingWebhookConfiguration.webhooks.http.deployment.replicas | int | `nil` | Desired number of pods |
| validatingWebhookConfiguration.webhooks.http.deployment.revisionHistoryLimit | int | `10` | The number of revisions to keep |
| validatingWebhookConfiguration.webhooks.http.deployment.annotations | object | `{}` | Deployment annotations. |
| validatingWebhookConfiguration.webhooks.http.deployment.updateStrategy | object | See [values.yaml](values.yaml) | Deployment update strategy. Ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| validatingWebhookConfiguration.webhooks.http.pod.labels | object | `{}` | Additional labels to add to each pod |
| validatingWebhookConfiguration.webhooks.http.pod.annotations | object | `{}` | Additional annotations to add to each pod |
| validatingWebhookConfiguration.webhooks.http.pod.imagePullSecrets | list | `[]` | Image pull secrets |
| validatingWebhookConfiguration.webhooks.http.pod.securityContext | object | `{}` | Security context |
| validatingWebhookConfiguration.webhooks.http.pod.nodeSelector | object | `{}` | Node labels for pod assignment |
| validatingWebhookConfiguration.webhooks.http.pod.tolerations | list | `[]` | List of node taints to tolerate |
| validatingWebhookConfiguration.webhooks.http.pod.topologySpreadConstraints | list | `[]` | Topology spread constraints. |
| validatingWebhookConfiguration.webhooks.http.pod.priorityClassName | string | `""` | Optional priority class |
| validatingWebhookConfiguration.webhooks.http.pod.hostNetwork | bool | `false` | Change `hostNetwork` to `true` when you want the pod to share its host's network namespace. Useful for situations like when you end up dealing with a custom CNI over Amazon EKS. Update the `dnsPolicy` accordingly as well to suit the host network mode. |
| validatingWebhookConfiguration.webhooks.http.pod.dnsPolicy | string | `"ClusterFirst"` | `dnsPolicy` determines the manner in which DNS resolution happens in the cluster. In case of `hostNetwork: true`, usually, the `dnsPolicy` is suitable to be `ClusterFirstWithHostNet`. For further reference: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy. |
| validatingWebhookConfiguration.webhooks.http.pod.antiAffinity | object | See [values.yaml](values.yaml) | Pod anti affinity constraints. |
| validatingWebhookConfiguration.webhooks.http.pod.affinity | object | `{}` | Pod affinity constraints. |
| validatingWebhookConfiguration.webhooks.http.pod.nodeAffinity | object | `{}` | Node affinity constraints. |
| validatingWebhookConfiguration.webhooks.http.container.image.registry | string | `"ghcr.io"` | Image registry |
| validatingWebhookConfiguration.webhooks.http.container.image.repository | string | `"kyverno/kyverno-envoy-plugin"` | Image repository |
| validatingWebhookConfiguration.webhooks.http.container.image.tag | string | `nil` | Image tag Defaults to appVersion in Chart.yaml if omitted |
| validatingWebhookConfiguration.webhooks.http.container.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| validatingWebhookConfiguration.webhooks.http.container.resources.limits | object | `{"memory":"384Mi"}` | Pod resource limits |
| validatingWebhookConfiguration.webhooks.http.container.resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | Pod resource requests |
| validatingWebhookConfiguration.webhooks.http.container.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Container security context |
| validatingWebhookConfiguration.webhooks.http.container.startupProbe | object | See [values.yaml](values.yaml) | Startup probe. The block is directly forwarded into the deployment, so you can use whatever startupProbes configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.http.container.livenessProbe | object | See [values.yaml](values.yaml) | Liveness probe. The block is directly forwarded into the deployment, so you can use whatever livenessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.http.container.readinessProbe | object | See [values.yaml](values.yaml) | Readiness Probe. The block is directly forwarded into the deployment, so you can use whatever readinessProbe configuration you want. ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/ |
| validatingWebhookConfiguration.webhooks.http.container.ports | list | `[{"containerPort":9080,"name":"probes","protocol":"TCP"},{"containerPort":9081,"name":"grpc","protocol":"TCP"},{"containerPort":9082,"name":"metrics","protocol":"TCP"}]` | Container ports. |
| validatingWebhookConfiguration.webhooks.http.container.args | list | `["serve","http","validation-webhook","--probes-address=:9080","--metrics-address=:9082"]` | Container args. |
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
