package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type flags struct {
	InputFiles []string // -f --input-files
	ResultFile string   // -r --result-file
	Composed   bool     // -c --composed (default false)
}

var opts = &flags{}

var lintHelp = `
To lint existing OSCAL files:
	lula tools lint -f <path1>,<path2>,<path3>
To lint composed OSCAL models:
	**WARNING** Only use if you are sure that the OSCAL model is already composed, (i.e. it has no imports or remote validations)
	lula tools lint -c true -f <path1>,<path2>,<path3>
`

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Validate OSCAL against schema",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Validate OSCAL documents are properly configured against the OSCAL schema",
		Example: lintHelp,
		Run: func(cmd *cobra.Command, args []string) {
			var validationResults []validation.ValidationResult
			var tmpDir string
			if len(opts.InputFiles) == 0 {
				message.Fatalf(nil, "No input files specified")
			}

			// Create a temporary directory to store the composed OSCAL models
			tmpDir, err := composition.CreateTempDir()
			if err != nil {
				message.Fatalf(err, "Failed to create temporary directory")
			}
			defer os.RemoveAll(tmpDir)

			for _, inputFile := range opts.InputFiles {

				spinner := message.NewProgressSpinner("Linting %s", inputFile)
				defer spinner.Stop()

				// ensure the input file is composed
				composedFile := inputFile
				if !opts.Composed {
					composedFile = path.Join(tmpDir, path.Base(inputFile))

					err := Compose(inputFile, composedFile)
					if err != nil {
						message.Fatalf(err, "Failed to compose %s", inputFile)
					}
				}

				validationResp, err := validation.ValidationCommand(composedFile)
				// fatal for non-validation errors
				if err != nil {
					message.Fatalf(err, "Failed to lint %s: %s", inputFile, err)
				}

				for _, warning := range validationResp.Warnings {
					message.Warn(warning)
				}

				// append the validation result to the results array
				validationResults = append(validationResults, validationResp.Result)

				// If result file is not specified, print the validation result
				if opts.ResultFile == "" {
					jsonBytes, err := json.MarshalIndent(validationResp.Result, "", "  ")
					if err != nil {
						message.Fatalf(err, "Failed to marshal validation result")
					}
					message.Infof("Validation result for %s: %s", inputFile, string(jsonBytes))
				}
				// New conditional for logging success or failed linting
				if validationResp.Result.Valid {
					message.Infof("Successfully validated %s is valid OSCAL version %s %s\n", inputFile, validationResp.Validator.GetSchemaVersion(), validationResp.Validator.GetModelType())
					spinner.Success()
				} else {
					message.WarnErrf(nil, "Failed to lint %s", inputFile)
					spinner.Stop()
				}
			}

			// If result file is specified, write the validation results to the file
			if opts.ResultFile != "" {
				// If there is only one validation result, write it to the file
				if len(validationResults) == 1 {
					validation.WriteValidationResult(validationResults[0], opts.ResultFile)
				} else {
					// If there are multiple validation results, write them to the file
					validation.WriteValidationResults(validationResults, opts.ResultFile)
				}
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

	toolsCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringSliceVarP(&opts.InputFiles, "input-files", "f", []string{}, "the paths to oscal json schema files (comma-separated)")
	lintCmd.Flags().StringVarP(&opts.ResultFile, "result-file", "r", "", "the path to write the validation result")
}
