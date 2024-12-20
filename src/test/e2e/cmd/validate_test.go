package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
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

		err := test(t, "-f", validInputFile, "-o", outputFile, "--save-resources")

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

	t.Run("Validate with valid oscal containing no control implementations - error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/oscal/valid-component-no-implementations.yaml")
		require.ErrorContains(t, err, "no control implementations found in component definition")
	})

	t.Run("Validate run tests", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		err := test(t, "-f", "./testdata/validate/component-composed.yaml", "-o", outputFile, "--run-tests")
		require.NoError(t, err)

		// Find all test-results-*.md files in tempDir
		files, err := os.ReadDir(tempDir)
		require.NoError(t, err)
		var testResultsFiles []string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "test-results-") {
				testResultsFiles = append(testResultsFiles, file.Name())
			}
		}

		// Should be only one test results file
		require.Equal(t, 1, len(testResultsFiles))

		// Check that the test results files contain the expected content
		filePath := filepath.Join(tempDir, testResultsFiles[0])
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)

		// Marshal into tests reports map
		var testReportsMap map[string]types.LulaValidationTestReport
		err = yaml.Unmarshal(data, &testReportsMap)
		require.NoError(t, err)

		// Check that the test reports map contains the expected number of test results
		assert.Equal(t, 2, len(testReportsMap))

		// Check that a test result contains the expected content
		testReport, ok := testReportsMap["82099492-0601-4287-a2d1-cc94c49dca9b"]
		require.True(t, ok)
		assert.Equal(t, 2, len(testReport.TestResults))
		assert.True(t, testReport.TestResults[0].Pass)
		assert.True(t, testReport.TestResults[1].Pass)
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

}
