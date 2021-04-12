package sdk

import (
	"fmt"
	"strings"
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

type keyValue struct {
	key   string
	value string
}

func (kv keyValue) String() string {
	return fmt.Sprintf(`"%s":%s`, kv.key, kv.value)
}

// AcceptRequest can be used inside of the `validate` function to accept the
// incoming request
func AcceptRequest() ([]byte, error) {
	return []byte(`{"accepted":true}`), nil
}

// RejectRequest can be used inside of the `validate` function to reject the
// incoming request
// * `message`: optional message to show to the user
// * `code`: optional error code to show to the user
func RejectRequest(message Message, code Code) ([]byte, error) {
	result := []keyValue{{key: "accepted", value: "false"}}
	if message != NoMessage {
		result = append(result, keyValue{key: "message", value: fmt.Sprintf(`"%s"`, string(message))})
	}
	if code != NoCode {
		result = append(result, keyValue{key: "code", value: fmt.Sprintf("%d", code)})
	}
	stringResult := []string{}
	for _, keyValue := range result {
		stringResult = append(stringResult, keyValue.String())
	}
	return []byte(fmt.Sprintf("{%s}", strings.Join(stringResult, ","))), nil
}

// AcceptSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as valid
func AcceptSettings() ([]byte, error) {
	return []byte(`{"valid":true}`), nil
}

// RejectSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as invalid
// * `message`: optional message to show to the user
func RejectSettings(message Message) ([]byte, error) {
	result := []keyValue{{key: "valid", value: "false"}}
	if message != NoMessage {
		result = append(result, keyValue{key: "message", value: fmt.Sprintf(`"%s"`, string(message))})
	}
	stringResult := []string{}
	for _, keyValue := range result {
		stringResult = append(stringResult, keyValue.String())
	}
	return []byte(fmt.Sprintf("{%s}", strings.Join(stringResult, ","))), nil
}
