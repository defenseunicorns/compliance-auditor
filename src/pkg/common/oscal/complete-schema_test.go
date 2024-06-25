package oscal_test

import (
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"testing"
)

func TestGetOscalModel(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		Model     oscalTypes_1_1_2.OscalModels
		ModelType string
	}

	testCases := []TestCase{
		{
			Model: oscalTypes_1_1_2.OscalModels{
				Catalog: &oscalTypes_1_1_2.Catalog{},
			},
			ModelType: "catalog",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				Profile: &oscalTypes_1_1_2.Profile{},
			},
			ModelType: "profile",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				ComponentDefinition: &oscalTypes_1_1_2.ComponentDefinition{},
			},
			ModelType: "component",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				SystemSecurityPlan: &oscalTypes_1_1_2.SystemSecurityPlan{},
			},
			ModelType: "system-security-plan",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				AssessmentPlan: &oscalTypes_1_1_2.AssessmentPlan{},
			},
			ModelType: "assessment-plan",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				AssessmentResults: &oscalTypes_1_1_2.AssessmentResults{},
			},
			ModelType: "assessment-results",
		},
		{
			Model: oscalTypes_1_1_2.OscalModels{
				PlanOfActionAndMilestones: &oscalTypes_1_1_2.PlanOfActionAndMilestones{},
			},
			ModelType: "poam",
		},
	}
	for _, testCase := range testCases {
		actual, err := oscal.GetOscalModel(&testCase.Model)
		if err != nil {
			t.Fatalf("unexpected error for model %s", testCase.ModelType)
		}
		expected := testCase.ModelType
		if expected != actual {
			t.Fatalf("error GetOscalModel: expected: %s | got: %s", expected, actual)
		}
	}

}
