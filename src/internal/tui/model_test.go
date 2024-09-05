package tui_test

import (
	"os"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
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

// func TestDeepCopy(t *testing.T) {
// 	oscalModel := oscalFromPath(t, "../../test/unit/common/oscal/valid-generated-component.yaml")
// 	model := tui.NewOSCALModel(oscalModel, "test.yaml", nil)

// 	if !model.isModelSaved() {
// 		t.Errorf("why aren't these the same?")
// 	}
// }
