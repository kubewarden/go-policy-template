package main

import (
	"fmt"

	onelog "github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	"github.com/mailru/easyjson"
	"github.com/tidwall/gjson"

	networkingv1 "github.com/kubewarden/k8s-objects/api/networking/v1"
)

func validate(payload []byte) ([]byte, error) {
	settings, err := NewSettingsFromValidationReq(payload)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	ingress_json := gjson.GetBytes(payload, "request.object").String()
	ingress := &networkingv1.Ingress{}
	if err := easyjson.Unmarshal([]byte(ingress_json), ingress); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode Ingress object: %s", err.Error())),
			kubewarden.NoCode)
	}

	logger.DebugWithFields("validating ingress object", func(e onelog.Entry) {
		e.String("name", ingress.Metadata.Name)
		e.String("namespace", ingress.Metadata.Namespace)
	})

	if settings.DeniedNames.Contains(ingress.Metadata.Name) {
		logger.InfoWithFields("rejecting ingress object", func(e onelog.Entry) {
			e.String("name", ingress.Metadata.Name)
			e.String("denied_names", settings.DeniedNames.String())
		})

		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' name is on the deny list", ingress.Metadata.Name)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
