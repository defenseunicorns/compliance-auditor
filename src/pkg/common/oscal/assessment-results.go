package oscal

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common/result"
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

	// Here we are going to add the threshold property by default
	props := []oscalTypes_1_1_2.Property{
		{
			Ns:    "https://docs.lula.dev/ns",
			Name:  "threshold",
			Value: "false",
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
	// This is used during evaluate to update the threshold prop automatically
	if original.UUID == latest.UUID {
		return latest, nil
	}

	original.Results = append(original.Results, latest.Results...)

	slices.SortFunc(original.Results, func(a, b oscalTypes_1_1_2.Result) int { return b.Start.Compare(a.Start) })
	// Update pertinent information
	original.Metadata.LastModified = time.Now()
	original.UUID = uuid.NewUUID()

	return original, nil
}

// IdentifyResults produces a map containing the threshold result and a result used for comparison
func IdentifyResults(assessmentMap map[string]*oscalTypes_1_1_2.AssessmentResults) (map[string]*oscalTypes_1_1_2.Result, error) {
	resultMap := make(map[string]*oscalTypes_1_1_2.Result)

	thresholds, sortedResults := findAndSortResults(assessmentMap)

	// Handle no results found in the assessment-results
	if len(sortedResults) == 0 {
		return nil, fmt.Errorf("less than 2 results found - no comparison possible")
	}

	// Handle single result found in the assessment-results
	if len(sortedResults) == 1 {
		// Only one result found - set latest and return
		resultMap["threshold"] = sortedResults[len(sortedResults)-1]
		return resultMap, fmt.Errorf("less than 2 results found - no comparison possible")
	}

	if len(thresholds) == 0 {
		// No thresholds identified but we have > 1 results - compare the preceding (threshold) against the latest
		resultMap["threshold"] = sortedResults[len(sortedResults)-2]
		resultMap["latest"] = sortedResults[len(sortedResults)-1]

		return resultMap, nil
	} else if len(thresholds) > 1 {
		// More than one threshold - likely the case with multiple assessment-results artifacts
		resultMap["threshold"] = thresholds[len(thresholds)-1]
		resultMap["latest"] = sortedResults[len(sortedResults)-1]

		if resultMap["threshold"] == resultMap["latest"] {
			// if threshold is latest here && we have > 1 threshold - make the threshold the older threshold
			resultMap["threshold"] = thresholds[len(thresholds)-2]
		}

		// Consider changing the namespace value to "false" here - only written if the command logic completes
		for _, result := range thresholds {
			UpdateProps("threshold", "https://docs.lula.dev/ns", "false", result.Props)
		}

		return resultMap, nil

	} else {
		// Otherwise we have a single threshold and we compare that against the latest result
		resultMap["threshold"] = thresholds[len(thresholds)-1]
		resultMap["latest"] = sortedResults[len(sortedResults)-1]

		if resultMap["threshold"] == resultMap["latest"] {
			return nil, fmt.Errorf("latest threshold is the latest result - no comparison possible")
		}

		return resultMap, nil
	}
}

func EvaluateResults(thresholdResult *oscalTypes_1_1_2.Result, newResult *oscalTypes_1_1_2.Result) (bool, map[string]result.ResultComparisonMap, error) {
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

func MakeAssessmentResultsDeterministic(assessment *oscalTypes_1_1_2.AssessmentResults) {

	// Sort Results
	slices.SortFunc(assessment.Results, func(a, b oscalTypes_1_1_2.Result) int { return b.Start.Compare(a.Start) })

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
			slices.SortFunc(observations, func(a, b oscalTypes_1_1_2.Observation) int { return a.Collected.Compare(b.Collected) })
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

// findAndSortResults takes a map of results and returns a list of thresholds and a sorted list of results in order of time
func findAndSortResults(resultMap map[string]*oscalTypes_1_1_2.AssessmentResults) ([]*oscalTypes_1_1_2.Result, []*oscalTypes_1_1_2.Result) {

	thresholds := make([]*oscalTypes_1_1_2.Result, 0)
	sortedResults := make([]*oscalTypes_1_1_2.Result, 0)

	for _, assessment := range resultMap {
		for _, result := range assessment.Results {
			if result.Props != nil {
				for _, prop := range *result.Props {
					if prop.Name == "threshold" && prop.Value == "true" {
						thresholds = append(thresholds, &result)
					}
				}
			}
			// Store all results in a non-sorted list
			sortedResults = append(sortedResults, &result)
		}
	}

	// Sort the results by start time
	slices.SortFunc(sortedResults, func(a, b *oscalTypes_1_1_2.Result) int { return a.Start.Compare(b.Start) })
	slices.SortFunc(thresholds, func(a, b *oscalTypes_1_1_2.Result) int { return a.Start.Compare(b.Start) })

	return thresholds, sortedResults
}

// Helper function to create observation
func CreateObservation(method string, relevantEvidence *[]oscalTypes_1_1_2.RelevantEvidence, descriptionPattern string, descriptionArgs ...any) oscalTypes_1_1_2.Observation {
	rfc3339Time := time.Now()
	uuid := uuid.NewUUID()
	return oscalTypes_1_1_2.Observation{
		Collected:        rfc3339Time,
		Methods:          []string{method},
		UUID:             uuid,
		Description:      fmt.Sprintf(descriptionPattern, descriptionArgs...),
		RelevantEvidence: relevantEvidence,
	}
}
