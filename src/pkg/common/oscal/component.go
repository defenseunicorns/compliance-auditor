package oscal

import (
	"fmt"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"

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
		return componentDefinition, fmt.Errorf("no Component Definition found in the provided data")
	}

	return *oscalModels.ComponentDefinition, nil
}

// Returns a map of the uuid - description of the back-matter resources
func BackMatterToMap(backMatter oscalTypes_1_1_2.BackMatter) (resourceMap map[string]string) {
	if backMatter.Resources == nil {
		return nil
	}

	resourceMap = make(map[string]string)
	for _, resource := range *backMatter.Resources {
		// perform a check to see if the key already exists (meaning duplicitive uuid use)
		_, exists := resourceMap[resource.UUID]
		if exists {
			message.Warnf("Duplicative UUID use detected - Overwriting UUID %s", resource.UUID)
		}

		resourceMap[resource.UUID] = resource.Description
	}
	return resourceMap

}
