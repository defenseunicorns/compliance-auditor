package oscal_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
	"gopkg.in/yaml.v3"
)

const validComponentPath = "../../../test/e2e/scenarios/resource-data/oscal-component.yaml"

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
	validBackMatterMapBytes := loadTestData(t, "../../../../test/validBackMatterMap.yaml")

	var validComponent oscalTypes.OscalCompleteSchema
	if err := yaml.Unmarshal(validComponentBytes, &validComponent); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}
	var validBackMatterMap map[string]types.Validation
	if err := yaml.Unmarshal(validBackMatterMapBytes, &validBackMatterMap); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	tests := []struct {
		name       string
		backMatter oscalTypes.BackMatter
		want       map[string]types.Validation
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
		want    oscalTypes.ComponentDefinition
		wantErr bool
	}{
		{
			name:    "Valid OSCAL Component Definition",
			data:    validBytes,
			want:    *validWantSchema.ComponentDefinition,
			wantErr: false,
		},
		{
			name:    "Invalid OSCAL Component Definition",
			data:    invalidBytes,
			wantErr: true,
		},
		{
			name:    "Empty Data",
			data:    []byte{},
			wantErr: true,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oscal.NewOscalComponentDefinition(tt.data)
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
