package tools

import (
	"github.com/defenseunicorns/go-oscal/src/cmd/validate"
	"github.com/spf13/cobra"
)

type flags struct {
	InputFile string // -f --input-file
	LogFile   string // -l --logger-file

}

var opts = &flags{}

var lintHelp = `
To lint an existing OSCAL file:
	lula tools lint <path to oscal>
`

func init() {
	lintCmd := &cobra.Command{
		Use:     "lint",
		Short:   "Validate OSCAL against schema",
		Long:    "Validate an OSCAL document is properly configured against the OSCAL schema",
		Example: lintHelp,
		RunE: func(cmd *cobra.Command, args []string) error {

			validate.ValidateCommand(opts.InputFile, opts.LogFile)

			return nil
		},
	}

	toolsCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	lintCmd.Flags().StringVarP(&opts.LogFile, "logger-file", "l", "", "the name of the file to write logs to (outputs to STDERR by default)")
}
