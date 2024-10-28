package component_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
)

func TestGetComponents(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefMultiValidations)
	components := component.GetComponents(oscalModel.ComponentDefinition)

	// Count controls
	require.Equal(t, 2, len(components))

	// Components should be sorted deterministically
	assert.Equal(t, "Component A", components[0].Name)

	// Count frameworks - rev4, rev5, source rev4 url, source rev5 url
	require.Equal(t, 4, len(components[0].Frameworks))

	// Frameworks should be sorted deterministically
	assert.Equal(t, "rev4", components[0].Frameworks[2].Name)

	// Check controls in framework
	assert.Equal(t, 3, len(components[0].Frameworks[2].Controls))

	// Check control ac-1 links to a validation
	for _, c := range components[0].Frameworks[2].Controls {
		if c.Name == "ac-1" {
			assert.Equal(t, 2, len(c.Validations))
		}
	}
}
