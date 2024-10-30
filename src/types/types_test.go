package types_test

import (
	"context"
	"encoding/json"
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

func TestRunTests(t *testing.T) {
	t.Parallel()

	runTest := func(t *testing.T, opaSpec opa.OpaSpec, validation types.LulaValidation, expectedTestReport []types.TestReport) {
		opaProvider, err := opa.CreateOpaProvider(context.Background(), &opaSpec)
		require.NoError(t, err)

		validation.Provider = &opaProvider

		testReports, err := validation.RunTests(context.Background())
		require.NoError(t, err)

		require.Equal(t, expectedTestReport, *testReports)
	}

	tests := []struct {
		name       string
		opaSpec    opa.OpaSpec
		validation types.LulaValidation
		want       []types.TestReport
	}{
		{
			name: "valid tests",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}",
			},
			validation: types.LulaValidation{
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				Tests: &[]types.Test{
					{
						Name: "test-modify-name",
						Changes: []types.Change{
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
			want: []types.TestReport{
				{
					TestName: "test-modify-name",
					Pass:     true,
					Remarks:  map[string]string{},
				},
			},
		},
		{
			name: "multiple tests",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}",
			},
			validation: types.LulaValidation{
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				Tests: &[]types.Test{
					{
						Name: "test-modify-name",
						Changes: []types.Change{
							{
								Path:     "test.metadata.name",
								Type:     transform.ChangeTypeUpdate,
								Value:    "another-resource",
								ValueMap: nil,
							},
						},
						ExpectedResult: "not-satisfied",
					},
					{
						Name: "test-add-another-field",
						Changes: []types.Change{
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
			want: []types.TestReport{
				{
					TestName: "test-modify-name",
					Pass:     true,
					Remarks:  map[string]string{},
				},
				{
					TestName: "test-add-another-field",
					Pass:     true,
					Remarks:  map[string]string{},
				},
			},
		},
		{
			name: "valid tests with remarks",
			opaSpec: opa.OpaSpec{
				Rego: "package validate\n\nvalidate {input.test.metadata.name == \"test-resource\"}\n\nmsg = input.test.metadata.name",
				Output: &opa.OpaOutput{
					Observations: []string{"validate.msg"},
				},
			},
			validation: types.LulaValidation{
				DomainResources: &types.DomainResources{
					"test": map[string]interface{}{
						"metadata": map[string]interface{}{
							"name": "test-resource",
						},
					},
				},
				Tests: &[]types.Test{
					{
						Name: "test-modify-name",
						Changes: []types.Change{
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
			want: []types.TestReport{
				{
					TestName: "test-modify-name",
					Pass:     true,
					Remarks: map[string]string{
						"validate.msg": "another-resource",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.opaSpec, tt.validation, tt.want)
		})
	}
}
