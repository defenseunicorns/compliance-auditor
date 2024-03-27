package oscal_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
	"gopkg.in/yaml.v3"
)

var (
	validComponentPath = "../../../test/e2e/scenarios/resource-data/oscal-component.yaml"
)

func TestBackMatterToMap(t *testing.T) {
	type args struct {
		backMatter oscalTypes_1_1_2.BackMatter
	}
	tests := []struct {
		name string
		args args
		want map[string]types.Validation
	}{
		{
			name: "Test No Resources",
			args: args{
				backMatter: oscalTypes_1_1_2.BackMatter{},
			},
		},
		// {
		// 	name: "Test "
		// }
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oscal.BackMatterToMap(tt.args.backMatter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackMatterToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOscalComponentDefinition(t *testing.T) {
	invalidConfig := oscalTypes_1_1_2.OscalCompleteSchema{}

	invalidBytes, err := yaml.Marshal(invalidConfig)
	if err != nil {
		t.Error(err)
	}

	validBytes, err := os.ReadFile(validComponentPath)
	if err != nil {
		t.Error(err)
	}

	var validWantSchema oscalTypes_1_1_2.OscalCompleteSchema
	err = yaml.Unmarshal(validBytes, &validWantSchema)
	if err != nil {
		t.Error(err)
	}

	validWant := *validWantSchema.ComponentDefinition

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    oscalTypes_1_1_2.ComponentDefinition
		wantErr bool
	}{
		{
			name: "Test NewOscalComponentDefinition",
			args: args{
				data: validBytes,
			},
			wantErr: false,
			want:    validWant,
		},
		{
			name: "Test NewOscalComponentDefinition Invalid",
			args: args{
				data: invalidBytes,
			},
			wantErr: true,
		},
		{
			name: "Test NewOscalComponentDefinition Empty Data",
			args: args{
				data: []byte{},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oscal.NewOscalComponentDefinition(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOscalComponentDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOscalComponentDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}
