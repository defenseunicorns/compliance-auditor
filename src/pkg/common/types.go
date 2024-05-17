package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	"sigs.k8s.io/yaml"
)

// Data structures for ingesting validation data
type Validation struct {
	LulaVersion string    `json:"lula-version" yaml:"lula-version"`
	Metadata    *Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Provider    *Provider `json:"provider,omitempty" yaml:"provider,omitempty"`
	Domain      *Domain   `json:"domain,omitempty" yaml:"domain,omitempty"`
}

// UnmarshalYaml is a convenience method to unmarshal a Validation object from a YAML byte array
func (v *Validation) UnmarshalYaml(data []byte) error {
	return yaml.Unmarshal(data, v)
}

// MarshalYaml is a convenience method to marshal a Validation object to a YAML byte array
func (v *Validation) MarshalYaml() ([]byte, error) {
	return yaml.Marshal(v)
}

// ToResource converts a Validation object to a Resource object
func (v *Validation) ToResource() (resource *oscalTypes_1_1_2.Resource, err error) {
	resource = &oscalTypes_1_1_2.Resource{}
	resource.Title = v.Metadata.Name
	if v.Metadata.UUID != "" {
		resource.UUID = v.Metadata.UUID
	} else {
		resource.UUID = uuid.NewUUID()
	}
	validationBytes, err := v.MarshalYaml()
	if err != nil {
		return nil, err
	}
	resource.Description = string(validationBytes)
	return resource, nil
}

// TODO: Perhaps extend this structure with other needed information, such as UUID or type of validation if workflow is needed
type Metadata struct {
	Name string `json:"name" yaml:"name"`
	UUID string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type Domain struct {
	Type           string               `json:"type" yaml:"type"`
	KubernetesSpec *kube.KubernetesSpec `json:"kubernetes-spec,omitempty" yaml:"kubernetes-spec,omitempty"`
	ApiSpec        *api.ApiSpec         `json:"api-spec,omitempty" yaml:"api-spec,omitempty"`
}

type Provider struct {
	Type        string               `json:"type" yaml:"type"`
	OpaSpec     *opa.OpaSpec         `json:"opa-spec,omitempty" yaml:"opa-spec,omitempty"`
	KyvernoSpec *kyverno.KyvernoSpec `json:"kyverno-spec,omitempty" yaml:"kyverno-spec,omitempty"`
}

// ToLulaValidation converts a Validation object to a LulaValidation object
func (validation *Validation) ToLulaValidation() (lulaValidation types.LulaValidation, err error) {
	// Do version checking here to establish if the version is correct/acceptable
	currentVersion := strings.Split(config.CLIVersion, "-")[0]

	versionConstraint := currentVersion
	if validation.LulaVersion != "" {
		versionConstraint = validation.LulaVersion
	}

	if validation.Domain == nil {
		return lulaValidation, fmt.Errorf("required domain is nil")
	}
	if validation.Provider == nil {
		return lulaValidation, fmt.Errorf("required provider is nil")
	}

	validVersion, versionErr := IsVersionValid(versionConstraint, currentVersion)
	if versionErr != nil {
		return lulaValidation, fmt.Errorf("version error: %s", versionErr.Error())
	} else if !validVersion {
		return lulaValidation, fmt.Errorf("version %s does not meet the constraint %s for this validation", currentVersion, versionConstraint)
	}

	// Construct the lulaValidation object
	// TODO: Is there a better location for context?
	ctx := context.Background()
	provider := GetProvider(validation.Provider, ctx)
	if provider == nil {
		return lulaValidation, fmt.Errorf("provider %s not found", validation.Provider.Type)
	}
	lulaValidation.Provider = &provider

	domain := GetDomain(validation.Domain, ctx)
	if domain == nil {
		return lulaValidation, fmt.Errorf("domain %s not found", validation.Domain.Type)
	}
	lulaValidation.Domain = &domain

	lulaValidation.LulaValidationType = types.DefaultLulaValidationType // TODO: define workflow/purpose for this

	return lulaValidation, nil
}
