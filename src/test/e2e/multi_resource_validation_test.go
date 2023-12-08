package test

import (
	"context"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"testing"
	"time"
)

func TestMultiResourceValidation(t *testing.T) {
	featureTrueAPIValidation := features.New("Check Multi-Resource API Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the configmap
			configMap, err := util.GetConfigMap("./scenarios/multi-resource/configmap.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-configmap", configMap)

			// Create the pod
			pod, err := util.GetPod("./scenarios/multi-resource/pod.yaml")
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
			ctx = context.WithValue(ctx, "api-field-pod", pod)
			// Create additional Namespace
			nsvt1, err := util.GetNamespace("validation-test1")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, nsvt1); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "validation-test1-ns", nsvt1)

			// Create the pod
			podvt1, err := util.GetPod("./scenarios/multi-resource/podvt1.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, podvt1); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					PodConditionMatch(podvt1, corev1.PodReady, corev1.ConditionTrue),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "podvt1", podvt1)
			// Create additional Namespace
			nsvt2, err := util.GetNamespace("validation-test2")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, nsvt2); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "validation-test2-ns", nsvt2)

			// Create the pod
			podvt2, err := util.GetPod("./scenarios/multi-resource/podvt2.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, podvt2); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					PodConditionMatch(podvt2, corev1.PodReady, corev1.ConditionTrue),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "podvt2", podvt2)

			return ctx
		}).
		Assess("Validate Multi-Resource Collections", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/multi-resource/oscal-component.yaml"

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
			podvt2 := ctx.Value("podvt2").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, podvt2); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(podvt2),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			nsvt2 := ctx.Value("validation-test2-ns").(*corev1.Namespace)
			if err := config.Client().Resources().Delete(ctx, nsvt2); err != nil {
				t.Fatal(err)
			}

			podvt1 := ctx.Value("podvt1").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, podvt1); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(podvt1),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			nsvt1 := ctx.Value("validation-test1-ns").(*corev1.Namespace)
			if err := config.Client().Resources().Delete(ctx, nsvt1); err != nil {
				t.Fatal(err)
			}

			pod := ctx.Value("api-field-pod").(*corev1.Pod)
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

			configMap := ctx.Value("api-field-configmap").(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMap),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueAPIValidation)
}
