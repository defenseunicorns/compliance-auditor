package types

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/defenseunicorns/lula/src/internal/transform"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// LulaValidationTestData is a struct that contains the details of the test performed on the LulaValidation
// as well as the result of the test
type LulaValidationTestData struct {
	Test   *LulaValidationTest
	Result *LulaValidationTestResult
}

// LulaValidationTest is a struct that contains the details of the test performed on the LulaValidation
// The 'test' is an evaluation of the Provider, comparing the actual result against expected,
// when provided with a changed input of Domain Resources
type LulaValidationTest struct {
	Name           string                     `json:"name" yaml:"name"`
	Changes        []LulaValidationTestChange `json:"changes,omitempty" yaml:"changes,omitempty"`
	ExpectedResult string                     `json:"expected-result" yaml:"expected-result"`
}

// ValidateData validates the data in the LulaValidationTest struct
func (l *LulaValidationTest) ValidateData() error {
	if l.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if l.ExpectedResult != "satisfied" && l.ExpectedResult != "not-satisfied" {
		return fmt.Errorf("expected-result must be satisfied or not-satisfied")
	}

	for _, change := range l.Changes {
		if err := change.validateData(); err != nil {
			return err
		}
	}

	return nil
}

// LulaValidationTestChange is a struct that contains the details of the changes that are to be made to the resources
// for a LulaValidationTest
type LulaValidationTestChange struct {
	Path     string                 `json:"path" yaml:"path"`
	Type     transform.ChangeType   `json:"type" yaml:"type"`
	Value    string                 `json:"value,omitempty" yaml:"value,omitempty"`
	ValueMap map[string]interface{} `json:"value-map,omitempty" yaml:"value-map,omitempty"`
}

// ValidateData validates the data in the LulaValidationTestChange struct
func (c *LulaValidationTestChange) validateData() error {
	if c.Path == "" {
		return fmt.Errorf("path is empty")
	}
	switch c.Type {
	case transform.ChangeTypeAdd, transform.ChangeTypeUpdate, transform.ChangeTypeDelete:
	default:
		return fmt.Errorf("invalid type")
	}

	return nil
}

// ExecuteTest executes a single LulaValidationTest
func (d *LulaValidationTestData) ExecuteTest(ctx context.Context, validation *LulaValidation, resources map[string]interface{}, saveResources bool) (*LulaValidationTestResult, error) {
	if d.Test == nil {
		return nil, fmt.Errorf("test is nil")
	}

	d.Result = &LulaValidationTestResult{
		TestName: d.Test.Name,
	}

	tt, err := transform.CreateTransformTarget(resources)
	if err != nil {
		d.Result.Pass = false
		d.Result.Remarks = map[string]string{
			"error creating transform target": err.Error(),
		}
		return d.Result, nil
	}

	for _, c := range d.Test.Changes {
		resources, err = tt.ExecuteTransform(c.Path, c.Type, c.Value, c.ValueMap)
		if err != nil {
			d.Result.Pass = false
			d.Result.Remarks = map[string]string{
				"error executing transform": err.Error(),
			}
			return d.Result, nil
		}
	}

	// save resources to validation directory
	if saveResources {
		workDir, ok := ctx.Value(LulaValidationWorkDir).(string)
		if !ok {
			workDir = "."
		}

		resourcesPath := filepath.Join(workDir, fmt.Sprintf("%s.json", d.Test.Name))

		err := WriteResources(resources, resourcesPath)
		if err != nil {
			message.Debugf("Error writing resource data to file: %v", err)
		} else {
			d.Result.TestResourcesPath = resourcesPath
		}
	}

	err = validation.Validate(ctx, WithStaticResources(resources))
	if err != nil {
		d.Result.Pass = false
		d.Result.Remarks = map[string]string{
			"error running validation": err.Error(),
		}
		return d.Result, nil
	}

	// Update test report
	result := "not-satisfied"
	if validation.Result.Passing > 0 {
		result = "satisfied"
	}
	d.Result.Result = result
	d.Result.Pass = d.Test.ExpectedResult == result
	d.Result.Remarks = validation.Result.Observations

	return d.Result, nil
}

// LulaValidationTestResult is a struct that contains the details of the results of the test performed
// on the LulaValidation
type LulaValidationTestResult struct {
	TestName          string            `json:"test-name" yaml:"test-name"`
	Pass              bool              `json:"pass" yaml:"pass"`
	Result            string            `json:"result" yaml:"result"`
	Remarks           map[string]string `json:"remarks,omitempty" yaml:"remarks,omitempty"`
	TestResourcesPath string            `json:"test-resources-path,omitempty" yaml:"test-resources-path,omitempty"`
}

// LulaValidationTestReport contains the report of all the tests performed on a LulaValidation
type LulaValidationTestReport struct {
	Name        string                      `json:"name" yaml:"name"`
	TestResults []*LulaValidationTestResult `json:"test-results" yaml:"test-results"`
}

// NewLulaValidationTestReport creates a new report for a Lula Validation
func NewLulaValidationTestReport(name string) *LulaValidationTestReport {
	return &LulaValidationTestReport{
		Name:        name,
		TestResults: make([]*LulaValidationTestResult, 0),
	}
}

func (r *LulaValidationTestReport) AddTestResult(result *LulaValidationTestResult) {
	r.TestResults = append(r.TestResults, result)
}

func (r *LulaValidationTestReport) PrintReport() {
	if r == nil {
		message.HeaderInfof("No tests found")
		return
	}
	message.HeaderInfof("Test results")
	for _, testResult := range r.TestResults {
		if testResult.Pass {
			message.Successf("Pass: %s", testResult.TestName)
		} else {
			var failMsg string
			if testResult.Result == "" {
				failMsg = "No Result"
			} else {
				failMsg = "Expected Result =/= Actual Result"
			}
			message.Failf("Fail: %s - %s", testResult.TestName, failMsg)
		}
		if testResult.Result != "" {
			message.Infof("Result: %s", testResult.Result)
		}
		for remark, value := range testResult.Remarks {
			message.Infof("--> %s: %s", remark, value)
		}
		if testResult.TestResourcesPath != "" {
			message.Infof("Test Resources File Path: %s", testResult.TestResourcesPath)
		}
	}
}

func (r *LulaValidationTestReport) TestFailed() bool {
	for _, testResult := range r.TestResults {
		if !testResult.Pass {
			return true
		}
	}
	return false
}

func SummarizeTestReport(testReportMap map[string]LulaValidationTestReport) (string, bool) {
	var testsFailing int
	var testsPassing int
	var missingTestResults int

	for _, testReport := range testReportMap {
		if len(testReport.TestResults) == 0 {
			missingTestResults++
		} else if testReport.TestFailed() {
			testsFailing++
		} else {
			testsPassing++
		}
	}

	noTestsRun := len(testReportMap) == missingTestResults

	return fmt.Sprintf("Test Results: %d passing, %d failing, %d missing", testsPassing, testsFailing, missingTestResults), noTestsRun
}
