package validation

import (
	"fmt"
	"os"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// ResultConsumer is the interface that must be implemented by any consumer of the validation
// store and results. It is responsible for evaluating the results and generating the output
// speific to the consumer.
type ResultConsumer interface {
	// Evaluate Results are the custom implementation for the consumer, which should take the
	// requirements, as specified by the producer, plus the data in the validation store
	// and evaluate them + generate the output
	EvaluateResults(store *ValidationStore) error

	// Generate Output is the custom implementation for the consumer that should create
	// a custom output
	GenerateOutput() error
}

// AssessmentResultsConsumer is an implementation of the ResultConsumer interface
// This consumer is responsible for generating an OSCAL Assessment Results model
type AssessmentResultsConsumer struct {
	assessmentResults *oscal.AssessmentResults
	path              string
}

func NewAssessmentResultsConsumer(path string) *AssessmentResultsConsumer {
	// Get asssessment results from file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	ar := oscal.NewAssessmentResults2()

	// Update the assessment results model if data is not nil
	if len(data) != 0 {
		err = ar.NewModel(data)
		if err != nil {
			return nil
		}
	}

	return &AssessmentResultsConsumer{
		assessmentResults: ar,
		path:              path,
	}
}

func (c *AssessmentResultsConsumer) EvaluateResults(store *ValidationStore) error {
	// Update the oscal.AssessmentResults with the results from the store
	// each requirement should be a finding
	// each validation in the requirement should be an observation

	// Create oscal results -> generate assessment results model (GenerateAssessmentResults)

	// If the existing assessment results are nil (c.assessmentResults == nil), set them

	// If they are populated, merge the results from the store into the existing assessment results

	return nil
}

func (c *AssessmentResultsConsumer) GenerateOutput() error {
	// Maybe should this consumer just create the results and then run a generate function
	// to create the assessment results model? I feel like if this is from an assessment plan
	// vs. a component definition, the assesment results model will be different... could/should
	// this be handled prior to the consumer being created?

	return oscal.WriteOscalModelNew(c.path, c.assessmentResults)
}

// SimpleConsumer is an implementation of the ResultConsumer interface
// The consumer determines "Pass" is true iff all requirements are satisfied
// Useful for quick determination of pass/fail status of the requirements
type SimpleConsumer struct {
	pass bool
	msg  string
}

func NewSimpleConsumer() *SimpleConsumer {
	return &SimpleConsumer{
		pass: false,
	}
}

func (c *SimpleConsumer) EvaluateResults(store *ValidationStore) error {
	var output strings.Builder
	requirements := store.GetRequirements()
	passCount := 0

	// Evaluate each requirement for pass/fail
	for _, requirement := range requirements {
		if requirement == nil {
			continue
		}

		pass, msg := requirement.EvaluateSuccess()
		if !pass {
			passCount++
		}
		output.WriteString(msg)
	}

	if passCount > 0 && passCount == len(requirements) {
		c.pass = true
		return nil
	}

	c.msg = output.String()

	return nil
}

func (c *SimpleConsumer) GenerateOutput() error {
	if !c.pass {
		return fmt.Errorf("requirements failed: %s", c.msg)
	}
	message.Infof("Requirements passed")
	return nil
}
