package validate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
)

type flags struct {
	OutputFile string // -o --output-file
	InputFile  string // -f --input-file
	Target     string // -t --target
}

var opts = &flags{}
var ConfirmExecution bool    // --confirm-execution
var RunNonInteractively bool // --non-interactive
var SaveResources bool       // --save-resources
var ResourcesDir string

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

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "validate an OSCAL component definition",
	Long:    "Lula Validation of an OSCAL component definition",
	Example: validateHelp,
	Run: func(cmd *cobra.Command, componentDefinitionPath []string) {
		outputFile := opts.OutputFile
		if outputFile == "" {
			outputFile = getDefaultOutputFile(opts.InputFile)
		}

		// Check if output file contains a valid OSCAL model
		_, err := oscal.ValidOSCALModelAtPath(outputFile)
		if err != nil {
			message.Fatalf(err, "Output file %s is not a valid OSCAL model: %v", outputFile, err)
		}

		if SaveResources {
			ResourcesDir = filepath.Join(filepath.Dir(outputFile))
		}

		if err := files.IsJsonOrYaml(opts.InputFile); err != nil {
			message.Fatalf(err, "Invalid file extension: %s, requires .json or .yaml", opts.InputFile)
		}

		ctx := context.WithValue(cmd.Context(), types.LulaValidationWorkDir, filepath.Dir(opts.InputFile))
		assessment, err := ValidateOnPath(ctx, opts.InputFile, opts.Target)
		if err != nil {
			message.Fatalf(err, "Validation error: %s", err)
		}

		var model = oscalTypes_1_1_2.OscalModels{
			AssessmentResults: assessment,
		}

		// Write the assessment results to file
		err = oscal.WriteOscalModel(outputFile, &model)
		if err != nil {
			message.Fatalf(err, "error writing component to file")
		}
	},
}

func init() {
	v := common.InitViper()

	validateCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write assessment results. Creates a new file or appends to existing files")
	validateCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	validateCmd.MarkFlagRequired("input-file")
	validateCmd.Flags().StringVarP(&opts.Target, "target", "t", v.GetString(common.VTarget), "the specific control implementations or framework to validate against")
	validateCmd.Flags().BoolVar(&ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of the validation")
	validateCmd.Flags().BoolVar(&RunNonInteractively, "non-interactive", false, "run the command non-interactively")
	validateCmd.Flags().BoolVar(&SaveResources, "save-resources", false, "saves the resources to 'resources' directory at assessment-results level")
}

func ValidateCommand() *cobra.Command {
	return validateCmd
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

// ValidateOnPath takes 1 -> N paths to OSCAL component-definition files
// It will then read those files to perform validation and return an ResultObject
func ValidateOnPath(ctx context.Context, path string, target string) (assessmentResult *oscalTypes_1_1_2.AssessmentResults, err error) {

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return assessmentResult, fmt.Errorf("path: %v does not exist - unable to digest document", path)
	}

	compositionCtx, err := composition.New(composition.WithModelFromLocalPath(path))
	if err != nil {
		return nil, fmt.Errorf("error creating composition context: %v", err)
	}

	oscalModel, err := compositionCtx.ComposeFromPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("error composing model: %v", err)
	}

	if oscalModel.ComponentDefinition == nil {
		return assessmentResult, fmt.Errorf("component definition is nil")
	}

	results, err := ValidateOnCompDef(ctx, oscalModel.ComponentDefinition, target)
	if err != nil {
		return assessmentResult, err
	}

	assessmentResult, err = oscal.GenerateAssessmentResults(results)
	if err != nil {
		return assessmentResult, err
	}

	return assessmentResult, nil

}

// ValidateOnCompDef takes a single ComponentDefinition object
// It will perform a validation and return a slice of results that can be written to an assessment-results object
func ValidateOnCompDef(ctx context.Context, compDef *oscalTypes_1_1_2.ComponentDefinition, target string) (results []oscalTypes_1_1_2.Result, err error) {
	if compDef == nil {
		return results, fmt.Errorf("cannot validate a component definition that is nil")
	}

	if *compDef.Components == nil {
		return results, fmt.Errorf("no components found in component definition")
	}

	// Create a validation store from the back-matter if it exists
	validationStore := validationstore.NewValidationStoreFromBackMatter(*compDef.BackMatter)

	// Create a map of control implementations from the component definition
	// This combines all same source/framework control implementations into an []Control-Implementation
	controlImplementations := oscal.FilterControlImplementations(compDef)

	if len(controlImplementations) == 0 {
		return results, fmt.Errorf("no control implementations found in component definition")
	}

	// target one specific controlImplementation
	// this could be either a framework or source property
	// this will only produce a single result
	if target != "" {
		if controlImplementation, ok := controlImplementations[target]; ok {
			findings, observations, err := ValidateOnControlImplementations(ctx, &controlImplementation, validationStore, target)
			if err != nil {
				return results, err
			}
			result, err := oscal.CreateResult(findings, observations)
			if err != nil {
				return results, err
			}
			// add/update the source to the result props - make source = framework or omit?
			oscal.UpdateProps("target", oscal.LULA_NAMESPACE, target, result.Props)
			results = append(results, result)
		} else {
			return results, fmt.Errorf("target %s not found", target)
		}
	} else {
		// default behavior - create a result for each unique source/framework
		// loop over the controlImplementations map & validate
		// we lose context of source if not contained within the loop
		for source, controlImplementation := range controlImplementations {
			findings, observations, err := ValidateOnControlImplementations(ctx, &controlImplementation, validationStore, source)
			if err != nil {
				return results, err
			}
			result, err := oscal.CreateResult(findings, observations)
			if err != nil {
				return results, err
			}
			// add/update the source to the result props
			oscal.UpdateProps("target", oscal.LULA_NAMESPACE, source, result.Props)
			results = append(results, result)
		}
	}

	return results, nil

}

func ValidateOnControlImplementations(ctx context.Context, controlImplementations *[]oscalTypes_1_1_2.ControlImplementationSet, validationStore *validationstore.ValidationStore, target string) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {

	// Create requirement store for all implemented requirements
	requirementStore := requirementstore.NewRequirementStore(controlImplementations)
	message.Title("\n🔍 Collecting Requirements and Validations for Target: ", target)
	requirementStore.ResolveLulaValidations(validationStore)
	reqtStats := requirementStore.GetStats(validationStore)
	message.Infof("Found %d Implemented Requirements", reqtStats.TotalRequirements)
	message.Infof("Found %d runnable Lula Validations", reqtStats.TotalValidations)

	// Check if validations perform execution actions
	if reqtStats.ExecutableValidations {
		message.Warnf(reqtStats.ExecutableValidationsMsg)
		if !ConfirmExecution {
			if !RunNonInteractively {
				ConfirmExecution = message.PromptForConfirmation(nil)
			}
			if !ConfirmExecution {
				// Break or just skip those those validations?
				message.Infof("Validations requiring execution will not be run")
				// message.Fatalf(errors.New("execution not confirmed"), "Exiting validation")
			}
		}
	}

	// Run Lula validations and generate observations & findings
	message.Title("\n📐 Running Validations", "")
	observations := validationStore.RunValidations(ctx, ConfirmExecution, SaveResources, ResourcesDir)
	message.Title("\n💡 Findings", "")
	findings := requirementStore.GenerateFindings(validationStore)

	// Print findings here to prevent repetition of findings in the output
	header := []string{"Control ID", "Status"}
	rows := make([][]string, 0)
	columnSize := []int{20, 25}

	for id, finding := range findings {
		rows = append(rows, []string{
			id, finding.Target.Status.State,
		})
	}

	if len(rows) != 0 {
		message.Table(header, rows, columnSize)
	}

	return findings, observations, nil
}

// getDefaultOutputFile returns the default output file name and checks if the file already exists
func getDefaultOutputFile(inputFile string) string {
	dirPath := filepath.Dir(inputFile)
	filename := "assessment-results" + filepath.Ext(inputFile)

	return filepath.Join(dirPath, filename)
}
