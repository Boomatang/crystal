package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConfigMap represents a mock Kubernetes ConfigMap with common fields and methods
type ConfigMap struct {
	// Standard Kubernetes metadata
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   ConfigMapMetadata `json:"metadata"`

	// ConfigMap specific data
	Data       map[string]string `json:"data"`
	BinaryData map[string][]byte `json:"binaryData,omitempty"`

	// Immutable flag
	Immutable *bool `json:"immutable,omitempty"`
}

// ConfigMapMetadata represents the metadata section of a ConfigMap
type ConfigMapMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int64             `json:"generation"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// NewConfigMap creates a new ConfigMap with default values
func NewConfigMap(name, namespace string) *ConfigMap {
	now := time.Now()
	immutable := false

	return &ConfigMap{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: ConfigMapMetadata{
			Name:              name,
			Namespace:         namespace,
			UID:               GenerateUID(),
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
			Labels:            make(map[string]string),
			Annotations:       make(map[string]string),
		},
		Data:       make(map[string]string),
		BinaryData: make(map[string][]byte),
		Immutable:  &immutable,
	}
}

// AddData adds a key-value pair to the ConfigMap data
func (cm *ConfigMap) AddData(key, value string) {
	cm.Data[key] = value
}

// AddBinaryData adds binary data to the ConfigMap
func (cm *ConfigMap) AddBinaryData(key string, value []byte) {
	cm.BinaryData[key] = value
}

// GetData retrieves a value from the ConfigMap data
func (cm *ConfigMap) GetData(key string) (string, bool) {
	value, exists := cm.Data[key]
	return value, exists
}

// GetBinaryData retrieves binary data from the ConfigMap
func (cm *ConfigMap) GetBinaryData(key string) ([]byte, bool) {
	value, exists := cm.BinaryData[key]
	return value, exists
}

// AddLabel adds a label to the ConfigMap
func (cm *ConfigMap) AddLabel(key, value string) {
	cm.Metadata.Labels[key] = value
}

// AddAnnotation adds an annotation to the ConfigMap
func (cm *ConfigMap) AddAnnotation(key, value string) {
	cm.Metadata.Annotations[key] = value
}

// SetImmutable sets the immutable flag
func (cm *ConfigMap) SetImmutable(immutable bool) {
	cm.Immutable = &immutable
}

// ToJSON converts the ConfigMap to JSON string
func (cm *ConfigMap) ToJSON() (string, error) {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal ConfigMap to JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON creates a ConfigMap from JSON string
func FromJSON(jsonStr string) (*ConfigMap, error) {
	var cm ConfigMap
	err := json.Unmarshal([]byte(jsonStr), &cm)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to ConfigMap: %w", err)
	}
	return &cm, nil
}

// Clone creates a deep copy of the ConfigMap
func (cm *ConfigMap) Clone() *ConfigMap {
	clone := &ConfigMap{
		APIVersion: cm.APIVersion,
		Kind:       cm.Kind,
		Metadata: ConfigMapMetadata{
			Name:              cm.Metadata.Name,
			Namespace:         cm.Metadata.Namespace,
			UID:               cm.Metadata.UID,
			ResourceVersion:   cm.Metadata.ResourceVersion,
			Generation:        cm.Metadata.Generation,
			CreationTimestamp: cm.Metadata.CreationTimestamp,
			Labels:            make(map[string]string),
			Annotations:       make(map[string]string),
		},
		Data:       make(map[string]string),
		BinaryData: make(map[string][]byte),
	}

	// Copy labels
	for k, v := range cm.Metadata.Labels {
		clone.Metadata.Labels[k] = v
	}

	// Copy annotations
	for k, v := range cm.Metadata.Annotations {
		clone.Metadata.Annotations[k] = v
	}

	// Copy data
	for k, v := range cm.Data {
		clone.Data[k] = v
	}

	// Copy binary data
	for k, v := range cm.BinaryData {
		clone.BinaryData[k] = make([]byte, len(v))
		copy(clone.BinaryData[k], v)
	}

	// Copy immutable flag
	if cm.Immutable != nil {
		immutable := *cm.Immutable
		clone.Immutable = &immutable
	}

	return clone
}

// GenerateUID generates a mock UID for testing purposes
func GenerateUID() string {
	return fmt.Sprintf("mock-uid-%d", time.Now().UnixNano())
}

// MockConfigMap creates a ConfigMap with sample data for testing
func MockConfigMap() *ConfigMap {
	cm := NewConfigMap("mock-config", "default")

	// Add sample data
	cm.AddData("DATABASE_URL", "postgresql://localhost:5432/mydb")
	cm.AddData("API_KEY", "mock-api-key-12345")
	cm.AddData("ENVIRONMENT", "development")
	cm.AddData("LOG_LEVEL", "debug")

	// Add sample binary data
	cm.AddBinaryData("cert.pem", []byte("-----BEGIN CERTIFICATE-----\nMOCK CERTIFICATE DATA\n-----END CERTIFICATE-----"))

	// Add sample labels
	cm.AddLabel("app", "crystal")
	cm.AddLabel("environment", "dev")
	cm.AddLabel("version", "1.0.0")

	// Add sample annotations
	cm.AddAnnotation("description", "Mock ConfigMap for testing")
	cm.AddAnnotation("maintainer", "crystal-team")

	return cm
}
