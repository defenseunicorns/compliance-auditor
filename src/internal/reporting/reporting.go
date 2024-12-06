package reporting

import (
	"context"
	"encoding/json"
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
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

	getOSCALModelsFile, err := network.Fetch(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get OSCAL file: %v", err)
	}

	spinner.Success()

	spinner = message.NewProgressSpinner("Reading OSCAL model from file")
	oscalModel, err := oscal.NewOscalModel(getOSCALModelsFile)
	if err != nil {
		return fmt.Errorf("failed to read OSCAL Model data: %v", err)
	}
	spinner.Success()

	// Set up the composer
	composer, err := composition.New(
		composition.WithRenderSettings("all", true),
		composition.WithTemplateRenderer("all", common.TemplateConstants, common.TemplateVariables, []string{}),
	)
	if err != nil {
		return fmt.Errorf("error creating new composer: %v", err)
	}

	err = handleOSCALModel(oscalModel, fileFormat, composer)
	if err != nil {
		return err
	}

	return nil
}

// Processes an OSCAL Model based on the model type
func handleOSCALModel(oscalModel *oscalTypes.OscalModels, format string, composer *composition.Composer) error {
	// Start a new spinner for the report generation process
	spinner := message.NewProgressSpinner("Determining OSCAL model type")
	modelType, err := oscal.GetOscalModel(oscalModel)
	if err != nil {
		return fmt.Errorf("unable to determine OSCAL model type: %v", err)
	}

	switch modelType {
	case "catalog", "profile", "assessment-plan", "assessment-results", "system-security-plan", "poam":
		// If the model type is not supported, stop the spinner with a warning
		return fmt.Errorf("reporting does not create reports for %s at this time", modelType)

	case "component":
		spinner.Updatef("Composing Component Definition")
		err := composer.ComposeComponentDefinitions(context.Background(), oscalModel.ComponentDefinition, "")
		if err != nil {
			return fmt.Errorf("failed to compose component definitions: %v", err)
		}

		spinner.Updatef("Processing Component Definition")
		// Process the component-definition model
		err = handleComponentDefinition(oscalModel.ComponentDefinition, format)
		if err != nil {
			// If an error occurs, stop the spinner and display the error
			return err
		}

	default:
		// For unknown model types, stop the spinner with a failure
		return fmt.Errorf("unknown OSCAL model type: %s", modelType)
	}

	spinner.Success()
	message.Info(fmt.Sprintf("Successfully processed OSCAL model: %s", modelType))
	return nil
}

// Handler for Component Definition OSCAL files to create the report
func handleComponentDefinition(componentDefinition *oscalTypes.ComponentDefinition, format string) error {

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
func ExtractControlIDs(controlMap map[string][]oscalTypes.ControlImplementationSet) *ComponentDefinitionReportData {
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
		// Use the Table function to print a formatted table
		message.Infof("Title: %s", data.ComponentDefinition.Title)

		// Prepare headers and data for Control ID By Source table
		sourceHeaders := []string{"Control Source", "Number of Controls"}
		sourceData := make([][]string, 0, len(data.ComponentDefinition.ControlIDBySource))
		for source, count := range data.ComponentDefinition.ControlIDBySource {
			sourceData = append(sourceData, []string{source, fmt.Sprintf("%d", count)})
		}
		// Print Control ID By Source using the Table function
		if err := message.Table(sourceHeaders, sourceData, []int{70, 30}); err != nil {
			// Handle the error, e.g., log or return
			return err
		}

		// Prepare headers and data for Control ID By Framework table
		frameworkHeaders := []string{"Framework", "Number of Controls"}
		frameworkData := make([][]string, 0, len(data.ComponentDefinition.ControlIDByFramework))
		for framework, count := range data.ComponentDefinition.ControlIDByFramework {
			frameworkData = append(frameworkData, []string{framework, fmt.Sprintf("%d", count)})
		}
		// Print Control ID By Framework using the Table function
		if err := message.Table(frameworkHeaders, frameworkData, []int{70, 30}); err != nil {
			// Handle the error, e.g., log or return
			return err
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
