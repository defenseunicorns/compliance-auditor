package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

func TestDevPrintResourcesCommand(t *testing.T) {
	message.NoProgress = true

	test := func(t *testing.T, args ...string) error {
		rootCmd := tools.PrintCommand()

		return runCmdTest(t, rootCmd, args...)
	}

	testAgainstGolden := func(t *testing.T, goldenFileName string, args ...string) error {
		rootCmd := tools.PrintCommand()

		return runCmdTestWithGolden(t, "tools/print/", goldenFileName, rootCmd, args...)
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

	t.Run("Print Validation non-composed component", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := filepath.Join(tempDir, "output.yaml")

		err := test(t, "--validation",
			"-a", "../scenarios/validation-composition/assessment-results.yaml",
			"-c", "../scenarios/validation-composition/component-definition.yaml",
			"-u", "d328a0a1-630b-40a2-9c9d-4818420a4126",
			"-o", outputFile,
		)

		require.NoError(t, err)

		// Check that the output file matches the expected validation
		var validation common.Validation
		validationBytes, err := os.ReadFile(outputFile)
		require.NoErrorf(t, err, "error reading validation output: %v", err)
		err = validation.UnmarshalYaml(validationBytes)
		require.NoErrorf(t, err, "error unmarshalling validation: %v", err)

		var expectedValidation common.Validation
		expectedValidationBytes, err := os.ReadFile("../scenarios/validation-composition/validation.opa.yaml")
		require.NoErrorf(t, err, "error reading expected validation: %v", err)
		err = expectedValidation.UnmarshalYaml(expectedValidationBytes)
		require.NoErrorf(t, err, "error unmarshalling expected validation: %v", err)

		require.Equalf(t, expectedValidation, validation, "expected validation does not match actual validation")
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
