package tools

import (
	"os"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/cmd/revise"
	"github.com/defenseunicorns/go-oscal/src/pkg/revision"
	"github.com/defenseunicorns/go-oscal/src/pkg/utils"
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
			// Check if the input file is specified and is a json or yaml file
			if upgradeOpts.InputFile == "" {
				message.Fatalf(nil, "No input file specified")
			} else if err := goOscalUtils.IsJsonOrYaml(upgradeOpts.InputFile); err != nil {
				message.Fatalf(err, "Invalid input file %s", upgradeOpts.InputFile)
			}

			// If no output file is specified, write to the input file
			if upgradeOpts.OutputFile == "" {
				upgradeOpts.OutputFile = upgradeOpts.InputFile
			} else if err := goOscalUtils.IsJsonOrYaml(upgradeOpts.OutputFile); err != nil {
				message.Fatalf(err, "Invalid output file %s", upgradeOpts.OutputFile)
			}

			// Get the output file extension
			split := strings.Split(upgradeOpts.OutputFile, ".")
			outputExt := split[len(split)-1]

			// If no version is specified, use the latest supported version
			if err := goOscalUtils.IsValidOscalVersion(upgradeOpts.Version); err != nil {
				message.Fatalf(err, "Invalid version %s", upgradeOpts.Version)
			}

			spinner := message.NewProgressSpinner("Upgrading %s to version %s", upgradeOpts.InputFile, upgradeOpts.Version)
			defer spinner.Stop()

			// Read the input file
			bytes, err := os.ReadFile(upgradeOpts.InputFile)
			if err != nil {
				message.Fatalf(err, "Failed to read %s", upgradeOpts.InputFile)
			}

			// Create Upgrader
			upgrader, err := revision.NewReviser(bytes, upgradeOpts.Version)
			if err != nil {
				message.Fatalf(err, "Failed to create reviser")
			}

			// Warn if the version is not the latest
			version := upgrader.GetSchemaVersion()
			warning := utils.VersionWarning(version)
			if warning != nil {
				message.Warn(warning.Error())
			}

			upgrader.SetDocumentPath(upgradeOpts.InputFile)

			// Upgrade the document
			revisionError := upgrader.Revise()

			// Write the validation result if it was specified and exists before handling the revision error
			result, err := upgrader.GetValidationResult()
			if err == nil && upgradeOpts.ValidationResult != "" {
				err = validation.WriteValidationResult(result, upgradeOpts.ValidationResult)
				if err != nil {
					message.Fatalf(err, "Failed to write the validation result to %s", upgradeOpts.ValidationResult)
				}
			}

			// Handle the revision error
			if revisionError != nil {
				message.Fatalf(revisionError, "Failed to upgrade %s to OSCAL version %s", upgradeOpts.InputFile, upgradeOpts.Version)
			}

			// Get the upgraded bytes
			upgradedBytes, err := upgrader.GetRevisedBytes(outputExt)
			if err != nil {
				message.Fatalf(err, "Failed to marshal the upgraded document")
			}

			// Write the upgraded document
			err = goOscalUtils.WriteOutput(upgradedBytes, upgradeOpts.OutputFile)
			if err != nil {
				message.Fatalf(err, "Failed to write the upgraded document to %s", upgradeOpts.OutputFile)
			}

			message.Infof("Successfully upgraded %s to OSCAL version %s %s\n", upgradeOpts.InputFile, upgrader.GetSchemaVersion(), upgrader.GetModelType())
			spinner.Success()
		},
	}

	toolsCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&upgradeOpts.InputFile, "input-file", "f", "", "the path to a oscal json schema file")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.OutputFile, "output-file", "o", "", "the path to write the linted oscal json schema file (default is the input file)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.Version, "version", "v", goOscalUtils.GetLatestSupportedVersion(), "the version of the oscal schema to validate against (default is the latest supported version)")
	upgradeCmd.Flags().StringVarP(&upgradeOpts.ValidationResult, "validation-result", "r", "", "the path to write the validation result file")
}
