package oscal

import (
	"fmt"
	"slices"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"gopkg.in/yaml.v3"
)

const OSCAL_VERSION = "1.1.2"

type ResultPair struct {
	ControlId       string
	oldFinding      *oscalTypes_1_1_2.Finding
	newFinding      *oscalTypes_1_1_2.Finding
	oldObservations []*oscalTypes_1_1_2.Observation
	newObservations []*oscalTypes_1_1_2.Observation
}

// newResultPair creates a new result pair object
func newResultPair(controlId string, oldFinding *oscalTypes_1_1_2.Finding, newFinding *oscalTypes_1_1_2.Finding, oldObservations []*oscalTypes_1_1_2.Observation, newObservations []*oscalTypes_1_1_2.Observation) ResultPair {
	return ResultPair{
		ControlId:       controlId,
		oldFinding:      oldFinding,
		newFinding:      newFinding,
		oldObservations: oldObservations,
		newObservations: newObservations,
	}
}

// GetResultPairDetails returns a string of the result pair observations
func (rp *ResultPair) GetResultPairDetails() string {
	var details string
	if rp.oldFinding != nil {
		details += fmt.Sprintf("Threshold Finding UUID: %s\n", rp.oldFinding.UUID)
		details += fmt.Sprintf("Threshold Observations: %s\n", getObservationsDetails(rp.oldObservations))
	}
	if rp.newFinding != nil {
		details += fmt.Sprintf("New Finding UUID: %s\n", rp.newFinding.UUID)
		details += fmt.Sprintf("New Observations: %s\n", getObservationsDetails(rp.newObservations))
	}
	return details
}

// GetMismatchedObservations returns a string of the mismatched observations
func (rp *ResultPair) GetMismatchedObservations() string {
	var details string
	if rp.oldFinding != nil && rp.newFinding != nil {
		details += fmt.Sprintf("Observations Deltas: %s\n", getObservationsDeltasDetails(rp.oldObservations, rp.newObservations))
	} else if rp.oldFinding != nil {
		details += "No New Findings\n"
		details += fmt.Sprintf("Threshold Observations: %s\n", getObservationsDetails(rp.oldObservations))
	} else if rp.newFinding != nil {
		details += "No Threshold Findings\n"
		details += fmt.Sprintf("New Observations: %s\n", getObservationsDetails(rp.oldObservations))
	}

	return details
}

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

func GenerateFindingsMap(findings []oscalTypes_1_1_2.Finding) map[string]oscalTypes_1_1_2.Finding {
	findingsMap := make(map[string]oscalTypes_1_1_2.Finding)
	for _, finding := range findings {
		findingsMap[finding.Target.TargetId] = finding
	}
	return findingsMap
}

// Helper function to create observations map
func GenerateObservationsMap(observations []oscalTypes_1_1_2.Observation) map[string]*oscalTypes_1_1_2.Observation {
	observationMap := make(map[string]*oscalTypes_1_1_2.Observation)
	for _, observation := range observations {
		observationMap[observation.UUID] = &observation
	}
	return observationMap
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

func EvaluateResults(thresholdResult *oscalTypes_1_1_2.Result, newResult *oscalTypes_1_1_2.Result) (bool, map[string][]ResultPair, error) {
	if thresholdResult.Findings == nil || newResult.Findings == nil {
		return false, nil, fmt.Errorf("results must contain findings to evaluate")
	}

	// Store unique findings for review here
	resultPairs := make(map[string][]ResultPair, 0)
	result := true
	var tempRp ResultPair

	findingMapThreshold := GenerateFindingsMap(*thresholdResult.Findings)
	findingMapNew := GenerateFindingsMap(*newResult.Findings)
	observationMapThreshold := GenerateObservationsMap(*thresholdResult.Observations)
	observationMapNew := GenerateObservationsMap(*newResult.Observations)

	// For a given oldResult - we need to prove that the newResult implements all of the oldResult findings/controls
	// We are explicitly iterating through the findings in order to collect a delta to display

	for targetId, finding := range findingMapThreshold {
		if _, ok := findingMapNew[targetId]; !ok {
			// If the new result does not contain the finding of the old result
			// set result to fail, add finding to the findings map and continue
			result = false
			tempRp = newResultPair(targetId, &finding, nil, getObservations(finding, observationMapThreshold), nil)
			resultPairs["no-longer-satisfied"] = append(resultPairs["no-longer-satisfied"], tempRp)
		} else {
			newFinding := findingMapNew[targetId]
			tempRp = newResultPair(targetId, &finding, &newFinding, getObservations(finding, observationMapThreshold), getObservations(newFinding, observationMapNew))

			// If the finding is present in each map - we need to check if the state has changed from "not-satisfied" to "satisfied"
			if finding.Target.Status.State == "satisfied" {
				// Was previously satisfied - compare state
				if newFinding.Target.Status.State == "not-satisfied" {
					// If the new finding is now not-satisfied - set result to false and add to findings
					result = false
					resultPairs["no-longer-satisfied"] = append(resultPairs["no-longer-satisfied"], tempRp)
				}
			} else {
				// was previously not-satisfied but now is satisfied
				if newFinding.Target.Status.State == "satisfied" {
					// If the new finding is now satisfied - add to new-passing-findings
					resultPairs["new-passing-findings"] = append(resultPairs["new-passing-findings"], tempRp)
				}
			}
			delete(findingMapNew, targetId)
		}
	}

	// All remaining findings in the new map are new findings
	for _, finding := range findingMapNew {
		tempRp = newResultPair(finding.Target.TargetId, nil, &finding, nil, getObservations(finding, observationMapNew))
		if finding.Target.Status.State == "satisfied" {
			resultPairs["new-passing-findings"] = append(resultPairs["new-passing-findings"], tempRp)
		} else {
			resultPairs["new-failing-findings"] = append(resultPairs["new-failing-findings"], tempRp)
		}

	}

	return result, resultPairs, nil
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

// getObservations gets a finding's observations
func getObservations(finding oscalTypes_1_1_2.Finding, observations map[string]*oscalTypes_1_1_2.Observation) []*oscalTypes_1_1_2.Observation {
	var relatedObservations []*oscalTypes_1_1_2.Observation
	if finding.RelatedObservations != nil {
		for _, observation := range *finding.RelatedObservations {
			if _, ok := observations[observation.ObservationUuid]; ok {
				relatedObservations = append(relatedObservations, observations[observation.ObservationUuid])
			}
		}
	}

	return relatedObservations
}

// getObservationsDetails returns a string of the observations details
func getObservationsDetails(observations []*oscalTypes_1_1_2.Observation) string {
	var details string
	for _, observation := range observations {
		details += fmt.Sprintf("\nObservation UUID: %s\n", observation.UUID)
		details += fmt.Sprintf("Observation Description: %s", observation.Description)
		if observation.RelevantEvidence != nil {
			for _, re := range *observation.RelevantEvidence {
				details += fmt.Sprintf("Observation %s", re.Description)
				details += fmt.Sprintf("Observation Evidence: %s", re.Remarks)
			}
		}
	}
	return details
}

// getObservationsDeltas retuns the observations that are different between the two observations arrays
func getObservationsDeltas(aObservations []*oscalTypes_1_1_2.Observation, bObservations []*oscalTypes_1_1_2.Observation) ([]*oscalTypes_1_1_2.Observation, []*oscalTypes_1_1_2.Observation) {
	var different []*oscalTypes_1_1_2.Observation
	var missing []*oscalTypes_1_1_2.Observation
	var foundInB bool

	for _, aOb := range aObservations {
		foundInB = false
		for _, bOb := range bObservations {
			// Get different A observations
			if aOb.Description == bOb.Description {
				if !isSameResult(aOb, bOb) {
					different = append(different, aOb)
				}
				foundInB = true
				break
			}
		}
		// Get missing A observations
		if !foundInB {
			missing = append(missing, aOb)
		}
	}

	return different, missing
}

// compareRelevantEvidence compares the relevant evidence of two observations
func isSameResult(oldObservation *oscalTypes_1_1_2.Observation, newObservation *oscalTypes_1_1_2.Observation) bool {
	oldRelevantEvidence := oldObservation.RelevantEvidence
	newRelevantEvidence := newObservation.RelevantEvidence

	if oldRelevantEvidence != nil && newRelevantEvidence != nil {
		for _, oldRe := range *oldRelevantEvidence {
			for _, newRe := range *newRelevantEvidence {
				if oldRe.Description == newRe.Description {
					return true
				}
			}
		}
	}
	return false
}

// getObservationsDeltasDetails returns a string of the observations details for the delta observations
func getObservationsDeltasDetails(oldObservations []*oscalTypes_1_1_2.Observation, newObservations []*oscalTypes_1_1_2.Observation) string {
	var details string
	oldDifferentFromNew, oldMissingFromNew := getObservationsDeltas(oldObservations, newObservations)
	newDifferentFromOld, newMissingFromOld := getObservationsDeltas(newObservations, oldObservations)

	// Print out the observations that are different between the two observations arrays
	if len(oldDifferentFromNew) > 0 {
		details += "\nThreshold Observations != New:\n"
		details += getObservationsDetails(oldDifferentFromNew)
	}
	if len(oldMissingFromNew) > 0 {
		details += "\nThreshold Observations missing from New:\n"
		details += getObservationsDetails(oldMissingFromNew)
	}
	if len(newDifferentFromOld) > 0 {
		details += "\nNew Observations != Threshold:\n"
		details += getObservationsDetails(newDifferentFromOld)
	}
	if len(newMissingFromOld) > 0 {
		details += "\nNew Observations missing from Threshold:\n"
		details += getObservationsDetails(newMissingFromOld)
	}

	return details
}
