package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/generate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestGenerateSSPCommand(t *testing.T) {

	test := func(t *testing.T, args ...string) error {
		t.Helper()
		rootCmd := generate.GenerateSSPCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := generate.GenerateSSPCommand()

		return runCmdTestWithGolden(t, "generate/", goldenFileName, rootCmd, args...)
	}

	t.Run("Generate SSP", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		args := []string{
			"--profile", "../../unit/common/oscal/valid-profile.yaml",
			"-o", outputFile,
		}
		err := test(t, args...)
		require.NoError(t, err, "executing lula generate ssp %v resulted in an error\n", args)

		// Check that the output file is valid OSCAL
		compiledBytes, err := os.ReadFile(outputFile)
		require.NoError(t, err, "error reading generated ssp")

		ssp := oscal.NewSystemSecurityPlan()

		// Create the new ssp object
		err = ssp.NewModel(compiledBytes)
		require.NoError(t, err, "error creating oscal model from ssp artifact")

		complete := ssp.GetCompleteModel()
		sspModel := complete.SystemSecurityPlan

		require.NotNil(t, sspModel, "expected the SystemSecurityPlan model to be non-nil")
		assert.Equal(t, 3, len(sspModel.ControlImplementation.ImplementedRequirements), "expected 3 controls")
	})

	t.Run("Generate SSP with components", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		// For the components, the valid-generated-component should have no source<>control intersections to the profile, but the
		// valid-multi-component-validations should have 3 controls and one component providing them
		args := []string{
			"--profile", "../../unit/common/oscal/valid-profile-remote-rev4.yaml",
			"-o", outputFile,
			"-c", "../../unit/common/oscal/valid-generated-component.yaml,../../unit/common/oscal/valid-multi-component-validations.yaml",
		}
		err := test(t, args...)
		require.NoError(t, err, "executing lula generate ssp %v resulted in an error\n", args)

		// Check that the output file is valid OSCAL
		compiledBytes, err := os.ReadFile(outputFile)
		require.NoError(t, err, "error reading generated ssp")

		ssp := oscal.NewSystemSecurityPlan()

		// Create the new ssp object
		err = ssp.NewModel(compiledBytes)
		require.NoError(t, err, "error creating oscal model from ssp artifact")

		complete := ssp.GetCompleteModel()
		sspModel := complete.SystemSecurityPlan

		require.NotNil(t, sspModel, "expected the SystemSecurityPlan model to be non-nil")
		assert.Equal(t, 3, len(sspModel.ControlImplementation.ImplementedRequirements), "expected 3 controls")

		for _, ir := range sspModel.ControlImplementation.ImplementedRequirements {
			require.NotNil(t, ir.ByComponents)
			assert.Equal(t, 1, len(*ir.ByComponents), "expected 1 component")
		}

		assert.Equal(t, "7c02500a-6e33-44e0-82ee-fba0f5ea0cae", sspModel.SystemImplementation.Components[0].UUID)
	})

	t.Run("Generate SSP on existing SSP", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		// Generate the SSP with just profile
		args := []string{
			"--profile", "../../unit/common/oscal/valid-profile-remote-rev4.yaml",
			"-o", outputFile,
		}
		err := test(t, args...)
		require.NoError(t, err, "executing lula generate ssp %v resulted in an error\n", args)

		// Re-generate the SSP with components
		args = []string{
			"--profile", "../../unit/common/oscal/valid-profile-remote-rev4.yaml",
			"-o", outputFile,
			"-c", "../../unit/common/oscal/valid-multi-component-validations.yaml",
		}

		err = test(t, args...)
		require.NoError(t, err, "executing re-generation of ssp %v resulted in an error\n", args)

		// Check output content is expected
		compiledBytes, err := os.ReadFile(outputFile)
		require.NoError(t, err, "error reading generated ssp")
		ssp := oscal.NewSystemSecurityPlan()

		// Create the new ssp object
		err = ssp.NewModel(compiledBytes)
		require.NoError(t, err, "error creating oscal model from ssp artifact")

		complete := ssp.GetCompleteModel()
		sspModel := complete.SystemSecurityPlan

		require.NotNil(t, sspModel, "expected the SystemSecurityPlan model to be non-nil")
		assert.Equal(t, 3, len(sspModel.ControlImplementation.ImplementedRequirements), "expected 3 controls")

		for _, ir := range sspModel.ControlImplementation.ImplementedRequirements {
			require.NotNil(t, ir.ByComponents)
			assert.Equal(t, 1, len(*ir.ByComponents), "expected 1 component")
		}

		assert.Equal(t, "7c02500a-6e33-44e0-82ee-fba0f5ea0cae", sspModel.SystemImplementation.Components[0].UUID)
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "ssp-help", "--help")
		require.NoError(t, err, "expected help message")
	})
}
