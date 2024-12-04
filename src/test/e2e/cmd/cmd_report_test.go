package cmd_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/report"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/stretchr/testify/require"
)

func TestLulaReportValidComponent(t *testing.T) {

	message.NoProgress = true

	// Helper function to test against golden files
	test := func(t *testing.T, goldenFileName string, args ...string) error {
		t.Helper()
		rootCmd := report.ReportCommand()
		return runCmdTest(t, rootCmd, args...)
	}

	t.Run("Valid YAML Report", func(t *testing.T) {
		err := test(t, "report_valid-multi-component-validations-yaml",
			"-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "yaml")
		require.NoError(t, err)
	})

	t.Run("Valid JSON Report", func(t *testing.T) {
		err := test(t, "report_valid-multi-component-validations-json",
			"-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "json")
		require.NoError(t, err)
	})

	t.Run("Valid TABLE Report", func(t *testing.T) {
		err := test(t, "report_valid-multi-component-validations",
			"-f", "../../unit/common/oscal/valid-multi-component-validations.yaml", "--file-format", "table")
		require.NoError(t, err)
	})

	t.Run("Unsupported OSCAL Report Model", func(t *testing.T) {
		err := test(t, "report_valid-multi-component-validations",
			"-f", "../../unit/common/oscal/catalog.yaml", "--file-format", "table")
		require.Error(t, err)
	})

	t.Run("invalid - file does not exist", func(t *testing.T) {
		err := test(t, "report_valid-multi-component-validations",
			"-f", "file-does-not-exist.yaml", "--file-format", "table")
		require.Error(t, err)
	})

	t.Run("Help Output", func(t *testing.T) {
		err := test(t, "report_help", "--help")
		require.NoError(t, err)
	})
}
