package report

import (
	"errors"
	"os"
	"fmt"
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
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
	ControlIDBySource map[string]int `json:"control ID mapped" yaml:"control ID mapped"`
	ControlIDByFramework map[string]int `json:"controlIDFramework" yaml:"controlIDFramework"`
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

		getOSCALModelsFile, err := fetchOrReadFile(opts.InputFile)
		if err != nil {
			message.Fatal(err, "Failed to get OSCAL file")
		}

		oscalModels, err := readOSCALModel(getOSCALModelsFile)
		if err != nil {
			message.Fatal(err, "Failed to get OSCAL Model data")
		}

		modelType, err := determineOSCALModel(oscalModels)
		if err != nil {
			message.Fatal(err, "Unable to determine OSCAL model type")
		}

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
	reportFlags.StringVar(&opts.FileFormat, "file-format", "json", "File format of output file")

}

func fetchOrReadFile(source string) ([]byte, error) {
	 // Check if the source is a URL
    if isURL(source) {
        return network.Fetch(source)  // Fetch data from the URL
    }

    // If it's not a URL, assume it's a file path and read the file
    return os.ReadFile(source)
}

// Reads OSCAL file in either YAML or JSON
func readOSCALModel(data []byte) (oscalTypes_1_1_2.OscalModels, error) {
    var oscalModels oscalTypes_1_1_2.OscalModels

    // Try to unmarshal as YAML
    err := yaml.Unmarshal(data, &oscalModels)
    if err == nil {
        return oscalModels, nil
    }

    // If YAML unmarshaling fails, try JSON
    err = json.Unmarshal(data, &oscalModels)
    if err == nil {
        return oscalModels, nil
    }

    return oscalModels, errors.New("data is neither valid YAML nor JSON")
}

// Checks the OSCAL file to determine what model the file is
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

// Handler for Component Definition OSCAL files to create the report
func handleComponentDefinition(componentDefinition *oscalTypes_1_1_2.ComponentDefinition, filePath string, format string) error {
	fmt.Println(componentDefinition, "This is a Component-Definition")

	// componentTitle := componentDefinition.Title

	controlMap := oscal.FilterControlImplementations(componentDefinition)

	extractedData := ExtractControlIDs(controlMap)

	extractedData.Title = componentDefinition.Metadata.Title



    // Create the final ReportData structure
    report := ReportData{
        ComponentDefinition: extractedData,
    }

    // Write the report to the specified file in the desired format if provided, defaults to yaml and ./oscal-report.yaml
    return WriteReport(report, filePath, format)
}

// Gets the unique Control IDs from each source and framework in the OSCAL Component Definition
func ExtractControlIDs(controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) *ComponentDefinitionReportData {
    // Split the control map into source and framework maps
    sourceMap, frameworkMap := SplitControlMap(controlMap)


    // Calculate the total number of unique control IDs by source
    sourceControlIDs := make(map[string]int)
    for source, controlMap := range sourceMap {
        total := 0
		for _, count := range controlMap {
			total += count
		}
		sourceControlIDs[source] = total
    }

    // Aggregate the control counts from frameworkMap
    aggregatedFrameworkCounts := make(map[string]int)
    for framework, controlCounts := range frameworkMap {
        total := 0
        for _, count := range controlCounts {
            total += count
        }
        aggregatedFrameworkCounts[framework] = total
    }

    // Create the report data
    reportData := &ComponentDefinitionReportData{
        ControlIDBySource: sourceControlIDs,              // Total count of unique Control IDs by Source
        ControlIDByFramework: aggregatedFrameworkCounts,   // Total count of unique Control IDs by Framework
    }

    return reportData

}

// Split the default controlMap into framework and source maps for further processing
func SplitControlMap(controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) (sourceMap map[string]map[string]int, frameworkMap map[string]map[string]int) {
    sourceMap = make(map[string]map[string]int)
    frameworkMap = make(map[string]map[string]int)

    for key, implementations := range controlMap {
        if isURL(key) {
            if _, exists := sourceMap[key]; !exists {
                sourceMap[key] = make(map[string]int)
            }
            for _, controlImplementation := range implementations {
                for _, implementedReq := range controlImplementation.ImplementedRequirements {
                    controlID := implementedReq.ControlId
                    sourceMap[key][controlID]++
                }
            }
        } else {
            if _, exists := frameworkMap[key]; !exists {
                frameworkMap[key] = make(map[string]int)
            }
            for _, controlImplementation := range implementations {
                for _, implementedReq := range controlImplementation.ImplementedRequirements {
                    controlID := implementedReq.ControlId
                    frameworkMap[key][controlID]++
                }
            }
        }
    }

    return sourceMap, frameworkMap
}

// Helper to determine in controlMap source from framework
func isURL(str string) bool {
    return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func GetSourceControlIDs(sourceMap map[string]map[string]int, frameworkMap map[string]map[string]int) (map[string]int, error) {
    result := make(map[string]int)

    // Process sourceMap
    for sourceURL := range sourceMap {
        sourceControls, err := fetchOrReadFile(sourceURL)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch controls from source %s: %v", sourceURL, err)
        }

        oscalModels, err := readOSCALModel(sourceControls)
        if err != nil {
            return nil, fmt.Errorf("failed to parse OSCAL data from source %s: %v", sourceURL, err)
        }

        modelType, err := determineOSCALModel(oscalModels)
        if err != nil {
            return nil, fmt.Errorf("failed to determine OSCAL model type for source %s: %v", sourceURL, err)
        }

        controlCount, err := extractControlIDsFromModel(oscalModels, modelType)
        if err != nil {
            return nil, fmt.Errorf("failed to extract control IDs from source %s: %v", sourceURL, err)
        }

        result[sourceURL] = controlCount
		fmt.Printf("Processed %s: found %d controls\n", sourceURL, controlCount)

    }

    // Process frameworkMap similarly
    for frameworkURL := range frameworkMap {
        frameworkControls, err := fetchOrReadFile(frameworkURL)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch controls from framework %s: %v", frameworkURL, err)
        }

        oscalModels, err := readOSCALModel(frameworkControls)
        if err != nil {
            return nil, fmt.Errorf("failed to parse OSCAL data from framework %s: %v", frameworkURL, err)
        }

        modelType, err := determineOSCALModel(oscalModels)
        if err != nil {
            return nil, fmt.Errorf("failed to determine OSCAL model type for framework %s: %v", frameworkURL, err)
        }

        controlCount, err := extractControlIDsFromModel(oscalModels, modelType)
        if err != nil {
            return nil, fmt.Errorf("failed to extract control IDs from framework %s: %v", frameworkURL, err)
        }

        result[frameworkURL] = controlCount
    }

    return result, nil
}

// Get the control ids from either the catalogs and/or profiles
func extractControlIDsFromModel(oscalModels oscalTypes_1_1_2.OscalModels, modelType string) (int, error) {
    switch modelType {
    case "catalog":
        return countCatalogControlIDs(*oscalModels.Catalog), nil
    case "profile":
        return countProfileControlIDs(*oscalModels.Profile), nil
    default:
        return 0, fmt.Errorf("unsupported OSCAL model type: %s", modelType)
    }
}

// Logic for counting the controls in an OSCAL catalog
func countCatalogControlIDs(catalog oscalTypes_1_1_2.Catalog) int {
    controlCount := 0

    // Iterate over each group in the catalog
    if catalog.Groups != nil {
        for _, group := range *catalog.Groups {
            // Iterate over each control in the group
            if group.Controls != nil {
                for _, control := range *group.Controls {
                    // Count the control itself
                    controlCount++

                    // Count each parameter's ID in the control
                    if control.Params != nil {
                        for _, param := range *control.Params {
                            fmt.Println("Found Param ID:", param.ID)
                            controlCount++ // Counting each param as well
                        }
                    }

                    // Recursively count any nested controls
                    if control.Controls != nil {
                        controlCount += countNestedControls(*control.Controls)
                    }
                }
            }
        }
    }

    return controlCount
}

// Logic for counting nested controls.Controls in an OSCAL Catalog
// TODO verify this, not 100% its accurate looking at the reference docs
func countNestedControls(controls []oscalTypes_1_1_2.Control) int {
    nestedCount := 0
    for _, control := range controls {
        nestedCount++
        // Count each parameter's ID in the control
        if control.Params != nil {
            for _, param := range *control.Params {
                fmt.Println("Found Param ID:", param.ID)
                nestedCount++ // Counting each param as well
            }
        }
        if control.Controls != nil {
            nestedCount += countNestedControls(*control.Controls)
        }
    }
    return nestedCount
}


func countProfileControlIDs(profile oscalTypes_1_1_2.Profile) int {
    controlCount := 0
    for _, imp := range profile.Imports {
        if imp.IncludeControls != nil {
            controlCount += len(*imp.IncludeControls)
        }
    }
    return controlCount
}

// Creates the report file
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

// TODO: Turn this into a print section using the message pkg
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
