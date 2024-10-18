package validation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Validator struct {
	composer                     *composition.Composer
	requestExecutionConfirmation bool
	runExecutableValidations     bool
	resourcesDir                 string
}

func New(opts ...Option) (*Validator, error) {
	var validator Validator

	for _, opt := range opts {
		if err := opt(&validator); err != nil {
			return nil, err
		}
	}

	return &validator, nil
}

func (v *Validator) ValidateOnPath(ctx context.Context, path, target string) (assessmentResult *oscalTypes_1_1_2.AssessmentResults, err error) {
	var oscalModel *oscalTypes_1_1_2.OscalCompleteSchema
	if v.composer == nil {
		path = filepath.Clean(path)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error getting path: %v", err)
		}
		oscalModel, err = oscal.NewOscalModel(data)
		if err != nil {
			return nil, fmt.Errorf("error creating oscal model from path: %v", err)
		}
	} else {
		oscalModel, err = v.composer.ComposeFromPath(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("error composing model: %v", err)
		}
	}

	if oscalModel.ComponentDefinition == nil {
		return assessmentResult, fmt.Errorf("component definition is nil")
	}

	results, err := v.ValidateOnCompDef(ctx, oscalModel.ComponentDefinition, target)
	if err != nil {
		return assessmentResult, err
	}

	assessmentResult, err = oscal.GenerateAssessmentResults(results)
	if err != nil {
		return assessmentResult, err
	}

	return assessmentResult, nil
}

func (v *Validator) ValidateOnCompDef(ctx context.Context, compDef *oscalTypes_1_1_2.ComponentDefinition, target string) (results []oscalTypes_1_1_2.Result, err error) {
	// TODO: Should we execute the validation even if there are no comp-def/components, e.g., create an empty assessment-results object?

	if compDef == nil {
		return nil, fmt.Errorf("cannot validate a component definition that is nil")
	}

	if *compDef.Components == nil {
		return nil, fmt.Errorf("no components found in component definition")
	}

	// Create a validation store from the back-matter if it exists
	validationStore := validationstore.NewValidationStoreFromBackMatter(*compDef.BackMatter)

	// Create a map of control implementations from the component definition
	// This combines all same source/framework control implementations into an []Control-Implementation
	controlImplementations := oscal.FilterControlImplementations(compDef)

	if len(controlImplementations) == 0 {
		return nil, fmt.Errorf("no control implementations found in component definition")
	}

	// Get results of validation execution
	results = make([]oscalTypes_1_1_2.Result, 0)
	if target != "" {
		if controlImplementation, ok := controlImplementations[target]; ok {
			findings, observations, err := v.ValidateOnControlImplementations(ctx, &controlImplementation, validationStore, target)
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
			findings, observations, err := v.ValidateOnControlImplementations(ctx, &controlImplementation, validationStore, source)
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

	return results, nil
}

func (v *Validator) ValidateOnControlImplementations(ctx context.Context, controlImplementations *[]oscalTypes_1_1_2.ControlImplementationSet, validationStore *validationstore.ValidationStore, target string) (map[string]oscalTypes_1_1_2.Finding, []oscalTypes_1_1_2.Observation, error) {
	// Create requirement store for all implemented requirements
	requirementStore := requirementstore.NewRequirementStore(controlImplementations)
	message.Title("\n🔍 Collecting Requirements and Validations for Target: ", target)
	requirementStore.ResolveLulaValidations(validationStore)
	reqtStats := requirementStore.GetStats(validationStore)
	message.Infof("Found %d Implemented Requirements", reqtStats.TotalRequirements)
	message.Infof("Found %d runnable Lula Validations", reqtStats.TotalValidations)

	// Check if validations perform execution actions
	if reqtStats.ExecutableValidations {
		if !v.runExecutableValidations && v.requestExecutionConfirmation {
			confirmExecution := message.PromptForConfirmation(nil)
			if !confirmExecution {
				message.Infof("Validations requiring execution will NOT be run")
			} else {
				v.runExecutableValidations = true
			}
		}
	}

	// Set values for saving resources
	saveResources := false
	if v.resourcesDir != "" {
		saveResources = true
	}

	// Run Lula validations and generate observations & findings
	message.Title("\n📐 Running Validations", "")
	observations := validationStore.RunValidations(ctx, v.runExecutableValidations, saveResources, v.resourcesDir)
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

	var err error
	if len(rows) != 0 {
		err = message.Table(header, rows, columnSize)
	}

	return findings, observations, err
}
