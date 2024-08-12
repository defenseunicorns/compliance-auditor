package report

import (
	"errors"
	"os"
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"gopkg.in/yaml.v3"
)

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
	FileFormat string // --file-format
}

type ReportData struct {
	ComponentDefinition *ComponentDefinitionReportData `json:"componentDefinition,omitempty" yaml:"componentDefinition,omitempty"`
}

type ComponentDefinitionReportData struct {
	Title string `json:"title" yaml:"title"`
	ControlIDMapped int `json:"control ID mapped" yaml:"control ID mapped"`
	ControlIDFramework map[string]int `json:"controlIDFramework" yaml:"controlIDFramework"`
	ControlIDSource int `json:"control ID from source" yaml:"control ID from source"`
}

var opts = &flags{}

var reportHelp = `
To create a new report:
lula report -f oscal-component-definition.yaml -o report.yaml
`

var reportCmd = &cobra.Command{
	Use:     "report",
	Hidden:  false,
	Aliases: []string{"r"},
	Short:   "Build a compliance report",
	Example: reportHelp,
	Run: func(_ *cobra.Command, args []string) {
		if opts.InputFile == "" {
			message.Fatal(errors.New("flag input-file is not set"),
				"Please specify an input file with the -f flag")
		}

		oscalModels, err := readOSCALFile(opts.InputFile)
		if err != nil {
			message.Fatal(err, "Failed to unmarshall OSCAL from file")
		}

		modelType, err := determineOSCALModel(oscalModels)
		if err != nil {
			message.Fatal(err, "Unable to determine OSCAL model type")
		}

		// REMOVE: Using for testing if I am able to determine models correctly
		fmt.Println(modelType)

		var reportModelErr error

		switch modelType {
		case "catalog":
			fmt.Println("reporting does not create reports for catalogs at this time")
		case "profile":
			fmt.Println("reporting does not create reports for profile at this time")
		case "component-definition":
			reportModelErr = handleComponentDefinition(oscalModels.ComponentDefinition, opts.OutputFile, opts.FileFormat)
		case "assessment-plan":
			fmt.Println("reporting does not create reports for assessment plans at this time")
		case "assessment-results":
			fmt.Println("reporting does not create reports for assessment results at this time")
		case "system-security-plan":
			fmt.Println("reporting does not create reports for system security plans at this time")
		case "plan-of-action-and-milestones":
			fmt.Println("reporting does not create reports for plan of action and milestones at this time")
		default:
			message.Fatal(fmt.Errorf("unknown OSCAL model type: %s", modelType), "Failed to process OSCAL file")
		}

		if reportModelErr != nil {
			message.Fatal(reportModelErr, "failed to create report")
		}

	},
}

func ReportCommand() *cobra.Command {

	reportFlags()

	return reportCmd
}

func reportFlags() {
	reportFlags := reportCmd.PersistentFlags()

	reportFlags.StringVarP(&opts.InputFile, "input-file", "f", "", "Path to an OSCAL file")
	reportFlags.StringVarP(&opts.OutputFile, "output-file", "o", "", "Path and Name to an output file")
	reportFlags.StringVar(&opts.FileFormat, "file-format", "yaml", "File format of output file")

}

func readOSCALFile(fileName string ) (oscalTypes_1_1_2.OscalModels, error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	// Read the YAML file using os.ReadFile
	yamlData, err := os.ReadFile(fileName)
	if err != nil {
		return oscalModels, err
	}

	err = yaml.Unmarshal(yamlData, &oscalModels)
    if err != nil {
        return oscalModels, err
    }

	return oscalModels, nil
}

func determineOSCALModel(oscalModels oscalTypes_1_1_2.OscalModels) (string, error) {
	// Determine which OSCAL model is present by checking non-nil fields
	switch {
	case oscalModels.AssessmentPlan != nil:
		return "assessment-plan", nil
	case oscalModels.AssessmentResults != nil:
		return "assessment-results", nil
	case oscalModels.Catalog != nil:
		return "catalog", nil
	case oscalModels.ComponentDefinition != nil:
		return "component-definition", nil
	case oscalModels.PlanOfActionAndMilestones != nil:
		return "plan-of-action-and-milestones", nil
	case oscalModels.Profile != nil:
		return "profile", nil
	case oscalModels.SystemSecurityPlan != nil:
		return "system-security-plan", nil
	default:
		return "", fmt.Errorf("unable to determine OSCAL model type")
	}
}

func handleComponentDefinition(componentDefinition *oscalTypes_1_1_2.ComponentDefinition, filePath string, format string) error {
	fmt.Println(componentDefinition, "This is a Component-Definition")

	// componentTitle := componentDefinition.Title

	controlMap := oscal.FilterControlImplementations(componentDefinition)

	extractedData := ExtractControlIDs(controlMap)

    // Call countControlIDs to get a count of mapped control ids in the component definition
	// componentReportData.ControlIDMapped = countControlIDs(oscalModels)

	// componentReportData.ControlIDSource = countControlIDsFromSource(oscalModels)

	// componentReportData.Title = getTitleFromComponentDefinition(oscalModels)
	extractedData.Title = componentDefinition.Metadata.Title



    // Create the final ReportData structure
    report := ReportData{
        ComponentDefinition: extractedData,
    }

    // Write the report to the specified file in the desired format if provided, defaults to yaml and ./oscal-report.yaml
    return WriteReport(report, filePath, format)
}

// func countControlIDs(oscalModels oscalTypes_1_1_2.OscalModels) int {
//     controlIDCounts := make(map[string]int)
//     if oscalModels.ComponentDefinition != nil && oscalModels.ComponentDefinition.Components != nil {
//         for _, component := range *oscalModels.ComponentDefinition.Components {
//             if component.ControlImplementations != nil {
//                 for _, controlImpl := range *component.ControlImplementations {
//                     for _, implementedReq := range controlImpl.ImplementedRequirements {
//                         controlID := implementedReq.ControlId
//                         controlIDCounts[controlID]++
//                     }
//                 }
//             }
//         }
//     }

//     // Return the count of unique Control IDs
// 	return len(controlIDCounts)
// }

// Testing
func ExtractControlIDs(controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) *ComponentDefinitionReportData {
    controlIDCountsSource := make(map[string]int)
    controlFrameworkCounts := make(map[string]int)

    for key, implementations := range controlMap {
        for _, controlImplementation := range implementations {
            // Iterate over ImplementedRequirements to count ControlIds
            for _, implementedReq := range controlImplementation.ImplementedRequirements {
                controlID := implementedReq.ControlId
                controlIDCountsSource[controlID]++  // Counting unique Control IDs by Source

                // Increment the framework count for the given framework (key)
                controlFrameworkCounts[key]++
            }
        }
    }

    // Create an instance of ComponentDefinitionReportData and set the fields
    reportData := &ComponentDefinitionReportData{
        ControlIDMapped:   len(controlIDCountsSource),  // Count of unique Control IDs by Source
        ControlIDFramework: controlFrameworkCounts,     // Map of framework name to count of controls
    }

    return reportData
}

func WriteReport(data ReportData, filePath string, format string) error {
    var err error
    var fileData []byte

	// Set default file and path if not provided
	if filePath == "" {
		if format == "yaml" {
			filePath = "oscal-report.yaml"
		} else {
			filePath = "oscal-report.json"
		}
	}

    // Determine the format
    if format == "yaml" {
        fileData, err = yaml.Marshal(data)
        if err != nil {
            return fmt.Errorf("failed to marshal data to YAML: %v", err)
        }
    } else {
        fileData, err = json.MarshalIndent(data, "", "  ")
        if err != nil {
            return fmt.Errorf("failed to marshal data to JSON: %v", err)
        }
    }

    // Write the serialized data to the file
    err = os.WriteFile(filePath, fileData, 0644)
    if err != nil {
        return fmt.Errorf("failed to write file: %v", err)
    }

    return nil
}

//
func PrintReport(data ReportData, format string) error {
    var err error
    var fileData []byte

    // Determine the format
    if format == "yaml" {
        fileData, err = yaml.Marshal(data)
        if err != nil {
            return fmt.Errorf("failed to marshal data to YAML: %v", err)
        }
    } else {
        fileData, err = json.MarshalIndent(data, "", "  ")
        if err != nil {
            return fmt.Errorf("failed to marshal data to JSON: %v", err)
        }
    }

    // Write the serialized data to the file
    fmt.Println(string(fileData))

    return nil
}
