package oscal_test

import (
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

var observations = []oscalTypes_1_1_2.Observation{
	{
		Collected:   time.Now(),
		Methods:     []string{"TEST"},
		UUID:        uuid.NewUUID(),
		Description: "test description",
	},
	{
		Collected:   time.Now(),
		Methods:     []string{"TEST"},
		UUID:        uuid.NewUUID(),
		Description: "test description",
	},
}

func TestIdentifyResults(t *testing.T) {
	t.Parallel()

	// Expecting an error when evaluating a single result
	t.Run("Handle valid assessment containing a single result", func(t *testing.T) {

		assessment, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml": assessment,
		}

		_, err = oscal.IdentifyResults(assessmentMap)
		if err == nil {
			t.Fatalf("Expected error for inability to identify multiple results : %v", err)
		}
	})

	// Identify threshold for multiple assessments and evaluate passing
	t.Run("Handle multiple threshold assessment containing a single result - pass", func(t *testing.T) {

		assessment, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml":   assessment,
			"invalid.yaml": assessment2,
		}

		resultMap, err := oscal.IdentifyResults(assessmentMap)
		if err != nil {
			t.Fatalf("Expected no error for inability to identify multiple results : %v", err)
		}

		if resultMap["threshold"] == nil || resultMap["latest"] == nil {
			t.Fatalf("Expected results to be identified")
		}

		if resultMap["threshold"].Start.After(resultMap["latest"].Start) {
			t.Fatalf("Expected threshold result to be before latest result")
		}

		status, _, err := oscal.EvaluateResults(resultMap["threshold"], resultMap["latest"])
		if err != nil {
			t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
		}

		if !status {
			t.Fatalf("Expected results to be evaluated as passing")
		}

	})

	// Identify threshold for multiple assessments and evaluate failing
	t.Run("Handle multiple threshold assessment containing a single result - fail", func(t *testing.T) {

		assessment, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults(findingMapFail, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		// key name does not matter here
		var assessmentMap = map[string]*oscalTypes_1_1_2.AssessmentResults{
			"valid.yaml":   assessment,
			"invalid.yaml": assessment2,
		}

		resultMap, err := oscal.IdentifyResults(assessmentMap)
		if err != nil {
			t.Fatalf("Expected error for inability to identify multiple results : %v", err)
		}

		if resultMap["threshold"] == nil || resultMap["latest"] == nil {
			t.Fatalf("Expected results to be identified")
		}

		if resultMap["threshold"].Start.After(resultMap["latest"].Start) {
			t.Fatalf("Expected threshold result to be before latest result")
		}

		status, _, err := oscal.EvaluateResults(resultMap["threshold"], resultMap["latest"])
		if err != nil {
			t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
		}

		if status {
			t.Fatalf("Expected results to be evaluated as failing")
		}
	})

	t.Run("Test merging two assessments - passing", func(t *testing.T) {

		assessment, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		assessment2, err := oscal.GenerateAssessmentResults(findingMapFail, observations)
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

		resultMap, err := oscal.IdentifyResults(assessmentMap)
		if err != nil {
			t.Fatalf("Expected error for inability to identify multiple results : %v", err)
		}

		if resultMap["threshold"] == nil || resultMap["latest"] == nil {
			t.Fatalf("Expected results to be identified")
		}

		if resultMap["threshold"].Start.After(resultMap["latest"].Start) {
			t.Fatalf("Expected threshold result to be before latest result")
		}

		status, _, err := oscal.EvaluateResults(resultMap["threshold"], resultMap["latest"])
		if err != nil {
			t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
		}

		if status {
			t.Fatalf("Expected results to be evaluated as failing")
		}
	})

	t.Run("Test merging two assessments - failing", func(t *testing.T) {

		assessment2, err := oscal.GenerateAssessmentResults(findingMapFail, observations)
		if err != nil {
			t.Fatalf("error generating assessment results: %v", err)
		}

		assessment, err := oscal.GenerateAssessmentResults(findingMapPass, observations)
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

		resultMap, err := oscal.IdentifyResults(assessmentMap)
		if err != nil {
			t.Fatalf("Expected error for inability to identify multiple results : %v", err)
		}

		if resultMap["threshold"] == nil || resultMap["latest"] == nil {
			t.Fatalf("Expected results to be identified")
		}

		if resultMap["threshold"].Start.After(resultMap["latest"].Start) {
			t.Fatalf("Expected threshold result to be before latest result")
		}

		status, _, err := oscal.EvaluateResults(resultMap["threshold"], resultMap["latest"])
		if err != nil {
			t.Fatalf("Expected error for inability to evaluate multiple results : %v", err)
		}

		if !status {
			t.Fatalf("Expected results to be evaluated as failing")
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

	if len(findings["new-passing-findings"]) != 1 {
		t.Fatal("error - expected 1 new finding, got ", len(findings["new-passing-findings"]))
	}

}
