package tui_test

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/muesli/termenv"
)

const timeout = time.Second * 20

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
	tea.Sequence()
}

// TestNewComponentDefinitionModel tests that the NewOSCALModel creates the expected model from component definition file
func TestNewComponentDefinitionModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel, "", tempLog)

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(timeout))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestMultiComponentDefinitionModel tests that the NewOSCALModel creates the expected model from component definition file
// and checks the component selection overlay -> new component section
func TestMultiComponentDefinitionModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-multi-component.yaml")
	model := tui.NewOSCALModel(oscalModel, "", tempLog)
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

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(timeout))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestNewAssessmentResultsModel tests that the NewOSCALModel creates the expected model from assessment results file
func TestNewAssessmentResultsModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-assessment-results.yaml")
	model := tui.NewOSCALModel(oscalModel, "", tempLog)
	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	testModel.Send(tea.KeyMsg{Type: tea.KeyTab})

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(timeout))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestComponentControlSelect tests that the user can navigate to a control, select it, and see expected
// remarks, description, and validations
func TestComponentControlSelect(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel, "", tempLog)
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

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(timeout))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}

// TestEditViewComponentDefinitionModel tests that the editing views of the component definition model are correct
func TestEditViewComponentDefinitionModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())
	tempOscalFile := testhelpers.CreateTempFile(t, "yaml")
	defer os.Remove(tempOscalFile.Name())

	oscalModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	model := tui.NewOSCALModel(oscalModel, tempOscalFile.Name(), tempLog)

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(common.DefaultWidth, common.DefaultHeight))

	// Edit the remarks
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight})                                    // Select component
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight})                                    // Select framework
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight})                                    // Select control
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})                                    // Open control
	testModel.Send(tea.KeyMsg{Type: tea.KeyRight})                                    // Navigate to remarks
	testModel.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})                // Edit remarks
	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlE})                                    // Newline
	testModel.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}}) // Add "test" to remarks
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter})                                    // Open control

	if err := testModel.Quit(); err != nil {
		t.Fatal(err)
	}

	if testModel == nil {
		t.Fatal("testModel is nil")
	}

	fm := testModel.FinalModel(t, teatest.WithFinalTimeout(timeout))

	teatest.RequireEqualOutput(t, []byte(fm.View()))
}
