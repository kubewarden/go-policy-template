package testing

import (
	"os"

	"encoding/json"

	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// BuildValidationRequestFromFixture creates the payload for the invocation of the `validate`
// function.
// * `reqFixture`: path to the json file with a recorded requst to evaluate
// * `settings`: instance of policy settings. Must be serializable to JSON using json.
func BuildValidationRequestFromFixture(reqFixture string, settings interface{}) ([]byte, error) {
	kubeAdmissionReqRaw, err := os.ReadFile(reqFixture)
	if err != nil {
		return nil, err
	}

	kubeAdmissionReq := kubewarden_protocol.KubernetesAdmissionRequest{}
	if err = json.Unmarshal(kubeAdmissionReqRaw, &kubeAdmissionReq); err != nil {
		return nil, err
	}

	settingsRaw, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	validationRequest := kubewarden_protocol.ValidationRequest{
		Request:  kubeAdmissionReq,
		Settings: settingsRaw,
	}

	return json.Marshal(validationRequest)
}

// BuildValidationRequest creates the payload for the invocation of the `validate`
// function.
// * `object`: instance of the object. Must be serializable to JSON using json
// * `settings`: instance of policy settings. Must be serializable to JSON using json.
func BuildValidationRequest(object, settings interface{}) ([]byte, error) {
	objectRaw, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	kubeAdmissionReq := kubewarden_protocol.KubernetesAdmissionRequest{
		Object: objectRaw,
	}

	settingsRaw, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	validationRequest := kubewarden_protocol.ValidationRequest{
		Request:  kubeAdmissionReq,
		Settings: settingsRaw,
	}

	return json.Marshal(validationRequest)
}
