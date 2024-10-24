package dev_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
)

func TestPrintValidation(t *testing.T) {
	t.Parallel()

	oscalAssessmentModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-assessment-results-with-resources.yaml")
	assessment := oscalAssessmentModel.AssessmentResults
	require.NotNil(t, assessment)

	oscalComponentModel := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-multi-component-validations.yaml")
	component := oscalComponentModel.ComponentDefinition
	require.NotNil(t, component)

	t.Run("Test print validation", func(t *testing.T) {
		tmpFile := testhelpers.CreateTempFile(t, ".json")
		defer os.Remove(tmpFile.Name())

		err := dev.PrintValidation(component, assessment, "92cb3cad-bbcd-431a-aaa9-cd47275a3982", tmpFile.Name())
		require.NoError(t, err)

		// get printed data
		printedData, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)

		// get actual data
		validationData, err := os.ReadFile("../../test/unit/common/validation/validation.resource-print.yaml")
		require.NoError(t, err)

		require.Equal(t, validationData, printedData)
	})

	t.Run("Test print validation with no validation prop", func(t *testing.T) {
		err := dev.PrintValidation(component, assessment, "e1ca2968-8652-41be-a19f-c32bc0b3086c", "")
		require.ErrorContains(t, err, "no validation linked to observation")
	})

	t.Run("Test print resources with validation not in backmatter", func(t *testing.T) {
		err := dev.PrintValidation(component, assessment, "af060637-2899-4f26-ae9d-2c1bbbddc4b0", "")
		require.ErrorContains(t, err, "validation not found in component definition")
	})

}
