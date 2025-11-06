# Example

## Setup

In this example we will deploy the Kyverno Authz server in a Kubernetes cluster.

### Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install the components
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

### Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.34.0

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

For more certificate management options, refer to [Certificates management](../../quick-start/kube-install.md#certificates-management).

### Deploy the Authz server

```bash
# deploy the kyverno authz server
helm install kyverno-authz-server                                             \
  --namespace kyverno --create-namespace                                      \
  --wait                                                                      \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server  \
  --set config.type=http                                                      \
  --set certManager.issuerRef.name=selfsigned-issuer \
  --set certManager.issuerRef.kind=ClusterIssuer \
  --set certManager.issuerRef.group=cert-manager.io
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
    expression: object.attributes.header[?"x-force-authorized"].orValue("") in ["enabled", "true"]
  validations:
  - expression: |
      !variables.force_authorized
        ? http.Denied("Forbidden").Response()
        : http.Allowed().Response()
EOF
```

This policy denies requests that don't contain the header `x-force-authorized` with the value `enabled` or `true`.

## Wrap Up

Congratulations on completing this guide!

You have successfully deployed the Kyverno HTTP Authorizer control plane, sidecar injector, and a sample ValidatingPolicy.
