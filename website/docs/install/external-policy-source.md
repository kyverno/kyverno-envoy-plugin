# External policy sources

The authz server can load policies from external sources via the `--external-policy-source` flag. You can pass this flag multiple times to combine several sources. These external sources are used in addition to, or instead of, the in-cluster source controlled by `--kube-policy-source`.

- **Flag**: `--external-policy-source=<url>` (repeatable)
- **Default**: none
- **Kubernetes source**: enabled by default. Disable it with `--kube-policy-source=false` if you only want external sources.

### Supported sources
External sources are backed by virtual filesystems. The following schemes are supported:

- **file**: Local or mounted filesystem directory containing policy YAML files.
  - Example (sidecar's default): `--external-policy-source=file:///data/kyverno-authz-server`

- **git**: Git repositories via `gitfs` (go-fsimpl). Useful to load policies from a repo.
  - Typical form: `--external-policy-source=git+https://<host>/<org>/<repo>.git`
  - You can point to a subdirectory or ref depending on your git URL; consult the `gitfs` docs for exact URL options (branch/tag/commit and subpaths).
  - Example: `--external-policy-source=git+https://github.com/acme/policies.git`

Notes:
- Each source is mounted read-only and scanned at startup; updates require restarting the server. 

### What files are expected
From each source directory, the server loads non-recursively all YAML files and parses embedded documents:
- `envoy.kyverno.io/v1alpha1` `AuthorizationPolicy`
- `policies.kyverno.io/v1alpha1` `ValidatingPolicy`

Invalid or non-policy documents are skipped. Compilation errors will make startup fail.

### Usage examples

- Run the authz server with a local directory and without in-cluster policies:
```bash
kyverno-envoy-plugin serve authz-server \
  --kube-policy-source=false \
  --external-policy-source=file:///policies
```

- Run with multiple sources (a local mount and a git repo):
```bash
kyverno-envoy-plugin serve authz-server \
  --external-policy-source=file:///policies/team-a \
  --external-policy-source=git+https://github.com/acme/policies.git
```

- Sidecar container (what the injector adds by default):
```text
--kube-policy-source=false
--external-policy-source=file:///data/kyverno-authz-server
```

### Adding external sources with Helm
When using the helm chart, you can configure additional sources to inject into sidecars:

Values example:
```yaml
externalPolicySources:
  - file:///data/kyverno-authz-server
  - git+https://github.com/acme/policies.git
```

These become repeated `--external-policy-source=...` arguments on the container.

### Troubleshooting
- Ensure the path or repo is reachable from the pod. For local sources, mount the directory into the sidecar at the correct path.
- Make sure the directory contains YAML files at its root (no recursive loading).
- Check logs for compilation errors if startup fails.

### See also
- `--kube-policy-source` to enable/disable in-cluster policy source (defaults to true)
- Git URL options and examples are provided by the upstream `gitfs` implementation (go-fsimpl).
