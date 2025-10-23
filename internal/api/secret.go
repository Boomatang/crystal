package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Secret represents a mock Kubernetes Secret with common fields and methods
type Secret struct {
	// Standard Kubernetes metadata
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   SecretMetadata `json:"metadata"`

	// Secret specific data
	Data       map[string][]byte `json:"data,omitempty"`
	StringData map[string]string `json:"stringData,omitempty"`

	// Secret type
	Type string `json:"type"`
}

// SecretMetadata represents the metadata section of a Secret
type SecretMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int64             `json:"generation"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// NewSecret creates a new Secret with default values
func NewSecret(name, namespace, secretType string) *Secret {
	now := time.Now()

	return &Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata: SecretMetadata{
			Name:              name,
			Namespace:         namespace,
			UID:               GenerateUID(),
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
			Labels:            make(map[string]string),
			Annotations:       make(map[string]string),
		},
		Data:       make(map[string][]byte),
		StringData: make(map[string]string),
		Type:       secretType,
	}
}

// AddData adds a key-value pair to the Secret data (base64 encoded)
func (s *Secret) AddData(key, value string) {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	s.Data[key] = []byte(encoded)
}

// AddRawData adds raw byte data to the Secret
func (s *Secret) AddRawData(key string, value []byte) {
	s.Data[key] = value
}

// AddStringData adds a key-value pair to the Secret stringData (will be base64 encoded)
func (s *Secret) AddStringData(key, value string) {
	s.StringData[key] = value
}

// GetData retrieves a decoded value from the Secret data
func (s *Secret) GetData(key string) (string, bool) {
	encoded, exists := s.Data[key]
	if !exists {
		return "", false
	}

	decoded, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return "", false
	}

	return string(decoded), true
}

// GetRawData retrieves raw byte data from the Secret
func (s *Secret) GetRawData(key string) ([]byte, bool) {
	value, exists := s.Data[key]
	return value, exists
}

// GetStringData retrieves a value from the Secret stringData
func (s *Secret) GetStringData(key string) (string, bool) {
	value, exists := s.StringData[key]
	return value, exists
}

// AddLabel adds a label to the Secret
func (s *Secret) AddLabel(key, value string) {
	s.Metadata.Labels[key] = value
}

// AddAnnotation adds an annotation to the Secret
func (s *Secret) AddAnnotation(key, value string) {
	s.Metadata.Annotations[key] = value
}

// SetType sets the secret type
func (s *Secret) SetType(secretType string) {
	s.Type = secretType
}

// ToJSON converts the Secret to JSON string
func (s *Secret) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Secret to JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON creates a Secret from JSON string
func SecretFromJSON(jsonStr string) (*Secret, error) {
	var s Secret
	err := json.Unmarshal([]byte(jsonStr), &s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Secret: %w", err)
	}
	return &s, nil
}

// Clone creates a deep copy of the Secret
func (s *Secret) Clone() *Secret {
	clone := &Secret{
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Metadata: SecretMetadata{
			Name:              s.Metadata.Name,
			Namespace:         s.Metadata.Namespace,
			UID:               s.Metadata.UID,
			ResourceVersion:   s.Metadata.ResourceVersion,
			Generation:        s.Metadata.Generation,
			CreationTimestamp: s.Metadata.CreationTimestamp,
			Labels:            make(map[string]string),
			Annotations:       make(map[string]string),
		},
		Data:       make(map[string][]byte),
		StringData: make(map[string]string),
		Type:       s.Type,
	}

	// Copy labels
	for k, v := range s.Metadata.Labels {
		clone.Metadata.Labels[k] = v
	}

	// Copy annotations
	for k, v := range s.Metadata.Annotations {
		clone.Metadata.Annotations[k] = v
	}

	// Copy data
	for k, v := range s.Data {
		clone.Data[k] = make([]byte, len(v))
		copy(clone.Data[k], v)
	}

	// Copy string data
	for k, v := range s.StringData {
		clone.StringData[k] = v
	}

	return clone
}

// MockSecret creates a Secret with sample data for testing
func MockSecret() *Secret {
	s := NewSecret("mock-secret", "default", "Opaque")

	// Add sample data (will be base64 encoded)
	s.AddData("username", "admin")
	s.AddData("password", "password123")
	s.AddData("api-key", "abcdefghijk")
	s.AddData("database-url", "postgresql://localhost:5432/mydb")

	// Add sample string data
	s.AddStringData("redis-password", "redispassword")
	s.AddStringData("jwt-secret", "myjwtsecretkey123456")

	// Add sample labels
	s.AddLabel("app", "crystal")
	s.AddLabel("environment", "dev")
	s.AddLabel("version", "1.0.0")

	// Add sample annotations
	s.AddAnnotation("description", "Mock Secret for testing")
	s.AddAnnotation("maintainer", "crystal-team")

	return s
}

// Common Secret types
const (
	SecretTypeOpaque              = "Opaque"
	SecretTypeServiceAccountToken = "kubernetes.io/service-account-token"
	SecretTypeDockercfg           = "kubernetes.io/dockercfg"
	SecretTypeDockerConfigJson    = "kubernetes.io/dockerconfigjson"
	SecretTypeBasicAuth           = "kubernetes.io/basic-auth"
	SecretTypeSSHAuth             = "kubernetes.io/ssh-auth"
	SecretTypeTLS                 = "kubernetes.io/tls"
	SecretTypeBootstrapToken      = "bootstrap.kubernetes.io/token"
)
