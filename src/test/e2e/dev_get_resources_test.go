package test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestGetResources(t *testing.T) {
	featureTrueGetResources := features.New("Check dev get-resources").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the pod
			pod, err := util.GetPod("./scenarios/dev-get-resources/pod.yaml")
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
			ctx = context.WithValue(ctx, "pod-get-resources", pod)

			// Create the configmap
			configMap, err := util.GetConfigMap("./scenarios/dev-get-resources/configmap.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "configmap-get-resources", configMap)

			return ctx
		}).
		Assess("Validate dev get-resources", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationFile := "./scenarios/dev-get-resources/validation.yaml"
			message.NoProgress = true

			collection, err := dev.DevGetResources(ctx, validationFile)
			if err != nil {
				t.Fatalf("error testing dev get-resources: %v", err)
			}

			// Check that collection has the expected resources
			if collection["test-pod"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string) != "test-pod-name" {
				t.Fatal("The test-pod-name resource was not found in the collection")
			}

			var foundConfig bool
			for _, c := range collection["configs"].([]map[string]interface{}) {
				if c["metadata"].(map[string]interface{})["name"].(string) == "nginx-conf" {
					foundConfig = true
				}
			}
			if !foundConfig {
				t.Fatal("The nginx-conf resource was not found in the collection")
			}

			message.Infof("Successfully validated dev get-resources command")

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Delete the configmap
			configMap := ctx.Value("configmap-get-resources").(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMap),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			// Delete the pod
			pod := ctx.Value("pod-get-resources").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueGetResources)
}
