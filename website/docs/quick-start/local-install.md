# Local installation

You can install the pre-compiled binary (in several ways), compile from sources, or run with Docker.

## Install the pre-compiled binary

### Manually

Download the pre-compiled binaries for your system from the [releases page](https://github.com/kyverno/kyverno-envoy-plugin/releases) and copy them to the desired location.

### Install using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/kyverno-envoy-plugin@latest
```

## Run with Docker

Chainsaw is also available as a Docker image which you can pull and run:

```bash
docker pull ghcr.io/kyverno/kyverno-envoy-plugin:<version>
```

```bash
docker run --rm                                     \
    -v ${HOME}/.kube/:/etc/kubeconfig/              \
    -e KUBECONFIG=/etc/kubeconfig/config            \
    --network=host                                  \
    ghcr.io/kyverno/kyverno-envoy-plugin:<version>  \
    version
```

## Compile from sources

**clone:**

```bash
git clone https://github.com/kyverno/kyverno-envoy-plugin.git
```

**build the binaries:**

```bash
cd kyverno-envoy-plugin
go mod tidy
make build
```

**verify it works:**

```bash
./kyverno-envoy-plugin version
```
