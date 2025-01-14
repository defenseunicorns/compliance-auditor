package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/types"
)

func TestFileValidation(t *testing.T) {
	failDir := "./scenarios/file-validations/fail"
	passDir := "./scenarios/file-validations/pass"
	oscalFile := "/component-definition.yaml"
	kyvernoFile := "/oscal-component-kyverno.yaml"

	t.Run("success - opa", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validator, err := validation.New()
		require.NoError(t, err)

		assessment, err := validator.ValidateOnPath(ctx, passDir+oscalFile, "")
		require.NoError(t, err)

		if len(assessment.Model.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			if state != "satisfied" {
				t.Fatal("State should be satisfied, but got :", state)
			}
		}
	})
	t.Run("success - kyverno", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validator, err := validation.New()
		require.NoError(t, err)

		assessment, err := validator.ValidateOnPath(ctx, passDir+kyvernoFile, "")
		require.NoError(t, err)
		require.NotEmpty(t, assessment.Model.Results, "Expected greater than zero results")

		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			require.Equal(t, "satisfied", state, fmt.Sprintf("State should be satisfied, but got %s", state))
		}
	})
	t.Run("success - arbitrary file contents", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validator, err := validation.New()
		require.NoError(t, err)
		assessment, err := validator.ValidateOnPath(ctx, passDir+"/component-definition-string-file.yaml", "")
		require.NoError(t, err)
		require.NotEmpty(t, assessment.Model.Results, "Expected greater than zero results")

		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			require.Equal(t, "satisfied", state, fmt.Sprintf("State should be satisfied, but got %s", state))
		}
	})
	t.Run("fail - opa", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validator, err := validation.New()
		require.NoError(t, err)

		assessment, err := validator.ValidateOnPath(ctx, failDir+oscalFile, "")
		require.NoError(t, err)
		require.NotEmpty(t, assessment.Model.Results, "Expected greater than zero results")

		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			require.Equal(t, "not-satisfied", state, fmt.Sprintf("State should not be satisfied, but got %s", state))
		}
	})
	t.Run("fail - kyverno", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validator, err := validation.New()
		require.NoError(t, err)
		assessment, err := validator.ValidateOnPath(ctx, failDir+kyvernoFile, "")
		require.NoError(t, err)

		if len(assessment.Model.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			require.Equal(t, "not-satisfied", state, fmt.Sprintf("State should not be satisfied, but got %s", state))
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, "scenarios/file-validations/invalid")
		validator, err := validation.New()
		require.NoError(t, err)
		_, err = validator.ValidateOnPath(ctx, "scenarios/file-validations/invalid/oscal-component.yaml", "")
		require.Error(t, err)
	})

	// This test fixture is referencing a file on GIT, so if you're moving
	// things around here you should probably check that, too
	t.Run("remote file download", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validator, err := validation.New()
		require.NoError(t, err)
		assessment, err := validator.ValidateOnPath(ctx, "scenarios/file-validations/pass/component-definition-remote-files.yaml", "")
		require.NoError(t, err)
		require.Len(t, assessment.Model.Results, 1)
		result := assessment.Model.Results[0]
		require.NotNil(t, result, "Expected findings to be not nil")
		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			require.Equal(t, "satisfied", state, fmt.Sprintf("State should be satisfied, but got %s", state))
		}
	})
}
