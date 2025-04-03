[![Stable](https://img.shields.io/badge/status-stable-brightgreen?style=for-the-badge)](https://github.com/kubewarden/community/blob/main/REPOSITORIES.md#stable)
[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go)

> Don't forget to take a look at [Kubewarden's official documentation](https://docs.kubewarden.io).
> The docs cover step-by-step instructions about how to write policies.

# Kubewarden Go Policy SDK

This module provides a SDK that can be used to write [Kubewarden
Policies](https://github.com/kubewarden/) using the Go programming
language.

Due to current Go compiler limitations, Go policies must be built
using [TinyGo](https://github.com/tinygo-org/tinygo).

## Known limitations of TinyGo

TinyGo doesn't have full support of the Go Standard Library. However, this
shouldn't pose a limit to policy authors.

# Validation

This SDK provides helper methods to accept and reject validation requests.

A validation policy consists of these steps:

1. Extract the object to inspect from the incoming payload
2. (Optional) Extract the settings object from the incoming payload.
3. Validation code
4. Communicate the outcome of the validation: accept/reject

The 4th step is done using helper functions provided by this SDK.

As for the 1st step, there are two approaches that can be used.

## Perform "jq-like" searches

The policy receives as input a `payload` parameter of type `[]byte`. This
contains the JSON document described [here](https://docs.kubewarden.io/writing-policies/spec/validating-policies#the-validationrequest-object).

Policy authors can leverage the [`github.com/tidwall/gjson`](https://github.com/tidwall/gjson)
package to search for data inside of this JSON document.

For example, assume you want to validate the `labels` that are inside of
a Kubernetes object. This can be done with this snippet:

```go
data := gjson.GetBytes(
  payload,
  "request.object.metadata.labels")

labels := mapset.NewThreadUnsafeSet()
denied_labels_violations := []string{}
constrained_labels_violations := []string{}

data.ForEach(func(key, value gjson.Result) bool {
  label := key.String()
  labelValue := value.String()
  // do something with label and labelValue

  return true
})
```

This _"jq-like"_ approach can be pretty handy when the policy has to look
deep inside of a Kubernetes object. This is especially helpful when dealing with
inner objects that are optional.

## Use native Go types

The majority of policies target a specific type of Kubernetes resource, like
Pod, Ingress, Service and similar. Because of that, another possible approach
is to unmarshal the incoming object into a native Go type.

TinyGo doesn't yet support the full Go Standard Library, plus it has limited
support of Go reflection.
Because of that, it is not possible to import the official Kubernetes Go library
from upstream (e.g.: `k8s.io/api/core/v1`).
Importing these official Kubernetes types will result in a compilation failure.

Moreover, Kubewarden provides TinyGo friendly Go types for all the Kubernetes
types inside of the [`github.com/kubewarden/k8s-objects`](https://github.com/kubewarden/k8s-objects)
package.

Using this SDK requires **TinyGo 0.28.1 or later.**

> **Warning**
> Using an older version of TinyGo will result in runtime errors due to the limited support for Go reflection.

### Example

This snippet shows how to implement a `validation` function that uses the
"native Go types" approach:

```go
// Create a ValidationRequest instance from the incoming payload
validationRequest := kubewarden_protocol.ValidationRequest{}
err := json.Unmarshal(payload, &validationRequest)
if err != nil {
	return kubewarden.RejectRequest(
		kubewarden.Message(err.Error()),
		kubewarden.Code(400))
}

// Access the **raw** JSON that describes the object
ingressJSON := validationRequest.Request.Object

// Try to create an Ingress instance using the RAW JSON we got from the
// ValidationRequest.
// This policy works only against Ingress objects, if the creation fails
// we reject the request and provide a meaningful error.
ingress := &networkingv1.Ingress{}
if err := json.Unmarshal([]byte(ingressJSON), ingress); err != nil {
	return kubewarden.RejectRequest(
		kubewarden.Message(
		fmt.Sprintf("Cannot decode Ingress object: %s", err.Error())),
		kubewarden.Code(400))
}

// the validation logic
```

**Note:** the `github.com/kubewarden/k8s-objects` package is organized
in the same way as the official `k8s.io` one.

# Mutating policy

Mutation policies works exactly like the validation ones. The only difference
is that, when a request has to be accepted AND mutated, the policy must return
the input object with all the required changes applied.

Mutation policies can be done by leveraging the Kubernetes Go types
defined inside of the `github.com/kubewarden/k8s-objects` package and
the helper methods provided by this SDK.

The following example defines a mutating policy that always changes the name of
Ingress objects:

```go
import (
    "encoding/json"
	"fmt"

	networkingv1 "github.com/kubewarden/k8s-objects/api/networking/v1"
	kubewarden "github.com/kubewarden/policy-sdk-go"

)

func validate(payload []byte) ([]byte, error) {
  // Create a ValidationRequest instance from the incoming payload
  validationRequest := kubewarden_protocol.ValidationRequest{}
  err := json.Unmarshal(payload, &validationRequest)
  if err != nil {
    return kubewarden.RejectRequest(
      kubewarden.Message(err.Error()),
      kubewarden.Code(400))
  }

  // Access the **raw** JSON that describes the object
  ingressJSON := validationRequest.Request.Object

  // Try to create a Ingress instance using the RAW JSON we got from the
  // ValidationRequest.
  // This policy works only against Ingress objects, if the creation fails
  // we reject the request and provide a meaningful error.
  ingress := &networkingv1.Ingress{}
  if err := json.Unmarshal([]byte(ingressJSON), ingress); err != nil {
    return kubewarden.RejectRequest(
      kubewarden.Message(
      fmt.Sprintf("Cannot decode Ingress object: %s", err.Error())),
      kubewarden.Code(400))
  }

  ingress.Metadata.Name = fmt.Sprintf("%s-changed", ingress.Metadata.Name)

  return kubewarden.MutateRequest(ingress)
}
```

# Logging

Policies can generate log messages that are then propagated to the host
environment (eg: [kwctl](https://github.com/kubewarden/kwctl),
[policy-server](https://github.com/kubewarden/policy-server)).

This Go module provides logging capabilities that integrate with the
[onelog](https://github.com/francoispqt/onelog) project.

This logging solution has been chosen because:

- It works also with WebAssembly binaries. Other popular logging solutions
  cannot even be built to WebAssembly.
- It provides [good performance](https://github.com/francoispqt/onelog#benchmarks)
- It supports structured logging.

## Usage

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

# Host capabilities

The policy executor exposes additional capabilities that can be leveraged by the
guest.

These capabilities are exposed to the Go policies via this SDK, through the `Host`
type defined inside of `github.com/kubewarden/policy-sdk-go/capabilities`.

## Get OCI manifest digest

The policy can request the digest of an OCI manifest. This can be used to
get the immutable reference of a container Image or anything that is stored
inside of a container registry (e.g. Kubewarden Policies, Helm charts,...).

```go
host := capabilities.NewHost()
digest, err := host.GetOCIManifestDigest("busybox:latest")
```

## Hostname DNS lookup

The policy can lookup the addresses for a given hostname by using the
DNS resolvers of the host that is evaluating the policy.

```go
host := capabilities.NewHost()
ips, err := host.LookupHost("kubewarden.io")
```

## Sigstore verification

The policy can ask the host to perform verification operations against
objects stored inside of container registries (e.g. container image, kubewarden
policy, helm chart,...) leveraging the [Sigstore](https://sigstore.dev) primitives.

Currently this SDK exposes helper function that can perform verification using
public keys and using the Sigstore keyless mechanism.

# Testing

[![GoDoc](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing?status.svg)](https://godoc.org/github.com/kubewarden/policy-sdk-go/testing)

The `kubewarden/policy-sdk-go/testing` module provides some test helpers
that simplify the process of writing unit tests.

## Host capabilities

The Go unit tests of a policy are **not** run inside of a WebAssembly environment,
they are instead built into a native executable using the official Go compiler.

Because of that, at test time the host capabilities have to be mocked. This is also
useful to write ad-hoc tests that can handle different kind of responses coming
from the host.

The `Host` type described above relies on an internal `waPC` client that
interacts with the host. At test time, the client is an instance of
`MockWapcClient`.

Developers can create `MockWapcClient` instances using the `NewMockWapcClient`
helper method from the `capabilities` package. It is used as follows:

```go
mockWapcClient := &mocks.MockWapcClient{}
mockWapcClient.On("HostCall", "kubewarden", "kubernetes", "get_resource", request).Return(wapcResponse, nil)
```

# Project template

We provide a GitHub repository template that can be used to quickly
scaffold a new Kubewarden policy writing Go.

This can be found [here](https://github.com/kubewarden/go-policy-template).
