package test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestDevValidation(t *testing.T) {
	featureTrueDevValidate := features.New("Check dev validate").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the pod
			pod, err := util.GetPod("./scenarios/dev-validate/pod.pass.yaml")
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
			ctx = context.WithValue(ctx, "pod-dev-validate", pod)

			return ctx
		}).
		Assess("Validate DevValidate Opa", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationFile := "./scenarios/dev-validate/validation.yaml"

			message.NoProgress = true

			validationBytes, err := common.ReadFileToBytes(validationFile)
			if err != nil {
				t.Errorf("Error reading file: %v", err)
			}

			validation, err := dev.DevValidate(ctx, validationBytes)
			if err != nil {
				t.Errorf("Error testing dev validate: %v", err)
			}

			// Check the validation result has been evaluated
			if !validation.Evaluated {
				t.Errorf("Validation result has not been evaluated")
			}

			if validation.Result.Failing > 0 {
				t.Errorf("Validation failed")
			}

			message.Infof("Successfully validated dev validate command with OPA")

			return ctx
		}).
		Assess("Validate DevValidate kyverno", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			validationFile := "./scenarios/dev-validate/validation.kyverno.yaml"

			message.NoProgress = true

			validationBytes, err := common.ReadFileToBytes(validationFile)
			if err != nil {
				t.Errorf("Error reading file: %v", err)
			}

			validation, err := dev.DevValidate(ctx, validationBytes)
			if err != nil {
				t.Errorf("Error testing dev validate: %v", err)
			}

			// Check the validation result has been evaluated
			if !validation.Evaluated {
				t.Errorf("Validation result has not been evaluated")
			}

			if validation.Result.Failing > 0 {
				t.Errorf("Validation failed")
			}

			message.Infof("Successfully validated dev validate comman with kyverno")

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {

			// Delete the pod
			pod := ctx.Value("pod-dev-validate").(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueDevValidate)
}
