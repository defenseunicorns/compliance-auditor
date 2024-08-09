package oscal_test

import (
	"slices"
	"testing"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// Create re-usable findings and observations
// use those in tests to generate test assessment results
var findingMapPass = map[string]oscalTypes_1_1_2.Finding{
	"ID-1": {
		Target: oscalTypes_1_1_2.FindingTarget{
			TargetId: "ID-1",
			Status: oscalTypes_1_1_2.ObjectiveStatus{
				State: "satisfied",
			},
		},
	},
}

var findingMapFail = map[string]oscalTypes_1_1_2.Finding{
	"ID-1": {
		Target: oscalTypes_1_1_2.FindingTarget{
			TargetId: "ID-1",
			Status: oscalTypes_1_1_2.ObjectiveStatus{
				State: "not-satisfied",
			},
		},
	},
}

var findings = []oscalTypes_1_1_2.Finding{
	{
		Target: oscalTypes_1_1_2.FindingTarget{
			TargetId: "ID-1",
			Status: oscalTypes_1_1_2.ObjectiveStatus{
				State: "satisfied",
			},
		},
	},
	{
		Target: oscalTypes_1_1_2.FindingTarget{
			TargetId: "ID-2",
			Status: oscalTypes_1_1_2.ObjectiveStatus{
				State: "not-satisfied",
			},
		},
	},
}

// Delineate between these two observations based on the description
var observations = []oscalTypes_1_1_2.Observation{
	{
		Collected:   time.Now(),
		Methods:     []string{"TEST"},
		UUID:        "4344e734-63d7-4bda-81f1-b805f60fdbf5",
		Description: "test description first",
	},
	{
		Collected:   time.Now(),
		Methods:     []string{"TEST"},
		UUID:        "1ac95fcc-1adb-4a25-89a7-08a708def2f3",
		Description: "test description second",
	},
}

func TestFilterResults(t *testing.T) {
	t.Parallel()

	// Expecting an error when evaluating assessment without results
	t.Run("Handle invalid assessment containing no results", func(t *testing.T) {

		var assessment = &oscalTypes_1_1_2.AssessmentResults{
			UUID: uuid.NewUUID(),
		}
		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml": assessment,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		if len(resultMap) > 0 {
			t.Fatalf("Expected resultMap length to be 0")
		}
	})

	// Expecting an error when evaluating a single result
	t.Run("Handle valid assessment containing a single result", func(t *testing.T) {

		result, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{result})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml": assessment,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		if resultMap == nil {
			t.Fatalf("Expected resultMap to be non-nil")
		}

		if len(resultMap) == 0 {
			t.Fatalf("Expected resultMap to contain the single result")
		}
	})

	// Identify threshold for multiple assessments and evaluate passing
	t.Run("Handle multiple threshold assessment containing a single result - pass", func(t *testing.T) {
		result, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{result})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		resultSecond, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultSecond})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml":   assessment,
			"invalid.yaml": assessment2,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		for _, result := range resultMap {
			if result.Threshold == nil || result.Latest == nil {
				t.Fatalf("Expected results to be identified")
			}

			if result.Threshold.Start.After(result.Latest.Start) {
				t.Fatalf("Expected threshold result to be before latest result")
			}

			status, _, err := oscal.EvaluateResults(result.Threshold, result.Latest)
			if err != nil {
				t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
			}

			if !status {
				t.Fatalf("Expected results to be evaluated as passing")
			}
		}

	})

	// Identify threshold for multiple assessments and evaluate failing
	t.Run("Handle multiple threshold assessment containing a single result - fail", func(t *testing.T) {

		resultPass, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultPass})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		resultFail, err := oscal.CreateResult(findingMapFail, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultFail})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml":   assessment,
			"invalid.yaml": assessment2,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		for _, result := range resultMap {
			if result.Threshold == nil || result.Latest == nil {
				t.Fatalf("Expected results to be identified")
			}

			if result.Threshold.Start.After(result.Latest.Start) {
				t.Fatalf("Expected threshold result to be before latest result")
			}

			status, _, err := oscal.EvaluateResults(result.Threshold, result.Latest)
			if err != nil {
				t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
			}

			if status {
				t.Fatalf("Expected results to be evaluated as failing")
			}
		}
	})

	t.Run("Test merging two assessments - failing", func(t *testing.T) {

		resultPass, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultPass})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		resultFail, err := oscal.CreateResult(findingMapFail, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultFail})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// Update assessment 2 props so that we only have 1 threshold
		oscal.UpdateProps("threshold", "docs.lula.dev/ns", "false", assessment2.Results[0].Props)

		assessment, err = oscal.MergeAssessmentResults(assessment, assessment2)
		if err != nil {
			t.Fatalf("error merging assessment results: %v", err)
		}

		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml": assessment,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		for _, result := range resultMap {
			if result.Threshold == nil || result.Latest == nil {
				t.Fatalf("Expected results to be identified")
			}

			if result.Threshold.Start.After(result.Latest.Start) {
				t.Fatalf("Expected threshold result to be before latest result")
			}

			status, _, err := oscal.EvaluateResults(result.Threshold, result.Latest)
			if err != nil {
				t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
			}

			if status {
				t.Fatalf("Expected results to be evaluated as failing")
			}
		}
	})

	t.Run("Test merging two assessments - passing", func(t *testing.T) {

		resultFail, err := oscal.CreateResult(findingMapFail, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultFail})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		resultPass, err := oscal.CreateResult(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating result from findings and observations: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults([]oscalTypes_1_1_2.Result{resultPass})
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// Update assessment props so that we only have 1 threshold
		oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "false", assessment.Results[0].Props)

		// TODO: review assumptions made about order of assessments during merge
		assessment, err = oscal.MergeAssessmentResults(assessment, assessment2)
		if err != nil {
			t.Fatalf("error merging assessment results: %v", err)
		}

		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml": assessment,
		}

		resultMap := oscal.FilterResults(assessmentMap)

		for _, result := range resultMap {
			if result.Threshold == nil || result.Latest == nil {
				t.Fatalf("Expected results to be identified")
			}

			if result.Threshold.Start.After(result.Latest.Start) {
				t.Fatalf("Expected threshold result to be before latest result")
			}

			status, resultComparison, err := oscal.EvaluateResults(result.Threshold, result.Latest)
			if err != nil {
				t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
			}

			if !status {
				t.Fatalf("Expected results to be evaluated as passing")
			}

			if len(resultComparison["now-satisfied"]) == 0 {
				t.Fatalf("Expected new passing findings to be found")
			}

		}
	})

}

// Given two results - evaluate for passing
func TestEvaluateResultsPassing(t *testing.T) {
	message.NoProgress = true

	mockThresholdResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			findingMapPass["ID-1"],
		},
	}

	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			findingMapPass["ID-1"],
		},
	}

	status, _, err := oscal.EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
	if err != nil {
		t.Fatal(err)
	}

	// If status is false - then something went wrong
	if !status {
		t.Fatal("error - evaluation failed")
	}

}

func TestEvaluateResultsFailed(t *testing.T) {
	message.NoProgress = true
	mockThresholdResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			findingMapPass["ID-1"],
		},
	}

	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			findingMapFail["ID-1"],
		},
	}

	status, findings, err := oscal.EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
	if err != nil {
		t.Fatal(err)
	}

	// If status is true - then something went wrong
	if status {
		t.Fatal("error - evaluation was successful when it should have failed")
	}

	if len(findings["no-longer-satisfied"]) != 1 {
		t.Fatal("error - expected 1 finding, got ", len(findings["no-longer-satisfied"]))
	}

}

func TestEvaluateResultsNoFindings(t *testing.T) {
	message.NoProgress = true
	mockThresholdResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{},
	}

	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{},
	}

	status, _, err := oscal.EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
	if err != nil {
		t.Fatal(err)
	}

	// If status is false - then something went wrong
	if !status {
		t.Fatal("error - evaluation failed")
	}

}

func TestEvaluateResultsNoThreshold(t *testing.T) {
	message.NoProgress = true
	mockThresholdResult := oscalTypes_1_1_2.Result{}

	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "satisfied",
					},
				},
			},
		},
	}

	_, _, err := oscal.EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
	if err == nil {
		t.Fatal("error - expected error, got nil")
	}
}

func TestEvaluateResultsNewFindings(t *testing.T) {
	message.NoProgress = true
	mockThresholdResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "satisfied",
					},
				},
			},
		},
	}
	// Adding two new findings
	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "satisfied",
					},
				},
			},
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-2",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "satisfied",
					},
				},
			},
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-3",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "not-satisfied",
					},
				},
			},
		},
	}

	status, findings, err := oscal.EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
	if err != nil {
		t.Fatal(err)
	}

	// If status is false - then something went wrong
	if !status {
		t.Fatal("error - evaluation failed")
	}

	if len(findings["new-satisfied"]) != 1 {
		t.Fatal("error - expected 1 new finding, got ", len(findings["new-passing-findings"]))
	}

}

func TestMakeAssessmentResultsDeterministic(t *testing.T) {
	// reverse the order
	slices.Reverse(findings)
	slices.Reverse(observations)

	// Will already be in reverse order
	var results = []oscalTypes_1_1_2.Result{
		{
			Start:        time.Now(),
			UUID:         "d66c9509-cb92-4597-86f8-6e6623ea9154",
			Findings:     &findings,
			Observations: &observations,
		},
		{
			Start:        time.Now(),
			UUID:         "28174d67-06a7-4c7c-be04-1edf437d4ece",
			Findings:     &findings,
			Observations: &observations,
		},
	}

	var assessment = oscalTypes_1_1_2.AssessmentResults{
		Results: results,
	}

	oscal.MakeAssessmentResultsDeterministic(&assessment)

	if len(assessment.Results) < 2 {
		t.Fatalf("Expected 2 results, got %d", len(assessment.Results))
	}

	// Assessment-Results.Results are sorted newest to oldest
	var resultExpected = []string{"28174d67-06a7-4c7c-be04-1edf437d4ece", "d66c9509-cb92-4597-86f8-6e6623ea9154"}
	//Verify order

	for key, id := range resultExpected {

		if assessment.Results[key].UUID != id {
			t.Fatalf("Expected UUID %q, got %q", id, assessment.Results[key].UUID)
		}

		assessmentResult := assessment.Results[key]
		if assessmentResult.Findings == nil {
			t.Fatal("Expected findings, got nil")
		}

		assesmentFindings := *assessmentResult.Findings

		if len(assesmentFindings) != 2 {
			t.Fatalf("Expected 2 findings, got %d", len(findings))
		}

		var findingExpected = []string{"ID-1", "ID-2"}

		for key, id := range findingExpected {
			if assesmentFindings[key].Target.TargetId != id {
				t.Fatalf("Expected finding %q, got %q", id, assesmentFindings[key].Target.TargetId)
			}
		}

		if assessmentResult.Observations == nil {
			t.Fatal("Expected observations, got nil")
		}

		assessmentObservations := *assessmentResult.Observations

		if len(assessmentObservations) != 2 {
			t.Fatalf("Expected 2 observations, got %d", len(assessmentObservations))
		}

		var observationExpected = []string{"4344e734-63d7-4bda-81f1-b805f60fdbf5", "1ac95fcc-1adb-4a25-89a7-08a708def2f3"}

		for key, id := range observationExpected {
			if assessmentObservations[key].UUID != id {
				t.Fatalf("Expected observation %q, got %q", id, assessmentObservations[key].UUID)
			}
		}
	}
}

func TestCreateResult(t *testing.T) {
	tests := []struct {
		name         string
		findingMap   map[string]oscalTypes_1_1_2.Finding
		observations []oscalTypes_1_1_2.Observation
		expected     string
	}{
		{
			name:         "passing result",
			findingMap:   findingMapPass,
			observations: observations,
			expected:     "satisfied",
		},
		{
			name:         "failing result",
			findingMap:   findingMapFail,
			observations: observations,
			expected:     "not-satisfied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result, err := oscal.CreateResult(tt.findingMap, tt.observations)
			if err != nil {
				t.Fatalf("error generating result from findings and observations: %v", err)
			}

			if result.Findings == nil {
				t.Fatal("expected findings to be non-nil")
			}

			for _, finding := range *result.Findings {
				if finding.Target.Status.State != tt.expected {
					t.Fatalf("Expected %s state, got %s", tt.expected, finding.Target.Status.State)
				}
			}

			if result.Observations == nil {
				t.Fatal("expected observations to be non-nil")
			}

			if len(*result.Observations) == 0 {
				t.Fatal("expected observations to be greater than zero")
			}
		})
	}
}
