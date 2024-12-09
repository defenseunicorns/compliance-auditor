package component_test

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timeout    = time.Second * 20
	maxRetries = 3
	height     = common.DefaultHeight
	width      = common.DefaultWidth

	validCompDef                 = "../../../test/unit/common/oscal/valid-generated-component.yaml"
	validCompDefValidations      = "../../../test/unit/common/oscal/valid-component.yaml"
	validCompDefMulti            = "../../../test/unit/common/oscal/valid-multi-component.yaml"
	validCompDefMultiValidations = "../../../test/unit/common/oscal/valid-multi-component-validations.yaml"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

// TestComponentDefinitionBasicView tests that the model is created correctly from a component definition with validations
func TestComponentDefinitionBasicView(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDef)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.Open(height, width)

	msgs := []tea.Msg{}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestComponentDefinitionComponentSwitch tests that the component picker executes correctly
func TestComponentDefinitionComponentSwitch(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefMulti)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight}, // Select component
		tea.KeyMsg{Type: tea.KeyEnter}, // enter component selection overlay
		tea.KeyMsg{Type: tea.KeyDown},  // navigate down
		tea.KeyMsg{Type: tea.KeyEnter}, // select new component, exit overlay
		tea.KeyMsg{Type: tea.KeyRight}, // Select framework
		tea.KeyMsg{Type: tea.KeyRight}, // Select control
		tea.KeyMsg{Type: tea.KeyEnter}, // Open control
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestComponentControlSelect tests that the user can navigate to a control, select it, and see expected
// remarks, description, and validations
func TestComponentControlSelect(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefMulti)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight}, // Select component
		tea.KeyMsg{Type: tea.KeyRight}, // Select framework
		tea.KeyMsg{Type: tea.KeyRight}, // Select control
		tea.KeyMsg{Type: tea.KeyEnter}, // Open control
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestEditViewComponentDefinitionModel tests that the editing views of the component definition model are correct
func TestEditViewComponentDefinitionModel(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefValidations)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select component
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select framework
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select control
		tea.KeyMsg{Type: tea.KeyEnter},                                    // Open control
		tea.KeyMsg{Type: tea.KeyRight},                                    // Navigate to remarks
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},                // Edit remarks
		tea.KeyMsg{Type: tea.KeyCtrlE},                                    // Newline
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}}, // Add "test" to remarks
		tea.KeyMsg{Type: tea.KeyEnter},                                    // Confirm edit
	}

	reset := func() tea.Model {
		resetOscalModel := testhelpers.OscalFromPath(t, validCompDefValidations)
		resetModel := component.NewComponentDefinitionModel(resetOscalModel.ComponentDefinition)
		resetModel.Open(height, width)
		return resetModel
	}

	err := testhelpers.RunTestModelView(t, model, reset, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestEditViewComponentDefinitionModel tests that the editing views of the component definition model are correct
func TestDetailValidationViewComponentDefinitionModel(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefValidations)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight},                     // Select component
		tea.KeyMsg{Type: tea.KeyRight},                     // Select framework
		tea.KeyMsg{Type: tea.KeyRight},                     // Select control
		tea.KeyMsg{Type: tea.KeyEnter},                     // Open control
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to remarks
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to description
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to validations
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, // Detail validation
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestEditComponentDefinitionModel tests that a component definition model can be modified, written, and re-read
func TestEditComponentDefinitionModel(t *testing.T) {
	tempOscalFile := testhelpers.CreateTempFile(t, "yaml")
	defer os.Remove(tempOscalFile.Name())

	oscalModel := testhelpers.OscalFromPath(t, validCompDef)
	model := component.NewComponentDefinitionModel(oscalModel.ComponentDefinition)

	testControlId := "ac-1"
	testRemarks := "test remarks"
	testDescription := "test description"

	model.TestSetSelectedControl(testControlId)
	model.UpdateRemarks(testRemarks)
	model.UpdateDescription(testDescription)

	// Create OSCAL model
	mdl := &oscalTypes.OscalCompleteSchema{
		ComponentDefinition: model.GetComponentDefinition(),
	}

	// Write the model to a temp file
	err := oscal.OverwriteOscalModel(tempOscalFile.Name(), mdl)
	require.NoError(t, err)

	// Read the model from the temp file
	modifiedOscalModel := testhelpers.OscalFromPath(t, tempOscalFile.Name())
	compDefn := modifiedOscalModel.ComponentDefinition
	require.NotNil(t, compDefn)
	for _, c := range *compDefn.Components {
		require.NotNil(t, c.ControlImplementations)
		for _, f := range *c.ControlImplementations {
			for _, r := range f.ImplementedRequirements {
				if r.ControlId == testControlId {
					assert.Equal(t, testRemarks, r.Remarks)
					assert.Equal(t, testDescription, r.Description)
				}
			}
		}
	}
}
