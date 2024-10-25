#############
# VARIABLES #
#############

GIT_SHA                            := $(shell git rev-parse HEAD)
ORG                                ?= kyverno
PACKAGE                            ?= github.com/$(ORG)/kyverno-envoy-plugin
GOPATH_SHIM                        := ${PWD}/.gopath
PACKAGE_SHIM                       := $(GOPATH_SHIM)/src/$(PACKAGE)
CLI_BIN                            := kyverno-envoy-plugin
CGO_ENABLED                        ?= 0
GOOS                               ?= $(shell go env GOOS)
ifdef VERSION
LD_FLAGS                           := "-s -w -X $(PACKAGE)/pkg/version.BuildVersion=$(VERSION)"
else
LD_FLAGS                           := "-s -w"
endif
KIND_IMAGE                         ?= kindest/node:v1.29.2
REGISTRY                           ?= ghcr.io
REPO                               ?= kyverno
IMAGE                              ?= kyverno-envoy-plugin
KO_REGISTRY                        ?= ko.local
KO_TAGS                            ?= $(GIT_SHA)
KO_PLATFORMS                       ?= all

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
HELM                               ?= $(TOOLS_DIR)/helm
HELM_VERSION                       ?= v3.12.3
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.22.0
KO                                 ?= $(TOOLS_DIR)/ko
KO_VERSION                         ?= v0.15.1
TOOLS                              := $(HELM) $(KIND) $(KO)
PIP                                ?= "pip"
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

.PHONY: install-tools
install-tools: ## Install tools
install-tools: $(TOOLS)

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

$(GOPATH_SHIM):
	@echo Create gopath shim... >&2
	@mkdir -p $(GOPATH_SHIM)

.INTERMEDIATE: $(PACKAGE_SHIM)
$(PACKAGE_SHIM): $(GOPATH_SHIM)
	@echo Create package shim... >&2
	@mkdir -p $(GOPATH_SHIM)/src/github.com/$(ORG) && ln -s -f ${PWD} $(PACKAGE_SHIM)

.PHONY: codegen-mkdocs
codegen-mkdocs: ## Generate mkdocs website
	@echo Generate mkdocs website... >&2
	@$(PIP) install mkdocs
	@$(PIP) install --upgrade pip
	@$(PIP) install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin mike
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen
codegen: ## Rebuild all generated code and docs
codegen: codegen-mkdocs

.PHONY: verify-codegen
verify-codegen: ## Verify all generated code and docs are up to date
verify-codegen: codegen
	@echo Checking codegen is up to date... >&2
	@git --no-pager diff -- .
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen".' >&2
	@echo 'To correct this, locally run "make codegen", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- .

#########
# BUILD #
#########

.PHONY: fmt
fmt: ## Run go fmt
	@echo Go fmt... >&2
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo Go vet... >&2
	@go vet ./...

#########
# BUILD #
#########

.PHONY: build
build: ## Build
build: fmt
build: vet
build:
	@echo "Build..." >&2
	@LD_FLAGS=$(LD_FLAGS) go build .

##############
# BUILD (KO) #
##############

.PHONY: ko-login
ko-login: $(KO)
	@$(KO) login $(REGISTRY) --username $(REGISTRY_USERNAME) --password $(REGISTRY_PASSWORD)

.PHONY: build-ko
build-ko: ## Build Docker image with ko
build-ko: fmt
build-ko: vet
build-ko: $(KO)
	@echo "Build Docker image with ko..." >&2
	@LD_FLAGS=$(LD_FLAGS) KO_DOCKER_REPO=$(KO_REGISTRY) $(KO) build . --preserve-import-paths --tags=$(KO_TAGS)

.PHONY: publish-ko
publish-ko: ## Publish Docker image with ko
publish-ko: fmt
publish-ko: vet
publish-ko: ko-login
publish-ko: $(KO)
	@echo "Publish Docker image with ko..." >&2
	@LD_FLAGS=$(LD_FLAGS) KO_DOCKER_REPO=$(REGISTRY)/$(REPO)/$(IMAGE) $(KO) build . --bare --tags=$(KO_TAGS) --platform=$(KO_PLATFORMS)

########
# TEST #
########

.PHONY: tests
tests: ## Run tests
	@echo Running tests... >&2
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out

##########
# MKDOCS #
##########

.PHONY: mkdocs-serve
mkdocs-serve: ## Generate and serve mkdocs website
	@echo Generate and serve mkdocs website... >&2
	@$(PIP) install mkdocs
	@$(PIP) install --upgrade pip
	@$(PIP) install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin mike
	@mkdocs serve -f ./website/mkdocs.yaml

########
# KIND #
########

.PHONY: kind-create-cluster
kind-create-cluster: ## Create kind cluster
kind-create-cluster: $(KIND)
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --image $(KIND_IMAGE) --wait 1m

.PHONY: kind-load-image
kind-load-image: ## Build image and load it in kind cluster
kind-load-image: $(KIND)
kind-load-image: build-ko
	@echo Load image in kind... >&2
	@$(KIND) load docker-image $(KO_REGISTRY)/$(PACKAGE):$(GIT_SHA)

.PHONY: kind-load-taged-image
kind-load-taged-image: ## Build image and load it in kind cluster
kind-load-taged-image: $(KIND)
kind-load-taged-image: build-ko
	@echo Load image in kind... >&2
	docker tag $(KO_REGISTRY)/$(PACKAGE):$(GIT_SHA) $(KO_REGISTRY)/$(PACKAGE):latest
	@$(KIND) load docker-image $(KO_REGISTRY)/$(PACKAGE):latest

#########
# ISTIO #
#########

.PHONY: install-istio
install-istio: ## Install ISTIO
install-istio: $(HELM)
	@echo Install istio... >&2
	@$(HELM) upgrade --install istio-base --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts base
	@$(HELM) upgrade --install istiod --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts istiod

########
# HELM #
########

.PHONY: chart-install
chart-install: ## Install chart
chart-install: kind-load-image
chart-install: $(HELM)
	@echo Install helm chart... >&2
	@$(HELM) upgrade --install kyverno-envoy-plugin --namespace kyverno --create-namespace --wait ./charts/kyverno-envoy-plugin \
		--set sidecarInjector.containers.injector.image.registry=ko.local \
		--set sidecarInjector.containers.injector.image.repository=github.com/kyverno/kyverno-envoy-plugin \
		--set sidecarInjector.containers.injector.image.tag=$(GIT_SHA)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
