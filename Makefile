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
CRDS_PATH                          := .crds
ifdef VERSION
LD_FLAGS                           := "-s -w -X $(PACKAGE)/pkg/version.BuildVersion=$(VERSION)"
else
LD_FLAGS                           := "-s -w"
endif
KIND_IMAGE                         ?= kindest/node:v1.31.1
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
CONTROLLER_GEN                     ?= $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             := latest
REGISTER_GEN                       ?= $(TOOLS_DIR)/register-gen
REGISTER_GEN_VERSION               := v0.28.0
REFERENCE_DOCS                     := $(TOOLS_DIR)/genref
REFERENCE_DOCS_VERSION             := latest
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

$(CONTROLLER_GEN):
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

$(REGISTER_GEN):
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(REGISTER_GEN_VERSION)

$(REFERENCE_DOCS):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(REFERENCE_DOCS_VERSION)

.PHONY: install-tools
install-tools: ## Install tools
install-tools: $(HELM)
install-tools: $(KIND)
install-tools: $(KO)
install-tools: $(CONTROLLER_GEN)
install-tools: $(REGISTER_GEN)
install-tools: $(REFERENCE_DOCS)

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

.PHONY: codegen-crds
codegen-crds: ## Generate CRDs
codegen-crds: $(CONTROLLER_GEN)
codegen-crds: $(REGISTER_GEN)
	@echo Generate CRDs... >&2
	@$(CONTROLLER_GEN) paths=./apis/v1alpha1/... object
	@$(CONTROLLER_GEN) paths=./apis/v1alpha1/... crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false output:dir=$(CRDS_PATH)
	@$(REGISTER_GEN) --input-dirs=./apis/v1alpha1 --go-header-file=./.hack/boilerplate.go.txt --output-base=.

.PHONY: codegen-mkdocs
codegen-mkdocs: ## Generate mkdocs website
	@echo Generate mkdocs website... >&2
	@$(PIP) install -r requirements.txt
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen-helm-docs
codegen-helm-docs: ## Generate helm docs
	@echo Generate helm docs... >&2
	@docker run -v ${PWD}/charts:/work -w /work jnorwood/helm-docs:v1.11.0 -s file

.PHONY: codegen-helm-crds
codegen-helm-crds: codegen-crds ## Generate helm CRDs
	@echo Generate helm crds... >&2
	@cat $(CRDS_PATH)/* \
		| $(SED) -e '1i{{- if .Values.crds.install }}' \
		| $(SED) -e '$$a{{- end }}' \
		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- end }}' \
 		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- toYaml . | nindent 4 }}' \
		| $(SED) -e '/^  annotations:/a \ \ \ \ {{- with .Values.crds.annotations }}' \
 		| $(SED) -e '/^  annotations:/i \ \ labels:' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- end }}' \
 		| $(SED) -e '/^  labels:/a \ \ \ \ {{- toYaml . | nindent 4 }}' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- with .Values.crds.labels }}' \
		| $(SED) -e '/^  labels:/a \ \ \ \ {{- include "kyverno-authz-server.labels" . | nindent 4 }}' \
 		> ./charts/kyverno-authz-server/templates/crds.yaml

.PHONY: codegen-api-docs
codegen-api-docs: ## Generate markdown API docs
codegen-api-docs: $(REFERENCE_DOCS)
codegen-api-docs: codegen-crds
	@echo Generate api docs... >&2
	@rm -rf ./website/docs/reference/apis
	@cd ./website/apis && $(REFERENCE_DOCS) -c config.yaml -f markdown -o ../docs/reference/apis

.PHONY: codegen-schemas-openapi
codegen-schemas-openapi: ## Generate openapi schemas (v2 and v3)
codegen-schemas-openapi: CURRENT_CONTEXT = $(shell kubectl config current-context)
codegen-schemas-openapi: codegen-crds
codegen-schemas-openapi: $(KIND)
	@echo Generate openapi schema... >&2
	@rm -rf ./.temp/.schemas
	@mkdir -p ./.temp/.schemas/openapi/v2
	@mkdir -p ./.temp/.schemas/openapi/v3/apis/envoy.kyverno.io
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@kubectl create -f $(CRDS_PATH)
	@sleep 15
	@kubectl get --raw /openapi/v2 > ./.temp/.schemas/openapi/v2/schema.json
	@kubectl get --raw /openapi/v3/apis/envoy.kyverno.io/v1alpha1 > ./.temp/.schemas/openapi/v3/apis/envoy.kyverno.io/v1alpha1.json
	@$(KIND) delete cluster --name schema
	@kubectl config use-context $(CURRENT_CONTEXT) || true

.PHONY: codegen-schemas-json
codegen-schemas-json: ## Generate json schemas
codegen-schemas-json: codegen-schemas-openapi
	@echo Generate json schema... >&2
	@$(PIP) install -r requirements.txt
	@rm -rf ./.temp/.schemas/json
	@rm -rf ./.schemas/json
	@openapi2jsonschema .temp/.schemas/openapi/v3/apis/envoy.kyverno.io/v1alpha1.json --kubernetes --strict --stand-alone --expanded -o ./.temp/.schemas/json
	@mkdir -p ./.schemas/json
	@cp ./.temp/.schemas/json/authorizationpolicy-envoy-*.json ./.schemas/json

.PHONY: codegen
codegen: ## Rebuild all generated code and docs
codegen: codegen-mkdocs
codegen: codegen-crds
codegen: codegen-helm-crds
codegen: codegen-helm-docs
codegen: codegen-api-docs
codegen: codegen-schemas-openapi
codegen: codegen-schemas-json

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
	@echo Build... >&2
	@LD_FLAGS=$(LD_FLAGS) go build .

##############
# BUILD (KO) #
##############

.PHONY: ko-login
ko-login: $(KO)
	@$(KO) login $(REGISTRY) --username $(REGISTRY_USERNAME) --password $(REGISTRY_PASSWORD)

.PHONY: ko-build
ko-build: ## Build Docker image with ko
ko-build: fmt
ko-build: vet
ko-build: $(KO)
	@echo Build Docker image with ko... >&2
	@LD_FLAGS=$(LD_FLAGS) KO_DOCKER_REPO=$(KO_REGISTRY) $(KO) build . --preserve-import-paths --tags=$(KO_TAGS)

.PHONY: ko-publish
ko-publish: ## Publish Docker image with ko
ko-publish: fmt
ko-publish: vet
ko-publish: ko-login
ko-publish: $(KO)
	@echo Publish Docker image with ko... >&2
	@LD_FLAGS=$(LD_FLAGS) KO_DOCKER_REPO=$(REGISTRY)/$(REPO)/$(IMAGE) $(KO) build . --bare --tags=$(KO_TAGS) --platform=$(KO_PLATFORMS)

##########
# DOCKER #
##########

.PHONY: docker-save-image
docker-save-image: ## Save docker image in archive
	@docker save $(KO_REGISTRY)/$(PACKAGE):$(GIT_SHA) > image.tar

.PHONY: docker-load-image
docker-load-image: ## Load docker image in archive
	@docker load --input image.tar

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
	@echo Generate and servemkdocs website... >&2
	@$(PIP) install -r requirements.txt
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
	@echo Load image in kind... >&2
	@$(KIND) load docker-image $(KO_REGISTRY)/$(PACKAGE):$(GIT_SHA)

################
# CERTIFICATES #
################

.PHONY: generate-certs
generate-certs: ## Generate certificates
generate-certs:
	@echo Generating certificates... >&2
	@rm -rf .certs
	@mkdir -p .certs
	@openssl req -new -x509  \
        -subj "/CN=kyverno-sidecar-injector.kyverno.svc" \
        -addext "subjectAltName = DNS:kyverno-sidecar-injector.kyverno.svc" \
        -nodes -newkey rsa:4096 -keyout .certs/tls.key -out .certs/tls.crt 

################
# CERT MANAGER #
################

.PHONY: install-cert-manager
install-cert-manager: ## Install cert-manager
install-cert-manager: $(HELM)
	@echo Install cert-manager... >&2
	@$(HELM) upgrade --install cert-manager --namespace cert-manager --create-namespace --wait --repo https://charts.jetstack.io cert-manager \
		--set crds.enabled=true

.PHONY: install-cluster-issuer
install-cluster-issuer: ## Install cert-manager cluster issuer
install-cluster-issuer:
	@echo Install cert-manager cluster issuer... >&2
	@kubectl apply -f .manifests/cert-manager/cluster-issuer.yaml

#########
# ISTIO #
#########

.PHONY: install-istio
install-istio: ## Install istio
install-istio: $(HELM)
	@echo Install istio... >&2
	@$(HELM) upgrade --install istio-base --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts base
	@$(HELM) upgrade --install istiod --namespace istio-system --create-namespace --wait --repo https://istio-release.storage.googleapis.com/charts istiod

########
# HELM #
########

.PHONY: install-kyverno-sidecar-injector
install-kyverno-sidecar-injector: ## Install kyverno-sidecar-injector chart
install-kyverno-sidecar-injector: kind-load-image
install-kyverno-sidecar-injector: $(HELM)
	@echo Build kyverno-sidecar-injector dependecy... >&2
	@$(HELM) dependency build --skip-refresh ./charts/kyverno-sidecar-injector
	@echo Install kyverno-sidecar-injector chart... >&2
	@$(HELM) upgrade --install kyverno-sidecar-injector --namespace kyverno --create-namespace --wait ./charts/kyverno-sidecar-injector \
		--set containers.injector.image.registry=$(KO_REGISTRY) \
		--set containers.injector.image.repository=$(PACKAGE) \
		--set containers.injector.image.tag=$(GIT_SHA) \
		--set certificates.certManager.issuerRef.name=selfsigned-issuer \
		--set certificates.certManager.issuerRef.kind=ClusterIssuer \
		--set certificates.certManager.issuerRef.group=cert-manager.io

.PHONY: install-kyverno-authz-server
install-kyverno-authz-server: ## Install kyverno-authz-server chart
install-kyverno-authz-server: kind-load-image
install-kyverno-authz-server: $(HELM)
	@echo Build kyverno-authz-server dependecy... >&2
	@$(HELM) dependency build --skip-refresh ./charts/kyverno-authz-server
	@echo Install kyverno-authz-server chart... >&2
	@$(HELM) upgrade --install kyverno-authz-server --namespace kyverno --create-namespace --wait ./charts/kyverno-authz-server \
		--set containers.server.image.registry=$(KO_REGISTRY) \
		--set containers.server.image.repository=$(PACKAGE) \
		--set containers.server.image.tag=$(GIT_SHA)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
