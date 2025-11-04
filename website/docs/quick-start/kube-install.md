# Kubernetes installation

## Prerequisites

- A Kubernetes cluster
- [Helm](https://helm.sh/) to install the components
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to interact with the cluster

## Setup a cluster (optional)

If you don't have a cluster at hand, you can create a local one with [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

```bash
KIND_IMAGE=kindest/node:v1.34.0

# create cluster
kind create cluster --image $KIND_IMAGE --wait 1m
```

## ValidatingPolicy CRD

The Kyverno Authz Server leverages the Kyverno `ValidatingPolicy` generic CRD.

This resource definition is not specific to the Kyverno Authz Server and must be installed separately.

Before deploying the Kyverno Authz Server, make sure the Kyverno ValidatingPolicy CRD is installed.

```bash
kubectl apply \
  -f https://raw.githubusercontent.com/kyverno/kyverno/refs/heads/main/config/crds/policies.kyverno.io/policies.kyverno.io_validatingpolicies.yaml
```

## Certificates management

The Kyverno Authz Server comes with a validation webhook and needs a valid certificate to let the api server call into it.

At deployment time you can either provide your own certificate or use [cert-manager](https://cert-manager.io) to create one for the Kyverno Authz Server.

### Bring your own

If you want to bring your own certificate, you can set `certificates.static` values when installing the helm chart.

```bash
# create certificate
openssl req -new -x509  \
  -subj "/CN=kyverno-authz-server-validation-authorizationpolicy.kyverno.svc" \
  -addext "subjectAltName = DNS:kyverno-authz-server-validation-authorizationpolicy.kyverno.svc" \
  -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt

# install chart with static certificate
helm install kyverno-authz-server \
  --namespace kyverno --create-namespace \
  --wait \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server \
  --set-file certificates.static.crt=tls.crt \
  --set-file certificates.static.key=tls.key
```

### Use cert-manager

If you don't want to manage the certificate yourself you can rely on [cert-manager](https://cert-manager.io) to create the certificate for you and inject it in the webhook configuration.

```bash
# install cert-manager
helm install cert-manager \
  --namespace cert-manager --create-namespace \
  --wait \
  --repo https://charts.jetstack.io cert-manager \
  --values - <<EOF
crds:
  enabled: true
EOF

# create a certificate issuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF

# install chart with managed certificate
helm upgrade --install kyverno-authz-server \
  --namespace kyverno --create-namespace \
  --wait \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server \
  --values - <<EOF
certificates:
  certManager:
    issuerRef:
      group: cert-manager.io
      kind: ClusterIssuer
      name: selfsigned-issuer
EOF
```

### Deploy the Kyverno Authz Server

Now we can deploy the Kyverno Authz Server.

```bash
# create the kyverno namespace
kubectl create ns kyverno

# label the namespace to inject the envoy proxy
kubectl label namespace kyverno istio-injection=enabled

# deploy the kyverno authz server
helm install kyverno-authz-server \
  --namespace kyverno \
  --wait  \
  --repo https://kyverno.github.io/kyverno-envoy-plugin kyverno-authz-server \
  --values - <<EOF
certificates:
  certManager:
    issuerRef:
      group: cert-manager.io
      kind: ClusterIssuer
      name: selfsigned-issuer
EOF
```
