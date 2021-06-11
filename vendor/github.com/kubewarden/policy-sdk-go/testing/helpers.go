package testing

import (
	"encoding/json"
	"io/ioutil"
)

// ValidationRequest describes the payload given to the `validate` function
type ValidationRequest struct {
	Request  *json.RawMessage `json:"request"`
	Settings interface{}      `json:"settings"`
}

// ValidationRequest describes the response returned by the `validate` function
type ValidationResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message,omitempty"`
	Code     uint64 `json:"code,omitempty"`
}

// BuildValidationRequest creates the payload for the invocation of the `validate`
// function.
// * `req_fixture`: path to the json file with a recorded requst to evaluate
// * `settings`: instance of policy settings. Must be serializable to JSON
func BuildValidationRequest(req_fixture string, settings interface{}) ([]byte, error) {
	requestRaw, err := ioutil.ReadFile(req_fixture)
	if err != nil {
		return []byte{}, err
	}

	request := json.RawMessage(requestRaw)

	validation_request := ValidationRequest{
		Request:  &request,
		Settings: settings,
	}

	return json.Marshal(validation_request)
}

// SettingsValidationResponse describes the response returned by the
// `validate_settings` function
type SettingsValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}
