package validate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type flags struct {
	AssessmentFile string // -a --assessment-file
	InputFile      string // -f --input-file
}

var opts = &flags{}

var validateHelp = `
To validate on a cluster:
	lula validate -f ./oscal-component.yaml

To indicate a specific Assessment Results file to create or append to:
	lula validate -f ./oscal-component.yaml -a assessment-results.yaml
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
		// Primary expected path for validation of OSCAL documents
		findings, observations, err := ValidateOnPath(opts.InputFile)
		if err != nil {
			message.Fatalf(err, "Validation error: %s", err)
		}

		report, err := oscal.GenerateAssessmentResults(findings, observations)
		if err != nil {
			message.Fatalf(err, "Generate error")
		}

		// Write report(s) to file
		err = WriteReport(report, opts.AssessmentFile)
		if err != nil {
			message.Fatalf(err, "Write error")
		}
	},
}

func ValidateCommand() *cobra.Command {

	// insert flag options here
	validateCmd.Flags().StringVarP(&opts.AssessmentFile, "assessment-file", "a", "", "the path to write assessment results. Creates a new file or appends to existing files")
	validateCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to the target OSCAL component definition")
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
		return findingMap, observations, fmt.Errorf("Path: %v does not exist - unable to digest document\n", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return findingMap, observations, err
	}

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
func ValidateOnCompDef(compDef oscalTypes_1_1_2.ComponentDefinition) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {

	// Populate a map[uuid]Validation into the validations
	validations := oscal.BackMatterToMap(compDef.BackMatter)

	// TODO: Is there a better location for context?
	ctx := context.Background()
	// Loops all the way down

	findings := make(map[string]oscalTypes_1_1_2.Finding)
	observations := make([]oscalTypes_1_1_2.Observation, 0)

	for _, component := range compDef.Components {
		for _, controlImplementation := range component.ControlImplementations {
			rfc3339Time := time.Now()
			for _, implementedRequirement := range controlImplementation.ImplementedRequirements {
				spinner := message.NewProgressSpinner("Validating Implemented Requirement - %s", implementedRequirement.UUID)
				defer spinner.Stop()

				// This should produce a finding - check if an existing finding for the control-id has been processed
				var finding oscalTypes_1_1_2.Finding
				tempObservations := make([]oscalTypes_1_1_2.Observation, 0)
				relatedObservations := make([]oscalTypes_1_1_2.RelatedObservation, 0)

				if _, ok := findings[implementedRequirement.ControlId]; ok {
					finding = findings[implementedRequirement.ControlId]
				} else {
					finding = oscalTypes_1_1_2.Finding{
						UUID:        uuid.NewUUID(),
						Title:       fmt.Sprintf("Validation Result - Component:%s / Control Implementation: %s / Control:  %s", component.UUID, controlImplementation.UUID, implementedRequirement.ControlId),
						Description: implementedRequirement.Description,
					}
				}

				var pass, fail int
				// IF the implemented requirement contains a link - check for Lula Validation

				for _, link := range implementedRequirement.Links {
					var result types.Result
					var err error
					// Current identifier is the link text
					if link.Text == "Lula Validation" {
						sharedUuid := uuid.NewUUID()
						observation := oscalTypes_1_1_2.Observation{
							Collected: rfc3339Time,
							Methods:   []string{"TEST"},
							UUID:      sharedUuid,
						}
						// Remove the leading '#' from the UUID reference
						id := strings.Replace(link.Href, "#", "", 1)
						observation.Description = fmt.Sprintf("[TEST] %s - %s\n", implementedRequirement.ControlId, id)
						// Check if the link exists in our pre-populated map of validations
						if val, ok := validations[id]; ok {
							// If the validation has already been evaluated, use the result from the evaluation - otherwise perform the validation
							if val.Evaluated {
								result = val.Result
							} else {
								result, err = ValidateOnTarget(ctx, id, val.Description)
								if err != nil {
									return map[string]oscalTypes_1_1_2.Finding{}, []oscalTypes_1_1_2.Observation{}, err
								}
								// Store the result in the validation object
								val.Result = result
								val.Evaluated = true
								validations[id] = val
							}
						} else {
							return map[string]oscalTypes_1_1_2.Finding{}, []oscalTypes_1_1_2.Observation{}, fmt.Errorf("Back matter Validation %v not found", id)
						}

						// Individual result state
						if result.Passing > 0 && result.Failing <= 0 {
							result.State = "satisfied"
						} else {
							result.State = "not-satisfied"
						}

						// Add remarks if Result has Observations
						var remarks string
						if len(result.Observations) > 0 {
							for k, v := range result.Observations {
								remarks += fmt.Sprintf("%s: %s\n", k, v)
							}
						}

						observation.RelevantEvidence = []oscalTypes_1_1_2.RelevantEvidence{
							{
								Description: fmt.Sprintf("Result: %s\n", result.State),
								Remarks:     remarks,
							},
						}

						relatedObservation := oscalTypes_1_1_2.RelatedObservation{
							ObservationUuid: sharedUuid,
						}

						pass += result.Passing
						fail += result.Failing

						// Coalesce slices and objects
						relatedObservations = append(relatedObservations, relatedObservation)
						tempObservations = append(tempObservations, observation)
					}

				}

				// Using language from Assessment Results model for Target Objective Status State
				var state string
				if finding.Target.Status.State == "not-satisfied" {
					state = "not-satisfied"
				} else if pass > 0 && fail <= 0 {
					state = "satisfied"
				} else {
					state = "not-satisfied"
				}

				message.Infof("UUID: %v", finding.UUID)
				message.Infof("    Status: %v", state)

				finding.Target = oscalTypes_1_1_2.FindingTarget{
					Status: oscalTypes_1_1_2.ObjectiveStatus{
						State: state,
					},
					TargetId: implementedRequirement.ControlId,
					Type:     "objective-id",
				}

				finding.RelatedObservations = relatedObservations

				findings[implementedRequirement.ControlId] = finding
				observations = append(observations, tempObservations...)
				spinner.Success()
			}
		}
	}

	return findings, observations, nil
}

// ValidateOnTarget takes a map[string]interface{}
// It will return a single Result
func ValidateOnTarget(ctx context.Context, id string, target map[string]interface{}) (types.Result, error) {
	// simple conditional until more providers are introduced
	if provider, ok := target["provider"].(string); ok && provider == "opa" {
		message.Debugf("OPA provider validating %s", id)
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
func WriteReport(report oscalTypes_1_1_2.AssessmentResults, assessmentFilePath string) error {

	var fileName string
	var tempAssessment oscalTypes_1_1_2.AssessmentResults

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
