package cmd_test

import (
	"path/filepath"
	"testing"

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

	t.Run("Valid multi validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-files", "./testdata/dev/lint/multi.validation.yaml",
			"--result-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)
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
	})

	// t.Run("Invalid OPA validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "./testdata/dev/lint/invalid.opa.validation.yaml",
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

	// t.Run("Remote validation file", func(t *testing.T) {
	// 	tempDir := t.TempDir()
	// 	outputFile := filepath.Join(tempDir, "output.json")

	// 	args := []string{
	// 		"--input-files", "https://raw.githubusercontent.com/defenseunicorns/lula/main/src/test/e2e/cmd/testdata/dev/lint/validation.kyverno.yaml",
	// 		"--result-file", outputFile,
	// 	}

	// 	err := test(t, args...)
	// 	require.NoError(t, err)
	// })

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

	// t.Run("Test include/exclude mutually exclusive", func(t *testing.T) {
	// 	err := test(t, "--source", "catalog.yaml", "--include", "ac-1", "--exclude", "ac-2")
	// 	if err == nil {
	// 		t.Error("Expected error message for flags being mutually exclusive")
	// 	}
	// 	if !strings.Contains(err.Error(), "none of the others can be") {
	// 		t.Errorf("Expected error for mutually exclusive flags - received %v", err.Error())
	// 	}
	// })

}
