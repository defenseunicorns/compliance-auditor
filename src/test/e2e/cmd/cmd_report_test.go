package cmd_test

import (
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	"github.com/stretchr/testify/require"
)

// Helper function to test against golden files
func testAgainstGolden(t *testing.T, goldenFileName string, args ...string) error {
	rootCmd := cmd.RootCommand()
	return runCmdTestWithGolden(t, "report", goldenFileName, rootCmd, args...)
}

func testAgainstGoldenWithErrorCheck(t *testing.T, goldenFileName string, expectedError string, args ...string) error {
	rootCmd := cmd.RootCommand()
	_, _, err := util.ExecuteCommand(rootCmd, args...) // Ignore output

	// If an error is expected, check if it matches the expected error message
	if expectedError != "" {
		if err == nil {
			t.Fatalf("expected error %q but got none", expectedError)
		} else if !strings.Contains(err.Error(), expectedError) {
			t.Fatalf("expected error %q but got %q", expectedError, err.Error())
		}
		return err // Return early as we are only testing error handling here
	}

	// No error is expected, so proceed to compare output with the golden file
	return runCmdTestWithGolden(t, "report", goldenFileName, rootCmd, args...)
}

func TestLulaReportValidComponent2(t *testing.T) {
	// Disable progress indicators
	message.NoProgress = true

	t.Run("Valid YAML Report", func(t *testing.T) {
		err := testAgainstGolden(t, "report_valid-multi-component-validations-yaml.golden",
			"report", "-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "yaml")
		require.NoError(t, err)
	})

	t.Run("Valid JSON Report", func(t *testing.T) {
		err := testAgainstGolden(t, "report_valid-multi-component-validations-json.golden",
			"report", "-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "json")
		require.NoError(t, err)
	})

	t.Run("Valid TABLE Report", func(t *testing.T) {
		err := testAgainstGolden(t, "report_valid-multi-component-validations.golden",
			"report", "-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "table")
		require.NoError(t, err)
	})

	t.Run("Help Output", func(t *testing.T) {
		err := testAgainstGolden(t, "report_help", "--help")
		require.NoError(t, err)
	})
}
