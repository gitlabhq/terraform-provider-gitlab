default: reviewable

reviewable: build fmt generate test ## Run before committing.

GOBIN = $(shell pwd)/bin

build: ## Build the provider binary.
	go mod tidy
	GOBIN=$(GOBIN) go install

generate: tool-tfplugindocs ## Generate files to be checked in.
	@# Setting empty environment variables to work around issue: https://github.com/hashicorp/terraform-plugin-docs/issues/12
	GITLAB_TOKEN="" $(GOBIN)/tfplugindocs generate

ifdef RUN
TESTARGS += -test.run $(RUN)
endif

test: ## Run unit tests.
	go test $(TESTARGS) ./gitlab

TFPROVIDERLINTX_CHECKS = -XAT001=false -XR003=false -XS002=false

fmt: tool-golangci-lint tool-tfproviderlintx tool-terraform tool-shfmt ## Format files and fix issues.
	gofmt -w -s .
	$(GOBIN)/golangci-lint run --fix
	$(GOBIN)/tfproviderlintx $(TFPROVIDERLINTX_CHECKS) --fix ./...
	$(GOBIN)/terraform fmt -recursive -list ./examples
	$(GOBIN)/shfmt -l -s -w ./examples

lint-golangci: tool-golangci-lint ## Run golangci-lint linter (same as fmt but without modifying files).
	$(GOBIN)/golangci-lint run

lint-tfprovider: tool-tfproviderlintx ## Run tfproviderlintx linter (same as fmt but without modifying files).
	$(GOBIN)/tfproviderlintx $(TFPROVIDERLINTX_CHECKS) ./...

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

SERVICE ?= gitlab-ce
GITLAB_TOKEN ?= ACCTEST1234567890123
GITLAB_BASE_URL ?= http://127.0.0.1:8080/api/v4

testacc-up: ## Launch a GitLab instance.
	docker-compose up -d $(SERVICE)
	./scripts/await-healthy.sh

testacc-down: ## Teardown a GitLab instance.
	docker-compose down

testacc: ## Run acceptance tests against a GitLab instance.
	TF_ACC=1 GITLAB_TOKEN=$(GITLAB_TOKEN) GITLAB_BASE_URL=$(GITLAB_BASE_URL) go test -v ./gitlab $(TESTARGS) -timeout 40m

# TOOLS
# Tool dependencies are installed into a project-local /bin folder.

tool-golangci-lint:
	@$(call install-tool, github.com/golangci/golangci-lint/cmd/golangci-lint)

tool-tfproviderlintx:
	@$(call install-tool, github.com/bflad/tfproviderlint/cmd/tfproviderlintx)

tool-tfplugindocs:
	@$(call install-tool, github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs)

tool-shfmt:
	@$(call install-tool, mvdan.cc/sh/v3/cmd/shfmt)

define install-tool
	cd tools && GOBIN=$(GOBIN) go install $(1)
endef

TERRAFORM_VERSION = v1.1.4
tool-terraform:
	@# See https://github.com/hashicorp/terraform/issues/30356
	@[ -f $(GOBIN)/terraform ] || { mkdir -p tmp; cd tmp; rm -rf terraform; git clone --branch $(TERRAFORM_VERSION) --depth 1 https://github.com/hashicorp/terraform.git; cd terraform; GOBIN=$(GOBIN) go install; cd ..; rm -rf terraform; }
