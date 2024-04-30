package oscal_test

import (
	"fmt"
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
		want    oscalTypes.ComponentDefinition
		wantErr bool
	}{
		{
			name:    "Valid OSCAL Component Definition",
			source:  "test.yaml",
			data:    validBytes,
			want:    *validWantSchema.ComponentDefinition,
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
			want:         *validWantSchema.ComponentDefinition,
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
			want:         *validWantSchema.ComponentDefinition,
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
			fmt.Println(err != nil)
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

			if !reflect.DeepEqual(implementedRequirements, tt.requirements) {
				t.Errorf("Generated Requirements length mismatch - got = %v, want %v", implementedRequirements, tt.requirements)
			}

		})
	}

}
