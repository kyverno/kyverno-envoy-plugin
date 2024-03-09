#!/bin/bash

# Create Kind cluster
kind create cluster

# Add Helm repo
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

# Create istio-system namespace
kubectl create namespace istio-system

# Install Istio components
helm install istio-base istio/base -n istio-system
helm install istiod istio/istiod -n istio-system --wait
helm install istio-ingress istio/gateway -n istio-system

# Label default namespace to inject sidecar automatically
kubectl label namespace default istio-injection=enabled