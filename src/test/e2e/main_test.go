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
	// "sigs.k8s.io/e2e-framework/klient/k8s"
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
		envfuncs.CreateClusterWithConfig(kind.NewProvider(), kindClusterName, "kind-config.yaml"),
		envfuncs.CreateNamespace(namespace),
		func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
			ingressByte, err := os.ReadFile("nginx-ingress.yaml")
			if err != nil {
				log.Fatal(err)
			}
			ingressYAML := string(ingressByte)

			// var deployment appsv1.Deployment

			r, err := resources.New(cfg.Client().RESTConfig())
			if err != nil {
				return ctx, err
			}
			// decode and create a stream of YAML or JSON documents from an io.Reader
			decoder.DecodeEach(ctx, strings.NewReader(ingressYAML), decoder.CreateHandler(r))

			// Does the deployment "ingress-nginx-controller" exist?
			err = wait.For(conditions.New(cfg.Client().Resources()).DeploymentAvailable("ingress-nginx-controller", "ingress-nginx"), wait.WithTimeout(time.Minute*5))
			if err != nil {
				log.Fatal(err)
			}

			// // Get the k8s.object for the deployment to pass to the next function call
			// err = cfg.Client().Resources().Get(ctx, "ingress-nginx-controller", "ingress-nginx", deployment)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			deployment := appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "ingress-nginx-controller", Namespace: "ingress-nginx"}}

			// Wait until the deployment object is ready
			err = wait.For(conditions.New(cfg.Client().Resources()).DeploymentConditionMatch(&deployment, appsv1.DeploymentAvailable, corev1.ConditionTrue), wait.WithTimeout(time.Minute*5))
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
