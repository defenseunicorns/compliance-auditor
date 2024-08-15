package oscal_test

import (
	"reflect"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestUpdateProps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		propName      string
		propNamespace string
		propValue     string
		props         *[]oscalTypes_1_1_2.Property
		want          *[]oscalTypes_1_1_2.Property
	}{
		{
			name:          "Update existing property",
			propName:      "generation",
			propNamespace: oscal.LULA_NAMESPACE,
			propValue:     "lula gen component <updated-cmd>",
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/ns",
					Value: "lula gen component <original-cmd>",
				},
			},
			want: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/oscal/ns",
					Value: "lula gen component <updated-cmd>",
				},
			},
		},
		{
			name:          "Add new property",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			propValue:     "test",
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/ns",
					Value: "lula gen component <original-cmd>",
				},
			},
			want: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/ns",
					Value: "lula gen component <original-cmd>",
				},
				{
					Name:  "target",
					Ns:    "https://docs.lula.dev/oscal/ns",
					Value: "test",
				},
			},
		},
		{
			name:          "Add new property in different namespace",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			propValue:     "test",
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    "https://some-other-ns.com",
					Value: "test",
				},
			},
			want: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    "https://some-other-ns.com",
					Value: "test",
				},
				{
					Name:  "target",
					Ns:    "https://docs.lula.dev/oscal/ns",
					Value: "test",
				},
			},
		},
		{
			name:          "Add new property to empty slice",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			propValue:     "test",
			props:         &[]oscalTypes_1_1_2.Property{},
			want: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    "https://docs.lula.dev/oscal/ns",
					Value: "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oscal.UpdateProps(tt.propName, tt.propNamespace, tt.propValue, tt.props)
			if !reflect.DeepEqual(*tt.props, *tt.want) {
				t.Errorf("UpdateProps() got = %v, want %v", *tt.props, *tt.want)
			}
		})
	}
}

func TestGetProps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		propName      string
		propNamespace string
		props         *[]oscalTypes_1_1_2.Property
		want          bool
		wantValue     string
	}{
		{
			name:          "Get existing property",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    oscal.LULA_NAMESPACE,
					Value: "test",
				},
			},
			want:      true,
			wantValue: "test",
		},
		{
			name:          "Get existing property with old namespace",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    "https://docs.lula.dev/ns",
					Value: "test",
				},
			},
			want:      true,
			wantValue: "test",
		},
		{
			name:          "Don't get property",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			props: &[]oscalTypes_1_1_2.Property{
				{
					Name:  "target",
					Ns:    "https://some-other-ns.com",
					Value: "test",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotValue := oscal.GetProp(tt.propName, tt.propNamespace, tt.props)
			if got != tt.want {
				t.Errorf("GetProp() got = %v, want %v", got, tt.want)
			}
			if gotValue != tt.wantValue {
				t.Errorf("GetProp() got = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
