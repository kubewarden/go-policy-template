package testing

import (
	"os"

	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"
)

// BuildValidationRequestFromFixture creates the payload for the invocation of the `validate`
// function.
// * `req_fixture`: path to the json file with a recorded requst to evaluate
// * `settings`: instance of policy settings. Must be serializable to JSON using easyjson
func BuildValidationRequestFromFixture(req_fixture string, settings easyjson.Marshaler) ([]byte, error) {
	kubeAdmissionReqRaw, err := os.ReadFile(req_fixture)
	if err != nil {
		return nil, err
	}

	kubeAdmissionReq := kubewarden_protocol.KubernetesAdmissionRequest{}
	if err := easyjson.Unmarshal(kubeAdmissionReqRaw, &kubeAdmissionReq); err != nil {
		return nil, err
	}

	settingsRaw, err := easyjson.Marshal(settings)
	if err != nil {
		return nil, err
	}

	validationRequest := kubewarden_protocol.ValidationRequest{
		Request:  kubeAdmissionReq,
		Settings: settingsRaw,
	}

	return easyjson.Marshal(validationRequest)
}

// BuildValidationRequest creates the payload for the invocation of the `validate`
// function.
// * `object`: instance of the object. Must be serializable to JSON using easyjson
// * `settings`: instance of policy settings. Must be serializable to JSON using easyjson
func BuildValidationRequest(object, settings easyjson.Marshaler) ([]byte, error) {
	objectRaw, err := easyjson.Marshal(object)
	if err != nil {
		return nil, err
	}

	kubeAdmissionReq := kubewarden_protocol.KubernetesAdmissionRequest{
		Object: objectRaw,
	}

	settingsRaw, err := easyjson.Marshal(settings)
	if err != nil {
		return nil, err
	}

	validationRequest := kubewarden_protocol.ValidationRequest{
		Request:  kubeAdmissionReq,
		Settings: settingsRaw,
	}

	return easyjson.Marshal(validationRequest)
}
