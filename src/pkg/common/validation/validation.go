package validation

import (
	"context"
	"fmt"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type ValidationContext struct {
	bctx                         context.Context
	componentDefinition          *oscalTypes_1_1_2.ComponentDefinition
	componentPath                string
	templateRenderer             *template.TemplateRenderer
	renderTemplate               bool
	requestExecutionConfirmation bool
	runExecutableValidations     bool
	resourcesDir                 string
}

func New(ctx context.Context, opts ...Option) (*ValidationContext, error) {
	var validationCtx ValidationContext

	for _, opt := range opts {
		if err := opt(&validationCtx); err != nil {
			return nil, err
		}
	}

	validationCtx.bctx = ctx

	// If from path, run a compose + template/render
	oscalModel, err := composition.ComposeFromPath(validationCtx.componentPath)
	if err != nil {
		return nil, fmt.Errorf("error composing from path: %v", err)
	}

	if oscalModel.ComponentDefinition == nil {
		return nil, fmt.Errorf("component definition is nil")
	}

	validationCtx.componentDefinition = oscalModel.ComponentDefinition

	return &validationCtx, nil
}

func (ctx *ValidationContext) RenderComponentDefinition(fileExt string, renderType template.RenderType) ([]byte, error) {
	if ctx.componentDefinition == nil {
		return nil, fmt.Errorf("cannot render a component definition that is nil")
	}

	fullModel := &oscalTypes_1_1_2.OscalModels{
		ComponentDefinition: ctx.componentDefinition,
	}

	modelData, err := oscal.ConvertOSCALToBytes(fullModel, fileExt)
	if err != nil {
		return nil, fmt.Errorf("error converting component definition to bytes: %v", err)
	}

	if ctx.renderTemplate {
		ctx.templateRenderer.UpdateTemplateString(string(modelData))

		// Execute the template render
		modelData, err = ctx.templateRenderer.Render(renderType)
		if err != nil {
			return nil, fmt.Errorf("error rendering template: %v", err)
		}

		// Convert modelData back to the component definition
		newModel, err := oscal.NewOscalModel(modelData)
		if err != nil {
			return nil, fmt.Errorf("error creating new oscal model from rendered data: %v", err)
		}

		if newModel.ComponentDefinition == nil {
			return nil, fmt.Errorf("error creating new oscal model from rendered data: component definition is nil")
		}

		ctx.componentDefinition = newModel.ComponentDefinition
	}
	return modelData, nil
}

func (ctx *ValidationContext) ExecuteOSCALValidation(target string) (*oscalTypes_1_1_2.AssessmentResults, error) {
	// TODO: Should we execute the validation even if there are no comp-def/components, e.g., create an empty assessment-results object?

	if ctx.componentDefinition == nil {
		return nil, fmt.Errorf("cannot validate a component definition that is nil")
	}

	if *ctx.componentDefinition.Components == nil {
		return nil, fmt.Errorf("no components found in component definition")
	}

	// Create a validation store from the back-matter if it exists
	validationStore := validationstore.NewValidationStoreFromBackMatter(*ctx.componentDefinition.BackMatter)

	// Create a map of control implementations from the component definition
	// This combines all same source/framework control implementations into an []Control-Implementation
	controlImplementations := oscal.FilterControlImplementations(ctx.componentDefinition)

	if len(controlImplementations) == 0 {
		return nil, fmt.Errorf("no control implementations found in component definition")
	}

	// Get results of validation execution
	results := make([]oscalTypes_1_1_2.Result, 0)
	if target != "" {
		if controlImplementation, ok := controlImplementations[target]; ok {
			findings, observations, err := ctx.ValidateOnControlImplementations(&controlImplementation, validationStore, target)
			if err != nil {
				return nil, err
			}
			result, err := oscal.CreateResult(findings, observations)
			if err != nil {
				return nil, err
			}
			// add/update the source to the result props - make source = framework or omit?
			oscal.UpdateProps("target", oscal.LULA_NAMESPACE, target, result.Props)
			results = append(results, result)
		} else {
			return nil, fmt.Errorf("target %s not found", target)
		}
	} else {
		// default behavior - create a result for each unique source/framework
		// loop over the controlImplementations map & validate
		// we lose context of source if not contained within the loop
		for source, controlImplementation := range controlImplementations {
			findings, observations, err := ctx.ValidateOnControlImplementations(&controlImplementation, validationStore, source)
			if err != nil {
				return nil, err
			}
			result, err := oscal.CreateResult(findings, observations)
			if err != nil {
				return nil, err
			}
			// add/update the source to the result props
			oscal.UpdateProps("target", oscal.LULA_NAMESPACE, source, result.Props)
			results = append(results, result)
		}
	}

	return oscal.GenerateAssessmentResults(results)
}

func (ctx *ValidationContext) ValidateOnControlImplementations(controlImplementations *[]oscalTypes_1_1_2.ControlImplementationSet, validationStore *validationstore.ValidationStore, target string) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {
	// Create requirement store for all implemented requirements
	requirementStore := requirementstore.NewRequirementStore(controlImplementations)
	message.Title("\n🔍 Collecting Requirements and Validations for Target: ", target)
	requirementStore.ResolveLulaValidations(validationStore)
	reqtStats := requirementStore.GetStats(validationStore)
	message.Infof("Found %d Implemented Requirements", reqtStats.TotalRequirements)
	message.Infof("Found %d runnable Lula Validations", reqtStats.TotalValidations)

	// Check if validations perform execution actions
	if reqtStats.ExecutableValidations {
		if !ctx.runExecutableValidations && ctx.requestExecutionConfirmation {
			confirmExecution := message.PromptForConfirmation(nil)
			if !confirmExecution {
				message.Infof("Validations requiring execution will NOT be run")
			} else {
				ctx.runExecutableValidations = true
			}
		}
	}

	// Set values for saving resources
	saveResources := false
	if ctx.resourcesDir != "" {
		saveResources = true
	}

	// Run Lula validations and generate observations & findings
	message.Title("\n📐 Running Validations", "")
	observations := validationStore.RunValidations(ctx.runExecutableValidations, saveResources, ctx.resourcesDir)
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
