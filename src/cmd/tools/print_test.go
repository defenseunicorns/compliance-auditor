package tools_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/types"
)

const (
	assessmentPath = "../../test/unit/common/oscal/valid-assessment-results-with-resources.yaml"
	componentPath  = "../../test/unit/common/oscal/valid-multi-component-validations.yaml"
)

func TestPrintResources(t *testing.T) {
	t.Parallel()

	oscalModel := testhelpers.OscalFromPath(t, assessmentPath)
	assessment := oscalModel.AssessmentResults
	require.NotNil(t, assessment)

	t.Run("Test print resources", func(t *testing.T) {
		tmpFile := testhelpers.CreateTempFile(t, ".json")
		defer os.Remove(tmpFile.Name())

		err := tools.PrintResources(assessment, "92cb3cad-bbcd-431a-aaa9-cd47275a3982", "../../test/unit/common/oscal", tmpFile.Name())
		require.NoError(t, err)

		// get printed resources
		data, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)

		var obsResources types.DomainResources
		err = json.Unmarshal(data, &obsResources)
		require.NoError(t, err)

		// get actual resources
		data, err = os.ReadFile("../../test/unit/common/resources/valid-resources.json")
		require.NoError(t, err)

		var resources types.DomainResources
		err = json.Unmarshal(data, &resources)
		require.NoError(t, err)

		require.Equal(t, resources, obsResources)
	})

	t.Run("Test print resources with invalid resources", func(t *testing.T) {
		err := tools.PrintResources(assessment, "e1ca2968-8652-41be-a19f-c32bc0b3086c", "../../test/unit/common/oscal", "")
		require.ErrorContains(t, err, "error unmarshalling resource")
	})

	t.Run("Test print resources with no resources", func(t *testing.T) {
		err := tools.PrintResources(assessment, "af060637-2899-4f26-ae9d-2c1bbbddc4b0", "../../test/unit/common/oscal", "")
		require.ErrorContains(t, err, "observation does not contain a remote reference")
	})

}

func TestPrintValidation(t *testing.T) {
	t.Parallel()

	oscalAssessmentModel := testhelpers.OscalFromPath(t, assessmentPath)
	assessment := oscalAssessmentModel.AssessmentResults
	require.NotNil(t, assessment)

	oscalComponentModel := testhelpers.OscalFromPath(t, componentPath)
	component := oscalComponentModel.ComponentDefinition
	require.NotNil(t, component)

	t.Run("Test print validation", func(t *testing.T) {
		tmpFile := testhelpers.CreateTempFile(t, ".json")
		defer os.Remove(tmpFile.Name())

		err := tools.PrintValidation(component, assessment, "92cb3cad-bbcd-431a-aaa9-cd47275a3982", tmpFile.Name())
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
		err := tools.PrintValidation(component, assessment, "e1ca2968-8652-41be-a19f-c32bc0b3086c", "")
		require.ErrorContains(t, err, "no validation linked to observation")
	})

	t.Run("Test print resources with validation not in backmatter", func(t *testing.T) {
		err := tools.PrintValidation(component, assessment, "af060637-2899-4f26-ae9d-2c1bbbddc4b0", "")
		require.ErrorContains(t, err, "validation not found in component definition")
	})

}
