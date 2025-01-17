package test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/pkg/message"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestCreateResourceDataValidation(t *testing.T) {
	featureTrueDataValidation := features.New("Check Create Resource Data Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the secure namespace for blocking admission of failPods
			secureNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "secure-ns",
					Labels: map[string]string{
						"pod-security.kubernetes.io/enforce": "restricted",
					},
				},
			}
			if err := config.Client().Resources().Create(ctx, secureNamespace); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess("Validate Create Resource Data Collections", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/create-resources/oscal-component.yaml"
			message.NoProgress = true

			validator, err := validation.New(
				validation.WithComposition(nil, oscalPath),
				validation.WithAllowExecution(true, true),
			)
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

			// Check that resources in the cluster were destroyed
			if err := config.Client().Resources().Get(ctx, "success-1", "validation-test", &corev1.Pod{}); err == nil {
				t.Fatal("pod success-1 should not exist")
			}
			if err := config.Client().Resources().Get(ctx, "success-2", "validation-test", &corev1.Pod{}); err == nil {
				t.Fatal("pod success-2 should not exist")
			}
			if err := config.Client().Resources().Get(ctx, "test-job", "another-ns", &batchv1.Job{}); err == nil {
				t.Fatal("job test-job should not exist")
			}
			if err := config.Client().Resources().Get(ctx, "test-pod-label", "validation-test", &corev1.Pod{}); err == nil {
				t.Fatal("pod test-pod-label should not exist")
			}
			if err := config.Client().Resources().Get(ctx, "validation-test", "", &corev1.Namespace{}); err != nil {
				t.Fatal("namespace validation-test should still exist")
			}
			if err := config.Client().Resources().Get(ctx, "secure-ns", "", &corev1.Namespace{}); err != nil {
				t.Fatal("namespace secure-ns should still exist")
			}
			if err := config.Client().Resources().Get(ctx, "another-ns", "", &corev1.Namespace{}); err == nil {
				t.Fatal("namespace another-ns should not exist")
			}

			return ctx
		}).
		Assess("Validate Create Resource With Wait and Read", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/create-resources/oscal-component-wait-read.yaml"
			message.NoProgress = true

			validator, err := validation.New(
				validation.WithComposition(nil, oscalPath),
				validation.WithAllowExecution(true, true),
			)
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

			// Check that resources in the cluster were destroyed
			podList := &corev1.PodList{}
			err = config.Client().Resources().WithNamespace("validation-test").List(ctx, podList)
			if len(podList.Items) != 0 || err != nil {
				t.Fatal("pods should not exist in validation-test namespace")
			}
			if err := config.Client().Resources().Get(ctx, "test-deployment", "validation-test", &appsv1.Deployment{}); err == nil {
				t.Fatal("deployment test-deployment should not exist")
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Delete the secure namespace
			secureNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "secure-ns",
				},
			}
			if err := config.Client().Resources().Delete(ctx, secureNamespace); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Feature()

	testEnv.Test(t, featureTrueDataValidation)
}

func TestDeniedCreateResources(t *testing.T) {
	featureDeniedCreateResource := features.New("Check Create Resource Denied - Success").
		Assess("Validate Create Resource Denied", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/create-resources/oscal-component-denied.yaml"
			message.NoProgress = true

			validator, err := validation.New(
				validation.WithComposition(nil, oscalPath),
				validation.WithAllowExecution(false, true),
			)
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
				if state != "not-satisfied" {
					t.Fatal("State should be not-satisfied, but got :", state)
				}
			}

			return ctx
		}).
		Feature()

	testEnv.Test(t, featureDeniedCreateResource)
}
