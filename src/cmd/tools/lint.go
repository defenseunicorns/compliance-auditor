package tools

import (
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type flags struct {
	InputFile  string // -f --input-file
	ResultFile string // -r --result-file
}

var opts = &flags{}

var lintHelp = `
To lint an existing OSCAL file:
	lula tools lint -f <path to oscal>
`

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Validate OSCAL against schema",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Validate an OSCAL document is properly configured against the OSCAL schema",
		Example: lintHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinner := message.NewProgressSpinner("Linting %s", opts.InputFile)
			defer spinner.Stop()

			validationResp, err := validation.ValidationCommand(opts.InputFile)

			for _, warning := range validationResp.Warnings {
				message.Warn(warning)
			}

			if opts.ResultFile != "" {
				validation.WriteValidationResult(validationResp.Result, opts.ResultFile)
			}

			if err != nil {
				message.Fatalf(err, "Failed to lint %s", opts.InputFile)
			}
			message.Infof("Successfully validated %s is valid OSCAL version %s %s\n", opts.InputFile, validationResp.Validator.GetSchemaVersion(), validationResp.Validator.GetModelType())
			spinner.Success()
		},
	}

	toolsCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	lintCmd.Flags().StringVarP(&opts.ResultFile, "result-file", "r", "", "the path to write the validation result")
}
