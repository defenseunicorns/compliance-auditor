package tools

import (
	"os"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type templateFlags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

var templateOpts = &templateFlags{}

var templateHelp = `
To template an OSCAL Model:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

Data for the templating should be stored under the 'variables' configuration item in a lula-config.yaml file
`
var templateCmd = &cobra.Command{
	Use:     "template",
	Short:   "Template an artifact",
	Long:    "Resolving templated artifacts with configuration data",
	Args:    cobra.NoArgs,
	Example: templateHelp,
	Run: func(cmd *cobra.Command, args []string) {
		// Read file
		data, err := pkgCommon.ReadFileToBytes(templateOpts.InputFile)
		if err != nil {
			message.Fatal(err, err.Error())
		}

		// Get current viper pointer
		v := common.GetViper()
		// Get all viper settings
		// This will only return config file items and resolved environment variables
		// that have an associated key in the config file
		viperData := v.AllSettings()

		// Handles merging viper config file data + environment variables
		mergedMap := template.CollectTemplatingData(viperData)

		templatedData, err := template.ExecuteTemplate(mergedMap, string(data))
		if err != nil {
			message.Fatalf(err, "error templating validation: %v", err)
		}

		if templateOpts.OutputFile == "" {
			_, err := os.Stdout.Write(templatedData)
			if err != nil {
				message.Fatalf(err, "failed to write to stdout: %v", err)
			}
		} else {
			err = files.CreateFileDirs(templateOpts.OutputFile)
			if err != nil {
				message.Fatalf(err, "failed to create output file path: %s\n", err)
			}
			err = os.WriteFile(templateOpts.OutputFile, templatedData, 0644)
			if err != nil {
				message.Fatal(err, err.Error())
			}
		}

	},
}

func TemplateCommand() *cobra.Command {
	return templateCmd
}

func init() {
	common.InitViper()

	toolsCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVarP(&templateOpts.InputFile, "input-file", "f", "", "the path to the target artifact")
	templateCmd.MarkFlagRequired("input-file")
	templateCmd.Flags().StringVarP(&templateOpts.OutputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be directed to stdout")
}
