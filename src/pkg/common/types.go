package common

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common/schemas"
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	"sigs.k8s.io/yaml"
)

// Define base errors for validations
var (
	ErrInvalidSchema   = errors.New("schema is invalid")
	ErrInvalidYaml     = errors.New("error unmarshaling JSON")
	ErrInvalidVersion  = errors.New("version is invalid")
	ErrInvalidDomain   = errors.New("domain is invalid")
	ErrInvalidProvider = errors.New("provider is invalid")
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

// Metadata is a structure that contains the name and uuid of a validation
type Metadata struct {
	Name string `json:"name" yaml:"name"`
	UUID string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

// Domain is a structure that contains the domain type and the corresponding spec
type Domain struct {
	// Type is the type of domain: enum: kubernetes, api
	Type string `json:"type" yaml:"type"`
	// KubernetesSpec is the specification for a Kubernetes domain, required if type is kubernetes
	KubernetesSpec *kube.KubernetesSpec `json:"kubernetes-spec,omitempty" yaml:"kubernetes-spec,omitempty"`
	// ApiSpec is the specification for an API domain, required if type is api
	ApiSpec *api.ApiSpec `json:"api-spec,omitempty" yaml:"api-spec,omitempty"`
}

type Provider struct {
	Type        string               `json:"type" yaml:"type"`
	OpaSpec     *opa.OpaSpec         `json:"opa-spec,omitempty" yaml:"opa-spec,omitempty"`
	KyvernoSpec *kyverno.KyvernoSpec `json:"kyverno-spec,omitempty" yaml:"kyverno-spec,omitempty"`
}

// Lint is a convenience method to lint a Validation object
func (validation *Validation) Lint() oscalValidation.ValidationResult {
	validationBytes, err := validation.MarshalYaml()
	if err != nil {
		return *oscalValidation.NewNonSchemaValidationError(err, &oscalValidation.ValidationParams{ModelType: "validation"})
	}
	return schemas.Validate("validation", validationBytes)
}

// ToLulaValidation converts a Validation object to a LulaValidation object
func (validation *Validation) ToLulaValidation() (lulaValidation types.LulaValidation, err error) {
	// Do version checking here to establish if the version is correct/acceptable
	currentVersion := strings.Split(config.CLIVersion, "-")[0]

	versionConstraint := currentVersion
	if validation.LulaVersion != "" {
		versionConstraint = validation.LulaVersion
	}

	lintResult := validation.Lint()
	// If the validation is not valid, return the error
	if oscalValidation.IsNonSchemaValidationError(&lintResult) {
		return lulaValidation, fmt.Errorf("%w: %v", ErrInvalidSchema, oscalValidation.GetNonSchemaError(&lintResult))
	} else if !lintResult.Valid {
		return lulaValidation, fmt.Errorf("%w: %v", ErrInvalidSchema, lintResult.Errors)
	}

	validVersion, versionErr := IsVersionValid(versionConstraint, currentVersion)
	if versionErr != nil {
		return lulaValidation, fmt.Errorf("%w: %s", ErrInvalidVersion, versionErr.Error())
	} else if !validVersion {
		return lulaValidation, fmt.Errorf("%w: version %s does not meet the constraint %s for this validation", ErrInvalidVersion, currentVersion, versionConstraint)
	}

	// Construct the lulaValidation object
	// TODO: Is there a better location for context?
	ctx := context.Background()

	domain, err := GetDomain(validation.Domain, ctx)
	if domain == nil {
		return lulaValidation, fmt.Errorf("%w: %s", ErrInvalidDomain, validation.Domain.Type)
	} else if err != nil {
		return lulaValidation, fmt.Errorf("%w: %v", ErrInvalidDomain, err)
	}
	lulaValidation.Domain = &domain

	provider, err := GetProvider(validation.Provider, ctx)
	if provider == nil {
		return lulaValidation, fmt.Errorf("%w: %s", ErrInvalidProvider, validation.Provider.Type)
	} else if err != nil {
		return lulaValidation, fmt.Errorf("%w: %v", ErrInvalidProvider, err)
	}
	lulaValidation.Provider = &provider

	lulaValidation.LulaValidationType = types.DefaultLulaValidationType // TODO: define workflow/purpose for this

	if validation.Metadata == nil {
		lulaValidation.Name = "lula-validation"
	} else {
		lulaValidation.Name = validation.Metadata.Name
	}

	return lulaValidation, nil
}
