---
title: Policy Sources
description: Configure how the Kyverno Authorization Server loads policies from different sources such as Kubernetes, local files, Git repositories, or OCI registries.
---

# Policy Sources

The Kyverno Authorization Server can load policies from multiple sources.  
Policies may come from the Kubernetes API, local files, Git repositories, or OCI images.

You can specify these sources when starting the server using command-line flags.

## Overview

- **Flag:** `--external-policy-source=<url>` (repeatable)  
- **Default:** none  
- **Kubernetes source:** Enabled by default. Disable it with `--kube-policy-source=false` if you only want external sources.

External policy sources are mounted as **read-only virtual filesystems**.  
You can pass multiple `--external-policy-source` flags to combine several sources at once.

## 1. Kubernetes Policy Source

By default, the Authz Server watches for policies from the Kubernetes API.  
This source allows it to dynamically load and update `ValidatingPolicy` resources from your cluster.

### Disable the Kubernetes Policy Source

If you only want to load policies from external sources (for example, local files or Git), you can disable this behavior:

    --kube-policy-source=false

This flag is often used in **sidecar** or **local development** environments.

### Example

Run the server without using in-cluster policies:

    kyverno-envoy-plugin serve authz-server \
      --kube-policy-source=false \
      --external-policy-source=file://policies

## 2. External Policy Sources

External sources are configured using the `--external-policy-source` flag.  
Each source URL specifies a virtual filesystem scheme.

You can combine multiple sources by repeating the flag:

    --external-policy-source=file:///policies/team-a \
    --external-policy-source=git+https://github.com/acme/policies.git

The following types of sources are supported:

### 2.1 File Source

The **file** source loads policies from a local or mounted filesystem directory.

- **Scheme:** `file://`
- **Description:** Reads YAML policy files from a directory.
- **Example:**  

        --external-policy-source=file:///data/kyverno-authz-server

This is the **default configuration** used by injected sidecars.

#### Usage Example

Run the Authz Server with only local policies:

    kyverno-envoy-plugin serve authz-server \
      --kube-policy-source=false \
      --external-policy-source=file:///policies

### 2.2 Git Source

The **git** source lets the server clone a Git repository and read policies from it.

- **Scheme:** `git+https://`
- **Description:** Loads policies from a Git repository using the `gitfs` (go-fsimpl) backend.
- **Example:**  

        --external-policy-source=git+https://github.com/acme/policies.git

You can reference specific **branches**, **tags**, **commits**, or **subdirectories** using URL parameters supported by `gitfs`.

#### Example (with multiple sources)

    kyverno-envoy-plugin serve authz-server \
      --external-policy-source=file:///policies/team-a \
      --external-policy-source=git+https://github.com/acme/policies.git

#### Notes

- The Git repository must be publicly accessible or configured with appropriate credentials.
- Policies are read once at startup; to update them, restart the server.

### 2.3 OCI Source

The **oci** source loads policies from an OCI image stored in a container registry.

- **Scheme:** `oci://`
- **Description:** Pulls a container image and uses its filesystem contents as policy definitions.
- **Example:**  

        --external-policy-source=oci://ghcr.io/org/policies:tag

This is useful for distributing versioned policy bundles via registries such as GHCR or Docker Hub.

#### Notes

- The image must contain YAML policy files at its root.
- The contents are mounted read-only.
- Updates require republishing the image and restarting the server.

## 3. Policy File Expectations

From each mounted source directory, the Authz Server loads all YAML and JSON files recusively.  
It expects documents of the following kind:

- `#!yaml apiVersion: policies.kyverno.io/v1alpha1` and `#!yaml kind: ValidatingPolicy`

Invalid or non-policy documents are skipped.  

## 4. Helm Integration

When deploying via Helm, you can specify external policy sources through values:

```yaml
config:
  sources:
    external:
    - file:///data/kyverno-authz-server
    - git+https://github.com/acme/policies.git
```

These values are automatically converted into repeated `--external-policy-source=...` flags in the container arguments.

## 5. Troubleshooting

- Ensure each file or repository path is reachable from the pod.
- For file-based sources, verify that the directory is correctly mounted in the container.
- Check container logs for detailed error messages or policy compilation issues.

## See Also

- [`--kube-policy-source`](#1-kubernetes-policy-source) — enable or disable in-cluster policy loading.  
- [go-fsimpl `gitfs` documentation](https://pkg.go.dev/github.com/go-git/go-fsimpl) — for Git URL options (branch, ref, subpath).  
- [Kyverno Policies documentation](../policies/index.md) — for defining and testing ValidatingPolicies.
