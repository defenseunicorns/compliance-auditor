package test

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestPodWaitValidation(t *testing.T) {
	const ckTestPodLabel contextKey = "test-pod-label"
	featureTrueValidation := features.New("Check Pod Wait for Ready - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/wait-field/pod.yaml")
			if err != nil {
				t.Fatal(err)
			}

			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}

			// We are purposefully not going to wait until the pod is ready and start Assess

			return context.WithValue(ctx, ckTestPodLabel, pod)
		}).
		Assess("Validate pod label", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/wait-field/oscal-component.yaml"
			message.NoProgress = true

			validator, err := validation.New()
			if err != nil {
				t.Errorf("error creating validation context: %v", err)
			}

			assessment, err := validator.ValidateOnPath(context.Background(), oscalPath, "")
			if err != nil {
				t.Fatal(err)
			}

			if len(assessment.Model.Results) == 0 {
				t.Fatal("Expected greater than zero results")
			}

			result := assessment.Model.Results[0]

			if result.Findings == nil {
				t.Fatal("Expected findings to be not nil")
			}

			for _, finding := range *result.Findings {
				state := finding.Target.Status.State
				if state != "satisfied" {
					t.Fatal("State should be satisfied, but got :", state)
				}
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value(ckTestPodLabel).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueValidation)
}
