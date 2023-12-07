package validate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type flags struct {
	AssessmentFile string // -a --assessment-file

}

var opts = &flags{}

var validateHelp = `
To validate on a cluster:
	lula validate ./oscal-component.yaml

To indicate a specific Assessment Results file to create or append to:
	lula validate ./oscal-component.yaml -a assessment-results.yaml
`

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "validate an OSCAL component definition",
	Long:    "Lula Validation of an OSCAL component definition",
	Example: validateHelp,
	RunE: func(cmd *cobra.Command, componentDefinitionPaths []string) error {
		// Conduct further error checking here (IE flags/arguments)
		if len(componentDefinitionPaths) == 0 {
			fmt.Println(cmd.Long)
			return errors.New("Path to OSCAL component definition(s) required")
		}

		results := types.ReportObject{
			FilePaths: componentDefinitionPaths,
		}

		// Primary expected path for validation of OSCAL documents
		err := ValidateOnPaths(&results)
		if err != nil {
			return fmt.Errorf("Validation error: %w\n", err)
		}

		report, err := oscal.GenerateAssessmentResults(&results)
		if err != nil {
			return fmt.Errorf("Generate error: %w\n", err)
		}

		// Write report(s) to file
		err = WriteReport(report, opts.AssessmentFile)
		if err != nil {
			return fmt.Errorf("Write error: %w\n", err)
		}
		return nil
	},
}

func ValidateCommand() *cobra.Command {

	// insert flag options here
	validateCmd.Flags().StringVarP(&opts.AssessmentFile, "assessment-file", "a", "", "the path to write assessment results. Creates a new file or appends to existing files")
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
func ValidateOnPaths(obj *types.ReportObject) error {
	// for each path
	for _, path := range obj.FilePaths {

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Printf("Path: %v does not exist - unable to digest document\n", path)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		compDef, err := oscal.NewOscalComponentDefinition(data)
		if err != nil {
			return err
		}

		err = ValidateOnCompDef(obj, compDef)
		if err != nil {
			return err
		}

		obj.UUID = compDef.UUID
	}

	return nil

}

// ValidateOnCompDef takes a single ComponentDefinition object
// It will perform a validation and add data to a referenced report object
func ValidateOnCompDef(obj *types.ReportObject, compDef oscalTypes.ComponentDefinition) error {

	// Populate a map[uuid]Validation into the validations
	obj.Validations = oscal.BackMatterToMap(compDef.BackMatter)

	// TODO: Is there a better location for context?
	ctx := context.Background()
	// Loops all the way down
	// Keeps track of UUID's for later reporting and relation
	for _, component := range compDef.Components {
		comp := types.Component{
			UUID: component.UUID,
		}
		for _, controlImplementation := range component.ControlImplementations {
			control := types.ControlImplementation{
				UUID: controlImplementation.UUID,
			}
			for _, implementedRequirement := range controlImplementation.ImplementedRequirements {

				impReq := types.ImplementedReq{
					UUID:        implementedRequirement.UUID,
					ControlId:   implementedRequirement.ControlId,
					Description: implementedRequirement.Description,
				}
				var pass, fail int
				// IF the implemented requirement contains a link - check for Lula Validation
				for _, link := range implementedRequirement.Links {
					var result types.Result
					var err error
					// Current identifier is the link text
					if link.Text == "Lula Validation" {
						// Remove the leading '#' from the UUID reference
						id := strings.Replace(link.Href, "#", "", 1)
						// Check if the link exists in our pre-populated map of validations
						if val, ok := obj.Validations[id]; ok {
							// If the validation has already been evaluated, use the result from the evaluation
							// Otherwise perform the validation
							if val.Evaluated {
								result = val.Result
							} else {
								result, err = ValidateOnTarget(ctx, val.Description)
								if err != nil {
									return err
								}
								// Store the result in the validation object
								val.Result = result
								val.Evaluated = true
								obj.Validations[id] = val
							}
						} else {
							return fmt.Errorf("Back matter Validation %v not found", id)
						}

						if result.Passing > 0 && result.Failing <= 0 {
							result.State = "satisfied"
						} else {
							result.State = "not-satisfied"
						}

						result.UUID = id

						pass += result.Passing
						fail += result.Failing

						impReq.Results = append(impReq.Results, result)

					}

				}

				// Using language from Assessment Results model for Target Objective Status State
				if pass > 0 && fail <= 0 {
					impReq.State = "satisfied"
				} else {
					impReq.State = "not-satisfied"
				}

				// TODO: convert to logging
				fmt.Printf("UUID: %v\n\tStatus: %v\n", impReq.UUID, impReq.State)

				control.ImplementedReqs = append(control.ImplementedReqs, impReq)
			}
			comp.ControlImplementations = append(comp.ControlImplementations, control)
		}
		obj.Components = append(obj.Components, comp)
	}

	return nil
}

// ValidateOnTarget takes a map[string]interface{}
// It will return a single Result
func ValidateOnTarget(ctx context.Context, target map[string]interface{}) (types.Result, error) {
	// simple conditional until more providers are introduced
	if provider, ok := target["provider"].(string); ok && provider == "opa" {
		fmt.Println("OPA provider validating...")
		results, err := opa.Validate(ctx, target["domain"].(string), target["payload"].(map[string]interface{}))
		if err != nil {
			return types.Result{}, err
		}
		return results, nil
	} else {
		return types.Result{}, errors.New("Unsupported provider")
	}

}

// This is the OSCAL document generation for final output.
// This should include some ability to consolidate controls met in multiple input documents under single control entries
// This should include fields that reference the source of the control to the original document ingested
func WriteReport(report oscalTypes.AssessmentResults, assessmentFilePath string) error {

	var fileName string
	var tempAssessment oscalTypes.AssessmentResults

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

			results := make([]oscalTypes.Result, 0)
			// append new results first - unfurl so as to allow multiple results in the future
			results = append(results, report.Results...)
			results = append(results, tempAssessment.Results...)
			tempAssessment.Results = results
			fileName = assessmentFilePath

		} else if os.IsNotExist(err) {
			// File does not exist
			tempAssessment = report
			fileName = assessmentFilePath
		} else {
			// Some other error occurred (permission issues, etc.)
			return err
		}

	} else {
		tempAssessment = report
		currentTime := time.Now()
		fileName = "assessment-results-" + currentTime.Format("01-02-2006-15:04:05") + ".yaml"
	}

	var b bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(tempAssessment)

	err := os.WriteFile(fileName, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
