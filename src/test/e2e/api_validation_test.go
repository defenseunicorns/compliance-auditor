package test

import (
	"context"
	"testing"
	// "time"

	// "github.com/defenseunicorns/lula/src/cmd/validate"
	// "github.com/defenseunicorns/lula/src/test/util"
	// "github.com/defenseunicorns/lula/src/types"
	// corev1 "k8s.io/api/core/v1"
	// "sigs.k8s.io/e2e-framework/klient/wait"
	// "sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestApiValidation(t *testing.T) {
	featureTrueValidation := features.New("Check API Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).Feature()

	featureFalseValidation := features.New("Check API Validation - Failure").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueValidation, featureFalseValidation)
}
