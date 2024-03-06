package tools

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/cmd/revise"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var upgradeHelp = `
To Upgrade an existing OSCAL file:
	lula tools upgrade -f <path to oscal> -v <version>
`

type upgradeOptions struct {
	revise.ReviseOptions
}

var upgradeOpts upgradeOptions = upgradeOptions{}

func init() {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade OSCAL document to a new version if possible",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Validate an OSCAL document against the OSCAL schema version specified in the document. If the document is valid, upgrade it to the latest version of OSCAL. Otherwise, return an ValidationError.",
		Example: upgradeHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinner := message.NewProgressSpinner("Upgrading %s to version %s", upgradeOpts.InputFile, upgradeOpts.Version)
			defer spinner.Stop()

			if upgradeOpts.InputFile == "" {
				message.Fatalf(nil, "No input file specified")
			}

			// If no output file is specified, write to the input file
			if upgradeOpts.OutputFile == "" {
				upgradeOpts.OutputFile = upgradeOpts.InputFile
			}

			// The Revise command has some logging behavior that is not ideal for lula.
			revisor, err := revise.Revise(&upgradeOpts.ReviseOptions)
			if err != nil {
				fmt.Println(err)
				message.Fatalf(err, "Failed to upgrade %s", upgradeOpts.InputFile)
			}

			message.Infof("Successfully upgraded %s to OSCAL version %s %s\n", upgradeOpts.InputFile, revisor.GetSchemaVersion(), revisor.GetModelType())
			spinner.Success()
		},
	}

	toolsCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&upgradeOpts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.OutputFile, "output-file", "o", "", "the path to write the linted oscal json schema file")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.Version, "version", "v", "", "the version of the oscal schema to validate against")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.ValidationResult, "validation-result", "r", "", "the path to write the validation result file")
}
