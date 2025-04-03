//go:build !tinygo
// +build !tinygo

// note well: we have to use the tinygo wasi target, because the wasm one is
// meant to be used inside of the browser

package sdk

import (
	"fmt"
)

func (k *KubewardenLogWriter) Write(p []byte) (int, error) {
	n, err := k.buffer.Write(p)
	line, _ := k.buffer.ReadBytes('\n')
	//nolint:forbidigo // this is a debug print
	fmt.Printf("NATIVE: |%s|\n", string(line))
	return n, err
}
