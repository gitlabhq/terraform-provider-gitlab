TEST?=./gitlab
SERVICE?=gitlab-ce
GITLAB_TOKEN?=ACCTEST1234567890123
GITLAB_BASE_URL?=http://127.0.0.1:8080/api/v4
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

ifdef RUN
TESTARGS += -test.run $(RUN)
endif

default: build

build:
	go install

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc-up:
	docker-compose up -d $(SERVICE)
	./scripts/await-healthy.sh

testacc-down:
	docker-compose down

testacc:
	TF_ACC=1 GITLAB_TOKEN=$(GITLAB_TOKEN) GITLAB_BASE_URL=$(GITLAB_BASE_URL) go test -v $(TEST) $(TESTARGS) -timeout 40m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

tfproviderlint:
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
	-XAT001=false -XR003=false -XR005=false -XS001=false -XS002=false \
	./...

.PHONY: default build test testacc-up testacc-down testacc vet fmt tfproviderlint
