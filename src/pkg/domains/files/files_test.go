package files

import (
	"testing"

	"github.com/defenseunicorns/lula/src/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

var _ types.Domain = (*Domain)(nil)

func TestGetResource(t *testing.T) {
	t.Run("local files", func(t *testing.T) {
		d := Domain{Spec: &Spec{Filepaths: []string{
			"testdata/foo.yaml",
			"testdata/bar.json",
		}}}

		resources, err := d.GetResources()
		require.NoError(t, err)
		if diff := cmp.Diff(resources, types.DomainResources{"bar.json": map[string]interface{}{"cat": "Cheetarah"}, "foo.yaml": "cat = Li Shou"}); diff != "" {
			t.Fatalf("wrong result:\n%s\n", diff)
		}
	})

	//remote files
	// TODO
}
