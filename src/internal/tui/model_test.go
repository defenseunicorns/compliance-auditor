package tui_test

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
	tea.Sequence()
}

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

// TestNewComponentDefinitionModel tests that the NewOSCALModel creates the expected model from component definition file
func TestNewComponentDefinitionModel(t *testing.T) {
	oscalModel := oscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel)

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(time.Second*5))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestMultiComponentDefinitionModel tests that the NewOSCALModel creates the expected model from component definition file
// and checks the component selection overlay -> new component section
func TestMultiComponentDefinitionModel(t *testing.T) {
	oscalModel := oscalFromPath(t, "../../test/unit/common/oscal/valid-multi-component.yaml")
	model := tui.NewOSCALModel(oscalModel)
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select component
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter}) // enter component selection overlay
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown})  // navigate down
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter}) // select new component, exit overlay
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select framework
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select control
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Open control

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(time.Second*5))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestNewAssessmentResultsModel tests that the NewOSCALModel creates the expected model from assessment results file
func TestNewAssessmentResultsModel(t *testing.T) {
	oscalModel := oscalFromPath(t, "../../test/unit/common/oscal/valid-assessment-results.yaml")
	model := tui.NewOSCALModel(oscalModel)
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	testModel.Send(tea.KeyMsg{Type: tea.KeyTab})

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(time.Second*5))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestComponentControlSelect tests that the user can navigate to a control, select it, and see expected
// remarks, description, and validations
func TestComponentControlSelect(t *testing.T) {
	oscalModel := oscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel)
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	// Navigate to the control
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select component
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select framework
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight}) // Select control
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Open control

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(time.Second*5))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}
