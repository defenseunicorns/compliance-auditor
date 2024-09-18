package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type templateFlags struct {
	InputFile  string   // -f --input-file
	OutputFile string   // -o --output-file
	Set        []string // --set
	All        bool     // --all
	// some demo flags
	Const        bool // --const
	NonSensitive bool // --non-sensitive
	Sensitive    bool // --sensitive
}

var templateOpts = &templateFlags{}

var templateHelp = `
To template an OSCAL Model:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

To perform overrides on the template data:
	lula tools template -f ./oscal-component.yaml --set var.key1=value1 --set const.key2=value2

To perform the full template operation, including sensitive data:
	lula tools template -f ./oscal-component.yaml --all

Data for templating should be stored under 'constants' or 'variables' configuration items in a lula-config.yaml file
See documentation for more detail on configuration schema
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

		// Get constants and variables for templating from viper config
		var constants map[string]interface{}
		var variables []template.VariableConfig

		err = v.UnmarshalKey(common.VConstants, &constants)
		if err != nil {
			message.Fatalf(err, "unable to unmarshal constants into map: %v", err)
		}

		err = v.UnmarshalKey(common.VVariables, &variables)
		if err != nil {
			message.Fatalf(err, "unable to unmarshal variables into slice: %v", err)
		}

		// Handles merging viper config file data + environment variables
		templateData := template.CollectTemplatingData(constants, variables)

		// override anything that's --set
		overrideTemplateValues(templateData, templateOpts.Set)

		// Execute the template function based on the flags (TESTING ONLY)
		var templatedData []byte
		if templateOpts.Const {
			templatedData, err = template.ExecuteConstTemplate(templateData.Constants, string(data))
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}
		} else if templateOpts.NonSensitive {
			templatedData, err = template.ExecuteNonSensitiveTemplate(templateData, string(data))
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}
		} else if templateOpts.Sensitive {
			templatedData, err = template.ExecuteSensitiveTemplate(templateData, string(data))
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}
		} else if templateOpts.All {
			templatedData, err = template.ExecuteFullTemplate(templateData, string(data))
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}
		} else {
			templatedData, err = template.ExecuteMaskedTemplate(templateData, string(data))
			if err != nil {
				message.Fatalf(err, "error templating validation: %v", err)
			}
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
	templateCmd.Flags().StringSliceVarP(&templateOpts.Set, "set", "s", []string{}, "set a value in the template data")

	templateCmd.Flags().BoolVar(&templateOpts.Const, "const", false, "only include constants in the template")
	templateCmd.Flags().BoolVar(&templateOpts.NonSensitive, "non-sensitive", false, "only include non-sensitive variables in the template")
	templateCmd.Flags().BoolVar(&templateOpts.Sensitive, "sensitive", false, "only include sensitive variables in the template")
	templateCmd.Flags().BoolVar(&templateOpts.All, "all", false, "include all variables in the template")
}

func overrideTemplateValues(templateData *template.TemplateData, setFlags []string) {
	for _, flag := range setFlags {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			message.Fatalf(fmt.Errorf("invalid --set flag format, should be key.path=value"), "invalid --set flag format, should be key.path=value")
		}
		path, value := parts[0], parts[1]

		// for each set flag, check if .var or .const
		// if .var, set the value in the templateData.Variables
		// if .const, set the value in the templateData.Constants
		if strings.HasPrefix(path, "."+template.VAR+".") {
			// Set the value in the templateData.Variables
			key := strings.TrimPrefix(path, "."+template.VAR+".")
			templateData.Variables[key] = value
		} else if strings.HasPrefix(path, "."+template.CONST+".") {
			// Set the value in the templateData.Constants
			key := strings.TrimPrefix(path, "."+template.CONST+".")
			setNestedValue(templateData.Constants, key, value)
		}
	}
}

// Helper function to set a value in a map based on a JSON-like key path
func setNestedValue(m map[string]interface{}, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	lastKey := keys[len(keys)-1]

	// Traverse the map, creating intermediate maps if necessary
	for _, key := range keys[:len(keys)-1] {
		if _, exists := m[key]; !exists {
			m[key] = make(map[string]interface{})
		}
		if nestedMap, ok := m[key].(map[string]interface{}); ok {
			m = nestedMap
		} else {
			return fmt.Errorf("path %s contains a non-map value", key)
		}
	}

	// Set the final value
	m[lastKey] = value
	return nil
}
