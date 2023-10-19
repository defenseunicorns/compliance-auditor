package oscal

import (
	"encoding/json"
	"fmt"

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
