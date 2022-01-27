SOURCE_FILES := $(shell find . -type f -name '*.go')

policy.wasm: $(SOURCE_FILES) go.mod go.sum
	docker run --rm -v ${PWD}:/src -w /src tinygo/tinygo:0.18.0 tinygo build \
		-o policy.wasm -target=wasi -no-debug .

annotated-policy.wasm: policy.wasm metadata.yml
	kwctl annotate -m metadata.yml -o annotated-policy.wasm policy.wasm

.PHONY: test
test:
	go test -v

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats e2e.bats

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm
