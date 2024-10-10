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
			{Name: "baz", Path: "baz", Parser: "json"},
			{Name: "arbitraryname", Path: "nested-directory/baz.hcl2"},
			{Name: "stringtheory", Path: "arbitrary.file", Parser: "string"},
		}}}

		resources, err := d.GetResources(context.WithValue(context.Background(), types.LulaValidationWorkDir, "testdata"))
		require.NoError(t, err)
		if diff := cmp.Diff(resources, types.DomainResources{
			"bar.json": map[string]interface{}{"cat": "Cheetarah"},
			"foo.yaml": "cat = Li Shou",
			"baz":      map[string]interface{}{"lizard": "Snakob"},
			"arbitraryname": map[string]any{
				"resource": map[string]any{"catname": map[string]any{"blackcat": map[string]any{"name": "robin"}}},
			},
			"stringtheory": "hello there!",
		}); diff != "" {
			t.Fatalf("wrong result:\n%s\n", diff)
		}
	})
}
