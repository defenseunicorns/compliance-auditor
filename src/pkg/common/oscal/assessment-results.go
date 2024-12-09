package oscal

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/result"
	"github.com/defenseunicorns/lula/src/types"
)

const OSCAL_VERSION = "1.1.2"

type EvalResult struct {
	Threshold *oscalTypes.Result
	Results   []*oscalTypes.Result
	Latest    *oscalTypes.Result
}

// NewAssessmentResults creates a new assessment results object from the given data.
func NewAssessmentResults(data []byte) (*oscalTypes.AssessmentResults, error) {
	var oscalModels oscalTypes.OscalModels

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

func GenerateAssessmentResults(results []oscalTypes.Result) (*oscalTypes.AssessmentResults, error) {
	var assessmentResults = &oscalTypes.AssessmentResults{}

	// Single time used for all time related fields
	rfc3339Time := time.Now()

	// Always create a new UUID for the assessment results (for now)
	assessmentResults.UUID = uuid.NewUUID()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	assessmentResults.Metadata = oscalTypes.Metadata{
		Title:        "[System Name] Security Assessment Results (SAR)",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Assessment Results generated from Lula",
		Published:    &rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Create results object
	assessmentResults.Results = results

	return assessmentResults, nil
}

func MergeAssessmentResults(original *oscalTypes.AssessmentResults, latest *oscalTypes.AssessmentResults) (*oscalTypes.AssessmentResults, error) {
	// If UUID's are matching - this must be a prop update for threshold
	// This is used during evaluate to update the threshold prop automatically
	if original.UUID == latest.UUID {
		return latest, nil
	}

	original.Results = append(original.Results, latest.Results...)

	slices.SortFunc(original.Results, func(a, b oscalTypes.Result) int { return b.Start.Compare(a.Start) })
	// Update pertinent information
	original.Metadata.LastModified = time.Now()
	original.UUID = uuid.NewUUID()

	return original, nil
}

func EvaluateResults(thresholdResult *oscalTypes.Result, newResult *oscalTypes.Result) (bool, map[string]result.ResultComparisonMap, error) {
	var status bool = true

	if thresholdResult.Findings == nil || newResult.Findings == nil {
		return false, nil, fmt.Errorf("results must contain findings to evaluate")
	}

	// Compare threshold result to new result and vice versa
	comparedToThreshold := result.NewResultComparisonMap(*newResult, *thresholdResult)

	// Group by categories
	categories := []struct {
		name        string
		stateChange result.StateChange
		satisfied   bool
		status      bool
	}{
		{
			name:        "new-satisfied",
			stateChange: result.NEW,
			satisfied:   true,
			status:      true,
		},
		{
			name:        "new-not-satisfied",
			stateChange: result.NEW,
			satisfied:   false,
			status:      true,
		},
		{
			name:        "no-longer-satisfied",
			stateChange: result.SATISFIED_TO_NOT_SATISFIED,
			satisfied:   false,
			status:      false,
		},
		{
			name:        "now-satisfied",
			stateChange: result.NOT_SATISFIED_TO_SATISFIED,
			satisfied:   true,
			status:      true,
		},
		{
			name:        "unchanged-not-satisfied",
			stateChange: result.UNCHANGED,
			satisfied:   false,
			status:      true,
		},
		{
			name:        "unchanged-satisfied",
			stateChange: result.UNCHANGED,
			satisfied:   true,
			status:      true,
		},
		{
			name:        "removed-not-satisfied",
			stateChange: result.REMOVED,
			satisfied:   false,
			status:      false,
		},
		{
			name:        "removed-satisfied",
			stateChange: result.REMOVED,
			satisfied:   true,
			status:      false,
		},
	}

	categorizedResultComparisons := make(map[string]result.ResultComparisonMap)
	for _, c := range categories {
		results := result.GetResultComparisonMap(comparedToThreshold, c.stateChange, c.satisfied)
		categorizedResultComparisons[c.name] = results
		if len(results) > 0 && !c.status {
			status = false
		}
	}

	return status, categorizedResultComparisons, nil
}

func MakeAssessmentResultsDeterministic(assessment *oscalTypes.AssessmentResults) {

	// Sort Results
	slices.SortFunc(assessment.Results, func(a, b oscalTypes.Result) int { return b.Start.Compare(a.Start) })

	for _, result := range assessment.Results {
		// sort findings by target id
		if result.Findings != nil {
			findings := *result.Findings
			sort.Slice(findings, func(i, j int) bool {
				return findings[i].Target.TargetId < findings[j].Target.TargetId
			})
			result.Findings = &findings
		}
		// sort observations by collected time
		if result.Observations != nil {
			observations := *result.Observations
			slices.SortFunc(observations, func(a, b oscalTypes.Observation) int { return a.Collected.Compare(b.Collected) })
			result.Observations = &observations
		}

		// Sort the include-controls in the control selections
		controlSelections := result.ReviewedControls.ControlSelections
		for _, selection := range controlSelections {
			if selection.IncludeControls != nil {
				controls := *selection.IncludeControls
				sort.Slice(controls, func(i, j int) bool {
					return controls[i].ControlId < controls[j].ControlId
				})
				selection.IncludeControls = &controls
			}
		}
	}

	// sort backmatter
	if assessment.BackMatter != nil {
		backmatter := *assessment.BackMatter
		if backmatter.Resources != nil {
			resources := *backmatter.Resources
			sort.Slice(resources, func(i, j int) bool {
				return resources[i].Title < resources[j].Title
			})
			backmatter.Resources = &resources
		}
		assessment.BackMatter = &backmatter
	}

}

// filterResults consumes many assessment-results objects and builds out a map of EvalResults filtered by target
// this function looks at the target prop as the key in the map
func FilterResults(resultMap map[string]*oscalTypes.AssessmentResults) map[string]EvalResult {
	evalResultMap := make(map[string]EvalResult)

	for _, assessment := range resultMap {
		if assessment == nil {
			continue
		}
		for _, result := range assessment.Results {
			if result.Props != nil {
				var target string
				hasTarget, targetValue := GetProp("target", LULA_NAMESPACE, result.Props)
				hasThreshold, thresholdValue := GetProp("threshold", LULA_NAMESPACE, result.Props)

				if hasTarget {
					// existing target prop
					target = targetValue
				} else {
					// non-existent target prop
					target = "default"
				}

				var evalResult EvalResult
				// target identified
				if tmpResult, ok := evalResultMap[target]; ok {
					// EvalResult Exists - append
					tmpResult.Results = append(tmpResult.Results, &result)
					evalResult = tmpResult
				} else {
					// EvalResult Does Not Exist - create
					results := make([]*oscalTypes.Result, 0)
					results = append(results, &result)
					tmpResult = EvalResult{
						Results: results,
					}
					evalResult = tmpResult
				}

				if hasThreshold && thresholdValue == "true" {
					if evalResult.Threshold == nil {
						evalResult.Threshold = &result
					} else {
						// If threshold exists and this is a newer threshold
						if result.Start.Compare(evalResult.Threshold.Start) > 0 {
							UpdateProps("threshold", LULA_NAMESPACE, "false", evalResult.Threshold.Props)
							evalResult.Threshold = &result
						}
					}
				}
				evalResultMap[target] = evalResult
			}
		}
	}
	// Now that all results are processed - iterate through each EvalResult, sort and assign latest/threshold
	for key, evalResult := range evalResultMap {
		if len(evalResult.Results) > 0 {
			slices.SortFunc(evalResult.Results, func(a, b *oscalTypes.Result) int { return a.Start.Compare(b.Start) })
			evalResult.Latest = evalResult.Results[len(evalResult.Results)-1]
			if evalResult.Threshold == nil && len(evalResult.Results) > 1 {
				// length of results > 1 and no established threshold - set threshold to the preceding result of latest
				evalResult.Threshold = evalResult.Results[len(evalResult.Results)-2]
			}
			evalResultMap[key] = evalResult
		}

	}

	return evalResultMap
}

// Helper function to create observation
func CreateObservation(method string, relevantEvidence *[]oscalTypes.RelevantEvidence, validation *types.LulaValidation, resourcesHref string, descriptionPattern string, descriptionArgs ...any) oscalTypes.Observation {
	rfc3339Time := time.Now()
	observationUuid := uuid.NewUUID()

	observation := oscalTypes.Observation{
		Collected:        rfc3339Time,
		Methods:          []string{method},
		UUID:             observationUuid,
		Description:      fmt.Sprintf(descriptionPattern, descriptionArgs...),
		RelevantEvidence: relevantEvidence,
	}
	if validation != nil {
		observation.Props = &[]oscalTypes.Property{
			{
				Name:  "validation",
				Ns:    "https://docs.lula.dev/oscal/ns",
				Value: common.AddIdPrefix(validation.UUID),
			},
		}
	}
	if resourcesHref != "" {
		observation.Links = &[]oscalTypes.Link{
			{
				Href: resourcesHref,
				Rel:  "lula.resources",
			},
		}
	}
	return observation
}

// Creates a result from findings and observations
func CreateResult(findingMap map[string]oscalTypes.Finding, observations []oscalTypes.Observation) (oscalTypes.Result, error) {

	// Single time used for all time related fields
	rfc3339Time := time.Now()
	controlList := make([]oscalTypes.AssessedControlsSelectControlById, 0)
	findings := make([]oscalTypes.Finding, 0)

	// Convert control map to slice of SelectControlById
	for controlId, finding := range findingMap {
		control := oscalTypes.AssessedControlsSelectControlById{
			ControlId: controlId,
		}
		controlList = append(controlList, control)
		findings = append(findings, finding)
	}

	props := []oscalTypes.Property{
		{
			Ns:    LULA_NAMESPACE,
			Name:  "threshold",
			Value: "false",
		},
	}

	result := oscalTypes.Result{
		UUID:        uuid.NewUUID(),
		Title:       "Lula Validation Result",
		Start:       rfc3339Time,
		Description: "Assessment results for performing Validations with Lula version " + config.CLIVersion,
		Props:       &props,
		ReviewedControls: oscalTypes.ReviewedControls{
			Description: "Controls validated",
			Remarks:     "Validation performed may indicate full or partial satisfaction",
			ControlSelections: []oscalTypes.AssessedControls{
				{
					Description:     "Controls Assessed by Lula",
					IncludeControls: &controlList,
				},
			},
		},
		Findings: &findings,
	}

	// Observations are only present with evidence generated by a validation
	// Lula can operate on oscal without any validations
	if len(observations) > 0 {
		result.Observations = &observations
	}

	return result, nil
}

// GetObservationByUuid returns the observation with the given UUID
func GetObservationByUuid(assessmentResults *oscalTypes.AssessmentResults, observationUuid string) (*oscalTypes.Observation, error) {
	if assessmentResults == nil {
		return nil, fmt.Errorf("assessment results is nil")
	}

	for _, result := range assessmentResults.Results {
		if result.Observations != nil {
			for _, observation := range *result.Observations {
				if observation.UUID == observationUuid {
					return &observation, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("observation with uuid %s not found", observationUuid)
}
