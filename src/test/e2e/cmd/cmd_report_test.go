package cmd_test

import (
	"testing"
	"strings"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/stretchr/testify/require"
	"github.com/defenseunicorns/lula/src/test/util"
)

// //TestLulaReportValidComponent checks that the 'lula report' command works with a valid component definition with multiple components and framework prop.
// func TestLulaReportValidComponent(t *testing.T) {
// 	// Disable progress indicators and other extra formatting
// 	message.NoProgress = true

// 	// Setup the root command and arguments
// 	rootCmd := cmd.RootCommand()
// 	cmdArgs := []string{"report", "-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "table"}

// 	// Run the command and compare output with golden file
// 	err := runCmdTestWithGolden(t, "report", "report_valid-multi-component-validations.golden", rootCmd, cmdArgs...)

// 	// Check for errors in command execution.
// 	if err != nil {
// 		t.Fatalf("Error executing `lula report` with valid component definition: %v", err)
// 	}
// }

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

	// t.Run("Invalid YAML File Format", func(t *testing.T) {
	// 	err := testAgainstGoldenWithErrorCheck(t, "report_valid-multi-component-validations-yaml",
	// 		"not-a-file.json", "report", "-f", "--file-format", "yaml")
	// 	require.ErrorContains(t, err, "failed to fetch data from URL") // Adjust error message as appropriate
	// })

	// t.Run("Invalid JSON File Format", func(t *testing.T) {
	// 	err := testAgainstGoldenWithErrorCheck(t, "report_valid-multi-component-validations-yaml",
	// 		"not-a-file.json", "report", "-f", "--file-format", "json")
	// 	require.ErrorContains(t, err, "failed to fetch data from URL") // Adjust error message as appropriate
	// })

	t.Run("Help Output", func(t *testing.T) {
		err := testAgainstGolden(t, "report_help", "--help")
		require.NoError(t, err)
	})
}
