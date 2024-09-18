package oscal_test

import (
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func TestInjectIntoOSCALModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		path         string
		targetPath   string
		subsetPath   string
		expectedPath string
	}{
		{
			name:         "inject-component-definition",
			path:         "component-definition.metadata",
			targetPath:   "../../../test/unit/common/oscal/valid-component.yaml",
			subsetPath:   "../../../test/unit/common/oscal/valid-component-metadata.yaml",
			expectedPath: "../../../test/unit/common/oscal/valid-component-metadata-injected.yaml",
		},
		{
			name:         "inject-ssp",
			path:         "system-security-plan.metadata",
			targetPath:   "../../../test/unit/common/oscal/valid-ssp.yaml",
			subsetPath:   "../../../test/unit/common/oscal/valid-ssp-metadata.yaml",
			expectedPath: "../../../test/unit/common/oscal/valid-ssp-metadata-injected.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert paths to correct types
			targetBytes := loadTestData(t, tt.targetPath)
			subsetBytes := loadTestData(t, tt.subsetPath)
			expectedBytes := loadTestData(t, tt.expectedPath)

			// Convert the test data to expected types
			var targetModel oscalTypes_1_1_2.OscalCompleteSchema
			if err := yaml.Unmarshal(targetBytes, &targetModel); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}
			var subsetMap map[string]interface{}
			if err := yaml.Unmarshal(subsetBytes, &subsetMap); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}
			var expectedModel oscalTypes_1_1_2.OscalCompleteSchema
			if err := yaml.Unmarshal(expectedBytes, &expectedModel); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}

			result, err := oscal.InjectIntoOSCALModel(&targetModel, subsetMap, tt.path)
			if err != nil {
				t.Errorf("InjectIntoOSCALModel() error = %v", err)
			}
			assert.Equal(t, *result, expectedModel, "The OSCAL models should be equal")
		})
	}
}
