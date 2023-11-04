package test

import (
	"os"
	"testing"
	"log"
	"context"
	"strings"
	"time"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/klient/wait"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testEnv         env.Environment
	kindClusterName string
	namespace       string
)

func TestMain(m *testing.M) {
	cfg, _ := envconf.NewFromFlags()
	testEnv = env.NewWithConfig(cfg)
	kindClusterName = envconf.RandomName("validation-test", 32)
	namespace = "validation-test"

	testEnv.Setup(
		envfuncs.CreateClusterWithConfig(
			kind.NewProvider(),
			kindClusterName,
			"kind-config.yaml"),

		envfuncs.CreateNamespace(namespace),

		func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
			// load stream of nginx-ingress resources
			ingressBytes, err := os.ReadFile("nginx-ingress.yaml")
			if err != nil {
				log.Fatal(err)
			}
			ingressYAML := string(ingressBytes)
			resource, err := resources.New(cfg.Client().RESTConfig())
			if err != nil {
				return ctx, err
			}
			decoder.DecodeEach(ctx, strings.NewReader(ingressYAML), decoder.CreateHandler(resource))

			// wait for ingress controller deployment object to be ready
			deployment := appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ingress-nginx-controller",
					Namespace: "ingress-nginx",
				},
			}
			err = wait.For(conditions.New(cfg.Client().Resources()).DeploymentConditionMatch(&deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				log.Fatal(err)
			}

			// find nginx ingress controller pod
			var pods corev1.PodList
			err = cfg.Client().Resources().WithNamespace("ingress-nginx").List(
				ctx, &pods, resources.WithLabelSelector(
					"app.kubernetes.io/component=controller," +
					"app.kubernetes.io/instance=ingress-nginx," +
					"app.kubernetes.io/name=ingress-nginx"))
			if err != nil {
				log.Fatal(err)
			}
			pod := &pods.Items[0]

			// wait for ingress controller to be ready
			err = wait.For(conditions.New(cfg.Client().Resources()).PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
			if err != nil {
				log.Fatal(err)
			}

			return ctx, nil
		},
	)

	testEnv.Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.DestroyCluster(kindClusterName),
	)

	os.Exit(testEnv.Run(m))
}
