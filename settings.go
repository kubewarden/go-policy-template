package main

import (
	"github.com/deckarep/golang-set"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"

	"fmt"
)

type Settings struct {
	DeniedNames mapset.Set `json:"denied_names"`
}

// Builds a new Settings instance starting from a validation
// request payload:
// {
//    "request": ...,
//    "settings": {
//       "denied_names": [...]
//    }
// }
func NewSettingsFromValidationReq(payload []byte) (Settings, error) {
	// Note well: we don't validate the input JSON now, this has
	// already done inside of the `validate` function

	return newSettings(
		payload,
		"settings.denied_names")
}

// Builds a new Settings instance starting from a Settings
// payload:
// {
//    "denied_names": ...
// }
func NewSettingsFromValidateSettingsPayload(payload []byte) (Settings, error) {
	if !gjson.ValidBytes(payload) {
		return Settings{}, fmt.Errorf("denied JSON payload")
	}

	return newSettings(
		payload,
		"denied_names")
}

func newSettings(payload []byte, paths ...string) (Settings, error) {
	if len(paths) != 1 {
		return Settings{}, fmt.Errorf("wrong number of json paths")
	}

	data := gjson.GetManyBytes(payload, paths...)

	deniedNames := mapset.NewThreadUnsafeSet()
	data[0].ForEach(func(_, entry gjson.Result) bool {
		deniedNames.Add(entry.String())
		return true
	})

	return Settings{
		DeniedNames: deniedNames,
	}, nil
}

// No special check has to be done
func (s *Settings) Valid() bool {
	return true
}

func validateSettings(payload []byte) ([]byte, error) {
	settings, err := NewSettingsFromValidateSettingsPayload(payload)
	if err != nil {
		return []byte{}, err
	}

	if settings.Valid() {
		return kubewarden.AcceptSettings()
	}
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
