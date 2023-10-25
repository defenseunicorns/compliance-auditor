package oscal

import (
	"fmt"
	"time"

	assessmentResultsTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1/assessment-results"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/google/uuid"
)

const OSCAL_VERSION = "1.1.1"

func GenerateAssessmentResult(report *types.ReportObject) (assessmentResultsTypes.OscalAssessmentResultsModel, error) {
	var assessmentResults assessmentResultsTypes.OscalAssessmentResultsModel

	// Single time used for all time related fields
	rfc3339Time := time.Now().Format(time.RFC3339)

	// Create placeholders for data required in objects
	controlMap := make(map[string]bool)
	controlList := make([]assessmentResultsTypes.SelectControlById, 0)
	findings := make([]assessmentResultsTypes.Finding, 0)

	// Build the controlMap and Findings array
	for _, component := range report.Components {
		for _, controlImplementation := range component.ControlImplementations {
			for _, implementedRequirement := range controlImplementation.ImplementedReqs {

				if _, ok := controlMap[implementedRequirement.ControlId]; ok {
					continue
				} else {
					controlMap[implementedRequirement.ControlId] = true
				}
				// TODO: Need to add in the control implementation UUID
				finding := assessmentResultsTypes.Finding{
					UUID:        uuid.NewString(),
					Title:       fmt.Sprintf("Validation Result - Component:%s / Control Implementation: %s / Control:  %s", component.UUID, controlImplementation.UUID, implementedRequirement.ControlId),
					Description: implementedRequirement.Description,
					Target: assessmentResultsTypes.FindingTarget{
						Status: assessmentResultsTypes.Status{
							State: implementedRequirement.Status,
						},
						TargetId: implementedRequirement.ControlId,
						Type:     "objective-id",
					},
				}
				findings = append(findings, finding)
			}
		}
	}

	// Convert control map to slice of SelectControlById
	for controlId := range controlMap {
		control := assessmentResultsTypes.SelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
	}

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.AssessmentResults.UUID = uuid.NewString()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	assessmentResults.AssessmentResults.Metadata = assessmentResultsTypes.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Lula Metadata Remarks",
		Published:    rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Create results object
	assessmentResults.AssessmentResults.Results = []assessmentResultsTypes.Result{
		{
			UUID:        uuid.NewString(),
			Title:       "Lula Result Title",
			Start:       rfc3339Time,
			Description: "Lula Result Description",
			ReviewedControls: assessmentResultsTypes.ReviewedControls{
				Description: "Lula Control Description",
				Remarks:     "Lula Control Remarks",
				ControlSelections: []assessmentResultsTypes.AssessedControls{
					{
						Description:     "Lula Assessed Controls Description",
						IncludeControls: controlList,
					},
				},
			},
			Findings: findings,
		},
	}

	return assessmentResults, nil
}
