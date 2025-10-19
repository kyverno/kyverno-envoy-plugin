# Quick Start

## Setup

In this quick start guide we will deploy the Kyverno HTTP Authorizer components in a Kubernetes cluster.

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install the components
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.31.1

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

### Deploy cert-manager

The Kyverno HTTP Authorizer components need certificates for their webhooks.

```bash
# install cert-manager
helm install cert-manager \
  --namespace cert-manager --create-namespace \
  --wait \
  --repo https://charts.jetstack.io cert-manager \
  --set crds.enabled=true

# create a self-signed cluster issuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF
```

For more certificate management options, refer to [Certificates management](../install/certificates.md).

### Deploy the Control Plane

Deploy the control plane which manages policies:

```bash
# deploy the control plane
helm install kyverno-http-authorizer-control-plane \
  --namespace kyverno --create-namespace \
  --wait \
  --repo https://kyverno.github.io/kyverno-http-authorizer kyverno-http-authorizer-control-plane \
  --set certificates.certManager.issuerRef.group=cert-manager.io \
  --set certificates.certManager.issuerRef.kind=ClusterIssuer \
  --set certificates.certManager.issuerRef.name=selfsigned-issuer
```

### Deploy the Sidecar Injector

Deploy the sidecar injector:

```bash
# deploy the sidecar injector
helm install kyverno-sidecar-injector \
  --namespace kyverno \
  --wait \
  --repo https://kyverno.github.io/kyverno-http-authorizer kyverno-sidecar-injector \
  --set certificates.certManager.issuerRef.group=cert-manager.io \
  --set certificates.certManager.issuerRef.kind=ClusterIssuer \
  --set certificates.certManager.issuerRef.name=selfsigned-issuer \
  --set controlPlaneAddress=kyverno-http-authorizer-control-plane.kyverno.svc.cluster.local:9081
```

### Deploy a ValidatingPolicy

Deploy a sample policy:

```bash
# deploy validating policy
kubectl apply -f - <<EOF
apiVersion: policies.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: demo
spec:
  failurePolicy: Fail
  evaluation:
    mode: HTTP
  variables:
  - name: force_authorized
    expression: object.headers.get("x-force-authorized")
  - name: allowed
    expression: variables.force_authorized in ["enabled", "true"]
  validations:
  - expression: |
      !variables.allowed
        ? http.response().status(403).withBody("Forbidden")
        : null
  - expression: |
      http.response().status(200)
EOF
```

This policy denies requests that don't contain the header `x-force-authorized` with the value `enabled` or `true`.

## Next Steps

Now that you have the components deployed, check out the [tutorials](../tutorials/index.md) to learn how to integrate with:

- [Ingress NGINX](../tutorials/ingress-nginx/index.md)

## Wrap Up

Congratulations on completing the quick start guide!

You have successfully deployed the Kyverno HTTP Authorizer control plane, sidecar injector, and a sample ValidatingPolicy.
