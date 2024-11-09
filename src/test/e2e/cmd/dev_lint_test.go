package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/stretchr/testify/require"
)

func TestDevLintCommand(t *testing.T) {

	test := func(t *testing.T, args ...string) error {
		t.Helper()
		rootCmd := dev.DevLintCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := dev.DevLintCommand()

		return runCmdTestWithGolden(t, "dev/lint/", goldenFileName, rootCmd, args...)
	}

	parseValidationResultList := func(t *testing.T, filePath string) ([]oscalValidation.ValidationResult, error) {
		t.Helper()
		result := make(map[string][]oscalValidation.ValidationResult)

		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return result["results"], err
		}

		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return result["results"], err
		}

		return result["results"], err
	}

	parseValidationResult := func(t *testing.T, filePath string) (oscalValidation.ValidationResult, error) {
		t.Helper()
		var result oscalValidation.ValidationResult

		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return result, err
		}

		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return result, err
		}

		return result, err
	}

	t.Run("Valid multi validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/multi.validation.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		results, err := parseValidationResultList(t, outputFile)
		require.NoError(t, err)

		expected := []bool{true, true}

		for i, result := range results {
			require.Equal(t, result.Valid, expected[i])
		}
	})

	t.Run("Valid OPA validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/opa.validation.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		result, err := parseValidationResult(t, outputFile)
		require.NoError(t, err)
		require.Equal(t, result.Valid, true)

	})

	t.Run("Valid Kyverno validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/validation.kyverno.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		result, err := parseValidationResult(t, outputFile)
		require.NoError(t, err)
		require.Equal(t, result.Valid, true)
	})

	t.Run("Invalid OPA validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/invalid.opa.validation.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.ErrorContains(t, err, "the following files failed linting")

		result, err := parseValidationResult(t, outputFile)
		require.NoError(t, err)
		require.Equal(t, result.Valid, false)
	})

	t.Run("valid template OPA validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/opa.validation.tpl.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		result, err := parseValidationResult(t, outputFile)
		require.NoError(t, err)
		require.Equal(t, result.Valid, true)
	})

	t.Run("Multiple files", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/validation.kyverno.yaml",
			"--input-files", "./testdata/dev/lint/opa.validation.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		results, err := parseValidationResultList(t, outputFile)
		require.NoError(t, err)

		expected := []bool{true, true}

		for i, result := range results {
			require.Equal(t, result.Valid, expected[i])
		}
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

}
