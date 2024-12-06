package oscal_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestResolveCatalogControls(t *testing.T) {
	t.Parallel()

	runTest := func(t *testing.T, catalogPath string, include, exclude, expectedControls []string) {
		validCatalogBytes := loadTestData(t, catalogPath)

		var validCatalog oscalTypes.OscalCompleteSchema
		if err := yaml.Unmarshal(validCatalogBytes, &validCatalog); err != nil {
			t.Fatalf("yaml.Unmarshal failed: %v", err)
		}

		require.NotNil(t, validCatalog.Catalog)

		controlMap, err := oscal.ResolveCatalogControls(validCatalog.Catalog, include, exclude)
		require.NoError(t, err)

		foundControls := make([]string, 0)
		for id := range controlMap {
			foundControls = append(foundControls, id)
		}
		require.ElementsMatch(t, expectedControls, foundControls)
	}

	tests := []struct {
		name             string
		catalogPath      string
		include          []string
		exclude          []string
		expectedControls []string
	}{
		{
			name:        "valid-catalog-include",
			catalogPath: "../../../test/unit/common/oscal/catalog.yaml",
			include:     []string{"ac-1", "ac-3", "ac-3.2", "ac-4", "ac-4.4"},
			exclude:     []string{},
			expectedControls: []string{
				"ac-1",
				"ac-3",
				"ac-3.2",
				"ac-4",
				"ac-4.4",
			},
		},
		{
			name:        "valid-catalog-include-all",
			catalogPath: "../../../test/unit/common/oscal/subdir/basic-catalog.yaml",
			include:     []string{},
			exclude:     []string{},
			expectedControls: []string{
				"s1.1.1",
				"s1.1.2",
				"s2.1.1",
				"s2.1.2",
			},
		},
		{
			name:        "valid-catalog-exclude",
			catalogPath: "../../../test/unit/common/oscal/subdir/basic-catalog.yaml",
			include:     []string{},
			exclude:     []string{"s1.1.2", "s2.1.2"},
			expectedControls: []string{
				"s1.1.1",
				"s2.1.1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.catalogPath, tt.include, tt.exclude, tt.expectedControls)
		})
	}
}
