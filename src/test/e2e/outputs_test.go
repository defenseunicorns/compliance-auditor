package test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestOutputs(t *testing.T) {
	featureTrueOutputs := features.New("Check Outputs").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod, err := util.GetPod("./scenarios/outputs/pod.yaml")
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

			// read the data and unmarshall
			data, err := os.ReadFile(oscalPath)
			if err != nil {
				t.Fatal(err)
			}

			compDef, err := oscal.NewOscalComponentDefinition(data)
			if err != nil {
				t.Fatal(err)
			}

			if compDef.Components == nil {
				t.Fatal("Expected non-nil components")
			}

			components := *compDef.Components
			validationStore := validationstore.NewValidationStoreFromBackMatter(*compDef.BackMatter)

			findingMap, observations, err := validate.ValidateOnControlImplementations(context.Background(), components[0].ControlImplementations, validationStore, "")
			if err != nil {
				t.Fatal(err)
			}

			// Check validation results are expected
			if findingMap["ID-1"].Target.Status.State != "satisfied" {
				t.Fatal("Failed to validate payload.output validation for ID-1")
			}
			if findingMap["ID-2"].Target.Status.State != "not-satisfied" {
				t.Fatal("Failed to validate payload.output validation for ID-2")
			}

			// Check that remarks are valid
			for _, o := range observations {
				if strings.Contains(o.Description, "ID-1") {
					// Check remarks have printed the correct observation value
					if o.RelevantEvidence != nil && (*o.RelevantEvidence)[0].Remarks != "validate.test: hello world\n" {
						t.Fatal("Failed to validate payload.output observations for ID-1")
					}
				}
				if strings.Contains(o.Description, "ID-2") {
					// Check remarks are empty due to incorrect observation formats
					if o.RelevantEvidence != nil && (*o.RelevantEvidence)[0].Remarks != "" {
						t.Fatal("Failed to validate payload.output observations for ID-2")
					}
				}
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
