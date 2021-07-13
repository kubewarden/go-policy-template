package sdk

import (
	"bytes"
)

// KubewardenLogWriter is a simple log collector that can be used
// with a `onelog.Logger` instance.
//
// KubewardenLogWriter will send the logs from the WebAssembly guest (the policy)
// to the WebAssembly host (kwct, policy-server).
//
// KubewardenLogWriter will write the log events to the standard output when the
// binary is NOT built for the WebAssembly target. This is useful for running
// native unit tests of policies.
type KubewardenLogWriter struct {
	buffer bytes.Buffer
}
