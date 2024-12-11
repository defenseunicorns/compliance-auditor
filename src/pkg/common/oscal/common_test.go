package oscal_test

import (
	"reflect"
	"sort"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestUpdateProps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		propName      string
		propNamespace string
		propValue     string
		props         *[]oscalTypes.Property
		want          *[]oscalTypes.Property
	}{
		{
			name:          "Update existing property",
			propName:      "generation",
			propNamespace: oscal.LULA_NAMESPACE,
			propValue:     "lula gen component <updated-cmd>",
			props: &[]oscalTypes.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/ns",
					Value: "lula gen component <original-cmd>",
				},
			},
			want: &[]oscalTypes.Property{
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
			props: &[]oscalTypes.Property{
				{
					Name:  "generation",
					Ns:    "https://docs.lula.dev/ns",
					Value: "lula gen component <original-cmd>",
				},
			},
			want: &[]oscalTypes.Property{
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
			props: &[]oscalTypes.Property{
				{
					Name:  "target",
					Ns:    "https://some-other-ns.com",
					Value: "test",
				},
			},
			want: &[]oscalTypes.Property{
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
			props:         &[]oscalTypes.Property{},
			want: &[]oscalTypes.Property{
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
		props         *[]oscalTypes.Property
		want          bool
		wantValue     string
	}{
		{
			name:          "Get existing property",
			propName:      "target",
			propNamespace: oscal.LULA_NAMESPACE,
			props: &[]oscalTypes.Property{
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
			props: &[]oscalTypes.Property{
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
			props: &[]oscalTypes.Property{
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

func TestCompareControls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "Compare controls with XX-##.## format 1, a sorted before b",
			a:        "AC-1",
			b:        "AC-1.10",
			expected: true,
		},
		{
			name:     "Compare controls with XX-##.## format 2, a sorted before b",
			a:        "ac-1.10",
			b:        "ac-2",
			expected: true,
		},
		{
			name:     "Compare controls with only one XX-##.## format, a sorted before b",
			a:        "apple",
			b:        "AC-2",
			expected: true,
		},
		{
			name:     "Compare controls with XX-##.## format 1, b sorted before a",
			a:        "AC-1.10",
			b:        "AC-1.2",
			expected: false,
		},
		{
			name:     "Compare controls with XX-##.## format 2, b sorted before a",
			a:        "ac-1.10",
			b:        "ac-1.2",
			expected: false,
		},
		{
			name:     "Compare controls with one XX-##.## format, b sorted before a",
			a:        "AC-2",
			b:        "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := oscal.CompareControls(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("CompareControls() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortControls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		controls []oscalTypes.Control
		expected []oscalTypes.Control
	}{
		{
			name: "Sort controls with XX-##.## format 1",
			controls: []oscalTypes.Control{
				{
					Title: "ac-14",
				},
				{
					Title: "ac-4",
				},
				{
					Title: "ac-4.21",
				},
				{
					Title: "ac-4.4",
				},
			},
			expected: []oscalTypes.Control{
				{
					Title: "ac-4",
				},
				{
					Title: "ac-4.4",
				},
				{
					Title: "ac-4.21",
				},
				{
					Title: "ac-14",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Slice(tt.controls, func(i, j int) bool {
				return oscal.CompareControls(tt.controls[i].Title, tt.controls[j].Title)
			})
			if !reflect.DeepEqual(tt.controls, tt.expected) {
				t.Errorf("SortControls() got = %v, want %v", tt.controls, tt.expected)
			}
		})
	}
}

func FuzzCompareControls(f *testing.F) {
	f.Add("apple", "anotherword")
	f.Add("AC-1", "ac-1")
	f.Add("ac-4.4", "ac-4.21")

	f.Fuzz(func(t *testing.T, a string, b string) {
		oscal.CompareControls(a, b)
	})
}

func FuzzCompareControlsInt(f *testing.F) {
	f.Add("apple", "anotherword")
	f.Add("AC-1", "ac-1")
	f.Add("ac-4.4", "ac-4.21")

	f.Fuzz(func(t *testing.T, a string, b string) {
		oscal.CompareControlsInt(a, b)
	})
}
