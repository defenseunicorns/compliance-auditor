package validate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
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

func ValidateCommand() *cobra.Command {
	v := common.GetViper()

	var (
		outputFile          string
		inputFile           string
		target              string
		outputComponent     string
		renderTemplate      bool
		setOpts             []string
		confirmExecution    bool
		runNonInteractively bool
		saveResources       bool
		resourcesDir        string
	)

	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "validate an OSCAL component definition",
		Long:    "Lula Validation of an OSCAL component definition",
		Example: validateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: check if remote or local?
			_, err := os.Stat(inputFile)
			if os.IsNotExist(err) {
				return fmt.Errorf("input-file: %v does not exist - unable to digest document", inputFile)
			}
			// Update path if relative
			path := inputFile
			if filepath.IsLocal(inputFile) {
				path = filepath.Join(filepath.Dir(inputFile), filepath.Base(inputFile))
			}

			// If no output file is specified, get the default output file
			if outputFile == "" {
				outputFile = getDefaultOutputFile(inputFile)
			}

			// Set values if save-resources is desired
			if saveResources {
				resourcesDir = filepath.Join(filepath.Dir(outputFile))
			}

			// check if setOpts without rendertemplate
			if len(setOpts) > 0 && !renderTemplate {
				message.Warnf("set-opts are only valid when render-template is true")
			}

			// Check if output-component is set
			rt := template.MASKED
			if outputComponent != "" {
				if strings.ToLower(outputComponent) == "full" {
					rt = template.ALL
				}
			}

			opts := []validation.Option{
				validation.WithComponentDefinitionFromPath(path),
				validation.WithTemplateRenderer(renderTemplate, setOpts),
				validation.WithAllowExecution(confirmExecution, runNonInteractively),
				validation.WithResourcesDir(saveResources, resourcesDir),
			}

			// Currently not using context, but could/should be used for logging things, etc.
			validationCtx, err := validation.New(context.Background(), opts...)
			if err != nil {
				return fmt.Errorf("error creating validation context: %v", err)
			}

			// Get component definition data, update validation context with rendered component definition
			componentOut, err := validationCtx.RenderComponentDefinition(filepath.Ext(inputFile), rt)
			if err != nil {
				return fmt.Errorf("error rendering template: %v", err)
			}

			assessment, err := validationCtx.ExecuteOSCALValidation(target) //ValidateOnPath(inputFile, target)
			if err != nil {
				message.Fatalf(err, "OSCAL Validation error: %s", err)
			}

			var model = oscalTypes_1_1_2.OscalModels{
				AssessmentResults: assessment,
			}

			// Write the rendered/output component-definition to file
			if outputComponent != "" {
				outputComponentFile := filepath.Join(filepath.Dir(outputFile), "component-rendered"+filepath.Ext(inputFile))
				err := os.WriteFile(outputComponentFile, componentOut, 0644)
				if err != nil {
					return err
				}
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
	cmd.Flags().BoolVar(&renderTemplate, "render-template", false, "render any templated data from the input file")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")
	cmd.Flags().StringVar(&outputComponent, "output-component", "", "specify the desired output component string. Options are 'masked', 'full'")

	return cmd
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
func ValidateOnPath(path string, target string) (assessmentResult *oscalTypes_1_1_2.AssessmentResults, err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return assessmentResult, fmt.Errorf("path: %v does not exist - unable to digest document", path)
	}

	oscalModel, err := composition.ComposeFromPath(path)
	if err != nil {
		return assessmentResult, err
	}

	if oscalModel.ComponentDefinition == nil {
		return assessmentResult, fmt.Errorf("component definition is nil")
	}

	results, err := ValidateOnCompDef(oscalModel.ComponentDefinition, target)
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
func ValidateOnCompDef(compDef *oscalTypes_1_1_2.ComponentDefinition, target string) (results []oscalTypes_1_1_2.Result, err error) {
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
			findings, observations, err := ValidateOnControlImplementations(&controlImplementation, validationStore, target)
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
			findings, observations, err := ValidateOnControlImplementations(&controlImplementation, validationStore, source)
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

func ValidateOnControlImplementations(controlImplementations *[]oscalTypes_1_1_2.ControlImplementationSet, validationStore *validationstore.ValidationStore, target string) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {

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
	observations := validationStore.RunValidations(ConfirmExecution, SaveResources, ResourcesDir)
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

// GetDefaultOutputFile returns the default output file name
func getDefaultOutputFile(inputFile string) string {
	dirPath := filepath.Dir(inputFile)
	filename := "assessment-results" + filepath.Ext(inputFile)

	return filepath.Join(dirPath, filename)
}
