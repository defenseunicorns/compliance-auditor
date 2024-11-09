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

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestGetResources(t *testing.T) {
	const (
		ckPodGetResources contextKey = "pod-get-resources"
		ckCfgGetResources contextKey = "config-get-resources"
	)
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
			ctx = context.WithValue(ctx, ckPodGetResources, pod)

			// Create the configmap
			configMap, err := util.GetConfigMap("./scenarios/dev-get-resources/configmap.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckCfgGetResources, configMap)

			return ctx
		}).
		Assess("Validate dev get-resources", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationFile := "./scenarios/dev-get-resources/validation.yaml"
			message.NoProgress = true
			dev.RunInteractively = false

			validationBytes, err := common.ReadFileToBytes(validationFile)
			if err != nil {
				t.Errorf("Error reading file: %v", err)
			}

			collection, err := dev.DevGetResources(ctx, validationBytes, false, nil)
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

			if len(collection["empty"].([]map[string]interface{})) != 0 {
				t.Fatalf("expected 0 length items in the empty payload - got %v\n", len(collection["empty"].([]map[string]interface{})))
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Delete the configmap
			configMap := ctx.Value(ckCfgGetResources).(*corev1.ConfigMap)
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
			pod := ctx.Value(ckPodGetResources).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueGetResources)
}
