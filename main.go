package main

import (
	onelog "github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

// This is not a good practice in general. Policy authors should avoid using global variables in the final code
//
//nolint:gochecknoglobals // Allowing global variables just to make the template code simple.
var (
	logWriter = kubewarden.KubewardenLogWriter{}
	logger    = onelog.New(
		&logWriter,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL
	)
)

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"validate":          validate,
		"validate_settings": validateSettings,
	})
}
