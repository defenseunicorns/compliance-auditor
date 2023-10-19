package test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/test/util"
	"github.com/defenseunicorns/lula/src/types"
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
			oscalPath := []string{"./scenarios/pod-label/oscal-component.yaml"}

			results := types.ReportObject{
				FilePaths: oscalPath,
			}
			err := validate.ValidateOnPaths(&results)
			if err != nil {
				t.Fatal("Validation error, result:", results)
			}

			// TODO: maybe this brings to light modifying the
			result := results.Components[0].ControlImplementations[0].ImplementedReqs[0].Results[0]

			if result.Failing != 0 {
				t.Fatal("Failing resources should be 0, but got :", result.Failing)
			}

			if result.Passing <= 0 {
				t.Fatal("Passing resources should be 1, but got :", result.Failing)
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
			oscalPath := []string{"./scenarios/pod-label/oscal-component.yaml"}

			results := types.ReportObject{
				FilePaths: oscalPath,
			}
			err := validate.ValidateOnPaths(&results)
			if err != nil {
				t.Fatal("Validation error, result:", results)
			}

			// TODO: maybe this brings to light modifying the
			result := results.Components[0].ControlImplementations[0].ImplementedReqs[0].Results[0]

			if result.Failing <= 0 {
				t.Fatal("Failing resources should be 1, but got :", result.Failing)
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
