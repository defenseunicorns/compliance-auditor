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
	LulaVersion string   `json:"lula-version" yaml:"lula-version"`
	Metadata    Metadata `json:"metadata" yaml:"metadata"`
	Provider    Provider `json:"provider" yaml:"provider"`
	Domain      Domain   `json:"domain" yaml:"domain"`
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
	Type           string              `json:"type" yaml:"type"`
	KubernetesSpec kube.KubernetesSpec `json:"kubernetes-spec" yaml:"kubernetes-spec"`
	ApiSpec        api.ApiSpec         `json:"api-spec" yaml:"api-spec"`
}

type Provider struct {
	Type        string              `json:"type" yaml:"type"`
	OpaSpec     opa.OpaSpec         `json:"opa-spec" yaml:"opa-spec"`
	KyvernoSpec kyverno.KyvernoSpec `json:"kyverno-spec" yaml:"kyverno-spec"`
}

// ToLulaValidation converts a Validation object to a LulaValidation object
func (validation *Validation) ToLulaValidation() (lulaValidation types.LulaValidation, err error) {
	// Do version checking here to establish if the version is correct/acceptable
	currentVersion := strings.Split(config.CLIVersion, "-")[0]

	versionConstraint := currentVersion
	if validation.LulaVersion != "" {
		versionConstraint = validation.LulaVersion
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
	lulaValidation.Provider = GetProvider(validation.Provider, ctx)
	if lulaValidation.Provider == nil {
		return lulaValidation, fmt.Errorf("provider %s not found", validation.Provider.Type)
	}
	lulaValidation.Domain = GetDomain(validation.Domain, ctx)
	if lulaValidation.Domain == nil {
		return lulaValidation, fmt.Errorf("domain %s not found", validation.Domain.Type)
	}

	lulaValidation.LulaValidationType = types.DefaultLulaValidationType // TODO: define workflow/purpose for this

	return lulaValidation, nil
}
