package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

func TestToolsComposeCommand(t *testing.T) {
	message.NoProgress = true

	test := func(t *testing.T, args ...string) error {
		rootCmd := tools.ComposeCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := tools.ComposeCommand()

		return runCmdTestWithGolden(t, "tools/compose/", goldenFileName, rootCmd, args...)
	}

	testAgainstOutputFile := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := tools.ComposeCommand()

		return runCmdTestWithOutputFile(t, "tools/compose/", goldenFileName, "yaml", rootCmd, args...)
	}

	t.Run("Compose Validation", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		err := test(t, "composed-file",
			"-f", "../../unit/common/composition/component-definition-import-compdefs.yaml",
			"-o", outputFile,
		)

		require.NoError(t, err)

		// Check that the output file is valid OSCAL
		compiledBytes, err := os.ReadFile(outputFile)
		require.NoErrorf(t, err, "error reading composed component definition: %v", err)

		compiledModel, err := oscal.NewOscalModel(compiledBytes)
		require.NoErrorf(t, err, "error creating oscal model from composed component definition: %v", err)

		require.NotNilf(t, compiledModel.ComponentDefinition, "composed component definition is nil")

		require.Equalf(t, 3, len(*compiledModel.ComponentDefinition.BackMatter.Resources), "expected 3 resources, got %d", len(*compiledModel.ComponentDefinition.BackMatter.Resources))
	})

	t.Run("Compose Validation with templating - all", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated",
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations")
		require.NoError(t, err)
	})

	t.Run("Compose Validation with templating - non-sensitive", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-non-sensitive",
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "non-sensitive",
			"--render-validations")
		require.NoError(t, err)
	})

	t.Run("Compose Validation with templating - constants", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-constants",
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "constants",
			"--render-validations")
		require.NoError(t, err)
	})

	t.Run("Compose Validation with templating - masked", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-masked",
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "masked",
			"--render-validations")
		require.NoError(t, err)
	})

	t.Run("Compose Validation with templating and overrides", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-overrides",
			"-f", "../../unit/common/composition/component-definition-template.yaml",
			"-r", "all",
			"--render-validations",
			"--set", ".const.resources.name=foo,.var.some_lula_secret=my-secret")
		require.NoError(t, err)
	})

	t.Run("Compose Validation with no templating on validations for valid validation template", func(t *testing.T) {
		err := testAgainstOutputFile(t, "composed-file-templated-no-validation-templated-valid",
			"-f", "../../unit/common/composition/component-definition-template-valid-validation-tmpl.yaml",
			"-r", "all")
		require.NoError(t, err)
	})

	t.Run("Test help", func(t *testing.T) {
		err := testAgainstGolden(t, "help", "--help")
		require.NoError(t, err)
	})

	t.Run("Test Compose - invalid file error", func(t *testing.T) {
		err := test(t, "-f", "not-a-file.yaml")
		require.ErrorContains(t, err, "error creating composition context")
	})

	t.Run("Test Compose - invalid file schema error", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/composition/component-definition-template.yaml")
		require.ErrorContains(t, err, "error composing model from path")
	})

	t.Run("Test Compose - invalid output file", func(t *testing.T) {
		err := test(t, "-f", "../../unit/common/composition/component-definition-multi.yaml", "-o", "../../unit/common/validation/validation.opa.yaml")
		require.ErrorContains(t, err, "invalid OSCAL model at output file")
	})
}
