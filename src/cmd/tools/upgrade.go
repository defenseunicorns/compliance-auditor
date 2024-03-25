package tools

import (
	"github.com/defenseunicorns/go-oscal/src/pkg/revision"
	goOscalUtils "github.com/defenseunicorns/go-oscal/src/pkg/utils"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var upgradeHelp = `
To Upgrade an existing OSCAL file:
	lula tools upgrade -f <path to oscal> -v <version>
`

type upgradeOptions struct {
	revision.RevisionOptions
}

var upgradeOpts upgradeOptions = upgradeOptions{}

func init() {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade OSCAL document to a new version if possible.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Validate an OSCAL document against the OSCAL schema version provided. If the document is valid, upgrade it to the provided OSCAL version. Otherwise, return or write as ValidationError. Yaml formatting handled by gopkg/yaml.v3 and while objects will maintain deep equality, visual representation may be different than the input file.",
		Example: upgradeHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinner := message.NewProgressSpinner("Upgrading %s to version %s", upgradeOpts.InputFile, upgradeOpts.Version)
			defer spinner.Stop()

			// If no output file is specified, write to the input file
			if upgradeOpts.OutputFile == "" {
				upgradeOpts.OutputFile = upgradeOpts.InputFile
			}

			revisionResponse, err := revision.RevisionCommand(&upgradeOpts.RevisionOptions)

			if upgradeOpts.ValidationResult != "" {
				validation.WriteValidationResult(revisionResponse.Result, upgradeOpts.ValidationResult)
			}

			if len(revisionResponse.Warnings) > 0 {
				for _, warning := range revisionResponse.Warnings {
					message.Warn(warning)
				}
			}

			err = goOscalUtils.WriteOutput(revisionResponse.RevisedBytes, upgradeOpts.OutputFile)
			if err != nil {
				message.Fatalf(err, "Failed to write upgraded %s with: %s", upgradeOpts.OutputFile, err)
			}

			if err != nil {
				message.Fatalf(err, "Failed to upgrade %s to OSCAL version %s %s", upgradeOpts.InputFile, revisionResponse.Reviser.GetSchemaVersion(), revisionResponse.Reviser.GetModelType())
			}
			message.Infof("Successfully upgraded %s to OSCAL version %s %s\n", upgradeOpts.InputFile, revisionResponse.Reviser.GetSchemaVersion(), revisionResponse.Reviser.GetModelType())
			spinner.Success()
		},
	}

	toolsCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&upgradeOpts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.OutputFile, "output-file", "o", "", "the path to write the linted oscal json schema file (default is the input file)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.Version, "version", "v", goOscalUtils.GetLatestSupportedVersion(), "the version of the oscal schema to validate against (default is the latest supported version)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.ValidationResult, "validation-result", "r", "", "the path to write the validation result file")
}
