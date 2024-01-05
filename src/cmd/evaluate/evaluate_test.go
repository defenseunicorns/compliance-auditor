package evaluate

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// Given two results - evaluate for passing
func TestEvaluateResultsPassing(t *testing.T) {
	message.NoProgress = true

	mockThresholdResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
		},
	}

	mockEvaluationResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
		},
	}

	status, _, err := EvaluateResults(mockThresholdResult, mockEvaluationResult)
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
	mockThresholdResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
		},
	}

	mockEvaluationResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "not-satisfied",
					},
				},
			},
		},
	}

	status, findings, err := EvaluateResults(mockThresholdResult, mockEvaluationResult)
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

func TestEvaluateResultsNewFindings(t *testing.T) {
	message.NoProgress = true
	mockThresholdResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
		},
	}
	// Adding two new findings
	mockEvaluationResult := oscalTypes.Result{
		Findings: []oscalTypes.Finding{
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-1",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-2",
					Status: oscalTypes.Status{
						State: "satisfied",
					},
				},
			},
			{
				Target: oscalTypes.FindingTarget{
					TargetId: "ID-3",
					Status: oscalTypes.Status{
						State: "not-satisfied",
					},
				},
			},
		},
	}

	status, findings, err := EvaluateResults(mockThresholdResult, mockEvaluationResult)
	if err != nil {
		t.Fatal(err)
	}

	// If status is false - then something went wrong
	if !status {
		t.Fatal("error - evaluation failed")
	}

	if len(findings["new-findings"]) != 2 {
		t.Fatal("error - expected 1 new finding, got ", len(findings["new-findings"]))
	}

}
