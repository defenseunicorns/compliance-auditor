package validate

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
)

var validateHelp = `
To validate on a cluster:
	lula validate -f ./oscal-component.yaml
To indicate a specific Assessment Results file to create or append to:
	lula validate -f ./oscal-component.yaml -o assessment-results.yaml
To target a specific control-implementation source / standard/ framework
	lula validate -f ./oscal-component.yaml -t critical
To run validations and automatically confirm execution
	lula dev validate -f ./oscal-component.yaml --confirm-execution
To run validations non-interactively (no execution)
	lula dev validate -f ./oscal-component.yaml --non-interactive
`

var (
	ErrValidating       = errors.New("error validating")
	ErrInvalidOut       = errors.New("error invalid OSCAL model at output")
	ErrWritingComponent = errors.New("error writing component to file")
	ErrCreatingVCtx     = errors.New("error creating validation context")
	ErrCreatingCCtx     = errors.New("error creating composition context")
)

func ValidateCommand() *cobra.Command {
	v := common.GetViper()

	var (
		outputFile          string
		inputFile           string
		target              string
		setOpts             []string
		confirmExecution    bool
		runNonInteractively bool
		saveResources       bool
	)

	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "validate an OSCAL component definition",
		Long:    "Lula Validation of an OSCAL component definition",
		Example: validateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {

			// If no output file is specified, get the default output file
			if outputFile == "" {
				outputFile = getDefaultOutputFile(inputFile)
			}

			// Check if output file contains a valid OSCAL model
			_, err := oscal.ValidOSCALModelAtPath(outputFile)
			if err != nil {
				return fmt.Errorf("invalid OSCAL model at output: %v", err)
			}

			// Set up the composition context
			compositionCtx, err := composition.New(
				composition.WithModelFromLocalPath(inputFile),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", common.TemplateConstants, common.TemplateVariables, setOpts),
			)
			if err != nil {
				return fmt.Errorf("error creating composition context: %v", err)
			}

			// Set up the validation context
			validationCtx, err := validation.New(
				validation.WithCompositionContext(compositionCtx, inputFile),
				validation.WithResourcesDir(saveResources, filepath.Dir(outputFile)),
				validation.WithAllowExecution(confirmExecution, runNonInteractively),
			)
			if err != nil {
				return fmt.Errorf("error creating validation context: %v", err)
			}

			ctx := context.WithValue(cmd.Context(), types.LulaValidationWorkDir, filepath.Dir(inputFile))
			assessmentResults, err := validationCtx.ValidateOnPath(ctx, inputFile, target)
			if err != nil {
				return fmt.Errorf("error validating on path: %v", err)
			}

			if assessmentResults == nil {
				return fmt.Errorf("assessment results are nil")
			}

			var model = oscalTypes_1_1_2.OscalModels{
				AssessmentResults: assessmentResults,
			}

			// Write the assessment results to file
			err = oscal.WriteOscalModel(outputFile, &model)
			if err != nil {
				return fmt.Errorf("error writing component to file: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to write assessment results. Creates a new file or appends to existing files")
	cmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	cmd.MarkFlagRequired("input-file")
	cmd.Flags().StringVarP(&target, "target", "t", v.GetString(common.VTarget), "the specific control implementations or framework to validate against")
	cmd.Flags().BoolVar(&confirmExecution, "confirm-execution", false, "confirm execution scripts run as part of the validation")
	cmd.Flags().BoolVar(&runNonInteractively, "non-interactive", false, "run the command non-interactively")
	cmd.Flags().BoolVar(&saveResources, "save-resources", false, "saves the resources to 'resources' directory at assessment-results level")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")

	return cmd
}

// getDefaultOutputFile returns the default output file name
func getDefaultOutputFile(inputFile string) string {
	dirPath := filepath.Dir(inputFile)
	filename := "assessment-results" + filepath.Ext(inputFile)

	return filepath.Join(dirPath, filename)
}

/*
	To tell the validation story:
		Lula is currently evaluating controls identified in the Implemented-Requirements of a component-definition.
		We would then be looking to retain information that may be required for relation of component-definition (input) to an assessment-results (output).
		In order to get there - we have to traverse and possibly track UUIDs at a minimum:

		Lula accepts 1 -> N paths to OSCAL component-definition files
		For each component definition:
			There are 1 -> N Components
			For each component:
				There are 1 -> N control-Implementations
				For each control-implementation:
					There are 1-> N implemented-requirements
					For each implemented-requirement:
						There are 1 -> N validations
							This allows for breaking complex query and policy into smaller  chunks
						Validations are evaluated individually with passing/failing resources
					Pass/Fail results from all validations is evaluated for a pass/fail status in the report

	As such, building a ReportObject to collect and retain the relational information could be preferred

*/
