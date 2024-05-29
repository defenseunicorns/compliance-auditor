package tools

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

func init() {
	composeCmd := &cobra.Command{
		Use:     "compose",
		Short:   "compose an OSCAL component definition",
		Long:    "Lula Composition of an OSCAL component definition. Used to compose remote validations within a component definition in order to resolve any references for portability.",
		Example: composeHelp,
		Run: func(cmd *cobra.Command, args []string) {
			if composeOpts.InputFile == "" {
				message.Fatal(errors.New("flag input-file is not set"),
					"Please specify an input file with the -f flag")
			}
			err := Compose(composeOpts.InputFile, composeOpts.OutputFile)
			if err != nil {
				message.Fatalf(err, "Composition error: %s", err)
			}
		},
	}

	toolsCmd.AddCommand(composeCmd)

	composeCmd.Flags().StringVarP(&composeOpts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	composeCmd.Flags().StringVarP(&composeOpts.OutputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be the original filename with `-composed` appended")
}

func Compose(inputFile, outputFile string) error {
	_, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file: %v does not exist - unable to compose document", inputFile)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Change Cwd to the directory of the component definition
	dirPath := filepath.Dir(inputFile)
	message.Infof("changing cwd to %s", dirPath)
	resetCwd, err := common.SetCwdToFileDir(dirPath)
	if err != nil {
		return err
	}

	model, err := oscal.NewOscalModel(data)
	if err != nil {
		return err
	}

	err = composition.ComposeComponentDefinitions(model.ComponentDefinition)
	if err != nil {
		return err
	}

	// Reset Cwd to original before outputting
	resetCwd()

	var b bytes.Buffer
	// Format the output
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(model)

	outputFileName := outputFile
	if outputFileName == "" {
		outputFileName = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-composed" + filepath.Ext(inputFile)
	}

	message.Infof("Writing Composed OSCAL Component Definition to: %s", outputFileName)

	err = files.WriteOutput(b.Bytes(), outputFileName)
	if err != nil {
		return err
	}

	return nil
}
