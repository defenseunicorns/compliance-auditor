package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

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
			json.Unmarshal(tt.want, &jsonWant)
			var jsonGot map[string]interface{}
			json.Unmarshal(got, &jsonGot)
			if !reflect.DeepEqual(jsonGot, jsonWant) {
				t.Errorf("GetDomainResourcesAsJSON() got = %v, want %v", jsonGot, jsonWant)
			}
		})
	}
}
