package oscal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/google/uuid"
)

const OSCAL_VERSION = "1.1.1"

func GenerateAssessmentResults(report *types.ReportObject) (oscalTypes.OscalModels, error) {
	var assessmentResults oscalTypes.OscalModels

	// Single time used for all time related fields
	rfc3339Time := time.Now().Format(time.RFC3339)

	// Create placeholders for data required in objects
	controlMap := make(map[string]bool)
	controlList := make([]oscalTypes.SelectControlById, 0)
	findings := make([]oscalTypes.Finding, 0)
	observations := make([]oscalTypes.Observation, 0)

	// Build the controlMap and Findings array
	for _, component := range report.Components {
		for _, controlImplementation := range component.ControlImplementations {
			for _, implementedRequirement := range controlImplementation.ImplementedReqs {
				tempObservations := make([]oscalTypes.Observation, 0)
				relatedObservations := make([]oscalTypes.RelatedObservation, 0)
				// For each result - there may be many observations
				for _, result := range implementedRequirement.Results {

					sharedUuid := uuid.NewString()
					observation := oscalTypes.Observation{
						Collected:   rfc3339Time,
						Description: fmt.Sprintf("[TEST] %s - %s\n", implementedRequirement.ControlId, result.UUID),
						Methods:     []string{"TEST"},
						UUID:        sharedUuid,
						RelevantEvidence: []oscalTypes.RelevantEvidence{
							{
								Description: fmt.Sprintf("Result: %s - Passing Resources: %s - Failing Resources %s\n", result.State, strconv.Itoa(result.Passing), strconv.Itoa(result.Failing)),
							},
						},
					}

					relatedObservation := oscalTypes.RelatedObservation{
						ObservationUuid: sharedUuid,
					}

					relatedObservations = append(relatedObservations, relatedObservation)
					tempObservations = append(observations, observation)
				}

				if _, ok := controlMap[implementedRequirement.ControlId]; ok {
					continue
				} else {
					controlMap[implementedRequirement.ControlId] = true
				}
				// TODO: Need to add in the control implementation UUID
				finding := oscalTypes.Finding{
					UUID:        uuid.NewString(),
					Title:       fmt.Sprintf("Validation Result - Component:%s / Control Implementation: %s / Control:  %s", component.UUID, controlImplementation.UUID, implementedRequirement.ControlId),
					Description: implementedRequirement.Description,
					Target: oscalTypes.FindingTarget{
						Status: oscalTypes.Status{
							State: implementedRequirement.State,
						},
						TargetId: implementedRequirement.ControlId,
						Type:     "objective-id",
					},
					RelatedObservations: relatedObservations,
				}
				findings = append(findings, finding)
				observations = append(observations, tempObservations...)
			}
		}
	}

	// Convert control map to slice of SelectControlById
	for controlId := range controlMap {
		control := oscalTypes.SelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
	}

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.AssessmentResults.UUID = uuid.NewString()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	assessmentResults.AssessmentResults.Metadata = oscalTypes.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Lula Metadata Remarks",
		Published:    rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Create results object
	assessmentResults.AssessmentResults.Results = []oscalTypes.Result{
		{
			UUID:        uuid.NewString(),
			Title:       "Lula Result Title",
			Start:       rfc3339Time,
			Description: "Lula Result Description",
			ReviewedControls: oscalTypes.ReviewedControls{
				Description: "Lula Control Description",
				Remarks:     "Lula Control Remarks",
				ControlSelections: []oscalTypes.AssessedControls{
					{
						Description:     "Lula Assessed Controls Description",
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
