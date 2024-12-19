package reporting

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// Mock functions for each OSCAL model

func MockAssessmentPlan() *oscalTypes.AssessmentPlan {
	return &oscalTypes.AssessmentPlan{
		UUID: "mock-assessment-plan-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock Assessment Plan",
			Version: "1.0",
		},
	}
}

func MockAssessmentResults() *oscalTypes.AssessmentResults {
	return &oscalTypes.AssessmentResults{
		UUID: "mock-assessment-results-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock Assessment Results",
			Version: "1.0",
		},
	}
}

func MockCatalog() *oscalTypes.Catalog {
	return &oscalTypes.Catalog{
		UUID: "mock-catalog-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock Catalog",
			Version: "1.0",
		},
	}
}

func MockComponentDefinition() *oscalTypes.ComponentDefinition {
	return &oscalTypes.ComponentDefinition{
		UUID: "mock-component-definition-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock Component Definition",
			Version: "1.0",
		},
		Components: &[]oscalTypes.DefinedComponent{
			{
				UUID:        "7c02500a-6e33-44e0-82ee-fba0f5ea0cae",
				Description: "Mock Component Description A",
				Title:       "Component A",
				Type:        "software",
				ControlImplementations: &[]oscalTypes.ControlImplementationSet{
					{
						Description: "Control Implementation Description",
						ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
							{
								ControlId:   "ac-1",
								Description: "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>",
								Remarks:     "STATEMENT: Implementation details for ac-1.",
								UUID:        "67dd59c4-0340-4aed-a49d-002815b50157",
							},
						},
						Source: "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev4/yaml/NIST_SP-800-53_rev4_HIGH-baseline-resolved-profile_catalog.yaml",
						UUID:   "0631b5b8-e51a-577b-8a43-2d3d0bd9ced8",
						Props: &[]oscalTypes.Property{
							{
								Name:  "framework",
								Ns:    "https://docs.lula.dev/ns",
								Value: "rev4",
							},
						},
					},
				},
			},
			{
				UUID:        "4cb1810c-d0d8-404e-b346-5a12c9629ed5",
				Description: "Mock Component Description B",
				Title:       "Component B",
				Type:        "software",
				ControlImplementations: &[]oscalTypes.ControlImplementationSet{
					{
						Description: "Control Implementation Description",
						ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
							{
								ControlId:   "ac-1",
								Description: "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>",
								Remarks:     "STATEMENT: Implementation details for ac-1.",
								UUID:        "857121b1-2992-412c-b34a-504ead86e117",
							},
						},
						Source: "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.yaml",
						UUID:   "b1723ecd-a15a-5daf-a8e0-a7dd20a19abf",
						Props: &[]oscalTypes.Property{
							{
								Name:  "framework",
								Ns:    "https://docs.lula.dev/ns",
								Value: "rev5",
							},
						},
					},
				},
			},
		},
	}
}

func MockPoam() *oscalTypes.PlanOfActionAndMilestones {
	return &oscalTypes.PlanOfActionAndMilestones{
		UUID: "mock-poam-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock POAM",
			Version: "1.0",
		},
	}
}

func MockProfile() *oscalTypes.Profile {
	return &oscalTypes.Profile{
		UUID: "mock-profile-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock Profile",
			Version: "1.0",
		},
	}
}

func MockSystemSecurityPlan() *oscalTypes.SystemSecurityPlan {
	return &oscalTypes.SystemSecurityPlan{
		UUID: "mock-system-security-plan-uuid",
		Metadata: oscalTypes.Metadata{
			Title:   "Mock System Security Plan",
			Version: "1.0",
		},
	}
}

func MockOscalModels() *oscalTypes.OscalCompleteSchema {
	return &oscalTypes.OscalCompleteSchema{
		AssessmentPlan:            MockAssessmentPlan(),
		AssessmentResults:         MockAssessmentResults(),
		Catalog:                   MockCatalog(),
		ComponentDefinition:       MockComponentDefinition(),
		PlanOfActionAndMilestones: MockPoam(),
		Profile:                   MockProfile(),
		SystemSecurityPlan:        MockSystemSecurityPlan(),
	}
}

// Test function for handleOSCALModel
func TestHandleOSCALModel(t *testing.T) {
	// Disable the spinner for this test function to work properly
	message.NoProgress = true

	// Define the test cases
	testCases := []struct {
		name       string
		oscalModel *oscalTypes.OscalCompleteSchema
		fileFormat string
		expectErr  bool
	}{
		{
			name:       "Component Definition Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{ComponentDefinition: MockComponentDefinition()},
			fileFormat: "table",
			expectErr:  false,
		},
		{
			name:       "Catalog Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{Catalog: MockCatalog()},
			fileFormat: "table",
			expectErr:  true,
		},
		{
			name:       "Assessment Plan Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{AssessmentPlan: MockAssessmentPlan()},
			fileFormat: "table",
			expectErr:  true,
		},
		{
			name:       "Assessment Results Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{AssessmentResults: MockAssessmentResults()},
			fileFormat: "table",
			expectErr:  true,
		},
		{
			name:       "POAM Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{PlanOfActionAndMilestones: MockPoam()},
			fileFormat: "table",
			expectErr:  true,
		},
		{
			name:       "Profile Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{Profile: MockProfile()},
			fileFormat: "table",
			expectErr:  true,
		},
		{
			name:       "System Security Plan Model",
			oscalModel: &oscalTypes.OscalCompleteSchema{SystemSecurityPlan: MockSystemSecurityPlan()},
			fileFormat: "table",
			expectErr:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Initialize CompositionContext
			compCtx, err := composition.New()
			if err != nil {
				t.Fatalf("failed to create composition context: %v", err)
			}

			// Call handleOSCALModel with compCtx
			err = handleOSCALModel(tc.oscalModel, tc.fileFormat, compCtx)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none for test case: %s", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got one for test case: %s, error: %v", tc.name, err)
				}
			}
		})
	}
}
