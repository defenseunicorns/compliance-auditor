package dev

import (
	"context"
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
)

type flags struct {
	InputFile        string // -f --input-file
	OutputFile       string // -o --output-file
	ConfirmExecution bool   // --confirm-execution
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
			spinnerMessage := fmt.Sprintf("Getting Resources from %s", opts.InputFile)
			spinner := message.NewProgressSpinner(spinnerMessage)
			defer spinner.Stop()

			ctx := context.Background()
			var validationBytes []byte
			var err error

			// Read the validation data from STDIN or provided file
			validationBytes, err = ReadValidation(cmd, spinner, validateOpts.InputFile, validateOpts.Timeout)
			if err != nil {
				message.Fatalf(err, "error reading validation: %v", err)
			}

			// Reset the spinner message
			spinner.Updatef(spinnerMessage)

			collection, err := DevGetResources(ctx, validationBytes)
			if err != nil {
				message.Fatalf(err, "error running dev get-resources: %v", err)
			}

			writeResources(collection, opts.OutputFile)

			spinner.Success()
		},
	}

	devCmd.AddCommand(getResourcesCmd)

	getResourcesCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to a validation manifest file")
	getResourcesCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write the resources json")
	getResourcesCmd.Flags().BoolVar(&opts.ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of getting resources")
}

func DevGetResources(ctx context.Context, validationBytes []byte) (types.DomainResources, error) {
	lulaValidation, err := RunSingleValidation(validationBytes, types.ExecutionAllowed(opts.ConfirmExecution), types.Interactive(true))
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
