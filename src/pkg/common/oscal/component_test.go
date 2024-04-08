package oscal_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
	"sigs.k8s.io/yaml"
)

const validComponentPath = "../../../../test/common/oscal/valid-component.yaml"

// Helper function to load test data
func loadTestData(t *testing.T, path string) []byte {
	t.Helper() // Marks this function as a test helper
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file '%s': %v", path, err)
	}
	return data
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

func TestBackMatterToMap(t *testing.T) {

	tests := []struct {
		name string
		args oscalTypes_1_1_2.BackMatter
		want map[string]types.Validation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oscal.BackMatterToMap(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackMatterToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
