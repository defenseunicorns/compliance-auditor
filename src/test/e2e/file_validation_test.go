package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		assert.NotNil(t, result, "Expected findings to be not nil")

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
		assert.NoError(t, err)
		assert.NotEmpty(t, assessment.Results, "Expected greater than zero results")

		result := assessment.Results[0]
		assert.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			assert.Equal(t, "satisfied", state, fmt.Sprintf("State should be satisfied, but got %s", state))
		}
	})
	t.Run("success - arbitrary file contexnts", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validator, err := validation.New()
		if err != nil {
			t.Errorf("error creating validator: %v", err)
		}
		assessment, err := validator.ValidateOnPath(ctx, passDir+"/component-definition-string-file.yaml", "")
		assert.NoError(t, err)
		assert.NotEmpty(t, assessment.Results, "Expected greater than zero results")

		result := assessment.Results[0]
		assert.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			assert.Equal(t, "satisfied", state, fmt.Sprintf("State should be satisfied, but got %s", state))
		}
	})
	t.Run("fail - opa", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validator, err := validation.New()
		require.NoError(t, err)

		assessment, err := validator.ValidateOnPath(ctx, failDir+oscalFile, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, assessment.Results, "Expected greater than zero results")

		result := assessment.Results[0]
		assert.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			assert.Equal(t, "not-satisfied", state, fmt.Sprintf("State should not be satisfied, but got %s", state))
		}
	})
	t.Run("fail - kyverno", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validator, err := validation.New()
		require.NoError(t, err)
		assessment, err := validator.ValidateOnPath(ctx, failDir+kyvernoFile, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		assert.NotNil(t, result, "Expected findings to be not nil")

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			assert.Equal(t, "not-satisfied", state, fmt.Sprintf("State should not be satisfied, but got %s", state))
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, "scenarios/file-validations/invalid")
		validator, err := validation.New()
		require.NoError(t, err)
		_, err = validator.ValidateOnPath(ctx, "scenarios/file-validations/invalid/oscal-component.yaml", "")
		if err == nil {
			t.Fatal("expected error, got success")
		}
	})
}
