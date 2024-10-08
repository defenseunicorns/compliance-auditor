package cmd_test

import (
	"os"
	"path/filepath"
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

	t.Run("Generate Profile", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		args := []string{
			"--source", "../unit/common/oscal/catalog.yaml",
			"--include", "ac-1,ac-2,ac-3",
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

	})

}
