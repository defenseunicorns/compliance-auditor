package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
)

// Data structures for ingesting validation data
type Validation struct {
	LulaVersion string   `json:"lula-version" yaml:"lula-version"`
	Metadata    Metadata `json:"metadata" yaml:"metadata"`
	Provider    Provider `json:"provider" yaml:"provider"`
	Domain      Domain   `json:"domain" yaml:"domain"`
}

// TODO: Perhaps extend this structure with other needed information, such as UUID or type of validation if workflow is needed
type Metadata struct {
	Name string `json:"name" yaml:"name"`
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
	var result types.Result
	var evaluated bool
	currentVersion := strings.Split(config.CLIVersion, "-")[0]

	versionConstraint := currentVersion
	if validation.LulaVersion != "" {
		versionConstraint = validation.LulaVersion
	}

	validVersion, versionErr := IsVersionValid(versionConstraint, currentVersion)
	if versionErr != nil {
		result.Failing = 1
		result.Observations = map[string]string{"Lula Version Error": versionErr.Error()}
		evaluated = true
	} else if !validVersion {
		result.Failing = 1
		result.Observations = map[string]string{"Version Constraint Incompatible": "Lula Version does not meet the constraint for this validation."}
		evaluated = true
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
	lulaValidation.Evaluated = evaluated
	lulaValidation.Result = result

	return lulaValidation, nil
}
