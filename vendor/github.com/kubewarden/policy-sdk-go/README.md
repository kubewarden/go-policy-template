# Kubewarden Go Policy SDK

This module provides a SDK that can be used to write [Kubewarden
Policies](https://github.com/kubewarden/) using the Go programming
language.

Due to current Go compiler limitations, Go policies must be built
using [TinyGo](https://github.com/tinygo-org/tinygo).

## Building your policy with TinyGo

Assuming your policy tree looks like:

```
.
├── go.mod
├── go.sum
├── main.go
```

You can build the project with TinyGo from within that folder:

```shell
docker run --rm -v ${PWD}:/src -w /src tinygo/tinygo:0.17.0 tinygo build -o my-policy.wasm -target=wasi -no-debug .
```

Resulting in:

```
.
├── go.mod
├── go.sum
├── main.go
├── my-policy.wasm
```
