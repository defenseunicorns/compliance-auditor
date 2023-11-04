package test

import (
	"context"
	"testing"
	"time"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/test/util"
	"github.com/defenseunicorns/lula/src/types"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestApiValidation(t *testing.T) {
	featureTrueValidation := features.New("Check API Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			configMap, err := util.GetConfigMap("./scenarios/api-field/configmap.pass.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-configmap", configMap)

			pod, err := util.GetPod("./scenarios/api-field/pod.yaml")
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

			service, err := util.GetService("./scenarios/api-field/service.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, service); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-service", service)

			ingress, err := util.GetIngress("./scenarios/api-field/ingress.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, ingress); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
				ResourceMatch(ingress, func(object k8s.Object) bool {
					ing, _ := object.(*netv1.Ingress)
					if len(ing.Status.LoadBalancer.Ingress) < 1 { return false }
					return ing.Status.LoadBalancer.Ingress[0].Hostname == "localhost"
				}),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-ingress", ingress)

			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := []string{"./scenarios/api-field/oscal-component.yaml"}

			results := types.ReportObject{
				FilePaths: oscalPath,
			}
			err := validate.ValidateOnPaths(&results)
			if err != nil {
				t.Fatal(err)
			}

			result := results.Components[0].ControlImplementations[0].ImplementedReqs[0].Results[0]

			if result.Failing != 0 {
				t.Fatal("Failing resources should be 0, but got :", result.Failing)
			}

			if result.Passing <= 0 {
				t.Fatal("Passing resources should be 1, but got :", result.Failing)
			}

			if result.State != "satisfied" {
				t.Fatal("State should be satisfied, but got :", result.State)
			}
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			ingress := ctx.Value("api-field-ingress").(*netv1.Ingress)
			if err := config.Client().Resources().Delete(ctx, ingress); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
				ResourceDeleted(ingress),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			service := ctx.Value("api-field-service").(*corev1.Service)
			if err := config.Client().Resources().Delete(ctx, service); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
				ResourceDeleted(service),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
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

	featureFalseValidation := features.New("Check API Validation - Failure").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			configMap, err := util.GetConfigMap("./scenarios/api-field/configmap.fail.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-configmap", configMap)

			pod, err := util.GetPod("./scenarios/api-field/pod.yaml")
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

			service, err := util.GetService("./scenarios/api-field/service.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, service); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-service", service)

			ingress, err := util.GetIngress("./scenarios/api-field/ingress.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, ingress); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
				ResourceMatch(ingress, func(object k8s.Object) bool {
					ing, _ := object.(*netv1.Ingress)
					if len(ing.Status.LoadBalancer.Ingress) < 1 { return false }
					return ing.Status.LoadBalancer.Ingress[0].Hostname == "localhost"
				}),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, "api-field-ingress", ingress)

			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := []string{"./scenarios/api-field/oscal-component.yaml"}

			results := types.ReportObject{
				FilePaths: oscalPath,
			}
			err := validate.ValidateOnPaths(&results)
			if err != nil {
				t.Fatal("Validation error, result:", results)
			}

			// TODO: maybe this brings to light modifying the
			result := results.Components[0].ControlImplementations[0].ImplementedReqs[0].Results[0]

			if result.Failing <= 0 {
				t.Fatal("Failing resources should be 1, but got :", result.Failing)
			}

			if result.State != "not-satisfied" {
				t.Fatal("State should be not-satisfied, but got :", result.State)
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			ingress := ctx.Value("api-field-ingress").(*netv1.Ingress)
			if err := config.Client().Resources().Delete(ctx, ingress); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
				ResourceDeleted(ingress),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			service := ctx.Value("api-field-service").(*corev1.Service)
			if err := config.Client().Resources().Delete(ctx, service); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
				ResourceDeleted(service),
				wait.WithTimeout(time.Minute*5))
			if err != nil {
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

	testEnv.Test(t, featureTrueValidation, featureFalseValidation )
}
