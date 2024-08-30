package dev

import (
	"context"
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
)

var getResourcesOpts = &flags{}

var getResourcesHelp = `
To get resources from lula validation manifest:
	lula dev get-resources -f /path/to/validation.yaml
To get resources from lula validation manifest and write to file:
	lula dev get-resources -f /path/to/validation.yaml -o /path/to/output.json
To get resources from lula validation and automatically confirm execution
	lula dev get-resources -f /path/to/validation.yaml --confirm-execution
To run validations using stdin:
	cat /path/to/validation.yaml | lula dev get-resources
To hang indefinitely for stdin:
	lula get-resources -t -1
To hang for timeout of 5 seconds:
	lula get-resources -t 5
`

var getResourcesCmd = &cobra.Command{
	Use:     "get-resources",
	Short:   "Get Resources from a Lula Validation Manifest",
	Long:    "Get the JSON resources specified in a Lula Validation Manifest",
	Example: getResourcesHelp,
	Run: func(cmd *cobra.Command, args []string) {
		spinnerMessage := fmt.Sprintf("Getting Resources from %s", getResourcesOpts.InputFile)
		spinner := message.NewProgressSpinner(spinnerMessage)
		defer spinner.Stop()

		ctx := context.Background()
		var validationBytes []byte
		var err error

		// Read the validation data from STDIN or provided file
		validationBytes, err = ReadValidation(cmd, spinner, getResourcesOpts.InputFile, getResourcesOpts.Timeout)
		if err != nil {
			message.Fatalf(err, "error reading validation: %v", err)
		}

		collection, err := DevGetResources(ctx, validationBytes, spinner)
		if err != nil {
			message.Fatalf(err, "error running dev get-resources: %v", err)
		}

		writeResources(collection, getResourcesOpts.OutputFile)

		spinner.Success()
	},
}

func init() {

	common.InitViper()

	devCmd.AddCommand(getResourcesCmd)

	getResourcesCmd.Flags().StringVarP(&getResourcesOpts.InputFile, "input-file", "f", STDIN, "the path to a validation manifest file")
	getResourcesCmd.Flags().StringVarP(&getResourcesOpts.OutputFile, "output-file", "o", "", "the path to write the resources json")
	getResourcesCmd.Flags().IntVarP(&getResourcesOpts.Timeout, "timeout", "t", DEFAULT_TIMEOUT, "the timeout for stdin (in seconds, -1 for no timeout)")
	getResourcesCmd.Flags().BoolVar(&getResourcesOpts.ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of getting resources")
}

func DevGetResources(ctx context.Context, validationBytes []byte, spinner *message.Spinner) (types.DomainResources, error) {
	lulaValidation, err := RunSingleValidation(
		validationBytes,
		types.ExecutionAllowed(getResourcesOpts.ConfirmExecution),
		types.Interactive(RunInteractively),
		types.WithSpinner(spinner),
		types.GetResourcesOnly(true),
	)
	if err != nil {
		return nil, err
	}

	return *lulaValidation.DomainResources, nil
}

func writeResources(data types.DomainResources, filepath string) {
	jsonData := message.JSONValue(data)

	// If a filepath is provided, write the JSON data to the file.
	if filepath != "" {
		err := os.WriteFile(filepath, []byte(jsonData), 0644)
		if err != nil {
			message.Fatalf(err, "error writing resource JSON to file: %v", err)
		}
	} else {
		// Else print to stdout
		fmt.Println(jsonData)
	}
}
