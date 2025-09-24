# CRDs

The Kyverno Authz Server leverages the Kyverno `ValidatingPolicy` generic CRD.

This resource definition is not specific to the Kyverno Authz Server and must be installed separately.

## Install Kyverno ValidatingPolicy CRD

Before deploying the Kyverno Authz Server, make sure the Kyverno ValidatingPolicy CRD is installed.

```bash
kubectl apply \
  -f https://raw.githubusercontent.com/kyverno/kyverno/refs/heads/main/config/crds/policies.kyverno.io/policies.kyverno.io_validatingpolicies.yaml
```
