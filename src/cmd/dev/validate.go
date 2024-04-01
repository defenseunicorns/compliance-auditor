package dev

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/utils"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var validateHelp = `
To run validations using a lula validation manifest:
	lula dev validate -f <path to manifest>
`

type ValidateFlags struct {
	flags
	ExpectedResult bool // -e --expected-result
}

var validateOpts = &ValidateFlags{}

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate Resources from a Lula Validation Manifest",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "TODO: Add long description here.",
		Example: validateHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinner := message.NewProgressSpinner("Validating %s", validateOpts.InputFile)
			defer spinner.Stop()

			ctx := context.Background()

			validation, err := DevValidate(ctx, validateOpts.InputFile)
			if err != nil {
				message.Fatalf(err, "error running dev validate: %v", err)
			}

			err = writeValidation(validation, validateOpts.OutputFile)
			if err != nil {
				message.Fatalf(err, "error writing result: %v", err)
			}

			result := validation.Result.Failing == 0
			if validateOpts.ExpectedResult != result {
				message.Fatalf(fmt.Errorf("validation failed"), "expected result to be %t got %t", validateOpts.ExpectedResult, result)
			}
			message.Infof("Validation completed with %d passing and %d failing results", validation.Result.Passing, validation.Result.Failing)
		},
	}

	devCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&validateOpts.InputFile, "input-file", "f", "", "the path to a validation manifest file")
	validateCmd.Flags().StringVarP(&validateOpts.OutputFile, "output-file", "o", "", "the path to write the validation with results")
	validateCmd.Flags().BoolVarP(&validateOpts.ExpectedResult, "expected-result", "e", true, "the expected result of the validation")

	validateCmd.MarkFlagRequired("input-file")
}

// DevValidate runs a validation using a lula validation manifest instead of an oscal file.
// It returns the validation result.
// or an error if the validation fails.
func DevValidate(ctx context.Context, inputFile string) (validation types.Validation, err error) {
	// Return an error if the input file is empty or not a yaml file
	if inputFile == "" || !strings.HasSuffix(inputFile, ".yaml") {
		return types.Validation{}, fmt.Errorf("input file must be a yaml file")
	}

	validationBytes, err := common.ReadFileToBytes(inputFile)
	if err != nil {
		return types.Validation{}, err
	}

	err = yaml.Unmarshal(validationBytes, &validation)
	if err != nil {
		return types.Validation{}, err
	}

	result, err := validate.ValidateOnTarget(ctx, validation.Title, validation.Target)
	if err != nil {
		return types.Validation{}, err
	}
	validation.Result = result
	validation.Evaluated = true

	return validation, nil
}

func writeValidation(result types.Validation, outputFile string) error {
	var resultBytes []byte
	var err error

	if outputFile == "" || strings.HasSuffix(outputFile, ".json") {
		resultBytes, err = json.Marshal(result)
	} else {
		resultBytes, err = yaml.Marshal(result)
	}
	if err != nil {
		return err
	}

	if outputFile == "" {
		message.Debug(string(resultBytes))
	} else {
		err = utils.WriteOutput(resultBytes, outputFile)
		if err != nil {
			return err
		}
	}

	return nil
}
