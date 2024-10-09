package kube

import (
	"context"
	"errors"
	"fmt"

	"github.com/defenseunicorns/lula/src/types"
)

type KubernetesDomain struct {
	// Spec is the specification of the Kubernetes resources
	Spec *KubernetesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func CreateKubernetesDomain(spec *KubernetesSpec) (types.Domain, error) {
	// Check validity of spec
	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	if spec.Resources == nil && spec.CreateResources == nil && spec.Wait == nil {
		return nil, fmt.Errorf("one of resources, create-resources, or wait must be specified")
	}

	if spec.Resources != nil {
		for _, resource := range spec.Resources {
			if resource.Name == "" {
				return nil, fmt.Errorf("resource name cannot be empty")
			}
			if resource.ResourceRule == nil {
				return nil, fmt.Errorf("resource rule cannot be nil")
			}
			if resource.ResourceRule.Resource == "" {
				return nil, fmt.Errorf("resource rule resource cannot be empty")
			}
			if resource.ResourceRule.Version == "" {
				return nil, fmt.Errorf("resource rule version cannot be empty")
			}
			if resource.ResourceRule.Name != "" && len(resource.ResourceRule.Namespaces) > 1 {
				return nil, fmt.Errorf("named resource requested cannot be returned from multiple namespaces")
			}
			if resource.ResourceRule.Field != nil {
				if resource.ResourceRule.Field.Type == "" {
					resource.ResourceRule.Field.Type = DefaultFieldType
				}
				err := resource.ResourceRule.Field.Validate()
				if err != nil {
					return nil, err
				}
				if resource.ResourceRule.Name == "" {
					return nil, fmt.Errorf("field cannot be specified without resource name")
				}
			}
		}
	}

	if spec.Wait != nil {
		if spec.Wait.Resource == "" {
			return nil, fmt.Errorf("wait resource cannot be empty")
		}
		if spec.Wait.Version == "" {
			return nil, fmt.Errorf("wait version cannot be empty")
		}
		if spec.Wait.Name == "" {
			return nil, fmt.Errorf("wait name cannot be empty")
		}
	}

	if spec.CreateResources != nil {
		for _, resource := range spec.CreateResources {
			if resource.Name == "" {
				return nil, fmt.Errorf("resource name cannot be empty")
			}
			if resource.Manifest == "" && resource.File == "" {
				return nil, fmt.Errorf("resource manifest or file must be specified")
			}
			if resource.Manifest != "" && resource.File != "" {
				return nil, fmt.Errorf("only resource manifest or file can be specified")
			}
		}
	}

	return KubernetesDomain{
		Spec: spec,
	}, nil
}

// GetResources returns the resources from the Kubernetes domain
// Evaluates the `create-resources` first, `wait` second, and finally `resources` last
func (k KubernetesDomain) GetResources(ctx context.Context) (resources types.DomainResources, err error) {
	cluster, err := GetCluster()
	if err != nil {
		return nil, err
	}

	// Evaluate the wait condition
	if k.Spec.Wait != nil {
		err := EvaluateWait(ctx, cluster, *k.Spec.Wait)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Return both?
	if k.Spec.Resources != nil {
		resources, err = QueryCluster(ctx, cluster, k.Spec.Resources)
		if err != nil {
			return nil, err
		}
	} else if k.Spec.CreateResources != nil {
		resources, err = CreateE2E(ctx, cluster, k.Spec.CreateResources)
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

func (k KubernetesDomain) IsExecutable() bool {
	// Domain is only executable if create-resources is not nil
	return len(k.Spec.CreateResources) > 0
}

type KubernetesSpec struct {
	Resources       []Resource       `json:"resources" yaml:"resources"`
	Wait            *Wait            `json:"wait,omitempty" yaml:"wait,omitempty"`
	CreateResources []CreateResource `json:"create-resources" yaml:"create-resources"`
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
	Name      string `json:"name" yaml:"name"`
	Group     string `json:"group" yaml:"group"`
	Version   string `json:"version" yaml:"version"`
	Resource  string `json:"resource" yaml:"resource"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Timeout   string `json:"timeout" yaml:"timeout"`
}

type CreateResource struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Manifest  string `json:"manifest" yaml:"manifest"`
	File      string `json:"file" yaml:"file"`
}
