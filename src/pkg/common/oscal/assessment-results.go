package oscal

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/config"
	"gopkg.in/yaml.v3"
)

const OSCAL_VERSION = "1.1.1"

func NewAssessmentResults(data []byte) (oscalTypes.AssessmentResults, error) {
	var assessmentResults oscalTypes.AssessmentResults

	err := yaml.Unmarshal(data, &assessmentResults)
	if err != nil {
		fmt.Printf("Error marshalling yaml: %s\n", err.Error())
		return oscalTypes.AssessmentResults{}, err
	}

	return assessmentResults, nil
}

func GenerateAssessmentResults(findingMap map[string]oscalTypes.Finding, observations []oscalTypes.Observation) (oscalTypes.AssessmentResults, error) {
	var assessmentResults oscalTypes.AssessmentResults

	// Single time used for all time related fields
	rfc3339Time := time.Now().Format(time.RFC3339)
	controlList := make([]oscalTypes.SelectControlById, 0)
	findings := make([]oscalTypes.Finding, 0)

	// Convert control map to slice of SelectControlById
	for controlId, finding := range findingMap {
		control := oscalTypes.SelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
		findings = append(findings, finding)
	}

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.UUID = uuid.NewUUID()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	assessmentResults.Metadata = oscalTypes.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Assessment Results generated from Lula",
		Published:    rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Create results object
	assessmentResults.Results = []oscalTypes.Result{
		{
			UUID:        uuid.NewUUID(),
			Title:       "Lula Validation Result",
			Start:       rfc3339Time,
			Description: "Assessment results for performing Validations with Lula version " + config.CLIVersion,
			ReviewedControls: oscalTypes.ReviewedControls{
				Description: "Controls validated",
				Remarks:     "Validation performed may indicate full or partial satisfaction",
				ControlSelections: []oscalTypes.AssessedControls{
					{
						Description:     "Controls Assessed by Lula",
						IncludeControls: controlList,
					},
				},
			},
			Findings:     findings,
			Observations: observations,
		},
	}

	return assessmentResults, nil
}

func GenerateFindingsMap(findings []oscalTypes.Finding) map[string]oscalTypes.Finding {
	findingsMap := make(map[string]oscalTypes.Finding)
	for _, finding := range findings {
		findingsMap[finding.Target.TargetId] = finding
	}
	return findingsMap
}
