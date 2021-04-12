wasm: go.mod go.sum *.go
	docker run --rm -v ${PWD}:/src -w /src tinygo/tinygo-dev:latest tinygo build -o policy.wasm -target=wasi -no-debug .

test:
	go test -v

e2e-tests:
	bats e2e.bats

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm
