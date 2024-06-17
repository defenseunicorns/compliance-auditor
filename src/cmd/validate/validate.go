package validate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type flags struct {
	OutputFile string // -o --output-file
	InputFile  string // -f --input-file
}

var opts = &flags{}
var ConfirmExecution bool    // --confirm-execution
var RunNonInteractively bool // --non-interactive

var validateHelp = `
To validate on a cluster:
	lula validate -f ./oscal-component.yaml
To indicate a specific Assessment Results file to create or append to:
	lula validate -f ./oscal-component.yaml -o assessment-results.yaml
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
		if opts.InputFile == "" {
			message.Fatal(errors.New("flag input-file is not set"),
				"Please specify an input file with the -f flag")
		}

		if err := files.IsJsonOrYaml(opts.InputFile); err != nil {
			message.Fatalf(err, "Invalid file extension: %s, requires .json or .yaml", opts.InputFile)
		}

		findings, observations, err := ValidateOnPath(opts.InputFile)
		if err != nil {
			message.Fatalf(err, "Validation error: %s", err)
		}

		report, err := oscal.GenerateAssessmentResults(findings, observations)
		if err != nil {
			message.Fatalf(err, "Generate error")
		}

		var model = oscalTypes_1_1_2.OscalModels{
			AssessmentResults: report,
		}

		// Write the assessment results to file
		err = oscal.WriteOscalModel(opts.OutputFile, &model)
		if err != nil {
			message.Fatalf(err, "error writing component to file")
		}
	},
}

func ValidateCommand() *cobra.Command {

	// insert flag options here
	validateCmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "", "the path to write assessment results. Creates a new file or appends to existing files")
	validateCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
	validateCmd.Flags().BoolVar(&ConfirmExecution, "confirm-execution", false, "confirm execution scripts run as part of the validation")
	validateCmd.Flags().BoolVar(&RunNonInteractively, "non-interactive", false, "run the command non-interactively")
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
func ValidateOnPath(path string) (findingMap map[string]oscalTypes_1_1_2.Finding, observations []oscalTypes_1_1_2.Observation, err error) {

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return findingMap, observations, fmt.Errorf("path: %v does not exist - unable to digest document", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return findingMap, observations, err
	}

	// Change Cwd to the directory of the component definition
	dirPath := filepath.Dir(path)
	message.Infof("changing cwd to %s", dirPath)
	resetCwd, err := common.SetCwdToFileDir(dirPath)
	if err != nil {
		return findingMap, observations, err
	}
	defer resetCwd()

	compDef, err := oscal.NewOscalComponentDefinition(data)
	if err != nil {
		return findingMap, observations, err
	}

	findingMap, observations, err = ValidateOnCompDef(compDef)
	if err != nil {
		return findingMap, observations, err
	}

	return findingMap, observations, err

}

// ValidateOnCompDef takes a single ComponentDefinition object
// It will perform a validation and add data to a referenced report object
func ValidateOnCompDef(compDef *oscalTypes_1_1_2.ComponentDefinition) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {
	err := composition.ComposeComponentDefinitions(compDef)
	if err != nil {
		return nil, nil, err

	}

	// Create a validation store from the back-matter if it exists
	validationStore := validationstore.NewValidationStoreFromBackMatter(*compDef.BackMatter)

	// Initialize findings and observations
	findings := make(map[string]oscalTypes_1_1_2.Finding)
	observations := make([]oscalTypes_1_1_2.Observation, 0)

	if *compDef.Components == nil {
		return findings, observations, fmt.Errorf("no components found in component definition")
	}

	// Create requirement store for all implemented requirements
	requirementStore := requirementstore.NewRequirementStore(compDef)
	message.Title("\n🔍 Collecting Requirements and Validations", "")
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
	observations = validationStore.RunValidations(ConfirmExecution)
	message.Title("\n💡 Findings", "")
	findings = requirementStore.GenerateFindings(validationStore)

	return findings, observations, nil
}

// This is the OSCAL document generation for final output.
// This should include some ability to consolidate controls met in multiple input documents under single control entries
// This should include fields that reference the source of the control to the original document ingested
// TODO: This is unused - remove?
func WriteReport(report oscalTypes_1_1_2.AssessmentResults, assessmentFilePath string) error {

	var fileName string
	var tempAssessment *oscalTypes_1_1_2.AssessmentResults

	if assessmentFilePath != "" {

		_, err := os.Stat(assessmentFilePath)
		if err == nil {
			// File does exist
			data, err := os.ReadFile(assessmentFilePath)
			if err != nil {
				return err
			}

			tempAssessment, err = oscal.NewAssessmentResults(data)
			if err != nil {
				return err
			}

			results := make([]oscalTypes_1_1_2.Result, 0)
			// append new results first - unfurl so as to allow multiple results in the future
			results = append(results, report.Results...)
			results = append(results, tempAssessment.Results...)
			tempAssessment.Results = results
			fileName = assessmentFilePath

		} else if os.IsNotExist(err) {
			// File does not exist
			tempAssessment = &report
			fileName = assessmentFilePath
		} else {
			// Some other error occurred (permission issues, etc.)
			return err
		}

	} else {
		tempAssessment = &report
		currentTime := time.Now()
		fileName = "assessment-results-" + currentTime.Format("01-02-206-15:04:05") + ".yaml"
	}

	var b bytes.Buffer

	var sar = oscalTypes_1_1_2.OscalModels{
		AssessmentResults: tempAssessment,
	}

	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(sar)

	message.Infof("Writing Security Assessment Results to: %s", fileName)

	err := os.WriteFile(fileName, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
