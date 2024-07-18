package result_test

import (
	"fmt"
	"testing"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/result"
)

func createTestResult(findingId, observationId, findingState, observationSatisfaction string) oscalTypes_1_1_2.Result {
	observationUuid := uuid.NewUUID()
	return oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: findingId,
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: findingState,
					},
				},
				RelatedObservations: &[]oscalTypes_1_1_2.RelatedObservation{
					{
						ObservationUuid: observationUuid,
					},
				},
			},
		},
		Observations: &[]oscalTypes_1_1_2.Observation{
			{
				UUID:        observationUuid,
				Description: observationId,
				RelevantEvidence: &[]oscalTypes_1_1_2.RelevantEvidence{
					{
						Description: fmt.Sprintf("Result: %s", observationSatisfaction),
						Remarks:     "Some remarks about this observation",
					},
				},
			},
		},
	}
}

func createTestResultNoObs(findingId, findingState string) oscalTypes_1_1_2.Result {
	return oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: findingId,
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: findingState,
					},
				},
			},
		},
	}
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func TestGetResultComparisonMap(t *testing.T) {
	// Tests both creating a results comparison map and testing getting the right comparisons
	tests := []struct {
		name                 string
		thresholdResult      oscalTypes_1_1_2.Result
		result               oscalTypes_1_1_2.Result
		expectedStateChange  result.StateChange
		expectedSatisfaction bool
		expectedId           string
	}{
		{
			name:                 "Unchanged, satisfied result",
			thresholdResult:      createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			result:               createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			expectedStateChange:  result.UNCHANGED,
			expectedSatisfaction: true,
			expectedId:           "id-1",
		},
		{
			name:                 "Changed, not satisfied to satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "not-satisfied", "not-satisfied"),
			result:               createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			expectedStateChange:  result.NOT_SATISFIED_TO_SATISFIED,
			expectedSatisfaction: true,
			expectedId:           "id-1",
		},
		{
			name:                 "Changed, satisfied to not-satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			result:               createTestResult("id-1", "test-1", "not-satisfied", "not-satisfied"),
			expectedStateChange:  result.SATISFIED_TO_NOT_SATISFIED,
			expectedSatisfaction: false,
			expectedId:           "id-1",
		},
		{
			name:                 "Removed finding, satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			result:               createTestResult("id-2", "test-2", "satisfied", "satisfied"),
			expectedStateChange:  result.REMOVED,
			expectedSatisfaction: false, // this is not-satisfied because it was removed, even though it was originally satisfied
			expectedId:           "id-1",
		},
		{
			name:                 "New finding, satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			result:               createTestResult("id-2", "test-2", "satisfied", "satisfied"),
			expectedStateChange:  result.NEW,
			expectedSatisfaction: true,
			expectedId:           "id-2",
		},
		{
			name:                 "Removed finding, not-satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "not-satisfied", "not-satisfied"),
			result:               createTestResult("id-2", "test-2", "satisfied", "satisfied"),
			expectedStateChange:  result.REMOVED,
			expectedSatisfaction: false,
			expectedId:           "id-1",
		},
		{
			name:                 "New finding, not-satisfied",
			thresholdResult:      createTestResult("id-1", "test-1", "satisfied", "satisfied"),
			result:               createTestResult("id-2", "test-2", "not-satisfied", "not-satisfied"),
			expectedStateChange:  result.NEW,
			expectedSatisfaction: false,
			expectedId:           "id-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultComparisonMap := result.NewResultComparisonMap(tt.result, tt.thresholdResult)
			subSetMap := result.GetResultComparisonMap(resultComparisonMap, tt.expectedStateChange, tt.expectedSatisfaction)

			if len(subSetMap) == 0 {
				t.Error("Expected subset populated, but it's empty")
			}
			if len(subSetMap) != 1 {
				t.Errorf("Expected subset to have 1 element, but has %d", len(subSetMap))
			}
			for id := range subSetMap {
				if id != tt.expectedId {
					t.Errorf("Expected id %s, but got %s", tt.expectedId, id)
				}
			}
		})
	}
}

func TestRefactorObservationsByControls(t *testing.T) {
	// create a bunch of result-comparisons for each ID...
	result1 := createTestResult("id-1", "test-1", "satisfied", "satisfied")
	thresholdResult1 := createTestResult("id-1", "test-1", "satisfied", "satisfied")
	resultComparisonMap1 := result.NewResultComparisonMap(result1, thresholdResult1)

	result2 := createTestResult("id-2", "test-1", "satisfied", "satisfied")
	thresholdResult2 := createTestResult("id-2", "test-1", "not-satisfied", "not-satisfied")
	resultComparisonMap2 := result.NewResultComparisonMap(result2, thresholdResult2)

	result3 := createTestResult("id-3", "test-2", "satisfied", "satisfied")
	thresholdResult3 := createTestResult("id-4", "test-2", "not-satisfied", "not-satisfied")
	resultComparisonMap3 := result.NewResultComparisonMap(result3, thresholdResult3)

	result4 := createTestResultNoObs("id-5", "satisfied")
	thresholdResult4 := createTestResultNoObs("id-5", "not-satisfied")
	resultComparisonMap4 := result.NewResultComparisonMap(result4, thresholdResult4)

	mapResultComparionMaps := map[string]result.ResultComparisonMap{
		"unchanged":       resultComparisonMap1,
		"changed":         resultComparisonMap2,
		"new-and-removed": resultComparisonMap3,
		"no-observations": resultComparisonMap4,
	}

	collapsedMap := result.Collapse(mapResultComparionMaps)
	observationPairMap, controlObservationMap, noObservations := result.RefactorObservationsByControls(collapsedMap)

	if len(observationPairMap) != 2 {
		t.Errorf("Expected 2 observation pairs, but got %d", len(observationPairMap))
	}
	for id := range observationPairMap {
		controls, ok := controlObservationMap[id]
		if !ok {
			t.Error("Expected controls to be in controlObservationMap, but it's not", id)
		}
		// check the observations are mapped to the controls correctly
		if id == "test-1" {
			if !contains(controls, "id-1") || !contains(controls, "id-2") {
				t.Errorf("Expected test-1 to contain id-1 and id-2, but got %v", controls)
			}
		}
		if id == "test-2" {
			if !contains(controls, "id-3") || !contains(controls, "id-4") {
				t.Errorf("Expected test-2 to contain id-3 and id-4, but got %v", controls)
			}
		}
	}
	if len(noObservations) != 1 {
		t.Errorf("Expected 1 no observation, but got %d", len(noObservations))
	}
	if !contains(noObservations, "id-5") {
		t.Errorf("Expected id-5 to be in no observations, but it's not")
	}
}
