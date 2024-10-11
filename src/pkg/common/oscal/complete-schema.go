package oscal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/inject"
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
			return fmt.Errorf("error reading file: %v", err)
		}
		existingModel, err := NewOscalModel(existingFileBytes)
		if err != nil {
			return fmt.Errorf("error getting existing model: %v", err)
		}

		existingModelType, err := GetOscalModel(existingModel)
		if err != nil {
			return fmt.Errorf("error getting existing model type: %v", err)
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
	// If the deterministic update is applied here - Lula will fix OSCAL that was previously written
	// or generated outside of Lula workflows
	switch modelType {
	case "component":
		MakeComponentDeterminstic(model.ComponentDefinition)
	case "assessment-results":
		MakeAssessmentResultsDeterministic(model.AssessmentResults)
	}

	b, err := ConvertOSCALToBytes(model, filepath.Ext(filePath))
	if err != nil {
		return fmt.Errorf("error converting OSCAL model to bytes: %v", err)
	}

	err = files.WriteOutput(b, filePath)
	if err != nil {
		return err
	}

	message.Infof("OSCAL artifact written to: %s", filePath)

	return nil

}

// OverwriteOscalModel takes a path and writes content to a file - does not check for existing content
// supports both json and yaml
func OverwriteOscalModel(filePath string, model *oscalTypes_1_1_2.OscalModels) error {

	// if no path or directory add default filename
	if filepath.Ext(filePath) == "" {
		filePath = filepath.Join(filePath, fmt.Sprintf("%s.yaml", "oscal"))
	} else {
		if err := files.IsJsonOrYaml(filePath); err != nil {
			return err
		}
	}

	// Make deterministic
	if model.ComponentDefinition != nil {
		MakeComponentDeterminstic(model.ComponentDefinition)
	}
	if model.AssessmentResults != nil {
		MakeAssessmentResultsDeterministic(model.AssessmentResults)
	}

	b, err := ConvertOSCALToBytes(model, filepath.Ext(filePath))
	if err != nil {
		return fmt.Errorf("error converting OSCAL model to bytes: %v", err)
	}

	err = files.WriteOutput(b, filePath)
	if err != nil {
		return err
	}

	return nil

}

func MergeOscalModels(existingModel *oscalTypes_1_1_2.OscalModels, newModel *oscalTypes_1_1_2.OscalModels, modelType string) (*oscalTypes_1_1_2.OscalModels, error) {
	var err error
	// Now to check each model type - currently only component definition and assessment-results apply

	// Component definition
	if modelType == "component" {

		if existingModel.ComponentDefinition == nil && newModel.ComponentDefinition != nil {
			return newModel, nil
		}

		merged, err := MergeComponentDefinitions(existingModel.ComponentDefinition, newModel.ComponentDefinition)
		if err != nil {
			return nil, err
		}
		// Re-assign after processing errors
		existingModel.ComponentDefinition = merged
	}

	// Assessment Results
	if modelType == "assessment-results" {

		if existingModel.AssessmentResults == nil && newModel.AssessmentResults != nil {
			return newModel, nil
		}

		merged, err := MergeAssessmentResults(existingModel.AssessmentResults, newModel.AssessmentResults)
		if err != nil {
			return existingModel, err
		}
		// Re-assign after processing errors
		existingModel.AssessmentResults = merged
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

// ValidOSCALModelAtPath takes a path and returns a bool indicating if the model exists/is valid
// bool = T/F that oscal model exists, error = if not nil OSCAL model is invalid
func ValidOSCALModelAtPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return true, err
	}

	_, err = NewOscalModel(data)
	if err != nil {
		return true, err
	}

	return true, nil
}

// InjectIntoOSCALModel takes a model target and a map[string]interface{} of values to inject into the model
func InjectIntoOSCALModel(target *oscalTypes_1_1_2.OscalModels, values map[string]interface{}, path string) (*oscalTypes_1_1_2.OscalModels, error) {
	// If the target is nil, return an error
	if target == nil {
		return nil, fmt.Errorf("target model is nil")
	}

	// Convert target to a map
	modelMap, err := convertOscalModelToMap(*target)
	if err != nil {
		return nil, err
	}

	// Inject the values into the map at the path
	newModelMap, err := inject.InjectMapData(modelMap, values, path)
	if err != nil {
		return nil, err
	}

	// Convert the new model map back to an OSCAL model
	newModel, err := convertMapToOscalModel(newModelMap)
	if err != nil {
		return nil, err
	}

	return newModel, nil
}

// ConvertOSCALToBytes returns a byte slice representation of an OSCAL model
func ConvertOSCALToBytes(model *oscalTypes_1_1_2.OscalModels, fileExt string) ([]byte, error) {
	var b bytes.Buffer

	if fileExt == ".json" {
		jsonEncoder := json.NewEncoder(&b)
		jsonEncoder.SetIndent("", "  ")
		err := jsonEncoder.Encode(model)
		if err != nil {
			return nil, err
		}
	} else {
		yamlEncoder := yamlV3.NewEncoder(&b)
		yamlEncoder.SetIndent(2)
		err := yamlEncoder.Encode(model)
		if err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

// convertOscalModelToMap converts an OSCAL model to a map[string]interface{}
func convertOscalModelToMap(model oscalTypes_1_1_2.OscalModels) (map[string]interface{}, error) {
	var modelMap map[string]interface{}
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(modelBytes, &modelMap)
	if err != nil {
		return nil, err
	}

	return modelMap, nil
}

// convertMapToOscalModel converts a map[string]interface{} to an OSCAL model
func convertMapToOscalModel(modelMap map[string]interface{}) (*oscalTypes_1_1_2.OscalModels, error) {
	var model oscalTypes_1_1_2.OscalModels
	modelBytes, err := json.Marshal(modelMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(modelBytes, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
