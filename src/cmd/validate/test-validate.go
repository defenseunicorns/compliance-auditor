package validate

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/validation"
)

// Take input file, create a new producer and consumer?
// If doing an compdef with several targets.. do that loop here?

func TestValidateCommand() *cobra.Command {
	v := common.GetViper()

	var (
		outputFile          string
		inputFile           string
		target              string
		setOpts             []string
		simple              bool
		silent              bool
		confirmExecution    bool
		runNonInteractively bool
		saveResources       bool
	)

	cmd := &cobra.Command{
		Use:     "test-validate",
		Aliases: []string{"tv"},
		Short:   "test-validate <anything>",
		Long:    "Lula Validation of <anything>",
		Example: validateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			consumerType := "assessment-results" // Default consumer type
			if simple {
				consumerType = "simple"
			}

			// Get the inputs, set up the producer/consumer
			inputFileBytes, err := os.ReadFile(inputFile)
			if err != nil {
				return err
			}

			if outputFile != "" {
				outputFile = getDefaultOutputFile(inputFile)
			}

			producer, err := validation.ResolveProducer(inputFileBytes, inputFile, target)
			if err != nil {
				return err
			}

			consumer, err := validation.ResolveConsumer(consumerType, outputFile)
			if err != nil {
				return err
			}

			// Set up the validation
			opts := []validation.Option{
				validation.WithAllowExecution(confirmExecution, runNonInteractively),
				// Other options...
			}

			validator, err := validation.New(producer, consumer, opts...)
			if err != nil {
				return err
			}

			// Get stats
			if !silent {
				numReqts, numVals, numExeVals := validator.GetStats()
				message.Title("\n🔍 Collecting Requirements and Validations for: ", inputFile)
				message.Infof("%d Requirements: ", numReqts)
				message.Infof("%d Validations: ", numVals)
				if numExeVals > 0 {
					message.Warnf("%d Executable Validations: ", numExeVals)
				}

			}

			// TODO: Request confirmation for execution, if needed
			runExecutableValidations := true

			err = validator.ExecuteValidations(cmd.Context(), runExecutableValidations)
			if err != nil {
				return err
			}

			return validator.YieldResults()
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to write assessment results. Creates a new file or appends to existing files")
	cmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	err := cmd.MarkFlagRequired("input-file")
	if err != nil {
		message.Fatal(err, "error initializing upgrade command flags")
	}
	cmd.Flags().StringVarP(&target, "target", "t", v.GetString(common.VTarget), "the specific control implementations or framework to validate against")
	cmd.Flags().BoolVar(&simple, "simple", false, "simple output")
	cmd.Flags().BoolVar(&silent, "silent", false, "if set, no output will be printed")
	cmd.Flags().BoolVar(&confirmExecution, "confirm-execution", false, "confirm execution scripts run as part of the validation")
	cmd.Flags().BoolVar(&runNonInteractively, "non-interactive", false, "run the command non-interactively")
	cmd.Flags().BoolVar(&saveResources, "save-resources", false, "saves the resources to 'resources' directory at assessment-results level")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")

	return cmd
}
