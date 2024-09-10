package test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestTemplateCommand(t *testing.T) {

	test := func(t *testing.T, expectError bool, args ...string) (string, error) {
		t.Helper()

		cmd := cmd.RootCommand()
		// cmdArgs := []string{"tools", "template"}
		// cmdArgs = append(cmdArgs, args...)

		cmd, output, err := util.ExecuteCommand(cmd, args...)
		if err != nil && !expectError {
			t.Fatal(err)
		}

		cmd.Execute()

		return output, err
	}

	t.Run("Test help", func(t *testing.T) {
		out, _ := test(t, false)

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
		out, _ := test(t, false, "-f", "../../unit/common/oscal/valid-component-template.yaml", "-o", "valid.yaml")

		fmt.Println(out)

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
