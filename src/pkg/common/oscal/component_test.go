package oscal_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"gopkg.in/yaml.v3"
)

const validComponentPath = "../../../test/unit/common/oscal/valid-component.yaml"
const catalogPath = "../../../test/unit/common/oscal/catalog.yaml"

// Helper function to load test data
func loadTestData(t *testing.T, path string) []byte {
	t.Helper() // Marks this function as a test helper
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file '%s': %v", path, err)
	}
	return data
}

func TestBackMatterToMap(t *testing.T) {
	validComponentBytes := loadTestData(t, validComponentPath)
	validBackMatterMapBytes := loadTestData(t, "../../../test/unit/common/oscal/valid-back-matter-map.yaml")

	var validComponent oscalTypes.OscalCompleteSchema
	if err := yaml.Unmarshal(validComponentBytes, &validComponent); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}
	var validBackMatterMap map[string]string
	if err := yaml.Unmarshal(validBackMatterMapBytes, &validBackMatterMap); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	tests := []struct {
		name       string
		backMatter oscalTypes.BackMatter
		want       map[string]string
	}{
		{
			name:       "Test No Resources",
			backMatter: oscalTypes.BackMatter{},
			want:       make(map[string]string),
		},
		{
			name:       "Test Valid Component",
			backMatter: *validComponent.ComponentDefinition.BackMatter,
			want:       validBackMatterMap,
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := oscal.BackMatterToMap(tc.backMatter)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("BackMatterToMap() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewOscalComponentDefinition(t *testing.T) {
	validBytes := loadTestData(t, validComponentPath)

	var validWantSchema oscalTypes.OscalCompleteSchema
	if err := yaml.Unmarshal(validBytes, &validWantSchema); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	invalidBytes, err := yaml.Marshal(oscalTypes.OscalCompleteSchema{})
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		source  string
		want    *oscalTypes.ComponentDefinition
		wantErr bool
	}{
		{
			name:    "Valid OSCAL Component Definition",
			source:  "test.yaml",
			data:    validBytes,
			want:    validWantSchema.ComponentDefinition,
			wantErr: false,
		},
		{
			name:    "Invalid OSCAL Component Definition",
			data:    invalidBytes,
			source:  "",
			wantErr: true,
		},
		{
			name:    "Invalid OSCAL source with valid data",
			data:    validBytes,
			source:  "test.go",
			wantErr: true,
		},
		{
			name:    "Empty Data",
			data:    []byte{},
			source:  "",
			wantErr: true,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oscal.NewOscalComponentDefinition(tt.source, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOscalComponentDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("NewOscalComponentDefinition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentFromCatalog(t *testing.T) {
	validBytes := loadTestData(t, "../../../test/unit/common/oscal/valid-generated-component.yaml")

	var validWantSchema oscalTypes.OscalCompleteSchema
	if err := yaml.Unmarshal(validBytes, &validWantSchema); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}
	// let's create a catalog from a test document
	catalogBytes := loadTestData(t, catalogPath)

	catalog, err := oscal.NewCatalog("https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json", catalogBytes)
	if err != nil {
		t.Errorf("error creating catalog from path %s", catalogPath)
	}

	tests := []struct {
		name         string
		data         oscalTypes.Catalog
		title        string
		source       string
		requirements []string
		remarks      []string
		want         oscalTypes.ComponentDefinition
		wantReqLen   int
		wantErr      bool
	}{
		{
			name:         "Valid test of component from Catalog",
			data:         catalog,
			title:        "Component Title",
			requirements: []string{"ac-1", "ac-3", "ac-3.2", "ac-4"},
			remarks:      []string{"statement"},
			source:       "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			wantReqLen:   4,
			wantErr:      false,
		},
		{
			name:         "Valid test of component from Catalog with malformed control",
			data:         catalog,
			title:        "Component Title",
			requirements: []string{"ac-1", "ac-3", "ac-3.2", "ac-4", "100"},
			remarks:      []string{"statement"},
			source:       "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			wantReqLen:   4,
			wantErr:      false,
		},
		{
			name:         "Invalid amount of requirements specified",
			data:         catalog,
			title:        "Component Test Title",
			requirements: []string{},
			remarks:      []string{"statement"},
			source:       "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			wantErr:      true,
		},
		{
			name:    "Invalid test of empty catalog",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oscal.ComponentFromCatalog(tt.source, tt.data, tt.title, tt.requirements, tt.remarks)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComponentFromCatalog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Given pointers below - let's return here if we've met the check above and wanted an error
			if tt.wantErr {
				return
			}

			// DeepEqual will be difficult with time/uuid generation
			component := (*got.Components)[0]
			if component.Title != tt.title {
				t.Errorf("ComponentFromCatalog() title = %v, want %v", component.Title, tt.title)
			}

			controlImplementation := (*component.ControlImplementations)[0]
			if controlImplementation.Source != tt.source {
				t.Errorf("ComponentFromCatalog() source = %v, want %v", controlImplementation.Source, tt.source)
			}

			implementedRequirements := make([]string, 0)
			for _, requirement := range controlImplementation.ImplementedRequirements {
				implementedRequirements = append(implementedRequirements, requirement.ControlId)
			}

			reqLen := len(implementedRequirements)
			if reqLen != tt.wantReqLen {
				t.Errorf("Generated Requirements length mismatch - got = %v, want %v", reqLen, tt.wantReqLen)
			}
		})
	}
}

func TestMergeComponentDefinitions(t *testing.T) {
	validBytes := loadTestData(t, "../../../test/unit/common/oscal/valid-generated-component.yaml")
	// generate a new artifact
	catalogBytes := loadTestData(t, catalogPath)
	source := "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json"
	catalog, err := oscal.NewCatalog(source, catalogBytes)
	if err != nil {
		t.Errorf("error creating catalog from path %s", catalogPath)
	}

	tests := []struct {
		name                                  string
		title                                 string
		source                                string
		requirements                          []string
		remarks                               []string
		expectedComponents                    int
		expectedControlImplementations        int
		expectedImplementedRequirements       int
		expectedTargetControlImplementations  int
		expectedTargetImplementedRequirements int
		wantErr                               bool
	}{
		{
			name:                                  "Valid test of component merge",
			title:                                 "Component Title",
			requirements:                          []string{"ac-1", "ac-3", "ac-3.2", "ac-4"},
			remarks:                               []string{"statement"},
			source:                                "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			expectedComponents:                    1,
			expectedControlImplementations:        1,
			expectedImplementedRequirements:       4,
			expectedTargetControlImplementations:  1,
			expectedTargetImplementedRequirements: 4,
			wantErr:                               false,
		},
		{
			name:                                  "Valid test of component merge with multiple unique controls",
			title:                                 "Component Title",
			requirements:                          []string{"ac-1", "ac-3", "ac-3.2", "ac-4", "ac-4.4", "au-5"},
			remarks:                               []string{"statement"},
			source:                                "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			expectedComponents:                    1,
			expectedControlImplementations:        1,
			expectedImplementedRequirements:       6,
			expectedTargetControlImplementations:  1,
			expectedTargetImplementedRequirements: 6,
			wantErr:                               false,
		},
		{
			name:                                  "Valid test of component merge with multiple unique components",
			title:                                 "Component Test Title",
			requirements:                          []string{"ac-1", "ac-3", "ac-3.2", "ac-4"},
			remarks:                               []string{"statement"},
			source:                                "https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json",
			expectedComponents:                    2,
			expectedImplementedRequirements:       8,
			expectedControlImplementations:        2,
			expectedTargetImplementedRequirements: 4,
			expectedTargetControlImplementations:  1,
			wantErr:                               false,
		},
		{
			name:    "Invalid test of empty existing component definition",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validComponent, _ := oscal.NewOscalComponentDefinition("valid-generated-component.yaml", validBytes)

			// Get the implemented requirements from existing for comparison
			existingComponent := (*validComponent.Components)[0]
			existingControlImplementation := (*existingComponent.ControlImplementations)[0]
			existingImplementedRequirementsMap := make(map[string]bool)
			for _, req := range existingControlImplementation.ImplementedRequirements {
				existingImplementedRequirementsMap[req.ControlId] = true
			}

			generated, _ := oscal.ComponentFromCatalog(tt.source, catalog, tt.title, tt.requirements, tt.remarks)

			merged, err := oscal.MergeComponentDefinitions(validComponent, generated)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeComponentDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Perform checks on quantities
			components := (*merged.Components)
			if len(components) != tt.expectedComponents {
				t.Errorf("MergeComponentDefinitions() expected %v components, got %v", tt.expectedComponents, len((*merged.Components)))
			}
			controlImplementations := make([]oscalTypes.ControlImplementationSet, 0)
			var targetComponent oscalTypes.DefinedComponent
			for _, component := range components {
				if component.ControlImplementations != nil {
					if component.Title == "Component Title" {
						targetComponent = component
					}
					controlImplementations = append(controlImplementations, (*component.ControlImplementations)...)
				}
			}
			if len(controlImplementations) != tt.expectedControlImplementations {
				t.Errorf("MergeComponentDefinitions() expected %v control implementations, got %v", tt.expectedControlImplementations, len(controlImplementations))
			}

			implementedRequirements := make([]oscalTypes.ImplementedRequirementControlImplementation, 0)
			for _, control := range controlImplementations {
				implementedRequirements = append(implementedRequirements, control.ImplementedRequirements...)
			}

			if len(implementedRequirements) != tt.expectedImplementedRequirements {
				t.Errorf("MergeComponentDefinitions() expected %v implemented requirements, got %v", tt.expectedImplementedRequirements, len(implementedRequirements))
			}

			// Now operate on the target existing component & items (should only be 1 component, 1 control-implementation and dynamic implemented-requirements)
			if targetComponent.ControlImplementations == nil {
				t.Errorf("MergeComponentDefinitions() missing control-implementations in component %s", targetComponent.Title)
			}
			targetControlImplementations := (*targetComponent.ControlImplementations)

			if len(targetControlImplementations) != tt.expectedTargetControlImplementations {
				t.Errorf("MergeComponentDefinitions() expected %v control-implementations in component %s, got %v", tt.expectedTargetControlImplementations, targetComponent.Title, len(controlImplementations))
			}
			var targetControlImp oscalTypes.ControlImplementationSet
			for _, item := range targetControlImplementations {
				if item.Source == source {
					targetControlImp = item
				}
			}
			// check implemented requirements for length
			if len(targetControlImp.ImplementedRequirements) != tt.expectedTargetImplementedRequirements {
				t.Errorf("MergeComponentDefinitions() expected %v implemented-requirements in component %s, got %v", tt.expectedTargetImplementedRequirements, targetComponent.Title, len(targetControlImp.ImplementedRequirements))
			}

			// check implemented requirements for existing content - add to the test artifact
			for _, req := range targetControlImp.ImplementedRequirements {
				if _, ok := existingImplementedRequirementsMap[req.ControlId]; ok {
					if req.Description != "This is the existing test string" {
						t.Errorf("MergeComponentDefinitions() expected description 'this is the existing test string' but got %s", req.Description)
					}
				}
			}
		})
	}

}
