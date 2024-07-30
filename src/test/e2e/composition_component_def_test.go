package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/test/util"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type compDefContextKey string

const componentDefinitionCompositionPodKey compDefContextKey = "component-definition-composition-pod"

func TestComponentDefinitionComposition(t *testing.T) {
	featureComponentDefinitionComposition := features.New("Check component definition composition").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the pod
			pod, err := util.GetPod("./scenarios/composition-component-definition/pod.pass.yaml")
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
			ctx = context.WithValue(ctx, componentDefinitionCompositionPodKey, pod)

			return ctx
		}).
		Assess("Validate local composition file", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			compDefPath := "../../test/unit/common/composition/component-definition-import-multi-compdef.yaml"
			compDefBytes, err := os.ReadFile(compDefPath)
			if err != nil {
				t.Error(err)
			}

			// Change Cwd to the directory of the component definition
			dirPath := filepath.Dir(compDefPath)
			resetCwd, err := common.SetCwdToFileDir(dirPath)
			if err != nil {
				t.Fatal(err)
			}

			compDef, err := oscal.NewOscalComponentDefinition(compDefBytes)
			if err != nil {
				t.Fatal(err)
			}

			results, err := validate.ValidateOnCompDef(compDef, "")
			if err != nil {
				t.Errorf("Error validating component definition: %v", err)
			}

			var expectedFindings, expectedObservations int
			expectedResults := len(results)

			for _, result := range results {
				expectedFindings += len(*result.Findings)
				expectedObservations += len(*result.Observations)
			}

			if expectedFindings == 0 {
				t.Errorf("Expected to find findings")
			}

			if expectedObservations == 0 {
				t.Errorf("Expected to find observations")
			}

			resetCwd()

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

			err = composition.ComposeComponentDefinitions(oscalModel.ComponentDefinition)
			if err != nil {
				t.Error(err)
			}

			composeResults, err := validate.ValidateOnCompDef(oscalModel.ComponentDefinition, "")
			if err != nil {
				t.Error(err)
			}

			if len(composeResults) != expectedResults {
				t.Errorf("Expected %v results, got %v", expectedResults, len(composeResults))
			}

			var composeFindings, composeObseravtions int
			for _, result := range composeResults {
				composeFindings += len(*result.Findings)
				composeObseravtions += len(*result.Observations)
			}

			if composeFindings != expectedFindings {
				t.Errorf("Expected %d findings, got %d", expectedFindings, composeFindings)
			}

			if composeObseravtions != expectedObservations {
				t.Errorf("Expected %d observations, got %d", expectedObservations, composeObseravtions)
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {

			// Delete the pod
			pod := ctx.Value(componentDefinitionCompositionPodKey).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureComponentDefinitionComposition)
}
