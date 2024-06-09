#!/bin/bash

KIND_IMAGE=kindest/node:v1.29.2
ISTIO_REPO=https://istio-release.storage.googleapis.com/charts
ISTIO_NS=istio-system

# Create Kind cluster
kind create cluster --image $KIND_IMAGE --wait 1m --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |-
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
        protocol: TCP
      - containerPort: 443
        hostPort: 443
        protocol: TCP
  - role: worker
EOF

# Install Istio components
helm upgrade --install istio-base       --namespace $ISTIO_NS           --create-namespace --wait --repo $ISTIO_REPO base
helm upgrade --install istiod           --namespace $ISTIO_NS           --create-namespace --wait --repo $ISTIO_REPO istiod
