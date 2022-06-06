// This package provides helper functions and structs for writing
// https://kubewarden.io policies using the Go programming
// language.
package sdk

import (
	"github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"
)

// Message is the optional string used to build validation responses
type Message string

// Code is the optional error code associated with validation responses
type Code uint16

const (
	// NoMessage can be used when building a response that doesn't have any
	// message to be shown to the user
	NoMessage Message = ""

	// NoCode can be used when building a response that doesn't have any
	// error code to be shown to the user
	NoCode Code = 0
)

// AcceptRequest can be used inside of the `validate` function to accept the
// incoming request
func AcceptRequest() ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted: true,
	}

	return easyjson.Marshal(response)
}

// RejectRequest can be used inside of the `validate` function to reject the
// incoming request
// * `message`: optional message to show to the user
// * `code`: optional error code to show to the user
func RejectRequest(message Message, code Code) ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted: false,
	}
	if message != NoMessage {
		msg := string(message)
		response.Message = &msg
	}
	if code != NoCode {
		c := uint16(code)
		response.Code = &c
	}

	return easyjson.Marshal(response)
}

// Accept the request and mutate the final object to match the
// one provided via the `newObject` param
func MutateRequest(newObject easyjson.Marshaler) ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted:      true,
		MutatedObject: newObject,
	}

	return easyjson.Marshal(response)
}

// AcceptSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as valid
func AcceptSettings() ([]byte, error) {
	response := protocol.SettingsValidationResponse{
		Valid: true,
	}
	return easyjson.Marshal(response)
}

// RejectSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as invalid
// * `message`: optional message to show to the user
func RejectSettings(message Message) ([]byte, error) {
	response := protocol.SettingsValidationResponse{
		Valid: false,
	}

	if message != NoMessage {
		msg := string(message)
		response.Message = &msg
	}
	return easyjson.Marshal(response)
}
