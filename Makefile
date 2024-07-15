ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := $(abspath $(ROOT_DIR)/bin)

SOURCE_FILES := $(shell find . -type f -name '*.go')
# It's necessary to call cut because kwctl command does not handle version
# starting with v.
VERSION ?= $(shell git describe | cut -c2-)

GOLANGCI_LINT_VER := v1.59.1
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(BIN_DIR)/$(GOLANGCI_LINT_BIN)


policy.wasm: $(SOURCE_FILES) go.mod go.sum
	docker run \
		--rm \
		-e GOFLAGS="-buildvcs=false" \
		-v ${PWD}:/src \
		-w /src tinygo/tinygo:0.31.2 \
		tinygo build -o policy.wasm -target=wasi -no-debug .

artifacthub-pkg.yml: metadata.yml go.mod
	$(warning If you are updating the artifacthub-pkg.yml file for a release, \
	  remember to set the VERSION variable with the proper value. \
	  To use the latest tag, use the following command:  \
	  make VERSION=$$(git describe --tags --abbrev=0 | cut -c2-) annotated-policy.wasm)
	kwctl scaffold artifacthub \
	  --metadata-path metadata.yml --version $(VERSION) \
	  --output artifacthub-pkg.yml

annotated-policy.wasm: policy.wasm metadata.yml
	kwctl annotate -m metadata.yml -u README.md -o annotated-policy.wasm policy.wasm

.PHONY: test
test:
	go test -v

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats e2e.bats

golangci-lint: $(GOLANGCI_LINT) ## Install a local copy of golang ci-lint.
$(GOLANGCI_LINT): ## Install golangci-lint.
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	go vet ./...
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm artifacthub-pkg.yml

.PHONY: fmt
fmt:
	go fmt ./...
