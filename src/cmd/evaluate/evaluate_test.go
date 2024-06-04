package evaluate

import (
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

var (
	validInputFile   = "../../test/unit/common/oscal/valid-assessment-result.yaml"
	invalidInputFile = "../../test/unit/common/oscal/invalid-assessment-result.yaml"
)

func TestEvaluateAssessmentResults(t *testing.T) {
	t.Parallel()

	// TODO: write logic to separate file read from core evaluation logic
	// TODO: move the core logic to library package
	// TODO: write a success test to receive assessments/results that we can verify intended prop change
	// t.Run("handles valid assessment result", func(t *testing.T) {
	// 	assessmentMap, err := EvaluateAssessmentResults([]string{validInputFile})
	// 	if err != nil {
	// 		t.Fatal("unexpected error for valid assessment result")
	// 	}

	// })

	t.Run("handles invalid path to assessment result file", func(t *testing.T) {
		_, err := EvaluateAssessmentResults([]string{"./invalid-path.yaml"})
		if err == nil {
			t.Fatal("expected error for invalid path")
		}
	})

	t.Run("handles invalid assessment result without any results", func(t *testing.T) {
		_, err := EvaluateAssessmentResults([]string{invalidInputFile})
		if err == nil {
			t.Fatal("expected error for invalid assessment result without results")
		}
	})

}

// Given two results - evaluate for passing
func TestEvaluateResultsPassing(t *testing.T) {
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

	status, _, err := EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
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

	mockEvaluationResult := oscalTypes_1_1_2.Result{
		Findings: &[]oscalTypes_1_1_2.Finding{
			{
				Target: oscalTypes_1_1_2.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: "not-satisfied",
					},
				},
			},
		},
	}

	status, findings, err := EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
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

	status, _, err := EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
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

	_, _, err := EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
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

	status, findings, err := EvaluateResults(&mockThresholdResult, &mockEvaluationResult)
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
