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

	test := func(t *testing.T, expectError bool, args ...string) (string, error) {
		t.Helper()

		cmdArgs := []string{"tools", "template"}
		cmdArgs = append(cmdArgs, args...)

		cmd := cmd.RootCommand()

		_, output, err := util.ExecuteCommand(cmd, cmdArgs...)
		if err != nil && !expectError {
			t.Fatal(err)
		}

		return output, err
	}

	t.Run("Template Valid File", func(t *testing.T) {

		_, err := test(t, false, "-f", "../../unit/common/oscal/valid-component-template.yaml", "-o", "valid.yaml")
		defer os.Remove("valid.yaml")
		if err != nil {
			t.Fatal(err)
		}

		// this comparison using golden files would make more sense
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

	})

	t.Run("Test help", func(t *testing.T) {
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

}
