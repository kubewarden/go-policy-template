package main

import (
	"testing"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/mailru/easyjson"
)

func TestEmptySettingsLeadsToApproval(t *testing.T) {
	settings := Settings{}
	pod := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&pod, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, *response.Code)
	}
}

func TestApproval(t *testing.T) {
	settings := Settings{
		DeniedNames: []string{"foo", "bar"},
	}
	pod := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&pod, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func TestApproveFixture(t *testing.T) {
	settings := Settings{
		DeniedNames: []string{},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture(
		"test_data/pod.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected rejection")
	}
}

func TestRejectionBecauseNameIsDenied(t *testing.T) {
	settings := Settings{
		DeniedNames: []string{"foo", "test-pod"},
	}

	pod := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&pod, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != false {
		t.Error("Unexpected approval")
	}

	expected_message := "The 'test-pod' name is on the deny list"
	if response.Message == nil {
		t.Errorf("expected response to have a message")
	}
	if *response.Message != expected_message {
		t.Errorf("Got '%s' instead of '%s'", *response.Message, expected_message)
	}
}
