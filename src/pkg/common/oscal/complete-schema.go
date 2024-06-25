package oscal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"
	yamlV3 "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

func NewOscalModel(data []byte) (*oscalTypes_1_1_2.OscalModels, error) {
	oscalModel := oscalTypes_1_1_2.OscalModels{}

	err := multiModelValidate(data)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &oscalModel)
	if err != nil {
		return nil, err
	}
	return &oscalModel, nil
}

// WriteOscalModel takes a path and writes content to a file while performing checks for existing content
// supports both json and yaml
func WriteOscalModel(filePath string, model *oscalTypes_1_1_2.OscalModels) error {

	modelType, err := GetOscalModel(model)
	if err != nil {
		return err
	}

	// if no path or directory add default filename
	if filepath.Ext(filePath) == "" {
		filePath = filepath.Join(filePath, fmt.Sprintf("%s.yaml", modelType))
	} else {
		if err := files.IsJsonOrYaml(filePath); err != nil {
			return err
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		// If the file exists - read the data into the model
		existingFileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		existingModel, err := NewOscalModel(existingFileBytes)
		if err != nil {
			return err
		}

		existingModelType, err := GetOscalModel(existingModel)
		if err != nil {
			return nil
		}

		if existingModelType != modelType {
			return fmt.Errorf("cannot merge model %s with existing model %s", modelType, existingModelType)
		}
		// Merge the existing model with the new model
		// re-assign to perform common operations below
		model, err = MergeOscalModels(existingModel, model, modelType)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer

	if filepath.Ext(filePath) == ".json" {
		jsonEncoder := json.NewEncoder(&b)
		jsonEncoder.SetIndent("", "  ")
		jsonEncoder.Encode(model)
	} else {
		yamlEncoder := yamlV3.NewEncoder(&b)
		yamlEncoder.SetIndent(2)
		yamlEncoder.Encode(model)
	}

	err = files.WriteOutput(b.Bytes(), filePath)
	if err != nil {
		return err
	}

	message.Infof("OSCAL artifact written to: %s", filePath)

	return nil

}

func MergeOscalModels(existingModel *oscalTypes_1_1_2.OscalModels, newModel *oscalTypes_1_1_2.OscalModels, modelType string) (*oscalTypes_1_1_2.OscalModels, error) {
	var err error
	// Now to check each model type - currently only component definition and assessment-results apply

	// Component definition
	if modelType == "component" {
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
	if modelType == "assessment-results" {
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

func GetOscalModel(model *oscalTypes_1_1_2.OscalModels) (modelType string, err error) {

	// Check if one model present and all other nil - is there a better way to do this?
	models := make([]string, 0)

	if model.Catalog != nil {
		models = append(models, "catalog")
	}

	if model.Profile != nil {
		models = append(models, "profile")
	}

	if model.ComponentDefinition != nil {
		models = append(models, "component")
	}

	if model.SystemSecurityPlan != nil {
		models = append(models, "system-security-plan")
	}

	if model.AssessmentPlan != nil {
		models = append(models, "assessment-plan")
	}

	if model.AssessmentResults != nil {
		models = append(models, "assessment-results")
	}

	if model.PlanOfActionAndMilestones != nil {
		models = append(models, "poam")
	}

	if len(models) > 1 {
		return "", fmt.Errorf("%v models identified when only oneOf is permitted", len(models))
	} else {
		return models[0], nil
	}

}
