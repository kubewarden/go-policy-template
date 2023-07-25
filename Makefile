SOURCE_FILES := $(shell find . -type f -name '*.go')
# It's necessary to call cut because kwctl command does not handle version
# starting with v.
VERSION ?= $(shell git describe | cut -c2-)


policy.wasm: $(SOURCE_FILES) go.mod go.sum
	docker run \
		--rm \
		-e GOFLAGS="-buildvcs=false" \
		-v ${PWD}:/src \
		-w /src tinygo/tinygo:0.28.1 \
		tinygo build -o policy.wasm -target=wasi -no-debug .

artifacthub-pkg.yml: metadata.yml go.mod
	$(warning If you are updating the artifacthub-pkg.yml file for a release, \
	  remember to set the VERSION variable with the proper value. \
	  To use the latest tag, use the following command:  \
	  make VERSION=$$(git describe --tags --abbrev=0 | cut -c2-) annotated-policy.wasm)
	kwctl scaffold artifacthub \
	  --metadata-path metadata.yml --version $(VERSION) \
	  --output artifacthub-pkg.yml

annotated-policy.wasm: policy.wasm metadata.yml artifacthub-pkg.yml
	kwctl annotate -m metadata.yml -u README.md -o annotated-policy.wasm policy.wasm

.PHONY: test
test:
	go test -v

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats e2e.bats

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm artifacthub-pkg.yml
