package dev

import (
	"context"

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
			spinner := message.NewProgressSpinner("Validating %s", opts.InputFile)
			defer spinner.Stop()

			ctx := context.Background()

			validation, err := DevValidate(ctx, opts.InputFile)
			if err != nil {
				message.Fatalf(err, "error running dev validate: %v", err)
			}

			message.Infof("Validation Result: %v", validation.Result)
			if validation.Result.Failing > 0 {
				message.Fatalf(nil, "Validation failed")
				spinner.Stop()
			} else {
				message.Infof("Validation Passed")
				spinner.Success()
			}
		},
	}

	devCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a validation manifest file")
	validateCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write the resources json")
}

func DevValidate(ctx context.Context, inputFile string) (validation types.Validation, err error) {
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
