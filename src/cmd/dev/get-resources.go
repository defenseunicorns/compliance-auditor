package dev

import (
	"context"
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
	OutputFile string // -o --output-file
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

			collection, err := DevGetResources(ctx, opts.InputFile)
			if err != nil {
				message.Fatalf(err, "error running dev get-resources: %v", err)
			}

			PrintJSON(collection, opts.OutputFile)

			spinner.Success()
		},
	}

	devCmd.AddCommand(getResourcesCmd)

	getResourcesCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a validation manifest file")
	getResourcesCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write the resources json")
}

func DevGetResources(ctx context.Context, inputFile string) (map[string]interface{}, error) {
	validationFile, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var validation types.Validation
	err = yaml.Unmarshal(validationFile, &validation)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	collection, err := GetResources(ctx, validation.Target)
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %v", err)
	}

	return collection, nil
}

func PrintJSON(data map[string]interface{}, filepath string) {
	jsonData := message.JSONValue(data)

	// If a filepath is provided, write the JSON data to the file.
	if filepath != "" {
		err := os.WriteFile(filepath, []byte(jsonData), 0644)
		if err != nil {
			message.Fatalf(err, "error writing JSON to file: %v", err)
		}
	} else {
		// Else print to stdout
		fmt.Println(jsonData)
	}
}

// TODO: refactor this to a common package
func GetResources(ctx context.Context, data types.Target) (map[string]interface{}, error) {
	if data.Domain == "kubernetes" {
		payload := data.Payload

		err := kube.EvaluateWait(payload.Wait)
		if err != nil {
			return nil, err
		}

		collection, err := kube.QueryCluster(ctx, payload.Resources)
		if err != nil {
			return nil, err
		}

		return collection, nil
	}

	return nil, fmt.Errorf("domain not supported: %s", data.Domain)
}
