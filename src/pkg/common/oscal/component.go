package oscal

import (
	"fmt"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
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

			err := yaml.Unmarshal([]byte(resource.Description), &validation)
			if err != nil {
				message.Fatalf(err, "error unmarshalling yaml: %s", err.Error())
				return nil
			}

			lulaValidation, err := validation.ToLulaValidation()
			if err != nil {
				message.Fatalf(err, "error converting validation to lula validation: %s", err.Error())
			}

			resourceMap[resource.UUID] = lulaValidation
		}

	}
	return resourceMap
}
