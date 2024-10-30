package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

func TestDevPrintResourcesCommand(t *testing.T) {
	message.NoProgress = true

	test := func(t *testing.T, args ...string) error {
		rootCmd := dev.PrintCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := dev.PrintCommand()

		return runCmdTestWithGolden(t, "dev/print/", goldenFileName, rootCmd, args...)
	}

	t.Run("Print Resources", func(t *testing.T) {
		err := testAgainstGolden(t, "resources", "--resources",
			"-a", "../../unit/common/oscal/valid-assessment-results-with-resources.yaml",
			"-u", "92cb3cad-bbcd-431a-aaa9-cd47275a3982",
		)
		require.NoError(t, err)
	})

	t.Run("Print Resources - invalid oscal", func(t *testing.T) {
		err := test(t, "--resources",
			"-a", "../../unit/common/validation/validation.opa.yaml",
			"-u", "92cb3cad-bbcd-431a-aaa9-cd47275a3982",
		)
		require.ErrorContains(t, err, "error creating oscal assessment results model")
	})

	t.Run("Print Resources - no uuid", func(t *testing.T) {
		err := test(t, "--resources",
			"-a", "../../unit/common/oscal/valid-assessment-results-with-resources.yaml",
			"-u", "foo",
		)
		require.ErrorContains(t, err, "error printing resources")
	})

	t.Run("Print Validation", func(t *testing.T) {
		err := testAgainstGolden(t, "validation", "--validation",
			"-a", "../../unit/common/oscal/valid-assessment-results-with-resources.yaml",
			"-c", "../../unit/common/oscal/valid-multi-component-validations.yaml",
			"-u", "92cb3cad-bbcd-431a-aaa9-cd47275a3982",
		)
		require.NoError(t, err)
	})

	t.Run("Print Validation - invalid assessment oscal", func(t *testing.T) {
		err := test(t, "--validation",
			"-a", "../../unit/common/validation/validation.opa.yaml",
			"-c", "../../unit/common/oscal/valid-multi-component-validations.yaml",
			"-u", "92cb3cad-bbcd-431a-aaa9-cd47275a3982",
		)
		require.ErrorContains(t, err, "error creating oscal assessment results model")
	})

	t.Run("Print Validation - no uuid", func(t *testing.T) {
		err := test(t, "--validation",
			"-a", "../../unit/common/oscal/valid-assessment-results-with-resources.yaml",
			"-c", "../../unit/common/oscal/valid-multi-component-validations.yaml",
			"-u", "foo",
		)
		require.ErrorContains(t, err, "error printing validation")
	})
}
