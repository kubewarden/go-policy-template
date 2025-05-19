package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	mapset "github.com/deckarep/golang-set/v2"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings is the structure that describes the policy settings.
type Settings struct {
	DeniedLabels      mapset.Set[string]            `json:"denied_labels"`
	ConstrainedLabels map[string]*RegularExpression `json:"constrained_labels"`
}

type RegularExpression struct {
	*regexp.Regexp
}

// UnmarshalText satisfies the encoding.TextMarshaler interface,
// also used by json.Unmarshal.
func (r *RegularExpression) UnmarshalText(text []byte) error {
	nativeRegExp, err := regexp.Compile(string(text))
	if err != nil {
		return err
	}
	r.Regexp = nativeRegExp
	return nil
}

func (r *RegularExpression) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	rawSettings := struct {
		DeniedLabels      []string                      `json:"denied_labels"`
		ConstrainedLabels map[string]*RegularExpression `json:"constrained_labels"`
	}{}

	err := json.Unmarshal(data, &rawSettings)
	if err != nil {
		return err
	}

	s.DeniedLabels = mapset.NewThreadUnsafeSet[string](rawSettings.DeniedLabels...)
	s.ConstrainedLabels = rawSettings.ConstrainedLabels

	return nil
}

func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}

// Valid is the structure that informs if the policy settings are valid. No
// special checks have to be done.
func (s *Settings) Valid() (bool, error) {
	constrainedLabels := mapset.NewThreadUnsafeSet[string]()

	for label := range s.ConstrainedLabels {
		constrainedLabels.Add(label)
	}

	constrainedAndDenied := constrainedLabels.Intersect(s.DeniedLabels)
	if constrainedAndDenied.Cardinality() != 0 {
		return false,
			fmt.Errorf("these labels cannot be constrained and denied at the same time: %v", constrainedAndDenied)
	}

	return true, nil
}

func validateSettings(payload []byte) ([]byte, error) {
	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(
			kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}

	valid, err := settings.Valid()
	if valid {
		return kubewarden.AcceptSettings()
	}

	return kubewarden.RejectSettings(
		kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
}
