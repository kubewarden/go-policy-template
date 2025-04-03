package protocol

import (
	"encoding/json"
)

// ValidationResponse represents the response given when validating a request.
type ValidationResponse struct {
	Accepted bool `json:"accepted"`
	// Optional - ignored if accepted
	Message *string `json:"message,omitempty"`
	// Optional - ignored if accepted
	Code *uint16 `json:"code,omitempty"`
	// Optional - used only by mutating policies
	MutatedObject interface{} `json:"mutated_object,omitempty"`
}

// SettingsValidationResponse repreents the response sent by a policy when validating its settings.
type SettingsValidationResponse struct {
	Valid bool `json:"valid"`
	// Optional - ignored if valid
	Message *string `json:"message,omitempty"`
}

// ValidationRequest represents the object received by the validate() function of Kubewarden policies.
type ValidationRequest struct {
	// The request to be evaluated
	Request KubernetesAdmissionRequest `json:"request"`
	// The settings to be used by the policy
	//
	// Note, the attributes holds the unmarshalled []bytes as found inside of
	// original JSON object.
	// This can then be parsed using `json.Unmarshal()` into a proper
	// type that implements the json interfaces.
	Settings json.RawMessage `json:"settings"`
}

// KubernetesAdmissionRequest represents Kubernetes' [AdmissionReview](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) request.
type KubernetesAdmissionRequest struct {
	// UID is an identifier for the individual request/response. It allows
	// us to distinguish instances of requests which are otherwise
	// identical (parallel requests, requests when earlier requests did not
	// modify etc) The UID is meant to track the round trip
	// (request/response) between the KAS and the WebHook, not the user
	// request. It is suitable for correlating log entries between the
	// webhook and apiserver, for either auditing or debugging.
	Uid string `json:"uid"` //nolint:revive // We cannot change this field name without breaking compatibility

	// Kind is the fully-qualified type of object being submitted (for
	// example, v1.Pod or autoscaling.v1.Scale)
	Kind GroupVersionKind `json:"kind"`

	// Resource is the fully-qualified resource being requested (for
	// example, v1.pods)
	Resource GroupVersionResource `json:"resource"`

	// SubResource is the subresource being requested, if any (for example,
	// "status" or "scale")
	SubResource string `json:"subResource"`

	// RequestKind is the fully-qualified type of the original API request
	// (for example, v1.Pod or autoscaling.v1.Scale). If this is specified
	// and differs from the value in "kind", an equivalent match and
	// conversion was performed.
	//
	// For example, if deployments can be modified via apps/v1 and
	// apps/v1beta1, and a webhook registered a rule of
	// `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`
	// and `matchPolicy: Equivalent`, an API request to apps/v1beta1
	// deployments would be converted and sent to the webhook with `kind:
	// {group:"apps", version:"v1", kind:"Deployment"}` (matching the rule
	// the webhook registered for), and `requestKind: {group:"apps",
	// version:"v1beta1", kind:"Deployment"}` (indicating the kind of the
	// original API request).
	//
	// See documentation for the "matchPolicy" field in the webhook
	// configuration type for more details.
	RequestKind GroupVersionKind `json:"requestKind"`

	// RequestResource is the fully-qualified resource of the original API
	// request (for example, v1.pods). If this is specified and differs
	// from the value in "resource", an equivalent match and conversion was
	// performed.
	//
	// For example, if deployments can be modified via apps/v1 and
	// apps/v1beta1, and a webhook registered a rule of
	// `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`
	// and `matchPolicy: Equivalent`, an API request to apps/v1beta1
	// deployments would be converted and sent to the webhook with
	// `resource: {group:"apps", version:"v1", resource:"deployments"}`
	// (matching the resource the webhook registered for), and
	// `requestResource: {group:"apps", version:"v1beta1",
	// resource:"deployments"}` (indicating the resource of the original
	// API request).
	//
	// See documentation for the "matchPolicy" field in the webhook
	// configuration type.
	RequestResource GroupVersionKind `json:"requestResource"`

	// RequestSubResource is the name of the subresource of the original
	// API request, if any (for example, "status" or "scale") If this is
	// specified and differs from the value in "subResource", an equivalent
	// match and conversion was performed. See documentation for the
	// "matchPolicy" field in the webhook configuration type.
	RequestSubResource string `json:"requestSubResource"`

	// Name is the name of the object as presented in the request.  On a
	// CREATE operation, the client may omit name and rely on the server to
	// generate the name.  If that is the case, this field will contain an
	// empty string.
	Name string `json:"name"`

	// Namespace is the namespace associated with the request (if any).
	Namespace string `json:"namespace"`

	// Operation is the operation being performed. This may be different
	// than the operation requested. e.g. a patch can result in either a
	// CREATE or UPDATE Operation.
	Operation string `json:"operation"`

	// UserInfo is information about the requesting user
	UserInfo UserInfo `json:"userInfo"`

	// Object is the object from the incoming request.
	//
	// Note, the attributes holds the unmarshalled []bytes as found inside of
	// original JSON object.
	// This can then be parsed using `json.Unmarshal()` into a proper
	// type that implements the json interfaces.
	Object json.RawMessage `json:"object"`

	// OldObject is the existing object. Only populated for DELETE and
	// UPDATE requests.
	//
	// Note, the attributes holds the unmarshalled []bytes as found inside
	// of original JSON object. This can then be parsed using
	// `json.Unmarshal()` into a proper type that implements the json
	// interfaces.
	OldObject json.RawMessage `json:"oldObject"`

	// DryRun indicates that modifications will definitely not be persisted
	// for this request. Defaults to false.
	DryRun bool `json:"dryRun"`

	// Options is the operation option structure of the operation being
	// performed. e.g. `meta.k8s.io/v1.DeleteOptions` or
	// `meta.k8s.io/v1.CreateOptions`. This may be different than the
	// options the caller provided. e.g. for a patch request the performed
	// Operation might be a CREATE, in which case the Options will a
	// `meta.k8s.io/v1.CreateOptions` even though the caller provided
	// `meta.k8s.io/v1.PatchOptions`.
	//
	// Note, the attributes holds the unmarshalled []bytes as found inside
	// of original JSON object. This can then be parsed using
	// `json.Unmarshal()` into a proper type that implements the json
	// interfaces.
	Options json.RawMessage `json:"options"`
}

// GroupVersionKind unambiguously identifies a kind.
type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// GroupVersionResource unambiguously identifies a resource.
type GroupVersionResource struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// UserInfo holds information about the user who made the request.
type UserInfo struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
	// Note, the attributes holds the unmarshalled []bytes as found inside
	// of original JSON object. This can then be parsed using
	// `json.Unmarshal()` into a proper type that implements the json
	// interfaces.
	Extra json.RawMessage `json:"extra,omitempty"`
}
