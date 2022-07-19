SOURCE_FILES := $(shell find . -type f -name '*.go')
BIN_DIR := $(abspath bin)

KWCTL_VER := v0.2.5
KWCTL_BIN := kwctl
KWCTL := $(BIN_DIR)/$(KWCTL_BIN)

## --------------------------------------
## Tools
## --------------------------------------

kwctl: $(KWCTL) ## Install a local copy of kwctl

$(KWCTL): ## Install kwctl.
	./hack/ensure-kwctl.sh $(KWCTL_VER)

## --------------------------------------
## Build
## --------------------------------------

policy.wasm: $(SOURCE_FILES) go.mod go.sum types_easyjson.go
	docker run \
		--rm \
		-e GOFLAGS="-buildvcs=false" \
		-v ${PWD}:/src \
		-w /src tinygo/tinygo:0.23.0 \
		tinygo build -o policy.wasm -target=wasi -no-debug .

annotated-policy.wasm: $(KWCTL) policy.wasm metadata.yml
	$(KWCTL) annotate -m metadata.yml -o annotated-policy.wasm policy.wasm

## --------------------------------------
## Tests
## --------------------------------------

.PHONY: generate-easyjson
types_easyjson.go: types.go
	docker run \
		--rm \
		-v ${PWD}:/src \
		-w /src \
		golang:1.17-alpine ./hack/generate-easyjson.sh

.PHONY: test
test: types_easyjson.go
	go test -v

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats e2e.bats

## --------------------------------------
## Cleanup
## --------------------------------------

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm bin/
