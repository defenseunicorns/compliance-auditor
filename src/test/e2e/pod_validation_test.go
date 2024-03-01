package test

import (
	"context"
	"os"
	"testing"
	"time"

	validator "github.com/defenseunicorns/go-oscal/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/cmd/validate"
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
			err = wait.For(conditions.New(config.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			return context.WithValue(ctx, "test-pod-label", pod)
		}).
		Assess("Validate pod label", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/pod-label/oscal-component.yaml"
			message.NoProgress = true

			findingMap, observations, err := validate.ValidateOnPath(oscalPath)
			if err != nil {
				t.Fatal(err)
			}

			for _, finding := range findingMap {
				state := finding.Target.Status.State
				if state != "satisfied" {
					t.Fatal("State should be satisfied, but got :", state)
				}
			}

			// Test report generation
			report, err := oscal.GenerateAssessmentResults(findingMap, observations)
			if err != nil {
				t.Fatal("Failed generation of Assessment Results object with: ", err)
			}

			// Write report(s) to file
			err = validate.WriteReport(report, "sar-test.yaml")
			if err != nil {
				t.Fatal("Failed to write report to file: ", err)
			}

			initialResultCount := len(report.Results)

			//Perform the write operation again and read the file to ensure result was appended
			report, err = oscal.GenerateAssessmentResults(findingMap, observations)
			if err != nil {
				t.Fatal("Failed generation of Assessment Results object with: ", err)
			}

			// Write report(s) to file
			err = validate.WriteReport(report, "sar-test.yaml")
			if err != nil {
				t.Fatal("Failed to write report to file: ", err)
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

			validator, err := validator.ValidateCommand("sar-test.yaml")
			if err != nil {
				t.Fatal("File failed linting")
			}
			message.Infof("Successfully validated %s is valid OSCAL version %s %s\n", "sar-test.yaml", validator.GetSchemaVersion(), validator.GetModelType())

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-label").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err := os.Remove("sar-test.yaml")
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
			message.NoProgress = true

			findingMap, _, err := validate.ValidateOnPath(oscalPath)
			if err != nil {
				t.Fatal(err)
			}

			for _, finding := range findingMap {
				state := finding.Target.Status.State
				if state != "not-satisfied" {
					t.Fatal("State should be not-satisfied, but got :", state)
				}
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-label").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueValidation, featureFalseValidation)
}
