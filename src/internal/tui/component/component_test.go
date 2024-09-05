package component_test

import (
	"os"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func oscalFromPath(t *testing.T, path string) *oscalTypes_1_1_2.OscalCompleteSchema {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	oscalModel, err := oscal.NewOscalModel(data)
	if err != nil {
		t.Fatalf("error creating oscal model from file: %v", err)
	}

	return oscalModel
}

func TestEditComponentDefinitionModel(t *testing.T) {
	oscalModel := oscalFromPath(t, "../../../test/unit/common/oscal/valid-generated-component.yaml")
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.TestSetControl()
	model.UpdateRemarks("test")

	// read remarks from the component definition
	compDefn := model.GetComponentDefinition()

	// try to write the model?
	mdl := &oscalTypes_1_1_2.OscalCompleteSchema{
		ComponentDefinition: compDefn,
	}
	oscal.WriteOscalModel("./oscal_test.yaml", mdl)
	for _, c := range *compDefn.Components {
		for _, f := range *c.ControlImplementations {
			for _, r := range f.ImplementedRequirements {
				if r.ControlId == "ac-1" {
					if r.Remarks != "test" {
						t.Errorf("Expected remarks to be 'test', got %s", r.Remarks)
					}
				}
			}
		}
	}
}

// func NewComponentDefinitionModel(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		component *component.Model
// 		want     *component.Model
// 	}{
// 	}
// 	// do a teatest thing here to update fields? and write?
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m := component.NewComponentDefinitionModel(tt.component)
// 			// test editing the pointers in the model
// 			m.
// 			assert.Equal(t, got, tt.want)
// 		})
// 	}
// }
