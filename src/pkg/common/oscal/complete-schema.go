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

func MergeOscalModels(existingModel *oscalTypes_1_1_2.OscalModels, model *oscalTypes_1_1_2.OscalModels) (*oscalTypes_1_1_2.OscalModels, error) {

	// Now to check each model type - currently only component definition and assessment-results apply

	return existingModel, nil
}
