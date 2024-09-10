package test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestTemplateCommand(t *testing.T) {

	test := func(t *testing.T, expectedError bool, args ...string) (string, error) {
		t.Helper()

		cmd := cmd.RootCommand()
		cmdArgs := []string{"tools", "template"}
		cmdArgs = append(cmdArgs, args...)

		output, err := util.ExecuteCommand(cmd, cmdArgs...)
		if err != nil && !expectedError {
			t.Fatal("erroring in the error location")
		}

		return output, err
	}

	t.Run("Test --help", func(t *testing.T) {
		out, _ := test(t, false, "--help")

		if !strings.Contains(out, "Resolving templated artifacts with configuration data") {
			t.Fatalf("Expected help string")
		}
	})

	// Tests that execute unhappy-paths will hit a fatal message which exits the runtime
	// TODO: review RunE command execution and ensure we don't prematurely exit where errors would still be valuable
	// t.Run("Test non-existent file", func(t *testing.T) {
	// 	out, _ := test(t, true, "-f", "non-existent.yaml")

	// 	if !strings.Contains(out, "Path: non-existent.yaml does not exist - unable to digest document") {
	// 		t.Fatalf("Expected error with unable to digest document error")
	// 	}
	// })

	t.Run("Template Valid File", func(t *testing.T) {
		test(t, false, "-f", "../../unit/common/oscal/valid-component-template.yaml", "-o", "valid.yaml")

		templated, err := os.ReadFile("valid.yaml")
		if err != nil {
			t.Fatal(err)
		}

		valid, err := os.ReadFile("../../unit/common/oscal/valid-component.yaml")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(templated, valid) {
			t.Fatalf("Expected: \n%s\n - Got \n%s\n", valid, templated)
		}

		// cleanup
		os.Remove("valid.yaml")

	})

}
