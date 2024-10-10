package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"gopkg.in/yaml.v3"
)

type ReportData struct {
	ComponentDefinition *ComponentDefinitionReportData `json:"componentDefinition,omitempty" yaml:"componentDefinition,omitempty"`
}

type ComponentDefinitionReportData struct {
	Title                string         `json:"title" yaml:"title"`
	ControlIDBySource    map[string]int `json:"control ID mapped" yaml:"control ID mapped"`
	ControlIDByFramework map[string]int `json:"controlIDFramework" yaml:"controlIDFramework"`
}

// Runs the logic of report generation
func GenerateReport(inputFile string, fileFormat string) error {
	spinner := message.NewProgressSpinner("Fetching or reading file %s", inputFile)
	getOSCALModelsFile, err := fetchOrReadFile(inputFile)
	if err != nil {
		spinner.Fatalf(fmt.Errorf("failed to get OSCAL file: %v", err), "failed to get OSCAL file")
	}
	spinner.Success()

	spinner = message.NewProgressSpinner("Reading OSCAL model from file")
	oscalModel, err := oscal.NewOscalModel(getOSCALModelsFile)
	if err != nil {
		spinner.Fatalf(fmt.Errorf("failed to read OSCAL Model data: %v", err), "failed to read OSCAL Model")
	}
	spinner.Success()

	// Initialize CompositionContext
	compCtx, err := composition.New()
	if err != nil {
		return fmt.Errorf("failed to create composition context: %v", err)
	}

	err = handleOSCALModel(oscalModel, fileFormat, compCtx)
	if err != nil {
		return err
	}

	return nil
}

// Processes an OSCAL Model based on the model type
func handleOSCALModel(oscalModel *oscalTypes_1_1_2.OscalModels, format string, compCtx *composition.CompositionContext) error {
	// Start a new spinner for the report generation process
	spinner := message.NewProgressSpinner("Determining OSCAL model type")
	modelType, err := oscal.GetOscalModel(oscalModel)
	if err != nil {
		spinner.Fatalf(fmt.Errorf("unable to determine OSCAL model type: %v", err), "unable to determine OSCAL model type")
		return err
	}

	switch modelType {
	case "catalog", "profile", "assessment-plan", "assessment-results", "system-security-plan", "poam":
		// If the model type is not supported, stop the spinner with a warning
		spinner.Warnf("reporting does not create reports for %s at this time", modelType)
		return fmt.Errorf("reporting does not create reports for %s at this time", modelType)

	case "component":
		spinner.Updatef("Composing Component Definition")
		err := compCtx.ComposeComponentDefinitions(context.Background(), oscalModel.ComponentDefinition, "")
		if err != nil {
			spinner.Fatalf(fmt.Errorf("failed to compose component definitions: %v", err), "failed to compose component definitions")
			return err
		}

		spinner.Updatef("Processing Component Definition")
		// Process the component-definition model
		err = handleComponentDefinition(oscalModel.ComponentDefinition, format, compCtx)
		if err != nil {
			// If an error occurs, stop the spinner and display the error
			spinner.Fatalf(err, "failed to process component-definition model")
			return err
		}

	default:
		// For unknown model types, stop the spinner with a failure
		spinner.Fatalf(fmt.Errorf("unknown OSCAL model type: %s", modelType), "failed to process OSCAL file")
		return fmt.Errorf("unknown OSCAL model type: %s", modelType)
	}

	spinner.Success()
	message.Info(fmt.Sprintf("Successfully processed OSCAL model: %s", modelType))
	return nil
}

// Handler for Component Definition OSCAL files to create the report
func handleComponentDefinition(componentDefinition *oscalTypes_1_1_2.ComponentDefinition, format string, compCtx *composition.CompositionContext) error {
	spinner := message.NewProgressSpinner("composing component definitions")

	err := compCtx.ComposeComponentDefinitions(context.Background(), componentDefinition, "")
	if err != nil {
		spinner.Fatalf(fmt.Errorf("failed to compose component definitions: %v", err), "failed to compose component definitions")
		return err
	}

	spinner.Success() // Mark the spinner as successful before moving forward

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
