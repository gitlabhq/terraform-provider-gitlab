default: reviewable

reviewable: build fmt generate test ## Run before committing.

GOBIN = $(shell pwd)/bin

build: ## Build the provider binary.
	go mod tidy
	GOBIN=$(GOBIN) go install

generate: tools ## Generate files to be checked in.
	@# Setting empty environment variables to work around issue: https://github.com/hashicorp/terraform-plugin-docs/issues/12
	GITLAB_TOKEN="" $(GOBIN)/tfplugindocs generate

ifdef RUN
TESTARGS += -test.run $(RUN)
endif

test: ## Run unit tests.
	go test $(TESTARGS) ./gitlab

TFPROVIDERLINTX_CHECKS = -XAT001=false -XR003=false -XS002=false

fmt: tools terraform ## Format files and fix issues.
	gofmt -w -s .
	$(GOBIN)/golangci-lint run --fix
	$(GOBIN)/tfproviderlintx $(TFPROVIDERLINTX_CHECKS) --fix ./...
	$(TERRAFORM) fmt -recursive -list ./examples
	$(GOBIN)/shfmt -l -s -w ./examples

lint-golangci: tools ## Run golangci-lint linter (same as fmt but without modifying files).
	$(GOBIN)/golangci-lint run

lint-tfprovider: tools ## Run tfproviderlintx linter (same as fmt but without modifying files).
	$(GOBIN)/tfproviderlintx $(TFPROVIDERLINTX_CHECKS) ./...

lint-examples-tf: terraform ## Run terraform linter on examples (same as fmt but without modifying files).
	$(TERRAFORM) fmt -recursive -check ./examples

lint-examples-sh: tools ## Run shell linter on examples (same as fmt but without modifying files).
	$(GOBIN)/shfmt -l -s -d ./examples

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

.PHONY: tools
tools:
	GOBIN=$(GOBIN) go generate ./tools/tools.go

TERRAFORM = $(GOBIN)/terraform
TERRAFORM_VERSION = v1.1.4
terraform:
	@# See https://github.com/hashicorp/terraform/issues/30356
	@[ -f $(TERRAFORM) ] || { mkdir -p tmp; cd tmp; rm -rf terraform; git clone --branch $(TERRAFORM_VERSION) --depth 1 https://github.com/hashicorp/terraform.git; cd terraform; GOBIN=$(GOBIN) go install; cd ..; rm -rf terraform; }
