package component_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
)

func TestGetComponents(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validCompDefMultiValidations)
	components := component.GetComponents(oscalModel.ComponentDefinition)

	// Count controls
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Components should be sorted deterministically
	if components[0].Name != "Component A" {
		t.Errorf("Expected component title to be Component A, got %s", components[0].Name)
	}

	// Count frameworks - rev4, rev5, source rev4 url, source rev5 url
	if len(components[0].Frameworks) != 4 {
		t.Errorf("Expected 4 frameworks, got %d", len(components[0].Frameworks))
	}

	// Frameworks should be sorted deterministically
	if components[0].Frameworks[2].Name != "rev4" {
		t.Errorf("Expected framework title to be rev4, got %s", components[0].Frameworks[0].Name)
	}

	// Check controls in framework
	if len(components[0].Frameworks[2].Controls) != 3 {
		t.Errorf("Expected 3 controls, got %d", len(components[0].Frameworks[2].Controls))
	}

	// Check control ac-1 links to a validation
	for _, c := range components[0].Frameworks[2].Controls {
		if c.Name == "ac-1" {
			if len(c.Validations) != 2 {
				t.Errorf("Expected 2 validations in AC-1, got %d", len(c.Validations))
			}
		}
	}
}
