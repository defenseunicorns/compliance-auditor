package validationstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

type ValidationStore struct {
	backMatterMap  map[string]string
	validationMap  map[string]*types.LulaValidation
	observationMap map[string]*oscalTypes.Observation
}

// NewValidationStore creates a new validation store
func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		backMatterMap:  make(map[string]string),
		validationMap:  make(map[string]*types.LulaValidation),
		observationMap: make(map[string]*oscalTypes.Observation),
	}
}

// NewValidationStoreFromBackMatter creates a new validation store from a back matter
func NewValidationStoreFromBackMatter(backMatter oscalTypes.BackMatter) *ValidationStore {
	return &ValidationStore{
		backMatterMap:  oscal.BackMatterToMap(backMatter),
		validationMap:  make(map[string]*types.LulaValidation),
		observationMap: make(map[string]*oscalTypes.Observation),
	}
}

// Number of validations in the store
func (v *ValidationStore) Count() int {
	return len(v.validationMap)
}

// AddValidation adds a validation to the store
func (v *ValidationStore) AddValidation(validation *common.Validation) (id string, err error) {
	if validation.Metadata == nil {
		validation.Metadata = &common.Metadata{}
	}

	if validation.Metadata.UUID == "" {
		validation.Metadata.UUID = uuid.NewUUID()
	}

	lulaValidation, err := validation.ToLulaValidation(validation.Metadata.UUID)
	v.validationMap[validation.Metadata.UUID] = &lulaValidation

	if err != nil {
		return "", err
	}

	return validation.Metadata.UUID, nil
}

// AddLulaValidation adds a LulaValidation to the store
func (v *ValidationStore) AddLulaValidation(validation *types.LulaValidation, id string) {
	trimmedId := common.TrimIdPrefix(id)
	v.validationMap[trimmedId] = validation
}

// GetLulaValidation gets the LulaValidation from the store
func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	trimmedId := common.TrimIdPrefix(id)

	if validation, ok := v.validationMap[trimmedId]; ok {
		return validation, nil
	}

	if validationString, ok := v.backMatterMap[trimmedId]; ok {
		lulaValidation, err := common.ValidationFromString(validationString, trimmedId)
		if err != nil {
			return &lulaValidation, err
		}
		v.validationMap[trimmedId] = &lulaValidation
		return &lulaValidation, nil
	}

	return validation, fmt.Errorf("validation #%s not found", trimmedId)
}

// DryRun checks if the validations are performing execution actions
func (v *ValidationStore) DryRun() (executable bool, msg string) {
	executableValidations := make([]string, 0)
	for k, val := range v.validationMap {
		if val != nil && val.Domain != nil {
			if (*val.Domain).IsExecutable() {
				executableValidations = append(executableValidations, k)
			}
		}
	}
	if len(executableValidations) > 0 {
		return true, fmt.Sprintf("The following validations are executable: %v", executableValidations)
	}
	return false, "No validation is executable"
}

// RunValidations runs the validations in the store
func (v *ValidationStore) RunValidations(ctx context.Context, confirmExecution, saveResources bool, outputsDir string) []oscalTypes.Observation {
	observations := make([]oscalTypes.Observation, 0, len(v.validationMap))

	for k, val := range v.validationMap {
		if val != nil {
			// Create observation for each non-nil validation
			completedText := "evaluated"
			spinnerMessage := fmt.Sprintf("Running validation %s", k)
			spinner := message.NewProgressSpinner("%s", spinnerMessage)
			defer spinner.Stop()
			err := val.Validate(ctx, types.ExecutionAllowed(confirmExecution))
			if err != nil {
				message.Debugf("Error running validation %s: %v", k, err)
				// Update validation with failed results
				val.Result.State = "not-satisfied"
				val.Result.Observations = map[string]string{
					"Error running validation": err.Error(),
				}
				completedText = "NOT evaluated"
			}

			// Update individual result state
			if val.Result.Passing > 0 && val.Result.Failing <= 0 {
				val.Result.State = "satisfied"
			} else {
				val.Result.State = "not-satisfied"
			}

			// Add the observation to the observation map
			var remarks string
			if len(val.Result.Observations) > 0 {
				for k, v := range val.Result.Observations {
					remarks += fmt.Sprintf("%s: %s\n", k, v)
				}
			}

			// Save Resources if specified
			var resourceHref string
			if saveResources {
				resourceUuid := uuid.NewUUID()
				// Create a remote resource file -> create directory 'resources' in the assessment-results directory -> create file with UUID as name
				filename := fmt.Sprintf("%s.json", resourceUuid)
				resourceFile := filepath.Join(outputsDir, "resources", filename)
				err := os.MkdirAll(filepath.Dir(resourceFile), os.ModePerm) // #nosec G301
				if err != nil {
					message.Debugf("Error creating directory for remote resource: %v", err)
				}
				jsonData := val.GetDomainResourcesAsJSON()
				err = files.WriteOutput(jsonData, resourceFile)
				if err != nil {
					message.Debugf("Error writing remote resource file: %v", err)
				}
				resourceHref = fmt.Sprintf("file://./resources/%s", filename)
			}

			// Create an observation
			relevantEvidence := &[]oscalTypes.RelevantEvidence{
				{
					Description: fmt.Sprintf("Result: %s\n", val.Result.State),
					Remarks:     remarks,
				},
			}
			observation := oscal.CreateObservation("TEST", relevantEvidence, val, resourceHref, "[TEST]: %s - %s\n", k, val.Name)
			v.observationMap[k] = &observation
			observations = append(observations, observation)

			spinner.Successf("%s -> %s -> %s", spinnerMessage, completedText, val.Result.State)
		}
	}
	return observations
}

// GetObservation returns the observation with the given ID as well as pass status
func (v *ValidationStore) GetRelatedObservation(id string) (oscalTypes.RelatedObservation, bool) {
	trimmedId := common.TrimIdPrefix(id)
	observation, ok := v.observationMap[trimmedId]
	if !ok {
		return oscalTypes.RelatedObservation{}, false
	}
	pass := false

	// check all descriptions in relevant evidence are satisfied
	if observation.RelevantEvidence != nil {
		for _, e := range *observation.RelevantEvidence {
			if e.Description == "Result: satisfied\n" {
				pass = true
			} else { // if any are not satisfied, return false
				pass = false
				break
			}
		}
	}

	return oscalTypes.RelatedObservation{
		ObservationUuid: observation.UUID,
	}, pass
}

// RunTests executes any tests defined on the validations in the validation store
func (v *ValidationStore) RunTests(ctx context.Context) map[string]types.LulaValidationTestReport {
	testReportMap := make(map[string]types.LulaValidationTestReport)

	for uuid, validation := range v.validationMap {
		// TODO: should test results be saved, e.g., if printResources is true?
		testReport, err := validation.RunTests(ctx, false)
		if err != nil {
			testReportMap[uuid] = types.LulaValidationTestReport{
				Name: validation.Name,
				TestResults: []*types.LulaValidationTestResult{
					{
						TestName: "Error running validation",
						Pass:     false,
						Result:   err.Error(),
					},
				},
			}
			continue
		}

		// If no tests are defined, return an empty report for the validation
		if testReport == nil {
			testReportMap[uuid] = types.LulaValidationTestReport{
				Name: validation.Name,
			}
			continue
		}

		testReportMap[uuid] = *testReport
	}
	return testReportMap
}
