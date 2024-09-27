package tools

import (
	"context"
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

	var cmd = &cobra.Command{
		Use:     "compose",
		Short:   "compose an OSCAL component definition",
		Long:    composeLong,
		Example: composeHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeSpinner := message.NewProgressSpinner("Composing %s", inputFile)
			defer composeSpinner.Stop()

			// Update input/output paths
			if filepath.IsLocal(inputFile) {
				inputFile = filepath.Join(filepath.Dir(inputFile), filepath.Base(inputFile))
			}

			if outputFile == "" {
				outputFile = GetDefaultOutputFile(inputFile)
			} else if filepath.IsLocal(outputFile) {
				outputFile = filepath.Join(filepath.Dir(outputFile), filepath.Base(outputFile))
			}

			// Check if output file contains a valid OSCAL model
			_, err := oscal.ValidOSCALModelAtPath(outputFile)
			if err != nil {
				message.Fatalf(err, "Output file %s is not a valid OSCAL model: %v", outputFile, err)
			}

			// Compose the OSCAL model
			constants, variables, err := common.GetTemplateConfig()
			if err != nil {
				return fmt.Errorf("error getting template config: %v", err)
			}

			opts := []composition.Option{
				composition.WithModelFromLocalPath(inputFile),
				composition.WithRenderSettings(renderTypeString, renderValidations),
				composition.WithTemplateRenderer(renderTypeString, constants, variables, setOpts),
			}

			err = Compose(cmd.Context(), inputFile, outputFile, opts...)
			if err != nil {
				return fmt.Errorf("error composing model: %v", err)
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
func Compose(ctx context.Context, inputFile, outputFile string, opts ...composition.Option) error {
	// Compose the OSCAL model
	compositionCtx, err := composition.New(opts...)
	if err != nil {
		return fmt.Errorf("error creating composition context: %v", err)
	}

	model, err := compositionCtx.ComposeFromPath(ctx, inputFile)
	if err != nil {
		return err
	}

	// Write the composed OSCAL model to a file
	err = oscal.WriteOscalModel(outputFile, model)
	if err != nil {
		return err
	}

	return nil
}

// GetDefaultOutputFile returns the default output file name
func GetDefaultOutputFile(inputFile string) string {
	return strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-composed" + filepath.Ext(inputFile)
}
