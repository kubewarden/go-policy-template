ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := $(abspath $(ROOT_DIR)/bin)

MOCKERY_VER := v2.43.2
MOCKERY_BIN := mockery
MOCKERY := $(BIN_DIR)/$(MOCKERY_BIN)

GOLANGCI_LINT_VER := v2.0.2
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(BIN_DIR)/$(GOLANGCI_LINT_BIN)

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

generate-mocks: mockery
	mockery

golangci-lint: $(GOLANGCI_LINT) ## Install a local copy of golang ci-lint.
$(GOLANGCI_LINT): ## Install golangci-lint.
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VER)

mockery: $(MOCKERY) ## Install a local copy of mockery.
$(MOCKERY): ## Install mockery.
	GOBIN=$(BIN_DIR) go install github.com/vektra/mockery/v2@$(MOCKERY_VER)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	go vet ./...
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

