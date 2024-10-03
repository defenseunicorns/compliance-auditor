package cmd_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/validate"
)

func TestValidateCommand(t *testing.T) {

	testWithGolden := func(t *testing.T, goldenFileName string, expectError bool, args ...string) error {
		rootCmd := validate.ValidateCommand()

		return runCmdTest(t, "validate/", goldenFileName, expectError, rootCmd, args...)
	}

	t.Run("Test help", func(t *testing.T) {
		err := testWithGolden(t, "help", false, "--help")
		if err != nil {
			t.Fatal(err)
		}
	})

	// Add more test cases for various flags... can't really compare any golden files but can test flags
	// I think based on the way the validate command prints logs...
	// I'll only be able to check the component-def-output file -> which is really just wrapped up in compose...
	// TODO: add more tests if/when the messaging/logging is improved
	// t.Run("Test validate command with flags?", func(t *testing.T) {

	// can I do a test for the interactive/confirm execution flags?
}
