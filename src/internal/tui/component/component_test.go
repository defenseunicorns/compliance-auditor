package component_test

import (
	"os"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

// TestEditComponentDefinitionModel tests that a component definition model can be modified, written, and re-read
func TestEditComponentDefinitionModel(t *testing.T) {
	tempOscalFile := testhelpers.CreateTempFile(t, "yaml")
	defer os.Remove(tempOscalFile.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../../test/unit/common/oscal/valid-generated-component.yaml")
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)

	testControlId := "ac-1"
	testRemarks := "test remarks"
	testDescription := "test description"

	model.TestSetSelectedControl(testControlId)
	model.UpdateRemarks(testRemarks)
	model.UpdateDescription(testDescription)

	// Create OSCAL model
	mdl := &oscalTypes_1_1_2.OscalCompleteSchema{
		ComponentDefinition: model.GetComponentDefinition(),
	}

	// Write the model to a temp file
	err := oscal.OverwriteOscalModel(tempOscalFile.Name(), mdl)
	if err != nil {
		t.Errorf("error overwriting oscal model: %v", err)
	}

	// Read the model from the temp file
	modifiedOscalModel := testhelpers.OscalFromPath(t, tempOscalFile.Name())
	compDefn := modifiedOscalModel.ComponentDefinition
	if compDefn == nil {
		t.Errorf("component definition is nil")
	}
	for _, c := range *compDefn.Components {
		if c.ControlImplementations == nil {
			t.Errorf("control implementations are nil")
		}
		for _, f := range *c.ControlImplementations {
			for _, r := range f.ImplementedRequirements {
				if r.ControlId == testControlId {
					if r.Remarks != testRemarks {
						t.Errorf("Expected remarks to be %s, got %s", testRemarks, r.Remarks)
					}
					if r.Description != testDescription {
						t.Errorf("Expected remarks to be %s, got %s", testDescription, r.Description)
					}
				}
			}
		}
	}
}
