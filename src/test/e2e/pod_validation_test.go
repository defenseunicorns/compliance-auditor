package test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/go-oscal/src/pkg/revision"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	"github.com/defenseunicorns/go-oscal/src/pkg/versioning"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestPodLabelValidation(t *testing.T) {
	featureTrueValidation := features.New("Check Pod Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/pod-label/pod.pass.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(config.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "test-pod-label", pod)
		}).
		Assess("Validate pod label", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component.yaml"
			return validatePodLabelPass(ctx, t, config, oscalPath)
		}).
		Assess("Validate pod label (Kyverno)", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component-kyverno.yaml"
			return validatePodLabelPass(ctx, t, config, oscalPath)
		}).
		Assess("Validate pod label (save-resources=backmatter)", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component.yaml"
			return validateSaveResources(ctx, t, oscalPath, "backmatter")
		}).
		Assess("Validate pod label (save-resources=remote)", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component.yaml"
			return validateSaveResources(ctx, t, oscalPath, "remote")
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-label").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}

			err := wait.For(conditions.New(config.Client().Resources()).ResourceDeleted(pod), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatal(err)
			}

			err = os.Remove("sar-test.yaml")
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	featureFalseValidation := features.New("Check Pod Validation - Failure").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/pod-label/pod.fail.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(config.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "test-pod-label", pod)
		}).
		Assess("Validate pod label", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component.yaml"
			return validatePodLabelFail(ctx, t, config, oscalPath)
		}).
		Assess("Validate pod label (Kyverno)", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component-kyverno.yaml"
			return validatePodLabelFail(ctx, t, config, oscalPath)
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-label").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err := wait.For(conditions.New(config.Client().Resources()).ResourceDeleted(pod), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	featureBadValidation := features.New("Check Graceful Failure - all not-satisfied without error").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/pod-label/pod.pass.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.For(conditions.New(config.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "test-pod-label", pod)
		}).
		Assess("All not-satisfied", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component-all-bad.yaml"
			return validatePodLabelFail(ctx, t, config, oscalPath)
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-label").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err := wait.For(conditions.New(config.Client().Resources()).ResourceDeleted(pod), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueValidation, featureFalseValidation, featureBadValidation)
}

func validatePodLabelPass(ctx context.Context, t *testing.T, config *envconf.Config, oscalPath string) context.Context {
	message.NoProgress = true

	tempDir := t.TempDir()

	// Upgrade the component definition to latest osscal version
	revisionOptions := revision.RevisionOptions{
		InputFile:  oscalPath,
		OutputFile: tempDir + "/oscal-component-upgraded.yaml",
		Version:    versioning.GetLatestSupportedVersion(),
	}
	revisionResponse, err := revision.RevisionCommand(&revisionOptions)
	if err != nil {
		t.Fatal("Failed to upgrade component definition with: ", err)
	}
	// Write the upgraded component definition to a temp file
	err = files.WriteOutput(revisionResponse.RevisedBytes, revisionOptions.OutputFile)
	if err != nil {
		t.Fatal("Failed to write upgraded component definition with: ", err)
	}
	message.Infof("Successfully upgraded %s to %s with OSCAL version %s %s\n", oscalPath, revisionOptions.OutputFile, revisionResponse.Reviser.GetSchemaVersion(), revisionResponse.Reviser.GetModelType())

	assessment, err := validate.ValidateOnPath(oscalPath, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(assessment.Results) == 0 {
		t.Fatal("Expected greater than zero results")
	}

	result := assessment.Results[0]

	if result.Findings == nil {
		t.Fatal("Expected findings to be not nil")
	}

	for _, finding := range *result.Findings {
		state := finding.Target.Status.State
		if state != "satisfied" {
			t.Fatal("State should be satisfied, but got :", state)
		}
	}

	// Test report generation
	report, err := oscal.GenerateAssessmentResults(assessment.Results, nil)
	if err != nil {
		t.Fatal("Failed generation of Assessment Results object with: ", err)
	}

	var model = oscalTypes_1_1_2.OscalModels{
		AssessmentResults: report,
	}

	// Write the assessment results to file
	err = oscal.WriteOscalModel("sar-test.yaml", &model)
	if err != nil {
		message.Fatalf(err, "error writing component to file")
	}

	initialResultCount := len(report.Results)

	//Perform the write operation again and read the file to ensure result was appended
	report, err = oscal.GenerateAssessmentResults(assessment.Results, nil)
	if err != nil {
		t.Fatal("Failed generation of Assessment Results object with: ", err)
	}

	// Get the UUID of the report results - there should only be one
	resultId := report.Results[0].UUID

	model = oscalTypes_1_1_2.OscalModels{
		AssessmentResults: report,
	}

	// Write the assessment results to file
	err = oscal.WriteOscalModel("sar-test.yaml", &model)
	if err != nil {
		message.Fatalf(err, "error writing component to file")
	}

	data, err := os.ReadFile("sar-test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	tempAssessment, err := oscal.NewAssessmentResults(data)
	if err != nil {
		t.Fatal(err)
	}

	// The number of results in the file should be more than initially
	if len(tempAssessment.Results) <= initialResultCount {
		t.Fatal("Failed to append results to existing report")
	}

	if resultId != tempAssessment.Results[0].UUID {
		t.Fatal("Failed to prepend results to existing report")
	}

	validatorResponse, err := validation.ValidationCommand("sar-test.yaml")
	if err != nil || validatorResponse.JsonSchemaError != nil {
		t.Fatal("File failed linting")
	}
	message.Infof("Successfully validated %s is valid OSCAL version %s %s\n", "sar-test.yaml", validatorResponse.Validator.GetSchemaVersion(), validatorResponse.Validator.GetModelType())

	return ctx
}

func validatePodLabelFail(ctx context.Context, t *testing.T, config *envconf.Config, oscalPath string) context.Context {
	message.NoProgress = true

	assessment, err := validate.ValidateOnPath(oscalPath, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(assessment.Results) == 0 {
		t.Fatal("Expected greater than zero results")
	}

	result := assessment.Results[0]

	if result.Findings == nil {
		t.Fatal("Expected findings to be not nil")
	}

	for _, finding := range *result.Findings {
		state := finding.Target.Status.State
		if state != "not-satisfied" {
			t.Fatal("State should be not-satisfied, but got :", state)
		}
	}

	return ctx
}

func validateSaveResources(ctx context.Context, t *testing.T, oscalPath, saveResources string) context.Context {
	message.NoProgress = true
	validate.SaveResources = saveResources
	tempDir := t.TempDir()
	validate.ResourcesDir = tempDir

	// Validate on path
	assessment, err := validate.ValidateOnPath(oscalPath, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(assessment.Results) == 0 {
		t.Fatal("Expected greater than zero results")
	}

	result := assessment.Results[0]

	// Should I call the dev commands here?
	if saveResources == "backmatter" {
		// Check that assessment results backmatter has the expected resources
		if assessment.BackMatter == nil {
			t.Fatal("Expected assessment backmatter, got nil")
		}
		if len(*assessment.BackMatter.Resources) != 2 {
			t.Fatal("Expected 2 resources, got ", len(*assessment.BackMatter.Resources))
		}
		// Check that the resources are the expected resources
		// helper function to convert oscalTypes_1_1_2.Resource to map[string]interface{}
		resourceStore := composition.NewResourceStoreFromBackMatter(assessment.BackMatter)

		for _, o := range *result.Observations {
			if o.Links == nil {
				t.Fatal("Expected observation links, got nil")
			}
			if len(*o.Links) != 1 {
				t.Fatal("Expected 1 link, got ", len(*o.Links))
			}
			link := (*o.Links)[0]
			resource, found := resourceStore.GetExisting(link.Href)
			if !found {
				t.Fatal("Expected resource to exist")
			}

			// Check that the resource has the expected data
			var data map[string]interface{}
			err := json.Unmarshal([]byte(resource.Description), &data)
			if err != nil {
				t.Fatal(err)
			}
			// Check that podvt exists - both should have this field, one is a struct and the other is an array
			if _, ok := data["podvt"]; !ok {
				t.Fatal("Expected podvt to exist")
			}
		}

	} else if saveResources == "remote" {
		// Check that remote files are created
		for _, o := range *result.Observations {
			if o.Links == nil {
				t.Fatal("Expected observation links, got nil")
			}
			if len(*o.Links) != 1 {
				t.Fatal("Expected 1 link, got ", len(*o.Links))
			}
			link := (*o.Links)[0]

			dataBytes, err := network.Fetch(link.Href)
			if err != nil {
				t.Fatal("Unable to fetch remote resource: ", err)
			}
			var data map[string]interface{}
			err = json.Unmarshal(dataBytes, &data)
			if err != nil {
				t.Fatal("Received invalid JSON: ", err)
			}
			// Check that podvt exists - both should have this field, one is a struct and the other is an array
			if _, ok := data["podvt"]; !ok {
				t.Fatal("Expected podvt to exist")
			}
		}
	}

	// Check that assessment results can be written to file
	var model = oscalTypes_1_1_2.OscalModels{
		AssessmentResults: assessment,
	}

	// Write the assessment results to file
	err = oscal.WriteOscalModel(filepath.Join(tempDir, "assessment-results.yaml"), &model)
	if err != nil {
		t.Fatal("error writing assessment results to file")
	}

	return ctx
}
