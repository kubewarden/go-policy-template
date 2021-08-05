package main

import (
	"fmt"

	onelog "github.com/francoispqt/onelog"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func validate(payload []byte) ([]byte, error) {
	settings, err := NewSettingsFromValidationReq(payload)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	logger.Info("validating request")

	data := gjson.GetBytes(
		payload,
		"request.object.metadata.name")

	if !data.Exists() {
		logger.Warn("cannot read object name from metadata: accepting request")
		return kubewarden.AcceptRequest()
	}
	name := data.String()

	logger.DebugWithFields("validating ingress object", func(e onelog.Entry) {
		namespace := gjson.GetBytes(payload, "request.object.metadata.namespace").String()
		e.String("name", name)
		e.String("namespace", namespace)
	})

	if settings.DeniedNames.Contains(name) {
		logger.InfoWithFields("rejecting ingress object", func(e onelog.Entry) {
			e.String("name", name)
			e.String("denied_names", settings.DeniedNames.String())
		})

		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' name is on the deny list", name)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
