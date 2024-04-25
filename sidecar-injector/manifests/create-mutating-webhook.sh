#!/bin/bash

# Get the base64 encoded tls.crt data
CA_BUNDLE=$(kubectl get secret kyverno-envoy-sidecar-certs -n kyverno-envoy-sidecar-injector -o jsonpath='{.data.tls\.crt}')

# Create the mutatingwebhook.yaml file
cat <<EOF > mutatingwebhook.yaml
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: kyverno-envoy-sidecar
  labels:
    app.kubernetes.io/name: sidecar-injector
    app.kubernetes.io/instance: sidecar-injector
webhooks:
  - name: kyverno-envoy-sidecar.kyverno-envoy-sidecar-injector.svc
    clientConfig:
      service:
        name: kyverno-envoy-sidecar
        namespace: kyverno-envoy-sidecar-injector
        path: "/mutate"
      caBundle: $CA_BUNDLE
    failurePolicy: Fail
    sideEffects: None
    admissionReviewVersions: ["v1"]
    rules:
      - apiGroups:
          - ""
        resources:
          - pods
        apiVersions:
          - "v1"
        operations:
          - CREATE
        scope: '*'  
    objectSelector:
      matchExpressions:
        - key: kyverno-envoy-sidecar/injection
          operator: In
          values:
          - enabled
EOF

# Apply the mutatingwebhook.yaml file
kubectl apply -f mutatingwebhook.yaml