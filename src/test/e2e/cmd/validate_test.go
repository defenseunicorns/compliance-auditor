package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

const (
	validInputFile    = "../../unit/common/oscal/valid-component.yaml"
	invalidOutputFile = "../../unit/common/validation/validation.opa.yaml"
)

func TestValidateCommand(t *testing.T) {

	message.NoProgress = true

	test := func(t *testing.T, args ...string) error {
		rootCmd := validate.ValidateCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := validate.ValidateCommand()

		return runCmdTestWithGolden(t, "validate/", goldenFileName, rootCmd, args...)
	}

	t.Run("Validate command", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		err := test(t, "-f", validInputFile, "-o", outputFile)

		require.NoError(t, err)

		// Check that the output file is valid OSCAL
		compiledBytes, err := os.ReadFile(outputFile)
		require.NoErrorf(t, err, "error reading assessment results file: %v", err)

		compiledModel, err := oscal.NewOscalModel(compiledBytes)
		require.NoErrorf(t, err, "error creating oscal model from assessment results: %v", err)

		require.NotNilf(t, compiledModel.AssessmentResults, "assessment results is nil")

		require.Equalf(t, 1, len(compiledModel.AssessmentResults.Results), "expected 1 result, got %d", len(compiledModel.AssessmentResults.Results))
	})

	t.Run("Validate with invalid input file - error", func(t *testing.T) {
		err := test(t, "-f", "invalid-file.yaml")
		require.ErrorContains(t, err, "error creating new composer")
	})

	t.Run("Validate with invalid output file - error", func(t *testing.T) {
		err := test(t, "-f", validInputFile, "-o", invalidOutputFile)
		require.ErrorContains(t, err, "invalid OSCAL model at output")
	})

	t.Run("Validate with invalid target - error", func(t *testing.T) {
		err := test(t, "-f", validInputFile, "-t", "invalid-target")
		require.ErrorContains(t, err, "error validating on path")
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

}
