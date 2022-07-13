default: reviewable

reviewable: build fmt generate test ## Run before committing.

GOBIN = $(shell pwd)/bin
PROVIDER_SRC_DIR := ./internal/provider
TERRAFORM_PLUGIN_DIR ?= ~/.terraform.d/plugins/gitlab.local/x/gitlab/99.99.99
TERRAFORM_PLATFORM_DIR ?= darwin_amd64

build: ## Build the provider binary.
	go mod tidy
	GOBIN=$(GOBIN) go install

local: build ## Build and Install the provider locally
	mkdir -p $(TERRAFORM_PLUGIN_DIR)/$(TERRAFORM_PLATFORM_DIR)
	cp -f $(GOBIN)/terraform-provider-gitlab $(TERRAFORM_PLUGIN_DIR)/$(TERRAFORM_PLATFORM_DIR)/terraform-provider-gitlab

generate: tool-tfplugindocs ## Generate files to be checked in.
	@# Setting empty environment variables to work around issue: https://github.com/hashicorp/terraform-plugin-docs/issues/12
	@# Setting the PATH so that tfplugindocs uses the same terraform binary as other targets here, and to resolve a "Error: Incompatible provider version" error on M1 macs.
	GITLAB_TOKEN="" PATH="$(GOBIN):$(PATH)" $(GOBIN)/tfplugindocs generate

ifdef RUN
TESTARGS += -test.run $(RUN)
endif

test: ## Run unit tests.
	go test $(TESTARGS) $(PROVIDER_SRC_DIR)

fmt: tool-golangci-lint tool-terraform tool-shfmt tfproviderlint-plugin ## Format files and fix issues.
	gofmt -w -s .
	$(GOBIN)/golangci-lint run --build-tags acceptance --fix
	$(GOBIN)/terraform fmt -recursive -list ./examples
	$(GOBIN)/shfmt -l -s -w ./examples

lint-golangci: tool-golangci-lint tfproviderlint-plugin ## Run golangci-lint linter (same as fmt but without modifying files).
	$(GOBIN)/golangci-lint run --build-tags acceptance

lint-examples-tf: tool-terraform ## Run terraform linter on examples (same as fmt but without modifying files).
	$(GOBIN)/terraform fmt -recursive -check ./examples

lint-examples-sh: tool-shfmt ## Run shell linter on examples (same as fmt but without modifying files).
	$(GOBIN)/shfmt -l -s -d ./examples

lint-generated: generate ## Check that "make generate" was called. Note this only works if the git workspace is clean.
	@echo "Checking git status"
	@[ -z "$(shell git status --short)" ] || { \
		echo "Error: Files should have been generated:"; \
		git status --short; echo "Diff:"; \
		git --no-pager diff HEAD; \
		echo "Run \"make generate\" and try again"; \
		exit 1; \
	}

lint-custom: ## Run custom checks and validations that do not fit into an existing lint framework.
	@./scripts/lint-custom.sh

apicovered: tool-apicovered ## Run an analysis tool to estimate the GitLab API coverage.
	@$(GOBIN)/apicovered ./gitlab

apiunused: tool-apiunused ## Run an analysis tool to output unused parts of the go-gitlab package.
	@$(GOBIN)/apiunused ./gitlab

SERVICE ?= gitlab-ce
GITLAB_TOKEN ?= ACCTEST1234567890123
GITLAB_BASE_URL ?= http://127.0.0.1:8080/api/v4

testacc-up: | certs ## Launch a GitLab instance.
	docker-compose up -d $(SERVICE)
	./scripts/await-healthy.sh

testacc-down: ## Teardown a GitLab instance.
	docker-compose down --volumes

testacc: ## Run acceptance tests against a GitLab instance.
	TF_ACC=1 GITLAB_TOKEN=$(GITLAB_TOKEN) GITLAB_BASE_URL=$(GITLAB_BASE_URL) go test --tags acceptance -v $(PROVIDER_SRC_DIR) $(TESTARGS) -timeout 40m

certs: ## Generate certs for the GitLab container registry
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes -keyout certs/gitlab-registry.key -out certs/gitlab-registry.crt -subj "/CN=gitlab-registry.com" -addext "subjectAltName=DNS:IP:127.0.0.1"

# TOOLS
# Tool dependencies are installed into a project-local /bin folder.

tool-golangci-lint:
	@$(call install-tool, github.com/golangci/golangci-lint/cmd/golangci-lint)

tool-tfplugindocs:
	@$(call install-tool, github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs)

tool-shfmt:
	@$(call install-tool, mvdan.cc/sh/v3/cmd/shfmt)

tool-apicovered:
	@$(call install-tool, ./cmd/apicovered)

tool-apiunused:
	@$(call install-tool, ./cmd/apiunused)

define install-tool
	cd tools && GOBIN=$(GOBIN) go install $(1)
endef

TERRAFORM_VERSION = v1.1.4
tool-terraform:
	@# See https://github.com/hashicorp/terraform/issues/30356
	@[ -f $(GOBIN)/terraform ] || { mkdir -p tmp; cd tmp; rm -rf terraform; git clone --branch $(TERRAFORM_VERSION) --depth 1 https://github.com/hashicorp/terraform.git; cd terraform; GOBIN=$(GOBIN) go install; cd ..; rm -rf terraform; }

clean: testacc-down
	@rm -rf certs/

tfproviderlint-plugin:
	@cd tools && go build -buildmode=plugin -o $(GOBIN)/tfproviderlint-plugin.so ./cmd/tfproviderlint-plugin
