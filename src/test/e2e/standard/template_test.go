package test

import (
	"flag"
	"os"
	"path/filepath"

	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/test/util"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestTemplateCommand(t *testing.T) {

	test := func(t *testing.T, goldenFileName string, expectError bool, args ...string) error {
		t.Helper()

		cmdArgs := []string{"tools", "template"}
		cmdArgs = append(cmdArgs, args...)

		cmd := cmd.RootCommand()

		_, output, err := util.ExecuteCommand(cmd, cmdArgs...)
		if err != nil {
			if !expectError {
				return err
			} else {
				return nil
			}
		}

		if !expectError {
			goldenFile := filepath.Join("testdata", goldenFileName+".golden")

			if *updateGolden && !expectError {
				err = os.WriteFile(goldenFile, []byte(output), 0644)
				if err != nil {
					return err
				}
			}

			expected, err := os.ReadFile(goldenFile)
			if err != nil {
				return err
			}

			if output != string(expected) {
				t.Fatalf("Expected:\n%s\n - Got \n%s\n", expected, output)
			}
		}

		return nil
	}

	t.Run("Template Validation", func(t *testing.T) {
		err := test(t, "validation", false, "-f", "../../unit/common/validation/validation.tmpl.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation with env vars", func(t *testing.T) {
		os.Setenv("LULA_VAR_SOME_ENV_VAR", "my-env-var")
		defer os.Unsetenv("LULA_VAR_SOME_ENV_VAR")
		err := test(t, "validation_with_env_vars", false, "-f", "../../unit/common/validation/validation.tmpl.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation with set", func(t *testing.T) {
		err := test(t, "validation_with_set", false, "-f", "../../unit/common/validation/validation.tmpl.yaml", "--set", ".const.resources.name=foo")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation for all", func(t *testing.T) {
		os.Setenv("LULA_VAR_SOME_LULA_SECRET", "env-secret")
		defer os.Unsetenv("LULA_VAR_SOME_LULA_SECRET")
		err := test(t, "validation_all", false, "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "all")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation for non-sensitive", func(t *testing.T) {
		err := test(t, "validation_non_sensitive", false, "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "non-sensitive")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation for constants", func(t *testing.T) {
		err := test(t, "validation_constants", false, "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "constants")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test help", func(t *testing.T) {
		err := test(t, "help", false, "--help")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation - invalid file error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "not-a-file.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Template Validation - invalid file schema error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "../../unit/common/validation/validation.bad.tmpl.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})
}
