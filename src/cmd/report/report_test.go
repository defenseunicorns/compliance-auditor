package report

import (
    "testing"

    "github.com/defenseunicorns/lula/src/pkg/message"
    oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

// Mock functions for each OSCAL model

func MockAssessmentPlan() *oscalTypes_1_1_2.AssessmentPlan {
    return &oscalTypes_1_1_2.AssessmentPlan{
        UUID: "mock-assessment-plan-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock Assessment Plan",
            Version: "1.0",
        },
    }
}

func MockAssessmentResults() *oscalTypes_1_1_2.AssessmentResults {
    return &oscalTypes_1_1_2.AssessmentResults{
        UUID: "mock-assessment-results-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock Assessment Results",
            Version: "1.0",
        },
    }
}

func MockCatalog() *oscalTypes_1_1_2.Catalog {
    return &oscalTypes_1_1_2.Catalog{
        UUID: "mock-catalog-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock Catalog",
            Version: "1.0",
        },
    }
}

func MockComponentDefinition() *oscalTypes_1_1_2.ComponentDefinition {
    return &oscalTypes_1_1_2.ComponentDefinition{
        UUID: "mock-component-definition-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock Component Definition",
            Version: "1.0",
        },
        Components: &[]oscalTypes_1_1_2.DefinedComponent{
            {
                UUID:        "7c02500a-6e33-44e0-82ee-fba0f5ea0cae",
                Description: "Mock Component Description A",
                Title:       "Component A",
                Type:        "software",
                ControlImplementations: &[]oscalTypes_1_1_2.ControlImplementationSet{
                    {
                        Description: "Control Implementation Description",
                        ImplementedRequirements: []oscalTypes_1_1_2.ImplementedRequirementControlImplementation{
                            {
                                ControlId:   "ac-1",
                                Description: "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>",
                                Remarks:     "STATEMENT: Implementation details for ac-1.",
                                UUID:        "67dd59c4-0340-4aed-a49d-002815b50157",
                            },
                        },
                        Source: "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev4/yaml/NIST_SP-800-53_rev4_HIGH-baseline-resolved-profile_catalog.yaml",
                        UUID:   "0631b5b8-e51a-577b-8a43-2d3d0bd9ced8",
                        Props: &[]oscalTypes_1_1_2.Property{
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
                ControlImplementations: &[]oscalTypes_1_1_2.ControlImplementationSet{
                    {
                        Description: "Control Implementation Description",
                        ImplementedRequirements: []oscalTypes_1_1_2.ImplementedRequirementControlImplementation{
                            {
                                ControlId:   "ac-1",
                                Description: "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>",
                                Remarks:     "STATEMENT: Implementation details for ac-1.",
                                UUID:        "857121b1-2992-412c-b34a-504ead86e117",
                            },
                        },
                        Source: "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.yaml",
                        UUID:   "b1723ecd-a15a-5daf-a8e0-a7dd20a19abf",
                        Props: &[]oscalTypes_1_1_2.Property{
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

func MockPoam() *oscalTypes_1_1_2.PlanOfActionAndMilestones {
    return &oscalTypes_1_1_2.PlanOfActionAndMilestones{
        UUID: "mock-poam-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock POAM",
            Version: "1.0",
        },
    }
}

func MockProfile() *oscalTypes_1_1_2.Profile {
    return &oscalTypes_1_1_2.Profile{
        UUID: "mock-profile-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock Profile",
            Version: "1.0",
        },
    }
}

func MockSystemSecurityPlan() *oscalTypes_1_1_2.SystemSecurityPlan {
    return &oscalTypes_1_1_2.SystemSecurityPlan{
        UUID: "mock-system-security-plan-uuid",
        Metadata: oscalTypes_1_1_2.Metadata{
            Title:   "Mock System Security Plan",
            Version: "1.0",
        },
    }
}

func MockOscalModels() *oscalTypes_1_1_2.OscalCompleteSchema {
    return &oscalTypes_1_1_2.OscalCompleteSchema{
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
        oscalModel *oscalTypes_1_1_2.OscalCompleteSchema
        fileFormat string
        expectErr  bool
    }{
        {
            name:       "Component Definition Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{ComponentDefinition: MockComponentDefinition()},
            fileFormat: "table",
            expectErr:  false,
        },
        {
            name:       "Catalog Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{Catalog: MockCatalog()},
            fileFormat: "table",
            expectErr:  true,
        },
        {
            name:       "Assessment Plan Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{AssessmentPlan: MockAssessmentPlan()},
            fileFormat: "table",
            expectErr:  true,
        },
        {
            name:       "Assessment Results Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{AssessmentResults: MockAssessmentResults()},
            fileFormat: "table",
            expectErr:  true,
        },
        {
            name:       "POAM Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{PlanOfActionAndMilestones: MockPoam()},
            fileFormat: "table",
            expectErr:  true,
        },
        {
            name:       "Profile Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{Profile: MockProfile()},
            fileFormat: "table",
            expectErr:  true,
        },
        {
            name:       "System Security Plan Model",
            oscalModel: &oscalTypes_1_1_2.OscalCompleteSchema{SystemSecurityPlan: MockSystemSecurityPlan()},
            fileFormat: "table",
            expectErr:  true,
        },
    }

    for _, tc := range testCases {
        // Run tests sequentially to avoid race conditions due to spinner
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            err := handleOSCALModel(tc.oscalModel, tc.fileFormat)
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
