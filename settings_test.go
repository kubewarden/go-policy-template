package main

import (
	"testing"
)

func TestParsingSettingsWithAllValuesProvidedFromValidationReq(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
			"denied_names": [ "foo", "bar" ]
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	expected := []string{"foo", "bar"}
	for _, exp := range expected {
		if !settings.DeniedNames.Contains(exp) {
			t.Errorf("Missing value %s", exp)
		}
	}
}

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if settings.DeniedNames.Cardinality() != 0 {
		t.Errorf("Expecpted DeniedNames to be empty")
	}
}

func TestSettingsAreValid(t *testing.T) {
	request := `
	{
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidateSettingsPayload(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if !settings.Valid() {
		t.Errorf("Settings are reported as not valid")
	}
}
