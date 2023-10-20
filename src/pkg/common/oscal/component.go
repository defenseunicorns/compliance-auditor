package oscal

import (
	"encoding/json"
	"fmt"

	"github.com/defenseunicorns/lula/src/types"
	oscalTypes "github.com/defenseunicorns/lula/src/types/oscal"
	yaml2 "github.com/ghodss/yaml"
)

// NewOscalComponentDefintion consumes a byte arrray and returns a new single OscalComponentDefinitionModel object
// Standard use is to read a file from the filesystem and pass the []byte to this function
func NewOscalComponentDefinition(data []byte) (oscalTypes.ComponentDefinition, error) {
	var oscalComponentDefinition oscalTypes.OscalComponentDefinitionModel

	// TODO: see if we unmarshall yaml data more effectively
	jsonDoc, err := yaml2.YAMLToJSON(data)
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
		return oscalComponentDefinition.ComponentDefinition, err
	}

	err = json.Unmarshal(jsonDoc, &oscalComponentDefinition)

	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s\n", err.Error())
	}

	return oscalComponentDefinition.ComponentDefinition, nil
}

// Map an array of resources to a map of UUID to validation object
func BackMatterToMap(backMatter oscalTypes.BackMatter) map[string]types.Validation {
	resourceMap := make(map[string]types.Validation)

	for _, resource := range backMatter.Resources {

		if resource.Title == "Lula Validation" {
			var lulaSelector map[string]interface{}

			data := []byte(resource.Description)
			jsonDoc, err := yaml2.YAMLToJSON(data)
			if err != nil {
				fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
				return nil
			}

			err = json.Unmarshal(jsonDoc, &lulaSelector)

			if err != nil {
				fmt.Printf("Error unmarshalling JSON: %s\n", err.Error())
			}
			validation := types.Validation{
				Title:       resource.Title,
				Description: lulaSelector,
				Evaluated:   false,
				Result:      types.Result{},
			}

			resourceMap[resource.UUID] = validation
		}

	}
	return resourceMap
}
