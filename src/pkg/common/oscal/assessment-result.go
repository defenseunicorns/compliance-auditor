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

	rfc3339Time := time.Now().Format(time.RFC3339)

	controlMap := make(map[string]bool)
	controlList := make([]assessmentResultsTypes.SelectControlById, 0)
	findings := make([]assessmentResultsTypes.Finding, 0)

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
					UUID:        implementedRequirement.UUID,
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

	for controlId := range controlMap {
		control := assessmentResultsTypes.SelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
	}

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.AssessmentResults.UUID = uuid.NewString()

	// Create metadata object with requires fields and a few extras
	assessmentResults.AssessmentResults.Metadata = assessmentResultsTypes.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1", // TODO: Set this to the version of the SAR
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
