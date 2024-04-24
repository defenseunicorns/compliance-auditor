package compile

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gooscalUtils "github.com/defenseunicorns/go-oscal/src/pkg/utils"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/compilation"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type compileFlags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

var compileOpts = &compileFlags{}

var compileHelp = `
To compile an OSCAL Model:
	lula compile -f ./oscal-component.yaml

To indicate a specific output file:
	lula compile -f ./oscal-component.yaml -o compiled-oscal-component.yaml
`

var compileCmd = &cobra.Command{
	Use:     "compile",
	Short:   "compile an OSCAL component definition",
	Long:    "Lula Compilation of an OSCAL component definition. Used to compile remote validations within a component definition in order to resolve any references for portability.",
	Example: compileHelp,
	Run: func(cmd *cobra.Command, componentDefinitionPath []string) {
		if compileOpts.InputFile == "" {
			message.Fatal(errors.New("flag input-file is not set"),
				"Please specify an input file with the -f flag")
		}
		// Primary expected path for compilation of OSCAL documents
		err := Compile(compileOpts.InputFile, compileOpts.OutputFile)
		if err != nil {
			message.Fatalf(err, "Compilation error: %s", err)
		}
	},
}

func CompileCommand() *cobra.Command {

	// insert flag options here
	compileCmd.Flags().StringVarP(&compileOpts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	compileCmd.Flags().StringVarP(&compileOpts.OutputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be the original filename with `-compiled` appended")

	return compileCmd
}

func Compile(inputFile, outputFile string) error {
	_, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("input file: %v does not exist - unable to compile document", inputFile)
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

	err = compilation.CompileComponentValidations(model.ComponentDefinition)
	if err != nil {
		return err
	}

	resetCwd()

	var b bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(model)

	outputFileName := outputFile
	if outputFileName == "" {
		outputFileName = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-compiled" + filepath.Ext(inputFile)
	}

	message.Infof("Writing Compiled OSCAL Component Definition to: %s", outputFileName)

	err = gooscalUtils.WriteOutput(b.Bytes(), outputFileName)
	if err != nil {
		return err
	}

	return nil
}
