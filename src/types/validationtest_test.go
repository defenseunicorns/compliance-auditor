package types_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/internal/transform"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
)

// TestExecuteTest tests the execution of a single LulaValidationTest
func TestExecuteTest(t *testing.T) {
	t.Parallel()
	
	runTest := func(t *testing.T, opaSpec opa.OpaSpec, validation types.LulaValidation, resources map[string]interface{}, expectedTestResult *types.LulaValidationTestResult) {
		opaProvider, err := opa.CreateOpaProvider(context.Background(), &opaSpec)
		require.NoError(t, err)

		validation.Provider = &opaProvider

		testResult := validation.ExecuteTest(context.Background(), resources, false)

		require.Equal(t, expectedTestResult, testResult)
	}

	tests := []struct {
		name       string
		opaSpec    opa.OpaSpec
		validation types.LulaValidation
		resources  map[string]interface{}
		want       *types.LulaValidationTestResult
	}{
		{
			name: "valid tests",
