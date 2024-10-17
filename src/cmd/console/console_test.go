package console_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/console"
)

func TestGetModelsByFiles(t *testing.T) {
	t.Run("Get output models by files", func(t *testing.T) {
		setOutputFiles := make(map[string]string)
		inputFiles := []string{"../../test/unit/common/oscal/valid-component.yaml", "../../test/unit/common/oscal/valid-generated-component.yaml", "../../test/unit/common/oscal/valid-assessment-results.yaml"}
		models, modelFiles, err := console.GetModelsByFiles(inputFiles, setOutputFiles)
		require.NoError(t, err)

		require.Len(t, models, 2)
		require.Len(t, modelFiles, 2)

		require.NotNil(t, models["component"].ComponentDefinition)
		require.NotNil(t, models["assessment-results"].AssessmentResults)

		require.Equal(t, modelFiles["component"], "component.yaml")
		require.Equal(t, modelFiles["assessment-results"], "../../test/unit/common/oscal/valid-assessment-results.yaml")
	})

	t.Run("Override output model files", func(t *testing.T) {
		setOutputFiles := make(map[string]string)
		setOutputFiles["component"] = "component-override.yaml"

		inputFiles := []string{"../../test/unit/common/oscal/valid-component.yaml", "../../test/unit/common/oscal/valid-generated-component.yaml", "../../test/unit/common/oscal/valid-assessment-results.yaml"}
		_, modelFiles, err := console.GetModelsByFiles(inputFiles, setOutputFiles)
		require.NoError(t, err)

		require.Equal(t, modelFiles["component"], "component-override.yaml")
		require.Equal(t, modelFiles["assessment-results"], "../../test/unit/common/oscal/valid-assessment-results.yaml")
	})
}
