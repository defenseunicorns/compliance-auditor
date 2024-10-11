package cmd_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/generate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestGenerateProfileCommand(t *testing.T) {

	test := func(t *testing.T, args ...string) error {
		t.Helper()
		rootCmd := generate.GenerateProfileCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := generate.GenerateProfileCommand()

		return runCmdTestWithGolden(t, "generate/", goldenFileName, rootCmd, args...)
	}

	testAgainstOutputFile := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := generate.GenerateProfileCommand()

		return runCmdTestWithOutputFile(t, "generate/", goldenFileName, "yaml", rootCmd, args...)
	}

	t.Run("Generate Profile", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		args := []string{
			"--source", "../../unit/common/oscal/catalog.yaml",
			"--include", "ac-1,ac-3,ac-2",
			"-o", outputFile,
		}
		err := test(t, args...)
		if err != nil {
			t.Errorf("executing lula generate profile %v resulted in an error\n", args)
		}

		// Check that the output file is valid OSCAL
		compiledBytes, err := os.ReadFile(outputFile)
		if err != nil {
			t.Errorf("error reading generated profile: %v\n", err)
		}

		profile := oscal.NewProfile()

		// Create the new profile object
		err = profile.NewModel(compiledBytes)
		if err != nil {
			t.Errorf("error creating oscal model from profile artifact: %v\n", err)
		}

		complete := profile.GetCompleteModel()
		if complete.Profile == nil {
			t.Error("expected the profile model to be non-nil")
		}

		profileModel := complete.Profile

		if len(profileModel.Imports) == 0 {
			t.Error("expected length of imports to be greater than 0")
		}

		// Target import item should be the only item in the list
		include := profileModel.Imports[0].IncludeControls
		controls := *include

		if len(controls) != 1 {
			t.Error("expected length of controls to be 1")
		}
		expected := []string{"ac-1", "ac-2", "ac-3"}
		ids := controls[0].WithIds
		if !reflect.DeepEqual(expected, *ids) {
			t.Errorf("expected control id slice to contain %+q, got %+q", expected, *ids)
		}
	})

	t.Run("Generate a profile with included controls", func(t *testing.T) {
		args := []string{
			"--source", "../../unit/common/oscal/catalog.yaml",
			"--include", "ac-1,ac-3,ac-2",
		}

		err := testAgainstOutputFile(t, "generate-profile", args...)
		if err != nil {
			t.Errorf("error executing: generate profile %v", strings.Join(args, " "))
		}
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		if err != nil {
			t.Errorf("Expected help message but received an error %v", err)
		}
	})

	t.Run("Test include/exclude mutually exclusive", func(t *testing.T) {
		err := test(t, "--source", "catalog.yaml", "--include", "ac-1", "--exclude", "ac-2")
		if err == nil {
			t.Error("Expected error message for flags being mutually exclusive")
		}
		if !strings.Contains(err.Error(), "none of the others can be") {
			t.Errorf("Expected error for mutually exclusive flags - received %v", err.Error())
		}
	})

}
