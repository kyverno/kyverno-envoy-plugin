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

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.22.0
TOOLS                              := $(KIND)
PIP                                ?= "pip"
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)
@
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
codegen-mkdocs: codegen-cli-docs
codegen-mkdocs: codegen-api-docs
codegen-mkdocs: codegen-jp-docs
	@echo Generate mkdocs website... >&2
	@$(PIP) install mkdocs
	@$(PIP) install --upgrade pip
	@$(PIP) install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin mike
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen
codegen: ## Rebuild all generated code and docs
codegen: codegen-api-docs
codegen: codegen-cli-docs
codegen: codegen-crds
codegen: codegen-deepcopy
codegen: codegen-mkdocs
codegen: codegen-register
codegen: codegen-schemas-json
codegen: codegen-tests-catalog

.PHONY: verify-codegen
verify-codegen: ## Verify all generated code and docs are up to date
verify-codegen: codegen
	@echo Checking codegen is up to date... >&2
	@git --no-pager diff -- .
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen".' >&2
	@echo 'To correct this, locally run "make codegen", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- .

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
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
