package tools

import (
	"context"
	"fmt"
	"os"
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

	var cmd = &cobra.Command{
		Use:     "compose",
		Short:   "compose an OSCAL component definition",
		Long:    composeLong,
		Example: composeHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeSpinner := message.NewProgressSpinner("Composing %s", inputFile)
			defer composeSpinner.Stop()

			// TODO: check if remote or local?
			_, err := os.Stat(inputFile)
			if os.IsNotExist(err) {
				return fmt.Errorf("input-file: %v does not exist - unable to digest document", inputFile)
			}

			// Update path if relative
			path := inputFile
			if filepath.IsLocal(inputFile) {
				path = filepath.Join(filepath.Dir(inputFile), filepath.Base(inputFile))
			}

			if outputFile == "" {
				outputFile = GetDefaultOutputFile(inputFile)
			}

			opts := []composition.Option{
				composition.WithModelFromPath(path),
				composition.WithTemplateRenderer(renderTypeString, renderValidations, setOpts),
			}

			compositionCtx, err := composition.New(context.Background(), opts...)
			if err != nil {
				return fmt.Errorf("error creating composition context: %v", err)
			}

			err = compositionCtx.ComposeFromPath(path)
			if err != nil {
				return fmt.Errorf("error composing model: %v", err)
			}

			// Write the composed OSCAL model to a file
			err = oscal.WriteOscalModel(outputFile, compositionCtx.GetModel())
			if err != nil {
				return fmt.Errorf("error writing composed model: %v", err)
			}

			message.Infof("Composed OSCAL Component Definition to: %s", outputFile)
			composeSpinner.Success()

			return nil
		},
	}
	cmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	cmd.MarkFlagRequired("input-file")
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be the original filename with `-composed` appended")
	cmd.Flags().StringVarP(&renderTypeString, "render", "r", "", "values to render the template with, options are: masked, constants, non-sensitive, all")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set value overrides for templated data")
	cmd.Flags().BoolVar(&renderValidations, "render-validations", false, "extend render to remote Lula Validations")

	return cmd
}

func init() {
	common.InitViper()
	toolsCmd.AddCommand(ComposeCommand())
}

// Compose composes an OSCAL model from a file path
// func Compose(inputFile, outputFile string, templateRenderer *template.TemplateRenderer, renderType template.RenderType) error {
// 	_, err := os.Stat(inputFile)
// 	if os.IsNotExist(err) {
// 		return fmt.Errorf("input file: %v does not exist - unable to compose document", inputFile)
// 	}

// 	// Compose the OSCAL model
// 	model, err := composition.ComposeFromPath(inputFile)
// 	if err != nil {
// 		return err
// 	}

// 	// Write the composed OSCAL model to a file
// 	err = oscal.WriteOscalModel(outputFile, model)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// GetDefaultOutputFile returns the default output file name
func GetDefaultOutputFile(inputFile string) string {
	return strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-composed" + filepath.Ext(inputFile)
}
