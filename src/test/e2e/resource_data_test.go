package test

import (
	"context"
	"testing"
	"time"

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

func TestResourceData(t *testing.T) {
	featureTrueOutputs := features.New("Check Outputs").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/resource-data/manifests.yaml")
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
			return context.WithValue(ctx, "test-pod-outputs", pod)
		}).
		Assess("Validate Outputs", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/outputs/oscal-component.yaml"
			message.NoProgress = true

			findingMap, observations, err := validate.ValidateOnPath(oscalPath)
			if err != nil {
				t.Fatal(err)
			}

			// Write report(s) to file to examine remarks
			report, err := oscal.GenerateAssessmentResults(findingMap, observations)
			if err != nil {
				t.Fatal("Failed generation of Assessment Results object with: ", err)
			}

			err = validate.WriteReport(report, "assessment-results-outputs.yaml")
			if err != nil {
				t.Fatal("Failed to write report to file: ", err)
			}

			message.Infof("Successfully validated payload.output structure")

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value("test-pod-outputs").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueOutputs)
}
