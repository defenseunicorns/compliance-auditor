package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/stretchr/testify/require"
)

func TestDevGetResourcesCommand(t *testing.T) {

	test := func(t *testing.T, args ...string) error {
		t.Helper()
		rootCmd := dev.DevGetResourcesCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := dev.DevGetResourcesCommand()

		return runCmdTestWithGolden(t, "dev/get-resources/", goldenFileName, rootCmd, args...)
	}

	parseOutput := func(t *testing.T, filePath string) (map[string]interface{}, error) {
		t.Helper()
		result := make(map[string]interface{})

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

	t.Run("Valid validation file", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-file", "./testdata/dev/get-resources/opa.validation.yaml",
			"--output-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		result, err := parseOutput(t, outputFile)
		require.NoError(t, err)
		name := result["pod"].(map[string]interface{})["metadata"].(map[string]interface{})["name"]
		require.Equal(t, name, "test-pod-name")
	})

	t.Run("Valid validation file - template", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.json")

		args := []string{
			"--input-file", "./testdata/dev/get-resources/opa.validation.tpl.yaml",
			"--output-file", outputFile,
		}

		err := test(t, args...)
		require.NoError(t, err)

		result, err := parseOutput(t, outputFile)
		require.NoError(t, err)
		name := result["pod"].(map[string]interface{})["metadata"].(map[string]interface{})["name"]
		require.Equal(t, name, "test-pod-name")
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

}
