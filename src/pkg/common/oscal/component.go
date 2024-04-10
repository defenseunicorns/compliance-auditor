package oscal

import (
	"context"
	"fmt"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"

	"sigs.k8s.io/yaml"
)

// NewOscalComponentDefinition consumes a byte array and returns a new single OscalComponentDefinitionModel object
// Standard use is to read a file from the filesystem and pass the []byte to this function
func NewOscalComponentDefinition(data []byte) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return componentDefinition, err
	}

	if oscalModels.ComponentDefinition == nil {
		return componentDefinition, fmt.Errorf("No Component Definition found in the provided data")
	}

	return *oscalModels.ComponentDefinition, nil
}

// Map an array of resources to a map of UUID to lulaValidation object
func BackMatterToMap(backMatter oscalTypes_1_1_2.BackMatter) map[string]types.LulaValidation {
	resourceMap := make(map[string]types.LulaValidation)

	if backMatter.Resources == nil {
		return nil
	}

	for _, resource := range *backMatter.Resources {
		// TODO: Possibly support different title values (e.g., "Placeholder", "Healthcheck")
		if resource.Title == "Lula Validation" {
			var validation common.Validation
			var lulaValidation types.LulaValidation

			err := yaml.Unmarshal([]byte(resource.Description), &validation)
			if err != nil {
				message.Fatalf(err, "error unmarshalling yaml: %s", err.Error())
				return nil
			}

			// Do version checking here to establish if the version is correct/acceptable
			var result types.Result
			var evaluated bool
			currentVersion := strings.Split(config.CLIVersion, "-")[0]

			versionConstraint := currentVersion
			if validation.LulaVersion != "" {
				versionConstraint = validation.LulaVersion
			}

			validVersion, versionErr := common.IsVersionValid(versionConstraint, currentVersion)
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
			lulaValidation.Provider = common.GetProvider(validation.Provider, ctx)
			if lulaValidation.Provider == nil {
				message.Fatalf(nil, "provider %s not found", validation.Provider.Type)
				return nil
			}
			lulaValidation.Domain = common.GetDomain(validation.Domain, ctx)
			if lulaValidation.Domain == nil {
				message.Fatalf(nil, "domain %s not found", validation.Domain.Type)
				return nil
			}
			lulaValidation.LulaValidationType = types.DefaultLulaValidationType // TODO: define workflow/purpose for this
			lulaValidation.Evaluated = evaluated
			lulaValidation.Result = result

			resourceMap[resource.UUID] = lulaValidation
		}

	}
	return resourceMap
}
