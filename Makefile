SOURCE_FILES := $(shell find . -type f -name '*.go')
VERSION := $(shell git describe --exact-match --tags $(git log -n1 --pretty='%h') | cut -c2-)

policy.wasm: $(SOURCE_FILES) go.mod go.sum types_easyjson.go
	docker run \
		--rm \
		-e GOFLAGS="-buildvcs=false" \
		-v ${PWD}:/src \
		-w /src tinygo/tinygo:0.23.0 \
		tinygo build -o policy.wasm -target=wasi -no-debug .

artifacthub-pkg.yml: metadata.yml go.mod
	kwctl scaffold artifacthub \
	    --metadata-path metadata.yml --version $(VERSION) \
		--questions-path questions-ui.yml > artifacthub-pkg.yml.tmp \
	&& mv artifacthub-pkg.yml.tmp artifacthub-pkg.yml \
	|| rm -f artifacthub-pkg.yml.tmp

annotated-policy.wasm: policy.wasm metadata.yml artifacthub-pkg.yml
	kwctl annotate -m metadata.yml -u README.md -o annotated-policy.wasm policy.wasm

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

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm
