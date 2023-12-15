package test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestPodWaitValidation(t *testing.T) {
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

			return context.WithValue(ctx, "test-pod-label", pod)
		}).
		Assess("Validate pod label", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/wait-field/oscal-component.yaml"

			findingMap, _, err := validate.ValidateOnPath(oscalPath)
			if err != nil {
				t.Fatal(err)
			}

			for _, finding := range findingMap {
				state := finding.Target.Status.State
				if state != "satisfied" {
					t.Fatal("State should be satisfied, but got :", state)
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

	testEnv.Test(t, featureTrueValidation)
}
