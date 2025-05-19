package main

import (
	"encoding/json"
	"fmt"
	"testing"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

func TestValidateLabel(t *testing.T) {
	cases := []struct {
		podLabels         map[string]string
		deniedLabels      []string
		constrainedLabels map[string]string
		expectedIsValid   bool
	}{
		{
			// Pod has no labels -> should be accepted
			podLabels:         map[string]string{},
			deniedLabels:      []string{"owner"},
			constrainedLabels: map[string]string{},
			expectedIsValid:   true,
		},
		{
			// Pod has labels, none is denied -> should be accepted
			podLabels: map[string]string{
				"hello": "world",
			},
			deniedLabels:      []string{"owner"},
			constrainedLabels: map[string]string{},
			expectedIsValid:   true,
		},
		{
			// Pod has labels, one is denied -> should be rejected
			podLabels: map[string]string{
				"hello": "world",
			},
			deniedLabels:      []string{"hello"},
			constrainedLabels: map[string]string{},
			expectedIsValid:   false,
		},
		{
			// Pod has labels, one has constraint that is respected -> should be accepted
			podLabels: map[string]string{
				"cc-center": "team-123",
			},
			deniedLabels: []string{"hello"},
			constrainedLabels: map[string]string{
				"cc-center": "team-\\d+",
			},
			expectedIsValid: true,
		},
		{
			// Pod has labels, one has constraint that are not respected -> should be rejected
			podLabels: map[string]string{
				"cc-center": "team-kubewarden",
			},
			deniedLabels: []string{"hello"},
			constrainedLabels: map[string]string{
				"cc-center": "team-\\d+",
			},
			expectedIsValid: false,
		},
		{
			// Settings have a constraint, pod doesn't have this label -> should be rejected
			podLabels: map[string]string{
				"owner": "team-kubewarden",
			},
			deniedLabels: []string{"hello"},
			constrainedLabels: map[string]string{
				"cc-center": "team-\\d+",
			},
			expectedIsValid: false,
		},
	}

	for _, testCase := range cases {
		settingsJSON := fmt.Sprintf(`{
			"denied_labels": %s,
			"constrained_labels": %s
		}`,
			mustMarshal(testCase.deniedLabels),
			mustMarshal(testCase.constrainedLabels),
		)

		settings := Settings{}
		if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
			t.Fatalf("Failed to unmarshal settings: %v", err)
		}

		pod := corev1.Pod{
			Metadata: &metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Labels:    testCase.podLabels,
			},
		}

		payload, buildErr := kubewarden_testing.BuildValidationRequest(&pod, &settings)
		if buildErr != nil {
			t.Errorf("Unexpected error: %+v", buildErr)
		}

		responsePayload, validateErr := validate(payload)
		if validateErr != nil {
			t.Errorf("Unexpected error: %+v", validateErr)
		}

		var response kubewarden_protocol.ValidationResponse
		if unmarshalErr := json.Unmarshal(responsePayload, &response); unmarshalErr != nil {
			t.Errorf("Unexpected error: %+v", unmarshalErr)
		}

		if testCase.expectedIsValid && !response.Accepted {
			t.Errorf(
				"Unexpected rejection: msg %s - code %d with pod labels: %v, denied labels: %v, constrained labels: %v",
				*response.Message,
				*response.Code,
				testCase.podLabels,
				testCase.deniedLabels,
				testCase.constrainedLabels,
			)
		}

		if !testCase.expectedIsValid && response.Accepted {
			t.Errorf("Unexpected acceptance with pod labels: %v, denied labels: %v, constrained labels: %v",
				testCase.podLabels, testCase.deniedLabels, testCase.constrainedLabels)
		}
	}
}

func mustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
