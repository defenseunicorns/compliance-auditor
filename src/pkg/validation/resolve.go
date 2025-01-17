package validation

import (
	"fmt"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func ResolveProducer(data []byte, path, target string) (ValidationProducer, error) {
	// Check if data matches a known producer type
	if compdef, ok := isComponentDefinition(data); ok {
		return NewComponentProducer(compdef, path, target), nil
	}

	// Add more producers here
	if validations, ok := isLulaValidation(data); ok {
		return NewSimpleProducer(validations), nil
	}

	return nil, fmt.Errorf("producer is not supported")

}

func ResolveConsumer(consumerType, path string) (ResultConsumer, error) {
	// Check consumer type specified - support multiple consumer, e.g., output forms?
	switch consumerType {
	case "assessment-results":
		return NewAssessmentResultsConsumer(path), nil
	case "simple":
		return NewSimpleConsumer(), nil
	default:
		return nil, fmt.Errorf("consumer is not supported")
	}
}

func isComponentDefinition(data []byte) (*oscal.ComponentDefinition, bool) {
	componentDefinition := oscal.NewComponentDefinition()

	err := componentDefinition.NewModel(data)
	if err != nil {
		return nil, false
	}

	return componentDefinition, true
}

func isLulaValidation(data []byte) ([]common.Validation, bool) {
	validations, err := common.ReadValidationsFromYaml(data)
	if err != nil {
		return nil, false
	}
	return validations, true
}
