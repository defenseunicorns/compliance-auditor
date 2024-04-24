package test

import (
	"context"
	"os"
	"testing"
	"time"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/test/util"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type contextKey string

const validationCompositionPodKey contextKey = "validation-composition-pod"

func TestValidationComposition(t *testing.T) {
	featureValidationComposition := features.New("Check validation composition").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the pod
			pod, err := util.GetPod("./scenarios/validation-composition/pod.pass.yaml")
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
			ctx = context.WithValue(ctx, validationCompositionPodKey, pod)

			return ctx
		}).
		Assess("Validate local composition file", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			compDefPath := "./scenarios/validation-composition/component-definition.yaml"
			compDefBytes, err := os.ReadFile(compDefPath)
			if err != nil {
				t.Error(err)
			}

			findings, observations, err := validate.ValidateOnPath(compDefPath)
			if err != nil {
				t.Errorf("Error validating component definition: %v", err)
			}
			expectedFindings := len(findings)
			expectedObservations := len(observations)

			if expectedFindings == 0 {
				t.Errorf("Expected to find findings")
			}

			if expectedObservations == 0 {
				t.Errorf("Expected to find observations")
			}

			var oscalModel oscalTypes_1_1_2.OscalCompleteSchema
			err = yaml.Unmarshal(compDefBytes, &oscalModel)
			if err != nil {
				t.Error(err)
			}
			reset, err := common.SetCwdToFileDir(compDefPath)
			if err != nil {
				t.Fatalf("Error setting cwd to file dir: %v", err)
			}
			defer reset()

			err = composition.ComposeComponentValidations(oscalModel.ComponentDefinition)
			if err != nil {
				t.Error(err)
			}

			findings, observations, err = validate.ValidateOnCompDef(*oscalModel.ComponentDefinition)
			if err != nil {
				t.Error(err)
			}

			if len(findings) != expectedFindings {
				t.Errorf("Expected %d findings, got %d", expectedFindings, len(findings))
			}

			if len(observations) != expectedObservations {
				t.Errorf("Expected %d observations, got %d", expectedObservations, len(observations))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {

			// Delete the pod
			pod := ctx.Value(validationCompositionPodKey).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureValidationComposition)
}
