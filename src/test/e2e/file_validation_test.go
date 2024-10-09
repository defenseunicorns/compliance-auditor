package test

import (
	"context"
	"testing"

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
		validationCtx, err := validation.New()
		if err != nil {
			t.Errorf("error creating validation context: %v", err)
		}
		assessment, err := validationCtx.ValidateOnPath(ctx, passDir+oscalFile, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		if result.Findings == nil {
			t.Fatal("Expected findings to be not nil")
		}

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			if state != "satisfied" {
				t.Fatal("State should be satisfied, but got :", state)
			}
		}
	})
	t.Run("success - kyverno", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, passDir)
		validationCtx, err := validation.New()
		if err != nil {
			t.Errorf("error creating validation context: %v", err)
		}
		assessment, err := validationCtx.ValidateOnPath(ctx, passDir+kyvernoFile, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		if result.Findings == nil {
			t.Fatal("Expected findings to be not nil")
		}

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			if state != "satisfied" {
				t.Fatal("State should be satisfied, but got :", state)
			}
		}
	})
	t.Run("fail - opa", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validationCtx, err := validation.New()
		if err != nil {
			t.Errorf("error creating validation context: %v", err)
		}
		assessment, err := validationCtx.ValidateOnPath(ctx, failDir+oscalFile, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		if result.Findings == nil {
			t.Fatal("Expected findings to be not nil")
		}

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			if state != "not-satisfied" {
				t.Fatal("State should be non-satisfied, but got :", state)
			}
		}
	})
	t.Run("fail - kyverno", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, failDir)
		validationCtx, err := validation.New()
		if err != nil {
			t.Errorf("error creating validation context: %v", err)
		}
		assessment, err := validationCtx.ValidateOnPath(ctx, failDir+kyvernoFile, "")
		if err != nil {
			t.Fatal(err)
		}

		if len(assessment.Results) == 0 {
			t.Fatal("Expected greater than zero results")
		}

		result := assessment.Results[0]
		if result.Findings == nil {
			t.Fatal("Expected findings to be not nil")
		}

		for _, finding := range *result.Findings {
			state := finding.Target.Status.State
			if state != "not-satisfied" {
				t.Fatal("State should be non-satisfied, but got :", state)
			}
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, "scenarios/file-validations/invalid")
		validationCtx, err := validation.New()
		if err != nil {
			t.Errorf("error creating validation context: %v", err)
		}
		_, err = validationCtx.ValidateOnPath(ctx, "scenarios/file-validations/invalid/oscal-component.yaml", "")
		if err == nil {
			t.Fatal("expected error, got success")
		}
	})
}
