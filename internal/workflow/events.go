package workflow

// AdmissionReview represents the structure sent by Kubernetes API server to webhooks
type AdmissionReview struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Request    *AdmissionRequest  `json:"request"`
	Response   *AdmissionResponse `json:"response,omitempty"`
}

// AdmissionRequest contains the information to create an admission request
type AdmissionRequest struct {
	UID             string               `json:"uid"`
	Kind            GroupVersionKind     `json:"kind"`
	Resource        GroupVersionResource `json:"resource"`
	RequestKind     GroupVersionKind     `json:"requestKind"`
	RequestResource GroupVersionResource `json:"requestResource"`
	Name            string               `json:"name"`
	Namespace       string               `json:"namespace"`
	Operation       string               `json:"operation"`
	UserInfo        UserInfo             `json:"userInfo"`
	Object          map[string]any       `json:"object"`
	OldObject       map[string]any       `json:"oldObject"`
	DryRun          bool                 `json:"dryRun"`
	Options         map[string]any       `json:"options"`
}

// AdmissionResponse contains the result of an admission request
type AdmissionResponse struct {
	UID       string  `json:"uid"`
	Allowed   bool    `json:"allowed"`
	Result    *Result `json:"result,omitempty"`
	Patch     []byte  `json:"patch,omitempty"`
	PatchType *string `json:"patchType,omitempty"`
}

// GroupVersionKind contains the group, version, and kind of an object
type GroupVersionKind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// GroupVersionResource contains the group, version, and resource of an object
type GroupVersionResource struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
}

// UserInfo contains information about the requesting user
type UserInfo struct {
	Username string   `json:"username"`
	UID      string   `json:"uid"`
	Groups   []string `json:"groups"`
}

// Result contains additional details about why a request was denied
type Result struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Event represents a simplified event structure for backward compatibility
type Event struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// ConfigMap represents a Kubernetes ConfigMap resource
type ConfigMap struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Data       map[string]string `json:"data"`
}

// Metadata contains the metadata for a Kubernetes resource
type Metadata struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	ResourceVersion string            `json:"resourceVersion,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
}
