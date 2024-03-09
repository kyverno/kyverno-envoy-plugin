#!/bin/bash

KIND_IMAGE=kindest/node:v1.29.2
ISTIO_REPO=https://istio-release.storage.googleapis.com/charts
ISTIO_NS=istio-system

# Create Kind cluster
kind create cluster --image $KIND_IMAGE

# Install Istio components
helm upgrade --install istio-base       --namespace $ISTIO_NS --create-namespace --wait --repo $ISTIO_REPO base
helm upgrade --install istiod           --namespace $ISTIO_NS --create-namespace --wait --repo $ISTIO_REPO istiod
helm upgrade --install istio-ingress    --namespace $ISTIO_NS --create-namespace --wait --repo $ISTIO_REPO gateway

# Label default namespace to inject sidecar automatically
kubectl label namespace default istio-injection=enabled
