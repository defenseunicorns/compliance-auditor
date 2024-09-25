package cmd_test

import (
	"testing"
)

func TestValidateCommand(t *testing.T) {

	test := func(t *testing.T, goldenFileName string, expectError bool, args ...string) error {
		t.Helper()

		cmdArgs := []string{"validate"}
		cmdArgs = append(cmdArgs, args...)

		return runCmdTest(t, "validate/"+goldenFileName, expectError, cmdArgs...)
	}

	t.Run("Test help", func(t *testing.T) {
		err := test(t, "help", false, "--help")
		if err != nil {
			t.Fatal(err)
		}
	})

	// Add more test cases as needed
	// I think based on the way the validate command prints logs, I'll only be able to check the component-def-output file
	// TODO: add more tests if/when the messaging/logging is improved
	// t.Run("Test validate command", func(t *testing.T) {
}
