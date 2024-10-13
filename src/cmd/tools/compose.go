package tools

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var composeHelp = `
To compose an OSCAL Model:
	lula tools compose -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools compose -f ./oscal-component.yaml -o composed-oscal-component.yaml
`

var composeLong = `
Lula Composition of an OSCAL component definition. Used to compose remote validations within a component definition in order to resolve any references for portability.

Supports templating of the composed component definition with the following configuration options:
- To compose with templating applied, specify '--render, -r' with values of 'all', 'non-sensitive', 'constants', or 'masked' (choice will depend on the use case for the composed content)
- To render Lula Validations include '--render-validations'
- To perform any manual overrides to the template data, specify '--set, -s' with the format '.const.key=value' or '.var.key=value'
`

func ComposeCommand() *cobra.Command {
	var (
		inputFile         string   // -f --input-file
		outputFile        string   // -o --output-file
		setOpts           []string // -s --set
		renderTypeString  string   // -r --render
		renderValidations bool     // --render-validations
	)

	var composeCmd = &cobra.Command{
		Use:     "compose",
		Short:   "compose an OSCAL component definition",
		Long:    composeLong,
		Example: composeHelp,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			composeSpinner := message.NewProgressSpinner("Composing %s", inputFile)
			defer composeSpinner.Stop()

			if outputFile == "" {
				outputFile = GetDefaultOutputFile(inputFile)
			}

			// Check if output file contains a valid OSCAL model
			_, err = oscal.ValidOSCALModelAtPath(outputFile)
			if err != nil {
				return fmt.Errorf("invalid OSCAL model at output file: %v", err)
			}

			opts := []composition.Option{
				composition.WithModelFromLocalPath(inputFile),
				composition.WithRenderSettings(renderTypeString, renderValidations),
				composition.WithTemplateRenderer(renderTypeString, common.TemplateConstants, common.TemplateVariables, setOpts),
			}

			// Compose the OSCAL model
			composer, err := composition.New(opts...)
			if err != nil {
				return fmt.Errorf("error creating new composer: %v", err)
			}

			model, err := composer.ComposeFromPath(cmd.Context(), inputFile)
			if err != nil {
				return fmt.Errorf("error composing model from path: %v", err)
			}

			// Write the composed OSCAL model to a file
			err = oscal.WriteOscalModel(outputFile, model)
			if err != nil {
				return fmt.Errorf("error writing composed model: %v", err)
			}

			message.Infof("Composed OSCAL Component Definition to: %s", outputFile)
			composeSpinner.Success()

			return nil
		},
	}
	composeCmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	composeCmd.MarkFlagRequired("input-file")
	composeCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be the original filename with `-composed` appended")
	composeCmd.Flags().StringVarP(&renderTypeString, "render", "r", "", "values to render the template with, options are: masked, constants, non-sensitive, all")
	composeCmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set value overrides for templated data")
	composeCmd.Flags().BoolVar(&renderValidations, "render-validations", false, "extend render to remote Lula Validations")

	return composeCmd
}

func init() {
	common.InitViper()
	toolsCmd.AddCommand(ComposeCommand())
}

// GetDefaultOutputFile returns the default output file name
func GetDefaultOutputFile(inputFile string) string {
	return strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-composed" + filepath.Ext(inputFile)
}
