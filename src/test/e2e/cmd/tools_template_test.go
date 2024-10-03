package cmd_test

import (
	"os"

	"testing"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/stretchr/testify/require"
)

// var updateGolden = flag.Bool("update", false, "update golden files")

func TestToolsTemplateCommand(t *testing.T) {

	test := func(t *testing.T, args ...string) error {
		rootCmd := tools.TemplateCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := tools.TemplateCommand()

		return runCmdTestWithGolden(t, "tools/template/", goldenFileName, rootCmd, args...)
	}

	t.Run("Template Validation", func(t *testing.T) {
		err := testAgainstGolden(t, "validation", "-f", "../../unit/common/validation/validation.tmpl.yaml")
		require.NoError(t, err)
	})

	t.Run("Template Validation with env vars", func(t *testing.T) {
		os.Setenv("LULA_VAR_SOME_ENV_VAR", "my-env-var")
		defer os.Unsetenv("LULA_VAR_SOME_ENV_VAR")
		err := testAgainstGolden(t, "validation_with_env_vars", "-f", "../../unit/common/validation/validation.tmpl.yaml")
		require.NoError(t, err)
	})

	t.Run("Template Validation with set", func(t *testing.T) {
		err := testAgainstGolden(t, "validation_with_set", "-f", "../../unit/common/validation/validation.tmpl.yaml", "--set", ".const.resources.name=foo")
		require.NoError(t, err)
	})

	t.Run("Template Validation for all", func(t *testing.T) {
		os.Setenv("LULA_VAR_SOME_LULA_SECRET", "env-secret")
		defer os.Unsetenv("LULA_VAR_SOME_LULA_SECRET")
		err := testAgainstGolden(t, "validation_all", "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "all")
		require.NoError(t, err)
	})

	t.Run("Template Validation for non-sensitive", func(t *testing.T) {
		err := testAgainstGolden(t, "validation_non_sensitive", "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "non-sensitive")
		require.NoError(t, err)
	})

	t.Run("Template Validation for constants", func(t *testing.T) {
		err := testAgainstGolden(t, "validation_constants", "-f", "../../unit/common/validation/validation.tmpl.yaml", "--render", "constants")
		require.NoError(t, err)
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

	t.Run("Template Validation - invalid file error", func(t *testing.T) {
		err := test(t, "-f", "not-a-file.yaml")
		require.ErrorContains(t, err, "error reading file")
	})

	t.Run("Template Validation - invalid set opts", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/validation/validation.tmpl.yaml", "--set", "not-valid")
		require.ErrorContains(t, err, "error parsing template overrides")
	})

	t.Run("Template Validation - invalid file schema error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/validation/validation.bad.tmpl.yaml")
		require.ErrorContains(t, err, "error rendering template")
	})
}
