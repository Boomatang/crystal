package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// Deployment represents a mock Kubernetes Deployment with common fields and methods
type Deployment struct {
	// Standard Kubernetes metadata
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   DeploymentMetadata `json:"metadata"`
	Spec       DeploymentSpec     `json:"spec"`
	Status     *DeploymentStatus  `json:"status,omitempty"`
}

// DeploymentMetadata represents the metadata section of a Deployment
type DeploymentMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int64             `json:"generation"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// DeploymentSpec represents the spec section of a Deployment
type DeploymentSpec struct {
	Replicas                *int32              `json:"replicas,omitempty"`
	Selector                *LabelSelector      `json:"selector,omitempty"`
	Template                PodTemplateSpec     `json:"template"`
	Strategy                *DeploymentStrategy `json:"strategy,omitempty"`
	MinReadySeconds         int32               `json:"minReadySeconds,omitempty"`
	RevisionHistoryLimit    *int32              `json:"revisionHistoryLimit,omitempty"`
	Paused                  bool                `json:"paused,omitempty"`
	ProgressDeadlineSeconds *int32              `json:"progressDeadlineSeconds,omitempty"`
}

// LabelSelector represents a label selector
type LabelSelector struct {
	MatchLabels      map[string]string          `json:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty"`
}

// LabelSelectorRequirement represents a label selector requirement
type LabelSelectorRequirement struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values,omitempty"`
}

// PodTemplateSpec represents a pod template specification
type PodTemplateSpec struct {
	Metadata ObjectMeta `json:"metadata,omitempty"`
	Spec     PodSpec    `json:"spec"`
}

// ObjectMeta represents basic object metadata
type ObjectMeta struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Generation        int64             `json:"generation,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
}

// PodSpec represents a pod specification
type PodSpec struct {
	Containers    []Container `json:"containers"`
	RestartPolicy string      `json:"restartPolicy,omitempty"`
	DNSPolicy     string      `json:"dnsPolicy,omitempty"`
	HostNetwork   bool        `json:"hostNetwork,omitempty"`
	HostPID       bool        `json:"hostPID,omitempty"`
	HostIPC       bool        `json:"hostIPC,omitempty"`
}

// Container represents a container specification
type Container struct {
	Name            string                `json:"name"`
	Image           string                `json:"image"`
	Command         []string              `json:"command,omitempty"`
	Args            []string              `json:"args,omitempty"`
	WorkingDir      string                `json:"workingDir,omitempty"`
	Ports           []ContainerPort       `json:"ports,omitempty"`
	Env             []EnvVar              `json:"env,omitempty"`
	EnvFrom         []EnvFromSource       `json:"envFrom,omitempty"`
	Resources       *ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts    []VolumeMount         `json:"volumeMounts,omitempty"`
	LivenessProbe   *Probe                `json:"livenessProbe,omitempty"`
	ReadinessProbe  *Probe                `json:"readinessProbe,omitempty"`
	Lifecycle       *Lifecycle            `json:"lifecycle,omitempty"`
	ImagePullPolicy string                `json:"imagePullPolicy,omitempty"`
}

// ContainerPort represents a container port
type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	HostPort      int32  `json:"hostPort,omitempty"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol,omitempty"`
	HostIP        string `json:"hostIP,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name      string        `json:"name"`
	Value     string        `json:"value,omitempty"`
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

// EnvVarSource represents a source for an environment variable
type EnvVarSource struct {
	FieldRef         *ObjectFieldSelector   `json:"fieldRef,omitempty"`
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty"`
	ConfigMapKeyRef  *ConfigMapKeySelector  `json:"configMapKeyRef,omitempty"`
	SecretKeyRef     *SecretKeySelector     `json:"secretKeyRef,omitempty"`
}

// ObjectFieldSelector represents a field selector for an object
type ObjectFieldSelector struct {
	APIVersion string `json:"apiVersion,omitempty"`
	FieldPath  string `json:"fieldPath"`
}

// ResourceFieldSelector represents a resource field selector
type ResourceFieldSelector struct {
	ContainerName string `json:"containerName,omitempty"`
	Resource      string `json:"resource"`
	Divisor       string `json:"divisor,omitempty"`
}

// ConfigMapKeySelector represents a config map key selector
type ConfigMapKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// SecretKeySelector represents a secret key selector
type SecretKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// EnvFromSource represents a source for environment variables
type EnvFromSource struct {
	Prefix       string              `json:"prefix,omitempty"`
	ConfigMapRef *ConfigMapEnvSource `json:"configMapRef,omitempty"`
	SecretRef    *SecretEnvSource    `json:"secretRef,omitempty"`
}

// ConfigMapEnvSource represents a config map environment source
type ConfigMapEnvSource struct {
	Name     string `json:"name"`
	Optional *bool  `json:"optional,omitempty"`
}

// SecretEnvSource represents a secret environment source
type SecretEnvSource struct {
	Name     string `json:"name"`
	Optional *bool  `json:"optional,omitempty"`
}

// ResourceRequirements represents resource requirements
type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// VolumeMount represents a volume mount
type VolumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath"`
	SubPath   string `json:"subPath,omitempty"`
}

// Probe represents a probe
type Probe struct {
	Handler             `json:",inline"`
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      int32 `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       int32 `json:"periodSeconds,omitempty"`
	SuccessThreshold    int32 `json:"successThreshold,omitempty"`
	FailureThreshold    int32 `json:"failureThreshold,omitempty"`
}

// Handler represents a probe handler
type Handler struct {
	Exec      *ExecAction      `json:"exec,omitempty"`
	HTTPGet   *HTTPGetAction   `json:"httpGet,omitempty"`
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
}

// ExecAction represents an exec action
type ExecAction struct {
	Command []string `json:"command,omitempty"`
}

// HTTPGetAction represents an HTTP GET action
type HTTPGetAction struct {
	Path        string       `json:"path,omitempty"`
	Port        IntOrString  `json:"port"`
	Host        string       `json:"host,omitempty"`
	Scheme      string       `json:"scheme,omitempty"`
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty"`
}

// TCPSocketAction represents a TCP socket action
type TCPSocketAction struct {
	Port IntOrString `json:"port"`
	Host string      `json:"host,omitempty"`
}

// HTTPHeader represents an HTTP header
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Lifecycle represents a container lifecycle
type Lifecycle struct {
	PostStart *Handler `json:"postStart,omitempty"`
	PreStop   *Handler `json:"preStop,omitempty"`
}

// DeploymentStrategy represents a deployment strategy
type DeploymentStrategy struct {
	Type          string                   `json:"type,omitempty"`
	RollingUpdate *RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

// RollingUpdateDeployment represents a rolling update deployment
type RollingUpdateDeployment struct {
	MaxUnavailable *IntOrString `json:"maxUnavailable,omitempty"`
	MaxSurge       *IntOrString `json:"maxSurge,omitempty"`
}

// DeploymentStatus represents the status of a deployment
type DeploymentStatus struct {
	ObservedGeneration  int64                 `json:"observedGeneration,omitempty"`
	Replicas            int32                 `json:"replicas,omitempty"`
	UpdatedReplicas     int32                 `json:"updatedReplicas,omitempty"`
	ReadyReplicas       int32                 `json:"readyReplicas,omitempty"`
	AvailableReplicas   int32                 `json:"availableReplicas,omitempty"`
	UnavailableReplicas int32                 `json:"unavailableReplicas,omitempty"`
	Conditions          []DeploymentCondition `json:"conditions,omitempty"`
}

// DeploymentCondition represents a deployment condition
type DeploymentCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastUpdateTime     time.Time `json:"lastUpdateTime,omitempty"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
}

// IntOrString is a type that can hold either an int32 or a string
type IntOrString struct {
	Type   Type   `json:"type"`
	IntVal int32  `json:"intVal,omitempty"`
	StrVal string `json:"strVal,omitempty"`
}

// Type represents the type of IntOrString
type Type int

const (
	Int Type = iota
	String
)

// NewDeployment creates a new Deployment with default values
func NewDeployment(name, namespace string) *Deployment {
	now := time.Now()
	replicas := int32(1)

	return &Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: DeploymentMetadata{
			Name:              name,
			Namespace:         namespace,
			UID:               GenerateUID(),
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
			Labels:            make(map[string]string),
			Annotations:       make(map[string]string),
		},
		Spec: DeploymentSpec{
			Replicas: &replicas,
			Selector: &LabelSelector{
				MatchLabels: make(map[string]string),
			},
			Template: PodTemplateSpec{
				Metadata: ObjectMeta{
					Labels: make(map[string]string),
				},
				Spec: PodSpec{
					Containers:    []Container{},
					RestartPolicy: "Always",
					DNSPolicy:     "ClusterFirst",
				},
			},
		},
	}
}

// AddContainer adds a container to the deployment
func (d *Deployment) AddContainer(container Container) {
	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)
}

// SetReplicas sets the number of replicas
func (d *Deployment) SetReplicas(replicas int32) {
	d.Spec.Replicas = &replicas
}

// AddLabel adds a label to the deployment
func (d *Deployment) AddLabel(key, value string) {
	d.Metadata.Labels[key] = value
}

// AddAnnotation adds an annotation to the deployment
func (d *Deployment) AddAnnotation(key, value string) {
	d.Metadata.Annotations[key] = value
}

// AddPodLabel adds a label to the pod template
func (d *Deployment) AddPodLabel(key, value string) {
	d.Spec.Template.Metadata.Labels[key] = value
}

// AddSelectorLabel adds a label to the selector
func (d *Deployment) AddSelectorLabel(key, value string) {
	if d.Spec.Selector == nil {
		d.Spec.Selector = &LabelSelector{
			MatchLabels: make(map[string]string),
		}
	}
	d.Spec.Selector.MatchLabels[key] = value
}

// ToJSON converts the Deployment to JSON string
func (d *Deployment) ToJSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Deployment to JSON: %w", err)
	}
	return string(data), err
}

// FromJSON creates a Deployment from JSON string
func DeploymentFromJSON(jsonStr string) (*Deployment, error) {
	var d Deployment
	err := json.Unmarshal([]byte(jsonStr), &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Deployment: %w", err)
	}
	return &d, nil
}

// MockDeployment creates a Deployment with sample data for testing
func MockDeployment() *Deployment {
	d := NewDeployment("mock-app", "default")

	// Add sample container
	container := Container{
		Name:  "mock-app",
		Image: "nginx:latest",
		EnvFrom: []EnvFromSource{
			{
				ConfigMapRef: &ConfigMapEnvSource{
					Name: "app-config",
				},
			},
		},
	}

	d.AddContainer(container)

	// Add sample labels
	d.AddLabel("app", "mock-app")
	d.AddLabel("environment", "dev")
	d.AddLabel("version", "1.0.0")

	// Add pod labels
	d.AddPodLabel("app", "mock-app")

	// Add selector labels
	d.AddSelectorLabel("app", "mock-app")

	// Add sample annotations
	d.AddAnnotation("description", "Mock Deployment for testing")
	d.AddAnnotation("maintainer", "crystal-team")

	return d
}
