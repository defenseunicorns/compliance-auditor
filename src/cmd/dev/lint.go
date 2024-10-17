package dev

import (
	"fmt"
	"strings"

	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/lula/src/cmd/common"
	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type LintFlags struct {
	InputFiles []string // -f --input-files
	ResultFile string   // -r --result-file
}

var lintOpts = &LintFlags{}

var lintHelp = `
To lint existing validation files:
	lula dev lint -f <path1>,<path2>,<path3> [-r <result-file>]
`

var lintCmd = &cobra.Command{
	Use:     "lint",
	Short:   "Lint validation files against schema",
	Long:    "Validate validation files are properly configured against the schema, file paths can be local or URLs (https://)",
	Example: lintHelp,
	Run: func(cmd *cobra.Command, args []string) {
		if len(lintOpts.InputFiles) == 0 {
			message.Fatalf(nil, "No input files specified")
		}

		validationResults := DevLintCommand(lintOpts.InputFiles)

		// If result file is specified, write the validation results to the file
		var err error
		if lintOpts.ResultFile != "" {
			// If there is only one validation result, write it to the file
			if len(validationResults) == 1 {
				err = oscalValidation.WriteValidationResult(validationResults[0], lintOpts.ResultFile)
			} else {
				// If there are multiple validation results, write them to the file
				err = oscalValidation.WriteValidationResults(validationResults, lintOpts.ResultFile)
			}
		}
		if err != nil {
			message.Fatal(err, "Error writing validation results")
		}

		// If there is at least one validation result that is not valid, exit with a fatal error
		failedFiles := []string{}
		for _, result := range validationResults {
			if !result.Valid {
				failedFiles = append(failedFiles, result.Metadata.DocumentPath)
			}
		}
		if len(failedFiles) > 0 {
			message.Fatal(nil, fmt.Sprintf("The following files failed linting: %s", strings.Join(failedFiles, ", ")))
		}
	},
}

func DevLintCommand(inputFiles []string) []oscalValidation.ValidationResult {
	var validationResults []oscalValidation.ValidationResult

	for _, inputFile := range inputFiles {
		var result oscalValidation.ValidationResult
		spinner := message.NewProgressSpinner("Linting %s", inputFile)

		// handleFail is a helper function to handle the case where the validation fails from
		// a non-schema error
		handleFail := func(err error) {
			result = *oscalValidation.NewNonSchemaValidationError(err, &oscalValidation.ValidationParams{ModelType: "validation"})
			validationResults = append(validationResults, result)
			message.WarnErrf(oscalValidation.GetNonSchemaError(&result), "Failed to lint %s, %s", inputFile, oscalValidation.GetNonSchemaError(&result).Error())
			spinner.Stop()
		}

		defer spinner.Stop()

		validationBytes, err := network.Fetch(inputFile)
		if err != nil {
			handleFail(err)
			break
		}

		validations, err := pkgCommon.ReadValidationsFromYaml(validationBytes)
		if err != nil {
			handleFail(err)
			break
		}

		allValid := true
		// Lint each validation in the file
		for _, validation := range validations {
			result = validation.Lint()
			result.Metadata.DocumentPath = inputFile
			validationResults = append(validationResults, result)

			// If any of the validations fail, set allValid to false
			if !result.Valid {
				allValid = false
			}
		}

		if allValid {
			message.Infof("Successfully linted %s", inputFile)
			spinner.Success()
		} else {
			message.Warnf("Validation failed for %s", inputFile)
			spinner.Stop()
		}
	}
	return validationResults
}

func init() {

	common.InitViper()

	devCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringSliceVarP(&lintOpts.InputFiles, "input-files", "f", []string{}, "the paths to validation files (comma-separated)")
	lintCmd.Flags().StringVarP(&lintOpts.ResultFile, "result-file", "r", "", "the path to write the validation result")
}
