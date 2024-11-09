package cmd_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/stretchr/testify/require"
)

func TestDevGetResourcesCommand(t *testing.T) {

	// test := func(t *testing.T, args ...string) error {
	// 	t.Helper()
	// 	rootCmd := dev.DevGetResourcesCommand()

	// 	return runCmdTest(t, rootCmd, args...)
	// }

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := dev.DevGetResourcesCommand()

		return runCmdTestWithGolden(t, "dev/get-resources/", goldenFileName, rootCmd, args...)
	}

	// t.Run("Valid multi validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/multi.validation.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	// t.Run("Valid OPA validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/opa.validation.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	// t.Run("Valid Kyverno validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/validation.kyverno.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	// t.Run("Invalid OPA validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/invalid.opa.validation.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.ErrorContains(t, err, "the following files failed linting")
	// })

	// t.Run("valid template OPA validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/opa.validation.tpl.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	// t.Run("Multiple files", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/validation.kyverno.yaml",
	// 		"--input-files", "./testdata/dev/lint/opa.validation.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

}
