package cmd_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

func TestToolsComposeCommand(t *testing.T) {
	message.NoProgress = true

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
			t.Error(err)
		}
	})

	t.Run("Compose Validation with templating - all", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Compose Validation with templating - non-sensitive", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-non-sensitive", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "non-sensitive",
			"--render-validations")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Compose Validation with templating - constants", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-constants", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "constants",
			"--render-validations")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Compose Validation with templating - masked", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-masked", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "masked",
			"--render-validations")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Compose Validation with templating and overrides", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-overrides", false,
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations",
			"--set", ".const.resources.name=foo,.var.some_lula_secret=my-secret")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Compose Validation with no templating on validations for valid validation template", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-no-validation-templated-valid", false,
			"-f", "../../unit/common/composition/component-definition-template-valid-validation-tmpl.yaml",
			"-r", "all")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test help", func(t *testing.T) {
		err := test(t, "help", false, "--help")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test Compose - invalid file error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "not-a-file.yaml")
		if err != nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Test Compose - invalid file schema error", func(t *testing.T) {
		err := test(t, "empty", true, "-f", "../../unit/common/composition/component-definition-template.yaml")
		if err != nil {
			t.Error(err)
		}
	})
}
