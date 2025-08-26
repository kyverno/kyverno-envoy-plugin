# Certificates management

The Kyverno Authz Server comes with a validation webhook and needs a valid certificate to let the api server call into it.

At deployment time you can either provide your own certificate or use [cert-manager](https://cert-manager.io) to create one for the Kyverno Authz Server.

## Bring your own

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

## Use cert-manager

If you don't want to manage the certificate yourself you can rely on [cert-manager](https://cert-manager.io) to create the certificate for you and inject it in the webhook configuration.

```bash
# install cert-manager
helm install cert-manager \
  --namespace cert-manager --create-namespace \
  --wait \
  --repo https://charts.jetstack.io cert-manager \
  --set crds.enabled=true

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
  --set certificates.certManager.issuerRef.group=cert-manager.io \
  --set certificates.certManager.issuerRef.kind=ClusterIssuer \
  --set certificates.certManager.issuerRef.name=selfsigned-issuer
```
