package tools

import (
	"errors"
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

type composeFlags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

var composeOpts = &composeFlags{}

var composeHelp = `
To compose an OSCAL Model:
	lula tools compose -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools compose -f ./oscal-component.yaml -o composed-oscal-component.yaml
`
var composeCmd = &cobra.Command{
	Use:     "compose",
	Short:   "compose an OSCAL component definition",
	Long:    "Lula Composition of an OSCAL component definition. Used to compose remote validations within a component definition in order to resolve any references for portability.",
	Example: composeHelp,
	Run: func(cmd *cobra.Command, args []string) {
		composeSpinner := message.NewProgressSpinner("Composing %s", composeOpts.InputFile)
		defer composeSpinner.Stop()

		if composeOpts.InputFile == "" {
			message.Fatal(errors.New("flag input-file is not set"),
				"Please specify an input file with the -f flag")
		}

		outputFile := composeOpts.OutputFile
		if outputFile == "" {
			outputFile = GetDefaultOutputFile(composeOpts.InputFile)
		}

		err := Compose(composeOpts.InputFile, outputFile)
		if err != nil {
			message.Fatalf(err, "Composition error: %s", err)
		}

		message.Infof("Composed OSCAL Component Definition to: %s", outputFile)
		composeSpinner.Success()
	},
}

func init() {
	common.InitViper()

	toolsCmd.AddCommand(composeCmd)

	composeCmd.Flags().StringVarP(&composeOpts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	composeCmd.Flags().StringVarP(&composeOpts.OutputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be the original filename with `-composed` appended")
}

// Compose composes an OSCAL model from a file path
func Compose(inputFile, outputFile string) error {
	_, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file: %v does not exist - unable to compose document", inputFile)
	}

	// Compose the OSCAL model
	model, err := composition.ComposeFromPath(inputFile)
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
