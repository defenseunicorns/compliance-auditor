package oscal

import (
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1/component-definition"
	"github.com/defenseunicorns/lula/src/types"
	"gopkg.in/yaml.v3"
)

// NewOscalComponentDefintion consumes a byte arrray and returns a new single OscalComponentDefinitionModel object
// Standard use is to read a file from the filesystem and pass the []byte to this function
func NewOscalComponentDefinition(data []byte) (oscalTypes.ComponentDefinition, error) {
	var oscalComponentDefinition oscalTypes.OscalComponentDefinitionModel

	err := yaml.Unmarshal(data, &oscalComponentDefinition)
	if err != nil {
		fmt.Printf("Error marshalling yaml: %s\n", err.Error())
		return oscalComponentDefinition.ComponentDefinition, err
	}

	return oscalComponentDefinition.ComponentDefinition, nil
}

// Map an array of resources to a map of UUID to validation object
func BackMatterToMap(backMatter oscalTypes.BackMatter) map[string]types.Validation {
	resourceMap := make(map[string]types.Validation)

	for _, resource := range backMatter.Resources {
		if resource.Title == "Lula Validation" {
			var lulaSelector map[string]interface{}

			err := yaml.Unmarshal([]byte(resource.Description), &lulaSelector)
			if err != nil {
				fmt.Printf("Error marshalling yaml: %s\n", err.Error())
				return nil
			}

			validation := types.Validation{
				Title:       resource.Title,
				Description: lulaSelector["target"].(map[string]interface{}),
				Evaluated:   false,
				Result:      types.Result{},
			}

			resourceMap[resource.UUID] = validation
		}

	}
	return resourceMap
}
