package oscal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/types"
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

func GenerateAssessmentResults(report *types.ReportObject) (oscalTypes.AssessmentResults, error) {
	var assessmentResults oscalTypes.AssessmentResults

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
					sharedUuid := uuid.NewUUID()
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
					tempObservations = append(tempObservations, observation)
				}

				if _, ok := controlMap[implementedRequirement.ControlId]; ok {
					continue
				} else {
					controlMap[implementedRequirement.ControlId] = true
				}
				// TODO: Need to add in the control implementation UUID
				finding := oscalTypes.Finding{
					UUID:        uuid.NewUUID(),
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
