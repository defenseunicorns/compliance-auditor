package cmd_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/tools"
)

func TestToolsComposeCommand(t *testing.T) {

	test := func(t *testing.T, goldenFileName string, expectError bool, args ...string) error {
		rootCmd := tools.ComposeCommand()

		return runCmdTest(t, "tools/compose/", goldenFileName, expectError, rootCmd, args...)
	}

	testAgainstOutputFile := func(t *testing.T, goldenFileName string, expectError bool, args ...string) error {
		rootCmd := tools.ComposeCommand()

		return runCmdTestWithOutputFile(t, "tools/compose/", goldenFileName, "yaml", expectError, rootCmd, args...)
	}

	t.Run("Compose Validation", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file", false,
			"-f", "../../unit/common/composition/component-definition-all-local.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Compose Validation with templating", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Compose Validation with templating and overrides", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-overrides", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations",
			"--set", ".const.resources.name=foo,.var.some_lula_secret=my-secret")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Compose Validation with no templating on validations", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-no-validation-templated-missing-validation", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Compose Validation with no templating on validations for valid validation template", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-no-validation-templated-valid", false,
			"-f", "../../unit/common/composition/component-definition-template-valid-validation-tmpl.yaml",
			"-r", "all")
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

	t.Run("Test Compose - invalid file error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "not-a-file.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test Compose - invalid file schema error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "../../unit/common/composition/component-definition-template.yaml")
		if err != nil {
			t.Fatal(err)
		}
	})
}
