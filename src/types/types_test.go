package types_test

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defenseunicorns/lula/src/internal/transform"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
)

func TestGetDomainResourcesAsJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		validation types.LulaValidation
		want       []byte
	}{
		{
			name: "valid validation",
			validation: types.LulaValidation{
				DomainResources: &types.DomainResources{
					"test-resource": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
			},
			want: []byte(`{"test-resource": {"metadata": {"name": "test-resource"}}}`),
		},
		{
			name: "nil validation",
			validation: types.LulaValidation{
				DomainResources: nil,
			},
			want: []byte(`{}`),
		},
		{
			name: "invalid validation",
			validation: types.LulaValidation{
				DomainResources: &types.DomainResources{
					"key": make(chan int),
				},
			},
			want: []byte(`{"Error":"Error marshalling to JSON"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validation.GetDomainResourcesAsJSON()
			var jsonWant map[string]interface{}
			err := json.Unmarshal(tt.want, &jsonWant)
			require.NoError(t, err)
			var jsonGot map[string]interface{}
			err = json.Unmarshal(got, &jsonGot)
			require.NoError(t, err)
			if !reflect.DeepEqual(jsonGot, jsonWant) {
				t.Errorf("GetDomainResourcesAsJSON() got = %v, want %v", jsonGot, jsonWant)
			}
		})
	}
}

// TestRunTests checks the execution of many tests on a single LulaValidation
func TestRunTests(t *testing.T) {
	t.Parallel()
	tmpDirName := "tmp-resources"

	runTest := func(t *testing.T, opaSpec opa.OpaSpec, validation types.LulaValidation, expectedTestReport *types.LulaValidationTestReport) {
		opaProvider, err := opa.CreateOpaProvider(context.Background(), &opaSpec)
		require.NoError(t, err)

		validation.Provider = &opaProvider

		testReport, err := validation.RunTests(context.Background(), false)
		require.NoError(t, err)

		require.Equal(t, expectedTestReport, testReport)
	}

	runTestWithPrint := func(t *testing.T, opaSpec opa.OpaSpec, validation types.LulaValidation, expectedTestReport *types.LulaValidationTestReport) {
		err := os.Mkdir(tmpDirName, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(tmpDirName)
		ctx := context.WithValue(context.Background(), types.LulaValidationWorkDir, tmpDirName)

		opaProvider, err := opa.CreateOpaProvider(context.Background(), &opaSpec)
		require.NoError(t, err)

		validation.Provider = &opaProvider

		testReport, err := validation.RunTests(ctx, true)
		require.NoError(t, err)

		require.Equal(t, expectedTestReport, testReport)
	}

	tests := []struct {
		name       string
		opaSpec    opa.OpaSpec
		validation types.LulaValidation
		want       *types.LulaValidationTestReport
		print      bool
	}{
		{
			name: "valid single test",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}",
			},
			validation: types.LulaValidation{
				Name: "test-validation",
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				ValidationTestData: []*types.LulaValidationTestData{
					{
						Test: &types.LulaValidationTest{
							Name: "test-modify-name",
							Changes: []types.LulaValidationTestChange{
								{
									Path:     "test.metadata.name",
									Type:     transform.ChangeTypeUpdate,
									Value:    "another-resource",
									ValueMap: nil,
								},
							},
							ExpectedResult: "not-satisfied",
						},
					},
				},
			},
			want: &types.LulaValidationTestReport{
				Name: "test-validation",
				TestResults: []*types.LulaValidationTestResult{
					{
						TestName: "test-modify-name",
						Result:   "not-satisfied",
						Pass:     true,
						Remarks:  map[string]string{},
					},
				},
			},
		},
		{
			name: "valid test with remarks",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}\n\nmsg = input.test.metadata.name",
				Output: &opa.OpaOutput{
					Observations: []string{"validate.msg"},
				},
			},
			validation: types.LulaValidation{
				Name: "test-validation",
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				ValidationTestData: []*types.LulaValidationTestData{
					{
						Test: &types.LulaValidationTest{
							Name: "test-modify-name",
							Changes: []types.LulaValidationTestChange{
								{
									Path:     "test.metadata.name",
									Type:     transform.ChangeTypeUpdate,
									Value:    "another-resource",
									ValueMap: nil,
								},
							},
							ExpectedResult: "not-satisfied",
						},
					},
				},
			},
			want: &types.LulaValidationTestReport{
				Name: "test-validation",
				TestResults: []*types.LulaValidationTestResult{
					{
						TestName: "test-modify-name",
						Result:   "not-satisfied",
						Pass:     true,
						Remarks: map[string]string{
							"validate.msg": "another-resource",
						},
					},
				},
			},
		},
		{
			name:  "valid test with printed resources",
			print: true,
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}",
			},
			validation: types.LulaValidation{
				Name: "test-validation",
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				ValidationTestData: []*types.LulaValidationTestData{
					{
						Test: &types.LulaValidationTest{
							Name: "test-modify-name",
							Changes: []types.LulaValidationTestChange{
								{
									Path:     "test.metadata.name",
									Type:     transform.ChangeTypeUpdate,
									Value:    "another-resource",
									ValueMap: nil,
								},
							},
							ExpectedResult: "not-satisfied",
						},
					},
				},
			},
			want: &types.LulaValidationTestReport{
				Name: "test-validation",
				TestResults: []*types.LulaValidationTestResult{
					{
						TestName:          "test-modify-name",
						Result:            "not-satisfied",
						Pass:              true,
						Remarks:           map[string]string{},
						TestResourcesPath: tmpDirName + "/test-modify-name.json",
					},
				},
			},
		},
		{
			name: "valid multiple tests",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}",
			},
			validation: types.LulaValidation{
				Name: "test-validation",
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				ValidationTestData: []*types.LulaValidationTestData{
					{
						Test: &types.LulaValidationTest{
							Name: "test-modify-name",
							Changes: []types.LulaValidationTestChange{
								{
									Path:     "test.metadata.name",
									Type:     transform.ChangeTypeUpdate,
									Value:    "another-resource",
									ValueMap: nil,
								},
							},
							ExpectedResult: "not-satisfied",
						},
					},
					{
						Test: &types.LulaValidationTest{
							Name: "test-add-another-field",
							Changes: []types.LulaValidationTestChange{
								{
									Path:     "test.metadata.anotherField",
									Type:     transform.ChangeTypeAdd,
									Value:    "new-resource",
									ValueMap: nil,
								},
							},
							ExpectedResult: "satisfied",
						},
					},
				},
			},
			want: &types.LulaValidationTestReport{
				Name: "test-validation",
				TestResults: []*types.LulaValidationTestResult{
					{
						TestName: "test-modify-name",
						Pass:     true,
						Result:   "not-satisfied",
						Remarks:  map[string]string{},
					},
					{
						TestName: "test-add-another-field",
						Pass:     true,
						Result:   "satisfied",
						Remarks:  map[string]string{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.print {
				runTestWithPrint(t, tt.opaSpec, tt.validation, tt.want)
			} else {
				runTest(t, tt.opaSpec, tt.validation, tt.want)
			}
		})
	}
}
