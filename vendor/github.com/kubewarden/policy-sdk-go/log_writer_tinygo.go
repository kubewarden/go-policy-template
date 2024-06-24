//go:build tinygo
// +build tinygo

// note well: we have to use the tinygo wasi target, because the wasm one is
// meant to be used inside of the browser

package sdk

import (
	wapc "github.com/wapc/wapc-guest-tinygo"
)

func (k *KubewardenLogWriter) Write(p []byte) (n int, err error) {
	n, err = k.buffer.Write(p)
	line, _ := k.buffer.ReadBytes('\n')
	wapc.HostCall("kubewarden", "tracing", "log", line)
	return
}
