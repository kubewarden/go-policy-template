package main

import (
	"fmt"

	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func validate(payload []byte) ([]byte, error) {
	if !gjson.ValidBytes(payload) {
		return kubewarden.RejectRequest(
			kubewarden.Message("Not a valid JSON document"),
			kubewarden.Code(400))
	}

	settings, err := NewSettingsFromValidationReq(payload)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	data := gjson.GetBytes(
		payload,
		"request.object.metadata.name")

	if !data.Exists() {
		return kubewarden.AcceptRequest()
	}
	name := data.String()

	if settings.DeniedNames.Contains(name) {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' name is on the deny list", name)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
