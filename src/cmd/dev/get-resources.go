package dev

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

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

func DevGetResourcesCommand() *cobra.Command {

	var (
		InputFile        string // -f --input-file
		OutputFile       string // -o --output-file
		Timeout          int    // -t --timeout
		ConfirmExecution bool   // --confirm-execution
	)

	cmd := &cobra.Command{
		Use:     "get-resources",
		Short:   "Get Resources from a Lula Validation Manifest",
		Long:    "Get the JSON resources specified in a Lula Validation Manifest",
		Example: getResourcesHelp,
		RunE: func(cmd *cobra.Command, args []string) error {

			spinnerMessage := fmt.Sprintf("Getting Resources from %s", InputFile)
			spinner := message.NewProgressSpinner("%s", spinnerMessage)
			defer spinner.Stop()

			ctx := context.Background()

			// Read the validation data from STDIN or provided file
			validationBytes, err := ReadValidation(cmd, spinner, InputFile, Timeout)
			if err != nil {
				return fmt.Errorf("error reading validation: %v", err)
			}

			config, _ := cmd.Flags().GetStringSlice("set")
			message.Debug("command line 'set' flags: %s", config)

			output, err := DevTemplate(validationBytes, config)
			if err != nil {
				return fmt.Errorf("error templating validation: %v", err)
			}

			// add to debug logs accepting that this will print sensitive information?
			message.Debug(string(output))

			collection, err := DevGetResources(ctx, output, ConfirmExecution, spinner)

			// do not perform the write if there is nothing to write (likely error)
			if collection != nil {
				errWrite := types.WriteResources(collection, OutputFile)
				if errWrite != nil {
					message.Fatalf(errWrite, "error writing resources: %v", err)
				}
			}

			if err != nil {
				message.Fatalf(err, "error running dev get-resources: %v", err)
			}

			spinner.Success()

			return nil
		},
	}

	cmd.Flags().StringVarP(&InputFile, "input-file", "f", STDIN, "the path to a validation manifest file")
	cmd.Flags().StringVarP(&OutputFile, "output-file", "o", "", "the path to write the resources json")
	cmd.Flags().IntVarP(&Timeout, "timeout", "t", DEFAULT_TIMEOUT, "the timeout for stdin (in seconds, -1 for no timeout)")
	cmd.Flags().BoolVar(&ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of getting resources")

	return cmd

}

func DevGetResources(ctx context.Context, validationBytes []byte, confirmExecution bool, spinner *message.Spinner) (types.DomainResources, error) {
	lulaValidation, err := RunSingleValidation(ctx,
		validationBytes,
		types.ExecutionAllowed(confirmExecution),
		types.Interactive(RunInteractively),
		types.WithSpinner(spinner),
		types.GetResourcesOnly(true),
	)
	if err != nil {
		if lulaValidation.DomainResources != nil {
			return *lulaValidation.DomainResources, err
		}
		return nil, err
	}

	return *lulaValidation.DomainResources, nil
}
