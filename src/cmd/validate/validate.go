package validate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	oscalTypes "github.com/defenseunicorns/lula/src/types/oscal"
	"github.com/spf13/cobra"
	yaml1 "sigs.k8s.io/yaml"
)

var validateHelp = `
To validate on a cluster:
	lula validate ./oscal-component.yaml
`

var cluster bool

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "validate",
	Long:    "Lula Validation for compliance with established policy",
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

		// Convert Results to expected report format
		report, err := GenerateReportFromResults(&results)
		if err != nil {
			return fmt.Errorf("Generate error: %w\n", err)
		}

		// Write report(s) to file
		err = WriteReport(report)
		if err != nil {
			return fmt.Errorf("Write error: %w\n", err)
		}
		return nil
	},
}

func ValidateCommand() *cobra.Command {

	// insert flag options here
	return ValidateCmd
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
		err = ValidateOnCompDef(obj, compDef)

		obj.UUID = compDef.UUID
	}

	return nil

}

// ValidateOnCompDef takes a single ComponentDefinition object
// It will perform a validation and add data to a referenced report object
func ValidateOnCompDef(obj *types.ReportObject, compDef oscalTypes.ComponentDefinition) error {

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
				for _, target := range implementedRequirement.Rules {
					result, err := ValidateOnTarget(ctx, target)

					if err != nil {
						return err
					}
					pass += result.Passing
					fail += result.Failing
					impReq.Results = append(impReq.Results, result)
				}

				if pass > 0 && fail <= 0 {
					impReq.Status = "Pass"
				} else if pass == 0 && fail == 0 {
					impReq.Status = "Not Evaluated"
				} else {
					impReq.Status = "Fail"
				}

				// TODO: convert to logging
				fmt.Printf("UUID: %v\n\tResources Passing: %v\n\tResources Failing: %v\n\tStatus: %v\n", impReq.UUID, pass, fail, impReq.Status)

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

// TODO: this needs to evolve quite a bit - should transform a ReportObject into an OSCAL Model
// Specifically re-traversing the layers
func GenerateReportFromResults(results *types.ReportObject) ([]types.ComplianceReport, error) {
	var complianceReports []types.ComplianceReport
	// component-definition -> component -> control-implementation -> implemented-requirements -> targets

	for _, component := range results.Components {
		for _, control := range component.ControlImplementations {
			for _, impReq := range control.ImplementedReqs {
				currentReport := types.ComplianceReport{
					UUID:        impReq.UUID,
					ControlId:   impReq.ControlId,
					Description: impReq.Description,
					Result:      impReq.Status,
				}

				complianceReports = append(complianceReports, currentReport)
			}
		}
	}

	return complianceReports, nil
}

// This is the OSCAL document generation for final output.
// This should include some ability to consolidate controls met in multiple input documents under single control entries
// This should include fields that reference the source of the control to the original document ingested
func WriteReport(compiledReport []types.ComplianceReport) error {
	reportData, err := yaml1.Marshal(&compiledReport)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	fileName := "compliance_report-" + currentTime.Format("01-02-2006-15:04:05") + ".yaml"

	err = os.WriteFile(fileName, reportData, 0644)
	if err != nil {
		return err
	}
	return nil
}
