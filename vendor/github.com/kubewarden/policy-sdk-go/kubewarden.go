// This package provides helper functions and structs for writing
// https://kubewarden.io policies using the Go programming
// language.
package sdk

import (
	"errors"

	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	batchv1 "github.com/kubewarden/k8s-objects/api/batch/v1"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	"github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"
)

// Message is the optional string used to build validation responses
type Message string

// Code is the optional error code associated with validation responses
type Code uint16

const (
	// NoMessage can be used when building a response that doesn't have any
	// message to be shown to the user
	NoMessage Message = ""

	// NoCode can be used when building a response that doesn't have any
	// error code to be shown to the user
	NoCode Code = 0
)

// AcceptRequest can be used inside of the `validate` function to accept the
// incoming request
func AcceptRequest() ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted: true,
	}

	return easyjson.Marshal(response)
}

// RejectRequest can be used inside of the `validate` function to reject the
// incoming request
// * `message`: optional message to show to the user
// * `code`: optional error code to show to the user
func RejectRequest(message Message, code Code) ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted: false,
	}
	if message != NoMessage {
		msg := string(message)
		response.Message = &msg
	}
	if code != NoCode {
		c := uint16(code)
		response.Code = &c
	}

	return easyjson.Marshal(response)
}

// Accept the request and mutate the final object to match the
// one provided via the `newObject` param
func MutateRequest(newObject easyjson.Marshaler) ([]byte, error) {
	response := protocol.ValidationResponse{
		Accepted:      true,
		MutatedObject: newObject,
	}

	return easyjson.Marshal(response)
}

// Update the pod spec from the resource defined in the original object and
// create an acceptance response.
// * `validation_request` - the original admission request
// * `pod_spec` - new PodSpec to be set in the response
func MutatePodSpecFromRequest(validationRequest protocol.ValidationRequest, podSepc corev1.PodSpec) ([]byte, error) {
	switch validationRequest.Request.Kind.Kind {
	case "apps.v1.Deployment":
		deployment := appsv1.Deployment{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &deployment); err != nil {
			return nil, err
		}
		deployment.Spec.Template.Spec = &podSepc
		return MutateRequest(deployment)
	case "apps.v1.ReplicaSet":
		replicaset := appsv1.ReplicaSet{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &replicaset); err != nil {
			return nil, err
		}
		replicaset.Spec.Template.Spec = &podSepc
		return MutateRequest(replicaset)
	case "apps.v1.StatefulSet":
		statefulset := appsv1.StatefulSet{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &statefulset); err != nil {
			return nil, err
		}
		statefulset.Spec.Template.Spec = &podSepc
		return MutateRequest(statefulset)
	case "apps.v1.DaemonSet":
		daemonset := appsv1.DaemonSet{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &daemonset); err != nil {
			return nil, err
		}
		daemonset.Spec.Template.Spec = &podSepc
		return MutateRequest(daemonset)
	case "v1.ReplicationController":
		replicationController := corev1.ReplicationController{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &replicationController); err != nil {
			return nil, err
		}
		replicationController.Spec.Template.Spec = &podSepc
		return MutateRequest(replicationController)
	case "batch.v1.CronJob":
		cronjob := batchv1.CronJob{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &cronjob); err != nil {
			return nil, err
		}
		cronjob.Spec.JobTemplate.Spec.Template.Spec = &podSepc
		return MutateRequest(cronjob)
	case "batch.v1.Job":
		job := batchv1.Job{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &job); err != nil {
			return nil, err
		}
		job.Spec.Template.Spec = &podSepc
		return MutateRequest(job)
	case "v1.Pod":
		pod := corev1.Pod{}
		if err := easyjson.Unmarshal(validationRequest.Request.Object, &pod); err != nil {
			return nil, err
		}
		pod.Spec = &podSepc
		return MutateRequest(pod)
	default:
		return RejectRequest("Object should be one of these kinds: Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, Job, CronJob, Pod", NoCode)
	}
}

// AcceptSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as valid
func AcceptSettings() ([]byte, error) {
	response := protocol.SettingsValidationResponse{
		Valid: true,
	}
	return easyjson.Marshal(response)
}

// RejectSettings can be used inside of the `validate_settings` function to
// mark the user provided settings as invalid
// * `message`: optional message to show to the user
func RejectSettings(message Message) ([]byte, error) {
	response := protocol.SettingsValidationResponse{
		Valid: false,
	}

	if message != NoMessage {
		msg := string(message)
		response.Message = &msg
	}
	return easyjson.Marshal(response)
}

// Extract PodSpec from high level objects. This method can be used to evaluate
// high level objects instead of just Pods. For example, it can be used to
// reject Deployments or StatefulSets that violate a policy instead of the Pods
// created by them. Objects supported are: Deployment, ReplicaSet, StatefulSet,
// DaemonSet, ReplicationController, Job, CronJob, Pod It returns an error if
// the object is not one of those. If it is a supported object it returns the
// PodSpec if present, otherwise returns an empty PodSpec.
// * `object`: the request to validate
func ExtractPodSpecFromObject(object protocol.ValidationRequest) (corev1.PodSpec, error) {
	switch object.Request.Kind.Kind {
	case "apps.v1.Deployment":
		deployment := appsv1.Deployment{}
		if err := easyjson.Unmarshal(object.Request.Object, &deployment); err != nil {
			return corev1.PodSpec{}, err
		}
		return *deployment.Spec.Template.Spec, nil
	case "apps.v1.ReplicaSet":
		replicaset := appsv1.ReplicaSet{}
		if err := easyjson.Unmarshal(object.Request.Object, &replicaset); err != nil {
			return corev1.PodSpec{}, err
		}
		return *replicaset.Spec.Template.Spec, nil
	case "apps.v1.StatefulSet":
		statefulset := appsv1.StatefulSet{}
		if err := easyjson.Unmarshal(object.Request.Object, &statefulset); err != nil {
			return corev1.PodSpec{}, err
		}
		return *statefulset.Spec.Template.Spec, nil
	case "apps.v1.DaemonSet":
		daemonset := appsv1.DaemonSet{}
		if err := easyjson.Unmarshal(object.Request.Object, &daemonset); err != nil {
			return corev1.PodSpec{}, err
		}
		return *daemonset.Spec.Template.Spec, nil
	case "v1.ReplicationController":
		replicationController := corev1.ReplicationController{}
		if err := easyjson.Unmarshal(object.Request.Object, &replicationController); err != nil {
			return corev1.PodSpec{}, err
		}
		return *replicationController.Spec.Template.Spec, nil
	case "batch.v1.CronJob":
		cronjob := batchv1.CronJob{}
		if err := easyjson.Unmarshal(object.Request.Object, &cronjob); err != nil {
			return corev1.PodSpec{}, err
		}
		return *cronjob.Spec.JobTemplate.Spec.Template.Spec, nil
	case "batch.v1.Job":
		job := batchv1.Job{}
		if err := easyjson.Unmarshal(object.Request.Object, &job); err != nil {
			return corev1.PodSpec{}, err
		}
		return *job.Spec.Template.Spec, nil
	case "v1.Pod":
		pod := corev1.Pod{}
		if err := easyjson.Unmarshal(object.Request.Object, &pod); err != nil {
			return corev1.PodSpec{}, err
		}
		return *pod.Spec, nil
	default:
		return corev1.PodSpec{}, errors.New("object should be one of these kinds: Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, Job, CronJob, Pod")
	}
}
