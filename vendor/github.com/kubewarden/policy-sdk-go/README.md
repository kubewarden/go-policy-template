[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go)

# Kubewarden Go Policy SDK

This module provides a SDK that can be used to write [Kubewarden
Policies](https://github.com/kubewarden/) using the Go programming
language.

Due to current Go compiler limitations, Go policies must be built
using [TinyGo](https://github.com/tinygo-org/tinygo).

## Testing

[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing)

The `kubewarden/policy-sdk-go/testing` module provides some test helpers
that simplify the process of writing unit tests.

# Project template

We provide a GitHub repository template that can be used to quickly
scaffold a new Kubewarden policy writing Go.

This can be found [here](https://github.com/kubewarden/go-policy-template).
