package test

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestResourceDataValidation(t *testing.T) {
	const (
		ckCfgMapYaml contextKey = "configmap-yaml"
		ckCfgMapJSON contextKey = "configmap-json"
		ckSecret     contextKey = "secret"
		ckPod        contextKey = "pod"
	)
	featureTrueDataValidation := features.New("Check Resource Data Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the json configmap
			configMapJson, err := util.GetConfigMap("./scenarios/resource-data/configmap_json.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMapJson); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckCfgMapJSON, configMapJson)

			// Create the configmap with yaml data
			configMapYaml, err := util.GetConfigMap("./scenarios/resource-data/configmap_yaml.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMapYaml); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckCfgMapYaml, configMapYaml)

			// Create the secret
			secret, err := util.GetSecret("./scenarios/resource-data/secret.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, secret); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckSecret, secret)

			// Create the pod
			pod, err := util.GetPod("./scenarios/resource-data/pod.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckPod, pod)

			return ctx
		}).
		Assess("Validate Resource Data Collections", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/resource-data/oscal-component.yaml"
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
			// Delete the json configmap
			configMapJson := ctx.Value(ckCfgMapJSON).(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMapJson); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMapJson),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			// Delete the yaml configmap
			configMapYaml := ctx.Value(ckCfgMapYaml).(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMapYaml); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMapYaml),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			// Delete the secret
			secret := ctx.Value(ckSecret).(*corev1.Secret)
			if err := config.Client().Resources().Delete(ctx, secret); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(secret),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			// Delete the pod
			pod := ctx.Value(ckPod).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(pod),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			return ctx

		}).Feature()

	testEnv.Test(t, featureTrueDataValidation)
}
