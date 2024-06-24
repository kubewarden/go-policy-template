//go:build tinygo
// +build tinygo

package protocol

import (
	"fmt"

	"github.com/kubewarden/policy-sdk-go/constants"

	wapc "github.com/wapc/wapc-guest-tinygo"
)

// Use the go module initialization function to automatically register
// the protocol version of the waPC module
func init() {
	wapc.RegisterFunctions(wapc.Functions{
		"protocol_version": func(payload []byte) ([]byte, error) {
			return []byte(fmt.Sprintf("%q", constants.ProtocolVersion)), nil
		},
	})
}
