[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go)

# Kubewarden Go Policy SDK

This module provides a SDK that can be used to write [Kubewarden
Policies](https://github.com/kubewarden/) using the Go programming
language.

Due to current Go compiler limitations, Go policies must be built
using [TinyGo](https://github.com/tinygo-org/tinygo).

## Validation

This SDK provides helper methods to accept and reject validation requests.

Mutation policies cannot be written using this SDK yet.

## Logging

Policies can generate log messages that are then propagated to the host
environment (eg: [kwctl](https://github.com/kubewarden/kwctl),
[policy-server](https://github.com/kubewarden/policy-server)).

This Go module provides logging capabilities that integrate with the
[onelog](https://github.com/francoispqt/onelog) project.

This logging solution has been chosen because:

  * It works also with WebAssembly binaries. Other popular logging solutions
    cannot even be built to WebAssembly.
  * It provides [good performance](https://github.com/francoispqt/onelog#benchmarks)
  * It supports structured logging.

### Usage

The instructions provided by the official
[onelog](https://github.com/francoispqt/onelog) project apply also to Kubewarden
policies.

The `onelog.Logger` instance must be configured to use a `KubewardenLogWriter`
object.

```go
	kl := kubewarden.KubewardenLogWriter{}
	logger := onelog.New(
		&kl,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
	)
	logger.Info("info message from tinygo")
	logger.DebugWithFields("i'm not sure what's going on", func(e onelog.Entry) {
		e.String("string", "foobar")
		e.Int("int", 12345)
		e.Int64("int64", 12345)
		e.Float("float64", 0.15)
		e.Bool("bool", true)
		e.Err("err", errors.New("someError"))
		e.ObjectFunc("user", func(e onelog.Entry) {
			e.String("name", "somename")
		})
	})
```


## Testing

[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing)

The `kubewarden/policy-sdk-go/testing` module provides some test helpers
that simplify the process of writing unit tests.

# Project template

We provide a GitHub repository template that can be used to quickly
scaffold a new Kubewarden policy writing Go.

This can be found [here](https://github.com/kubewarden/go-policy-template).
