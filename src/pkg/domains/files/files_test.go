package files

import (
	"context"
	"testing"

	"github.com/defenseunicorns/lula/src/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

var _ types.Domain = (*Domain)(nil)

func TestGetResource(t *testing.T) {
	t.Run("local files", func(t *testing.T) {
		d := Domain{Spec: &Spec{Filepaths: []FileInfo{
			{Name: "foo.yaml", Path: "foo.yaml"},
			{Name: "bar.json", Path: "bar.json"},
			{Name: "arbitraryname", Path: "nested-directory/baz.hcl2"},
		}}}

		resources, err := d.GetResources(context.WithValue(context.Background(), types.LulaValidationWorkDir, "testdata"))
		require.NoError(t, err)
		if diff := cmp.Diff(resources, types.DomainResources{
			"bar.json": map[string]interface{}{"cat": "Cheetarah"},
			"foo.yaml": "cat = Li Shou",
			"arbitraryname": map[string]any{
				"resource": map[string]any{"catname": map[string]any{"blackcat": map[string]any{"name": "robin"}}},
			},
		}); diff != "" {
			t.Fatalf("wrong result:\n%s\n", diff)
		}
	})
}
