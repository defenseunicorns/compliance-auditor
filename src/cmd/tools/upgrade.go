package tools

import (
	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/go-oscal/src/pkg/revision"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/go-oscal/src/pkg/versioning"
	"github.com/defenseunicorns/lula/src/cmd/common"
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

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade OSCAL document to a new version if possible.",
	Long:    "Validate an OSCAL document against the OSCAL schema version provided. If the document is valid, upgrade it to the provided OSCAL version. Otherwise, return or write as ValidationError. Yaml formatting handled by gopkg/yaml.v3 and while objects will maintain deep equality, visual representation may be different than the input file.",
	Example: upgradeHelp,
	Run: func(cmd *cobra.Command, args []string) {
		spinner := message.NewProgressSpinner("Upgrading %s to version %s", upgradeOpts.InputFile, upgradeOpts.Version)
		defer spinner.Stop()

		// If no output file is specified, write to the input file
		if upgradeOpts.OutputFile == "" {
			upgradeOpts.OutputFile = upgradeOpts.InputFile
		}

		revisionResponse, revisionErr := revision.RevisionCommand(&upgradeOpts.RevisionOptions)

		if upgradeOpts.ValidationResult != "" {
			err := validation.WriteValidationResult(revisionResponse.Result, upgradeOpts.ValidationResult)
			if err != nil {
				message.Fatalf("Failed to write validation result to %s: %s\n", upgradeOpts.ValidationResult, err)
			}
		}

		if revisionErr != nil {
			message.Fatalf(revisionErr, "Failed to upgrade %s: %s", upgradeOpts.InputFile, revisionErr)
		}

		if len(revisionResponse.Warnings) > 0 {
			for _, warning := range revisionResponse.Warnings {
				message.Warn(warning)
			}
		}

		err := files.WriteOutput(revisionResponse.RevisedBytes, upgradeOpts.OutputFile)
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

func init() {

	common.InitViper()
	toolsCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&upgradeOpts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	upgradeCmd.MarkFlagRequired("input-file")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.OutputFile, "output-file", "o", "", "the path to write the linted oscal json schema file (default is the input file)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.Version, "version", "v", versioning.GetLatestSupportedVersion(), "the version of the oscal schema to validate against (default is the latest supported version)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.ValidationResult, "validation-result", "r", "", "the path to write the validation result file")
}
