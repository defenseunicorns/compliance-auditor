package kube

import (
	"context"
	"errors"

	"github.com/defenseunicorns/lula/src/types"
)

type KubernetesDomain struct {
	// Context is the context that Kubernetes resources are being evaluated in
	Context context.Context `json:"context" yaml:"context"`

	// Spec is the specification of the Kubernetes resources
	Spec *KubernetesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (k KubernetesDomain) GetResources() (types.DomainResources, error) {
	var resources types.DomainResources

	if k.Spec.Wait != nil {
		err := EvaluateWait(*k.Spec.Wait)
		if err != nil {
			return nil, err
		}
	}

	resources, err := QueryCluster(k.Context, k.Spec.Resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

type KubernetesSpec struct {
	Resources []Resource `json:"resources" yaml:"resources"`
	Wait      *Wait      `json:"wait,omitempty" yaml:"wait,omitempty"`
}

type Resource struct {
	Name         string        `json:"name" yaml:"name"`
	Description  string        `json:"description" yaml:"description"`
	ResourceRule *ResourceRule `json:"resource-rule,omitempty" yaml:"resource-rule,omitempty"`
}

type ResourceRule struct {
	Name       string   `json:"name" yaml:"name"`
	Group      string   `json:"group" yaml:"group"`
	Version    string   `json:"version" yaml:"version"`
	Resource   string   `json:"resource" yaml:"resource"`
	Namespaces []string `json:"namespaces" yaml:"namespaces"`
	Field      *Field   `json:"field,omitempty" yaml:"field,omitempty"`
}

type FieldType string

const (
	FieldTypeJSON    FieldType = "json"
	FieldTypeYAML    FieldType = "yaml"
	DefaultFieldType FieldType = FieldTypeJSON
)

type Field struct {
	Jsonpath string    `json:"jsonpath" yaml:"jsonpath"`
	Type     FieldType `json:"type" yaml:"type"`
	Base64   bool      `json:"base64" yaml:"base64"`
}

// Validate the Field type if valid
func (f Field) Validate() error {
	switch f.Type {
	case FieldTypeJSON, FieldTypeYAML:
		return nil
	default:
		return errors.New("field Type must be 'json' or 'yaml'")
	}
}

type Wait struct {
	Condition string `json:"condition" yaml:"condition"`
	Jsonpath  string `json:"jsonpath" yaml:"jsonpath"`
	Kind      string `json:"kind" yaml:"kind"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Timeout   string `json:"timeout" yaml:"timeout"`
}
