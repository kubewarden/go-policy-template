package main

import (
	"regexp"
	"testing"

	"encoding/json"

	mapset "github.com/deckarep/golang-set/v2"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

func TestValidateLabel(t *testing.T) {
	// NOTE 1
	cases := []struct {
		podLabels         map[string]string
		deniedLabels      mapset.Set[string]
		constrainedLabels map[string]*RegularExpression
		expectedIsValid   bool
	}{
		{
			// ➀
			// Pod has no labels -> should be accepted
			podLabels:         map[string]string{},
			deniedLabels:      mapset.NewThreadUnsafeSet[string]("owner"),
			constrainedLabels: map[string]*RegularExpression{},
			expectedIsValid:   true,
		},
		{
			// ➁
			// Pod has labels, none is denied -> should be accepted
			podLabels: map[string]string{
				"hello": "world",
			},
			deniedLabels:      mapset.NewThreadUnsafeSet[string]("owner"),
			constrainedLabels: map[string]*RegularExpression{},
			expectedIsValid:   true,
		},
		{
			// ➂
			// Pod has labels, one is denied -> should be rejected
			podLabels: map[string]string{
				"hello": "world",
			},
			deniedLabels:      mapset.NewThreadUnsafeSet[string]("hello"),
			constrainedLabels: map[string]*RegularExpression{},
			expectedIsValid:   false,
		},
		{
			// ➃
			// Pod has labels, one has constraint that is respected -> should be accepted
			podLabels: map[string]string{
				"cc-center": "team-123",
			},
			deniedLabels: mapset.NewThreadUnsafeSet[string]("hello"),
			constrainedLabels: map[string]*RegularExpression{
				"cc-center": {
					Regexp: regexp.MustCompile(`team-\d+`),
				},
			},
			expectedIsValid: true,
		},
		{
			// ➄
			// Pod has labels, one has constraint that are not respected -> should be rejected
			podLabels: map[string]string{
				"cc-center": "team-kubewarden",
			},
			deniedLabels: mapset.NewThreadUnsafeSet[string]("hello"),
			constrainedLabels: map[string]*RegularExpression{
				"cc-center": {
					Regexp: regexp.MustCompile(`team-\d+`),
				},
			},
			expectedIsValid: false,
		},
		{
			// ➅
			// Settings have a constraint, pod doesn't have this label -> should be rejected
			podLabels: map[string]string{
				"owner": "team-kubewarden",
			},
			deniedLabels: mapset.NewThreadUnsafeSet[string]("hello"),
			constrainedLabels: map[string]*RegularExpression{
				"cc-center": {
					Regexp: regexp.MustCompile(`team-\d+`),
				},
			},
			expectedIsValid: false,
		},
	}

	// NOTE 2
	for _, testCase := range cases {
		settings := Settings{
			DeniedLabels:      testCase.deniedLabels,
			ConstrainedLabels: testCase.constrainedLabels,
		}

		pod := corev1.Pod{
			Metadata: &metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Labels:    testCase.podLabels,
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
		if err := json.Unmarshal(responsePayload, &response); err != nil {
			t.Errorf("Unexpected error: %+v", err)
		}

		if testCase.expectedIsValid && !response.Accepted {
			t.Errorf("Unexpected rejection: msg %s - code %d with pod labels: %v, denied labels: %v, constrained labels: %v",
				*response.Message, *response.Code, testCase.podLabels, testCase.deniedLabels, testCase.constrainedLabels)
		}

		if !testCase.expectedIsValid && response.Accepted {
			t.Errorf("Unexpected acceptance with pod labels: %v, denied labels: %v, constrained labels: %v",
				testCase.podLabels, testCase.deniedLabels, testCase.constrainedLabels)
		}
	}
}
