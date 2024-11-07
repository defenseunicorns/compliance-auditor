package dev

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var validateHelp = `
To run validation from a lula validation manifest:
	lula dev validate -f /path/to/validation.yaml
To run validation using a custom resources file:
	lula dev validate -f /path/to/validation.yaml -r /path/to/resources.json
To run validation and automatically confirm execution
	lula dev validate -f /path/to/validation.yaml --confirm-execution
To run validation from stdin:
	cat /path/to/validation.yaml | lula dev validate
To hang indefinitely for stdin:
	lula dev validate -t -1
To hang for timeout of 5 seconds:
	lula dev validate -t 5
`

func DevValidateCommand() *cobra.Command {

	var (
		InputFile        string // -f --input-file
		OutputFile       string // -o --output-file
		Timeout          int    // -t --timeout
		ConfirmExecution bool   // --confirm-execution
		ExpectedResult   bool   // -e --expected-result
		ResourcesFile    string // -r --resources-file
	)

	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Run an individual Lula validation.",
		Long:    "Run an individual Lula validation for quick testing and debugging of a Lula Validation. This command is intended for development purposes only.",
		Example: validateHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinnerMessage := fmt.Sprintf("Validating %s", InputFile)
			spinner := message.NewProgressSpinner("%s", spinnerMessage)
			defer spinner.Stop()

			ctx := context.Background()
			var validationBytes []byte
			var resourcesBytes []byte
			var err error

			// Read the validation data from STDIN or provided file
			validationBytes, err = ReadValidation(cmd, spinner, InputFile, Timeout)
			if err != nil {
				message.Fatalf(err, "error reading validation: %v", err)
			}

			// Reset the spinner message
			spinner.Updatef("%s", spinnerMessage)

			// If a resources file is provided, read the resources file
			if ResourcesFile != "" {
				if !strings.HasSuffix(ResourcesFile, ".json") {
					message.Fatalf(fmt.Errorf("resource file must be a json file"), "resource file must be a json file")
				} else {
					// Read the resources data
					resourcesBytes, err = pkgCommon.ReadFileToBytes(ResourcesFile)
					if err != nil {
						message.Fatalf(err, "error reading file: %v", err)
					}
				}
			}

			config, _ := cmd.Flags().GetStringSlice("set")
			message.Debug("command line 'set' flags: %s", config)

			output, err := DevTemplate(validationBytes, config)
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}

			// add to debug logs accepting that this will print sensitive information?
			message.Debug(string(output))

			validation, err := DevValidate(ctx, output, resourcesBytes, ConfirmExecution, spinner)
			if err != nil {
				message.Fatalf(err, "error running dev validate: %v", err)
			}

			// Write the validation result to a file if an output file is provided
			// Otherwise, print the result to the debug console
			err = writeValidation(validation, OutputFile)
			if err != nil {
				message.Fatalf(err, "error writing result: %v", err)
			}

			// Print observations if there are any
			if len(validation.Result.Observations) > 0 {
				message.Infof("Observations:")
				for key, observation := range validation.Result.Observations {
					message.Infof("--> %s: %s", key, observation)
				}
			}

			result := validation.Result.Passing > 0 && validation.Result.Failing <= 0
			// If the expected result is not equal to the actual result, return an error
			if ExpectedResult != result {
				message.Fatalf(fmt.Errorf("validation failed"), "expected result to be %t got %t", ExpectedResult, result)
			}
			// Print the number of passing and failing results
			message.Infof("Validation completed with %d passing and %d failing results", validation.Result.Passing, validation.Result.Failing)
		},
	}

	cmd.Flags().StringVarP(&InputFile, "input-file", "f", STDIN, "the path to a validation manifest file")
	cmd.Flags().StringVarP(&ResourcesFile, "resources-file", "r", "", "the path to an optional resources file")
	cmd.Flags().StringVarP(&OutputFile, "output-file", "o", "", "the path to write the validation with results")
	cmd.Flags().IntVarP(&Timeout, "timeout", "t", DEFAULT_TIMEOUT, "the timeout for stdin (in seconds, -1 for no timeout)")
	cmd.Flags().BoolVarP(&ExpectedResult, "expected-result", "e", true, "the expected result of the validation (-e=false for failing result)")
	cmd.Flags().BoolVar(&ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of the validation")

	return cmd
}

// DevValidate reads a validation manifest and converts it to a LulaValidation struct, then validates it
// Returns the LulaValidation struct and any error encountered
func DevValidate(ctx context.Context, validationBytes []byte, resourcesBytes []byte, confirmExecution bool, spinner *message.Spinner) (lulaValidation types.LulaValidation, err error) {
	// Set resources if resourcesBytes is not empty
	var resources types.DomainResources
	if len(resourcesBytes) > 0 {
		// Unmarshal the resources data to the DomainResources type
		err = json.Unmarshal(resourcesBytes, &resources)
		if err != nil {
			return lulaValidation, err
		}
	}

	lulaValidation, err = RunSingleValidation(ctx,
		validationBytes,
		types.WithStaticResources(resources),
		types.ExecutionAllowed(confirmExecution),
		types.Interactive(RunInteractively),
		types.WithSpinner(spinner),
	)
	if err != nil {
		return lulaValidation, err
	}

	return lulaValidation, nil
}

func writeValidation(result types.LulaValidation, outputFile string) error {
	var resultBytes []byte
	var err error

	// Marshal to json if the output file is empty or a json file
	if outputFile == "" || strings.HasSuffix(outputFile, ".json") {
		resultBytes, err = json.Marshal(result)
	} else {
		resultBytes, err = yaml.Marshal(result)
	}
	// Return an error if it fails to marshal the result
	if err != nil {
		return err
	}

	// Write the result to the output file if provided, otherwise print to the debug console
	if outputFile == "" {
		message.Debug(string(resultBytes))
	} else {
		err = files.WriteOutput(resultBytes, outputFile)
		if err != nil {
			return err
		}
	}

	return nil
}
