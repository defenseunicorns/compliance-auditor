package component_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/stretchr/testify/require"
)

// TestComponentDefinitionBasicView tests that the model is created correctly from a component definition with validations
func TestValidateInitialView(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefMultiValidations)
	model := component.NewValidateModel(oscalModel.ComponentDefinition)

	msgs := []tea.Msg{
		component.ValidateOpenMsg{
			Height: height,
			Width:  width,
			Target: "rev4",
		},
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}
