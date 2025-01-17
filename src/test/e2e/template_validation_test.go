package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

// Write validation template tests...
// To test
// 1. Templated comp-def
// check pass and fail?...

func TestTemplateValidation(t *testing.T) {
	const ckPodTmplValidation contextKey = "pod-template-validation"
	featureTemplateValidation := features.New("Check Template Validation").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the pod
			pod, err := util.GetPod("./scenarios/template-validation/pod.yaml")
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
			ctx = context.WithValue(ctx, ckPodTmplValidation, pod)

			return ctx
		}).
		Assess("Template to Pass", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/template-validation/component-definition.tmpl.yaml"
			// Set up the composition context
			composer, err := composition.New(
				composition.WithModelFromLocalPath(oscalPath),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", map[string]interface{}{
					"type":  interface{}("software"),
					"title": interface{}("lula"),
					"resources": interface{}(map[string]interface{}{
						"name":      interface{}("test-pod-label"),
						"namespace": interface{}("validation-test"),
					}),
				}, []template.VariableConfig{
					{
						Key:       "pod_label",
						Default:   "bar",
						Sensitive: false,
					},
					{
						Key:       "container_name",
						Default:   "nginx",
						Sensitive: true,
					},
				}, []string{}),
			)
			if err != nil {
				t.Errorf("error creating composition context: %v", err)
			}

			ctx = validateFindingsSatisfied(ctx, t, oscalPath, validation.WithComposition(composer, oscalPath))

			return ctx
		}).
		Assess("Template to Pass with env vars", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/template-validation/component-definition.tmpl.yaml"

			// Add env vars
			os.Setenv("LULA_VAR_POD_LABEL", "bar")
			os.Setenv("LULA_VAR_CONTAINER_NAME", "nginx")
			defer os.Unsetenv("LULA_VAR_POD_LABEL")
			defer os.Unsetenv("LULA_VAR_CONTAINER_NAME")

			// Set up the composition context - no default variable values now
			composer, err := composition.New(
				composition.WithModelFromLocalPath(oscalPath),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", map[string]interface{}{
					"type":  interface{}("software"),
					"title": interface{}("lula"),
					"resources": interface{}(map[string]interface{}{
						"name":      interface{}("test-pod-label"),
						"namespace": interface{}("validation-test"),
					}),
				}, []template.VariableConfig{
					{
						Key:       "pod_label",
						Sensitive: false,
					},
					{
						Key:       "container_name",
						Sensitive: true,
					},
				}, []string{}),
			)
			if err != nil {
				t.Errorf("error creating composition context: %v", err)
			}

			ctx = validateFindingsSatisfied(ctx, t, oscalPath, validation.WithComposition(composer, oscalPath))

			return ctx
		}).
		Assess("Template to Pass with overrides", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/template-validation/component-definition.tmpl.yaml"
			// Set up the composition context - bad resource name and no pod_label var - in overrides
			composer, err := composition.New(
				composition.WithModelFromLocalPath(oscalPath),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", map[string]interface{}{
					"type":  interface{}("software"),
					"title": interface{}("lula"),
					"resources": interface{}(map[string]interface{}{
						"name":      interface{}("bad-pod-name"),
						"namespace": interface{}("validation-test"),
					}),
				}, []template.VariableConfig{
					{
						Key:       "pod_label",
						Sensitive: false,
					},
					{
						Key:       "container_name",
						Default:   "nginx",
						Sensitive: true,
					},
				}, []string{".var.pod_label=bar", ".const.resources.name=test-pod-label"}),
			)
			if err != nil {
				t.Errorf("error creating composition context: %v", err)
			}

			ctx = validateFindingsSatisfied(ctx, t, oscalPath, validation.WithComposition(composer, oscalPath))

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value(ckPodTmplValidation).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTemplateValidation)
}

func validateFindingsSatisfied(ctx context.Context, t *testing.T, oscalPath string, opts ...validation.Option) context.Context {
	message.NoProgress = true

	validator, err := validation.New(opts...)
	if err != nil {
		t.Errorf("error creating validation context: %v", err)
	}

	assessment, err := validator.ValidateOnPath(context.Background(), oscalPath, "")
	if err != nil {
		t.Fatalf("Failed to validate oscal file: %s", oscalPath)
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
}
