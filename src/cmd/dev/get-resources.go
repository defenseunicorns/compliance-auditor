package dev

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/config"
	kube "github.com/defenseunicorns/lula/src/pkg/common/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -r --result-file
}

var opts = &flags{}

var getResourcesHelp = `
To get resources from lula validation manifest:
	lula dev get-resources -f <path to manifest>

	Example:
	lula dev get-resources -f /path/to/manifest.json -o /path/to/output.json
`

func init() {
	getResourcesCmd := &cobra.Command{
		Use:   "get-resources",
		Short: "Get Resources from a Lula Validation Manifest",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Get the JSON resources specified in a Lula Validation Manifest",
		Example: getResourcesHelp,
		Run: func(cmd *cobra.Command, args []string) {
			spinner := message.NewProgressSpinner("Getting Resources from %s", opts.InputFile)
			defer spinner.Stop()

			ctx := context.Background()

			validationFile, err := os.ReadFile(opts.InputFile)
			if err != nil {
				message.Fatalf(err, "error reading YAML file: %v", err)
			}

			var validation types.Validation
			err = yaml.Unmarshal(validationFile, &validation)
			if err != nil {
				message.Fatalf(err, "error unmarshaling YAML: %v", err)
			}

			var collection map[string]interface{}
			// Get resources from the validation manifest -> refactor this?
			if validation.Target.Domain == "kubernetes" {
				payload := validation.Target.Payload

				err := kube.EvaluateWait(payload.Wait)
				if err != nil {
					message.Fatalf(err, "error getting resources with wait: %v", err)
				}

				collection, err = kube.QueryCluster(ctx, payload.Resources)
				if err != nil {
					message.Fatalf(err, "error getting resources: %v", err)
				}
			}

			PrintJSON(collection, opts.OutputFile)

			// message.Infof("Successfully got resources %s is valid OSCAL version %s %s\n", opts.InputFile, validationResp.Validator.GetSchemaVersion(), validationResp.Validator.GetModelType())
			spinner.Success()
		},
	}

	devCmd.AddCommand(getResourcesCmd)

	getResourcesCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a validation manifest file")
	getResourcesCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write the resources json")
}

func PrintJSON(data map[string]interface{}, filepath string) {
	// Marshal the data into JSON bytes.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		message.Fatalf(err, "error marshaling data to JSON: %v", err)
	}

	// If a filepath is provided, write the JSON data to the file.
	if filepath != "" {
		err = os.WriteFile(filepath, jsonData, 0644)
		if err != nil {
			message.Fatalf(err, "error writing JSON to file: %v", err)
		}
	} else {
		// Print to stdout
		fmt.Println(string(jsonData))
	}
}
