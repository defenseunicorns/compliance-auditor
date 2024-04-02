package types

import (
	"errors"
	"fmt"
	"reflect"

	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type Validation struct {
	Title       string `json:"title" yaml:"title"`
	LulaVersion string `json:"lula-version" yaml:"lula-version"`
	Target      Target `json:"target" yaml:"target"`
	Evaluated   bool   `json:"evaluated" yaml:"evaluated"`
	Result      Result `json:"result" yaml:"result"`
}

// native type for conversion to targeted report format
type Result struct {
	UUID         string            `json:"uuid" yaml:"uuid"`
	ControlId    string            `json:"control-id" yaml:"control-id"`
	Description  string            `json:"description" yaml:"description"`
	Passing      int               `json:"passing" yaml:"passing"`
	Failing      int               `json:"failing" yaml:"failing"`
	State        string            `json:"state" yaml:"state"`
	Observations map[string]string `json:"observations" yaml:"observations"`
}

// Current placeholder for all requisite data in the payload
// Fields will be populated as required otherwise left empty
// This could be expanded as providers add more fields
type Payload struct {
	Resources []Resource              `json:"resources" yaml:"resources"`
	Requests  []Request               `mapstructure:"requests" json:"requests" yaml:"requests"`
	Wait      Wait                    `json:"wait" yaml:"wait"`
	Rego      string                  `json:"rego" yaml:"rego"`
	Kyverno   *kjson.ValidatingPolicy `json:"kyverno" yaml:"kyverno"`
	Output    Output                  `json:"output" yaml:"output"`
}

type Output struct {
	Validation   string   `json:"validation" yaml:"validation"`
	Observations []string `json:"observations" yaml:"observations"`
}

type Resource struct {
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	ResourceRule ResourceRule `json:"resource-rule" yaml:"resource-rule"`
}

type Wait struct {
	Condition string `json:"condition" yaml:"condition"`
	Jsonpath  string `json:"jsonpath" yaml:"jsonpath"`
	Kind      string `json:"kind" yaml:"kind"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Timeout   string `json:"timeout" yaml:"timeout"`
}

type Request struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type Target struct {
	Provider string  `json:"provider" yaml:"provider"`
	Domain   string  `json:"domain" yaml:"domain"`
	Payload  Payload `json:"payload" yaml:"payload"`
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

type ResourceRule struct {
	Name       string   `json:"name" yaml:"name"`
	Group      string   `json:"group" yaml:"group"`
	Version    string   `json:"version" yaml:"version"`
	Resource   string   `json:"resource" yaml:"resource"`
	Namespaces []string `json:"namespaces" yaml:"namespaces"`
	Field      Field    `json:"field" yaml:"field"`
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

// Lint the validation
func (validation *Validation) Lint() error {
	if validation.Title == "" {
		return fmt.Errorf("validation title is required")
	}

	// Requires a target
	if reflect.DeepEqual(validation.Target, Target{}) {
		return fmt.Errorf("validation target is required")
	}

	// Requires a payload
	if reflect.DeepEqual(validation.Target.Payload, Payload{}) {
		return fmt.Errorf("validation target payload is required")
	}

	// Requires resources
	if len(validation.Target.Payload.Resources) == 0 {
		return fmt.Errorf("validation target resources are required")
	}

	// Iterate through each resource and check if the rule has all the required fields
	for _, resource := range validation.Target.Payload.Resources {
		// get the resource rule
		rule := resource.ResourceRule

		// Requires a version
		if rule.Version == "" {
			return fmt.Errorf("resource %s has no version", rule.Name)
		}

		// Requires a resource
		if rule.Resource == "" {
			return fmt.Errorf("resource %s has no resource", rule.Name)
		}

		// Requires a namespace if the resource has a name
		if rule.Name != "" && len(rule.Namespaces) == 0 {
			return fmt.Errorf("resource %s has no namespaces", rule.Name)
		}

		// Requires a name if the resource has a field
		if !reflect.DeepEqual(rule.Field, Field{}) && rule.Name == "" {
			return fmt.Errorf("resource-rule with field must have a name")
		}
	}
	return nil
}
