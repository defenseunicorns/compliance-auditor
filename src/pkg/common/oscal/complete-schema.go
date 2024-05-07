package oscal

import (
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"sigs.k8s.io/yaml"
)

func NewOscalModel(data []byte) (*oscalTypes_1_1_2.OscalModels, error) {
	oscalModel := oscalTypes_1_1_2.OscalModels{}
	err := yaml.Unmarshal(data, &oscalModel)
	if err != nil {
		return nil, err
	}
	return &oscalModel, nil
}

func MergeOscalModels(existingModel *oscalTypes_1_1_2.OscalModels, newModel *oscalTypes_1_1_2.OscalModels) (*oscalTypes_1_1_2.OscalModels, error) {
	var err error
	// Now to check each model type - currently only component definition and assessment-results apply

	// Component definition
	if existingModel.ComponentDefinition != nil && newModel.ComponentDefinition != nil {
		merged, err := MergeComponentDefinitions(existingModel.ComponentDefinition, newModel.ComponentDefinition)
		if err != nil {
			return nil, err
		}
		// Re-assign after processing errors
		existingModel.ComponentDefinition = merged
	} else if existingModel.ComponentDefinition == nil && newModel.ComponentDefinition != nil {
		existingModel.ComponentDefinition = newModel.ComponentDefinition
	}

	// Assessment Results
	if existingModel.AssessmentResults != nil && newModel.AssessmentResults != nil {
		merged, err := MergeAssessmentResults(existingModel.AssessmentResults, newModel.AssessmentResults)
		if err != nil {
			return existingModel, err
		}
		// Re-assign after processing errors
		existingModel.AssessmentResults = merged
	} else if existingModel.AssessmentResults == nil && newModel.AssessmentResults != nil {
		existingModel.AssessmentResults = newModel.AssessmentResults
	}

	return existingModel, err
}
