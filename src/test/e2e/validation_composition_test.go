package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
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
			oscalPath := "./scenarios/validation-composition/component-definition.yaml"
			return validateComposition(ctx, t, oscalPath, "satisfied")
		}).
		Assess("Validate local composition file with bad links", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/validation-composition/component-definition-bad-href.yaml"
			return validateComposition(ctx, t, oscalPath, "not-satisfied")
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

func validateComposition(ctx context.Context, t *testing.T, oscalPath, expectedFindingState string) context.Context {
	compDefBytes, err := os.ReadFile(oscalPath)
	if err != nil {
		t.Error(err)
	}

	assessment, err := validate.ValidateOnPath(context.Background(), oscalPath, "")
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
		if state != expectedFindingState {
			t.Fatalf("State should be %s, but got: %s", expectedFindingState, state)
		}
	}

	var oscalModel oscalTypes_1_1_2.OscalCompleteSchema
	err = yaml.Unmarshal(compDefBytes, &oscalModel)
	if err != nil {
		t.Error(err)
	}

	compDef := oscalModel.ComponentDefinition

	compositionCtx, err := composition.New(composition.WithModelFromLocalPath(oscalPath))
	if err != nil {
		t.Errorf("error creating composition context: %v", err)
	}

	baseDir := filepath.Dir(oscalPath)
	err = compositionCtx.ComposeComponentValidations(ctx, compDef, baseDir)
	if err != nil {
		t.Error(err)
	}

	components := *compDef.Components

	// Create a validation store from the back-matter if it exists
	validationStore := validationstore.NewValidationStoreFromBackMatter(*compDef.BackMatter)

	findingMap, observations, err := validate.ValidateOnControlImplementations(ctx, components[0].ControlImplementations, validationStore, "")
	if err != nil {
		t.Fatalf("Error with validateOnControlImplementations: %v", err)
	}

	// For fun - create a result
	composeResult, err := oscal.CreateResult(findingMap, observations)
	if err != nil {
		t.Fatal(err)
	}

	// Compare results
	status, _, err := oscal.EvaluateResults(&result, &composeResult)
	if err != nil {
		t.Fatal(err)
	}

	if !status {
		t.Fatal("Expected Successful evaluate")
	}

	return ctx
}
