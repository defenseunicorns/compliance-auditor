package oscal_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestGetOscalModel(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		Model     oscalTypes.OscalModels
		ModelType string
	}

	testCases := []TestCase{
		{
			Model: oscalTypes.OscalModels{
				Catalog: &oscalTypes.Catalog{},
			},
			ModelType: "catalog",
		},
		{
			Model: oscalTypes.OscalModels{
				Profile: &oscalTypes.Profile{},
			},
			ModelType: "profile",
		},
		{
			Model: oscalTypes.OscalModels{
				ComponentDefinition: &oscalTypes.ComponentDefinition{},
			},
			ModelType: "component",
		},
		{
			Model: oscalTypes.OscalModels{
				SystemSecurityPlan: &oscalTypes.SystemSecurityPlan{},
			},
			ModelType: "system-security-plan",
		},
		{
			Model: oscalTypes.OscalModels{
				AssessmentPlan: &oscalTypes.AssessmentPlan{},
			},
			ModelType: "assessment-plan",
		},
		{
			Model: oscalTypes.OscalModels{
				AssessmentResults: &oscalTypes.AssessmentResults{},
			},
			ModelType: "assessment-results",
		},
		{
			Model: oscalTypes.OscalModels{
				PlanOfActionAndMilestones: &oscalTypes.PlanOfActionAndMilestones{},
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
			var targetModel oscalTypes.OscalCompleteSchema
			if err := yaml.Unmarshal(targetBytes, &targetModel); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}
			var subsetMap map[string]interface{}
			if err := yaml.Unmarshal(subsetBytes, &subsetMap); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}
			var expectedModel oscalTypes.OscalCompleteSchema
			if err := yaml.Unmarshal(expectedBytes, &expectedModel); err != nil {
				t.Fatalf("yaml.Unmarshal failed: %v", err)
			}

			result, err := oscal.InjectIntoOSCALModel(&targetModel, subsetMap, tt.path)
			if err != nil {
				t.Errorf("InjectIntoOSCALModel() error = %v", err)
			}
			assert.Equal(t, expectedModel, *result, "The OSCAL models should be equal")
		})
	}
}

func TestFetchOSCALModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		uri      string
		rootDir  string
		expected string
	}{
		{
			name:     "fetch-remote-catalog",
			uri:      "https://raw.githubusercontent.com/usnistgov/oscal-content/refs/heads/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_MODERATE-baseline-resolved-profile_catalog.yaml",
			rootDir:  "",
			expected: "catalog",
		},
		{
			name:     "fetch-local-component",
			uri:      "../../../test/unit/common/oscal/valid-component.yaml",
			rootDir:  "",
			expected: "component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, modelType, err := oscal.FetchOSCALModel(tt.uri, tt.rootDir)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, modelType)
			assert.NotNil(t, result)
		})
	}
}
