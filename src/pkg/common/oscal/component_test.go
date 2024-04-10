package oscal_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"gopkg.in/yaml.v3"
)

const validComponentPath = "../../../test/unit/common/oscal/valid-component.yaml"

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
	data := loadTestData(t, validComponentPath)
	var component oscalTypes_1_1_2.OscalCompleteSchema
	err := yaml.Unmarshal(data, &component)
	if err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	got := oscal.BackMatterToMap(*component.ComponentDefinition.BackMatter)
	if got == nil {
		t.Fatalf("BackMatterToMap returned nil")
	}

	if len(got) == 0 {
		t.Fatalf("BackMatterToMap returned empty map")
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
