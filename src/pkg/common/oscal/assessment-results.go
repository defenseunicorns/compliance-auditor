package oscal

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"gopkg.in/yaml.v3"
)

const OSCAL_VERSION = "1.1.2"

// NewAssessmentResults creates a new assessment results object from the given data.
func NewAssessmentResults(data []byte) (*oscalTypes_1_1_2.AssessmentResults, error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	err := multiModelValidate(data)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		fmt.Printf("Error marshalling yaml: %s\n", err.Error())
		return nil, err
	}

	return oscalModels.AssessmentResults, nil
}

func GenerateAssessmentResults(findingMap map[string]oscalTypes_1_1_2.Finding, observations []oscalTypes_1_1_2.Observation) (*oscalTypes_1_1_2.AssessmentResults, error) {
	var assessmentResults = &oscalTypes_1_1_2.AssessmentResults{}

	// Single time used for all time related fields
	rfc3339Time := time.Now()
	controlList := make([]oscalTypes_1_1_2.AssessedControlsSelectControlById, 0)
	findings := make([]oscalTypes_1_1_2.Finding, 0)

	// Convert control map to slice of SelectControlById
	for controlId, finding := range findingMap {
		control := oscalTypes_1_1_2.AssessedControlsSelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
		findings = append(findings, finding)
	}

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.UUID = uuid.NewUUID()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	assessmentResults.Metadata = oscalTypes_1_1_2.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Assessment Results generated from Lula",
		Published:    &rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Here we are going to add the threshold property
	props := []oscalTypes_1_1_2.Property{
		{
			Ns:    "https://docs.lula.dev/ns",
			Name:  "threshold",
			Value: "true",
		},
	}

	// Create results object
	assessmentResults.Results = []oscalTypes_1_1_2.Result{
		{
			UUID:        uuid.NewUUID(),
			Title:       "Lula Validation Result",
			Start:       rfc3339Time,
			Description: "Assessment results for performing Validations with Lula version " + config.CLIVersion,
			Props:       &props,
			ReviewedControls: oscalTypes_1_1_2.ReviewedControls{
				Description: "Controls validated",
				Remarks:     "Validation performed may indicate full or partial satisfaction",
				ControlSelections: []oscalTypes_1_1_2.AssessedControls{
					{
						Description:     "Controls Assessed by Lula",
						IncludeControls: &controlList,
					},
				},
			},
			Findings:     &findings,
			Observations: &observations,
		},
	}

	return assessmentResults, nil
}

func MergeAssessmentResults(original *oscalTypes_1_1_2.AssessmentResults, latest *oscalTypes_1_1_2.AssessmentResults) (*oscalTypes_1_1_2.AssessmentResults, error) {

	// If UUID's are matching - this must be a prop update for threshold
	// We should be able to return the latest results
	// This is used during evaluate to update the threshold prop automatically
	if original.UUID == latest.UUID {
		// Consider that this is a potential modification and this might be a good location to generate a new UUID
		return latest, nil
	}

	// Validate only ever creates one result
	// Assumed that there is always an original threshold
	result := latest.Results[0]
	for index, prop := range *result.Props {
		if prop.Name == "threshold" {
			prop.Value = "false"
			// Better way to update the prop?
			(*result.Props)[index] = prop
		}
	}

	results := make([]oscalTypes_1_1_2.Result, 0)
	// append newest to oldest results
	results = append(results, result)
	results = append(results, original.Results...)
	original.Results = results

	// Update pertinent information
	original.Metadata.LastModified = time.Now()
	original.UUID = uuid.NewUUID()

	return original, nil
}

func GenerateFindingsMap(findings []oscalTypes_1_1_2.Finding) map[string]oscalTypes_1_1_2.Finding {
	findingsMap := make(map[string]oscalTypes_1_1_2.Finding)
	for _, finding := range findings {
		findingsMap[finding.Target.TargetId] = finding
	}
	return findingsMap
}
