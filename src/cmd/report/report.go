package report

import (
	"errors"
	"os"
	"encoding/json"
	"strings"
    "fmt"

	"github.com/spf13/cobra"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"gopkg.in/yaml.v3"
)

type flags struct {
	InputFile  string // -f --input-file
	FileFormat string // --file-format
}

type ReportData struct {
	ComponentDefinition *ComponentDefinitionReportData `json:"componentDefinition,omitempty" yaml:"componentDefinition,omitempty"`
}

type ComponentDefinitionReportData struct {
	Title                string            `json:"title" yaml:"title"`
	ControlIDBySource    map[string]int    `json:"control ID mapped" yaml:"control ID mapped"`
	ControlIDByFramework map[string]int    `json:"controlIDFramework" yaml:"controlIDFramework"`
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

		spinner := message.NewProgressSpinner("Fetching or reading file %s", opts.InputFile)
		getOSCALModelsFile, err := fetchOrReadFile(opts.InputFile)
		if err != nil {
			spinner.Fatalf(err, "Failed to get OSCAL file")
		}
		spinner.Success()

		spinner = message.NewProgressSpinner("Reading OSCAL model from file")
		oscalModels, err := readOSCALModel(getOSCALModelsFile)
		if err != nil {
			spinner.Fatalf(err, "Failed to get OSCAL Model data")
		}
		spinner.Success()

		spinner = message.NewProgressSpinner("Determining OSCAL model type")
		modelType, err := determineOSCALModel(oscalModels)
		if err != nil {
			spinner.Fatalf(err, "Unable to determine OSCAL model type")
		}
		spinner.Success()

		spinner = message.NewProgressSpinner("Generating report")
		var reportModelErr error
		switch modelType {
		case "catalog", "profile", "assessment-plan", "assessment-results", "system-security-plan", "plan-of-action-and-milestones":
			message.Warnf("Reporting does not create reports for %s at this time", modelType)
		case "component-definition":
			reportModelErr = handleComponentDefinition(oscalModels.ComponentDefinition, opts.FileFormat)
		default:
			spinner.Fatalf(fmt.Errorf("unknown OSCAL model type: %s", modelType), "Failed to process OSCAL file")
		}
		spinner.Success()

		if reportModelErr != nil {
			message.Fatal(reportModelErr, "Failed to create report")
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
	reportFlags.StringVar(&opts.FileFormat, "file-format", "table", "File format of output file")
}

func fetchOrReadFile(source string) ([]byte, error) {
	if isURL(source) {
		spinner := message.NewProgressSpinner("Fetching data from URL: %s", source)
		defer spinner.Stop()
		data, err := network.Fetch(source)
		if err != nil {
			spinner.Fatalf(err, "Failed to fetch data from URL")
		}
		spinner.Success()
		return data, nil
	}
	spinner := message.NewProgressSpinner("Reading file: %s", source)
	defer spinner.Stop()
	data, err := os.ReadFile(source)
	if err != nil {
		spinner.Fatalf(err, "Failed to read file")
	}
	spinner.Success()
	return data, nil
}

// Reads OSCAL file in either YAML or JSON
func readOSCALModel(data []byte) (oscalTypes_1_1_2.OscalModels, error) {
	var oscalModels oscalTypes_1_1_2.OscalModels
	err := yaml.Unmarshal(data, &oscalModels)
	if err == nil {
		return oscalModels, nil
	}
	err = json.Unmarshal(data, &oscalModels)
	if err == nil {
		return oscalModels, nil
	}
	return oscalModels, errors.New("data is neither valid YAML nor JSON")
}

// Checks the OSCAL file to determine what model the file is
func determineOSCALModel(oscalModels oscalTypes_1_1_2.OscalModels) (string, error) {
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
func handleComponentDefinition(componentDefinition *oscalTypes_1_1_2.ComponentDefinition, format string) error {
	controlMap := oscal.FilterControlImplementations(componentDefinition)
	extractedData := ExtractControlIDs(controlMap)
	extractedData.Title = componentDefinition.Metadata.Title

	report := ReportData{
		ComponentDefinition: extractedData,
	}

	message.Info("Generating report...")
	return PrintReport(report, format)
}

// Gets the unique Control IDs from each source and framework in the OSCAL Component Definition
func ExtractControlIDs(controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) *ComponentDefinitionReportData {
	sourceMap, frameworkMap := SplitControlMap(controlMap)

	sourceControlIDs := make(map[string]int)
	for source, controlMap := range sourceMap {
		total := 0
		for _, count := range controlMap {
			total += count
		}
		sourceControlIDs[source] = total
	}

	aggregatedFrameworkCounts := make(map[string]int)
	for framework, controlCounts := range frameworkMap {
		total := 0
		for _, count := range controlCounts {
			total += count
		}
		aggregatedFrameworkCounts[framework] = total
	}

	return &ComponentDefinitionReportData{
		ControlIDBySource:    sourceControlIDs,
		ControlIDByFramework: aggregatedFrameworkCounts,
	}
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

// Helper to determine if the controlMap source is a URL
func isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func PrintReport(data ReportData, format string) error {
    if format == "table" {
        // Use the message package for printing table data
        message.Infof("Title: %s", data.ComponentDefinition.Title)

        // Print the Control ID By Source as a table
        message.Info("\nControl Source            | Number of Controls")
        message.Info(strings.Repeat("-", 60))

        for source, count := range data.ComponentDefinition.ControlIDBySource {
            message.Infof("%-40s | %-15d", source, count)
        }

        // Print the Control ID By Framework as a table
        message.Info("\nFramework                | Number of Controls")
        message.Info(strings.Repeat("-", 40))

        for framework, count := range data.ComponentDefinition.ControlIDByFramework {
            message.Infof("%-20s | %-15d", framework, count)
        }

    } else {
        var err error
        var fileData []byte

        if format == "yaml" {
            message.Info("Generating report in YAML format...")
            fileData, err = yaml.Marshal(data)
            if err != nil {
                message.Fatal(err, "Failed to marshal data to YAML")
            }
        } else {
            message.Info("Generating report in JSON format...")
            fileData, err = json.MarshalIndent(data, "", "  ")
            if err != nil {
                message.Fatal(err, "Failed to marshal data to JSON")
            }
        }

        message.Info(string(fileData))
    }

    return nil
}
