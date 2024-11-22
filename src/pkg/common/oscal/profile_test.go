package oscal_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

func TestGetType(t *testing.T) {
	test := func(t *testing.T, model *oscal.Profile, expected string) {
		t.Helper()

		got := model.GetType()

		if got != expected {
			t.Fatalf("Expected %s - got %s\n", expected, got)
		}
	}

	t.Run("Test populated model", func(t *testing.T) {

		profile := oscal.NewProfile()
		profile.Model = &oscalTypes.Profile{}

		test(t, profile, "profile")
	})

	t.Run("Test unpopulated model", func(t *testing.T) {

		profile := oscal.NewProfile()

		test(t, profile, "profile")
	})
}

func TestGetCompleteModel(t *testing.T) {
	test := func(t *testing.T, model *oscal.Profile, expectedNil bool) {
		t.Helper()

		result := model.GetCompleteModel()

		if result.Profile == nil && !expectedNil {
			t.Fatal("Expected profile to not return nil")
		}

	}

	t.Run("Test complete with non-nil model", func(t *testing.T) {
		profile := oscal.NewProfile()
		profile.Model = &oscalTypes.Profile{}
		test(t, profile, false)
	})

	t.Run("Test complete with no model declaration", func(t *testing.T) {
		// Expecting a nil model
		profile := oscal.NewProfile()
		test(t, profile, true)
	})
}

func TestMakeDeterministic(t *testing.T) {
	test := func(t *testing.T, model *oscal.Profile, expectedImports []string, expectedControls []string, expectNil bool) {
		t.Helper()

		// Make deterministic
		_ = model.MakeDeterministic()

		if model.Model == nil && expectNil {
			return
		} else if model.Model != nil && expectNil {
			t.Fatal("Expected nil model but model was not nil")
		} else {
			profile := *model.Model

			for i, importItem := range profile.Imports {
				// Check href
				if importItem.Href != expectedImports[i] {
					t.Fatalf("Expected href %s - got %s\n", expectedImports[i], importItem.Href)
				}

				// Check included
				if importItem.IncludeControls != nil {
					includes := *importItem.IncludeControls
					for _, include := range includes {
						if include.WithIds != nil {
							for j, id := range *include.WithIds {
								if expectedControls[j] != id {
									t.Fatalf("Expected id %s - got %s\n", expectedControls[j], id)
								}
							}
						}
					}
				}

				// Check excluded
				if importItem.ExcludeControls != nil {
					excludes := *importItem.ExcludeControls
					for _, exclude := range excludes {
						if exclude.WithIds != nil {
							for j, id := range *exclude.WithIds {
								if expectedControls[j] != id {
									t.Fatalf("Expected id %s - got %s\n", expectedControls[j], id)
								}
							}
						}
					}
				}
			}
		}
	}

	t.Run("Profile with included controls", func(t *testing.T) {
		profile, err := oscal.GenerateProfile("", "#a3fb260d-0b89-4a12-b65c-a2737500febc", []string{"ac-4", "ac-1", "ac-7", "ac-3", "ac-2"}, []string{}, false)
		if err != nil {
			t.Fatal(err)
		}

		test(t, profile, []string{"#a3fb260d-0b89-4a12-b65c-a2737500febc"}, []string{"ac-1", "ac-2", "ac-3", "ac-4", "ac-7"}, false)
	})

	t.Run("Profile with exclude controls", func(t *testing.T) {
		profile, err := oscal.GenerateProfile("", "#a3fb260d-0b89-4a12-b65c-a2737500febc", []string{}, []string{"ac-4", "ac-1", "ac-7", "ac-3", "ac-2"}, false)
		if err != nil {
			t.Fatal(err)
		}

		test(t, profile, []string{"#a3fb260d-0b89-4a12-b65c-a2737500febc"}, []string{"ac-1", "ac-2", "ac-3", "ac-4", "ac-7"}, false)
	})

	t.Run("Profile with empty model", func(t *testing.T) {
		profile := oscal.NewProfile()

		test(t, profile, []string{}, []string{}, true)
	})
}

func TestResolveProfileControls(t *testing.T) {
	runTest := func(t *testing.T, profilePath string, include, exclude, expectedControls []string) {
		validProfileBytes := loadTestData(t, profilePath)

		var validProfile oscalTypes.OscalCompleteSchema
		if err := yaml.Unmarshal(validProfileBytes, &validProfile); err != nil {
			t.Fatalf("yaml.Unmarshal failed: %v", err)
		}

		require.NotNil(t, validProfile.Profile)

		controlMap, err := oscal.ResolveProfileControls(validProfile.Profile, profilePath, "", include, exclude)
		require.NoError(t, err)

		foundControls := make([]string, 0)
		for id := range controlMap {
			foundControls = append(foundControls, id)
		}
		require.ElementsMatch(t, expectedControls, foundControls)
	}

	tests := []struct {
		name             string
		profilePath      string
		include          []string
		exclude          []string
		expectedControls []string
	}{
		{
			name:        "valid-profile",
			profilePath: "../../../test/unit/common/oscal/valid-profile.yaml",
			include:     []string{},
			exclude:     []string{},
			expectedControls: []string{
				"ac-1",
				"ac-2",
				"ac-3",
			},
		},
		{
			name:        "valid-profile-with-exclude",
			profilePath: "../../../test/unit/common/oscal/valid-profile.yaml",
			include:     []string{},
			exclude:     []string{"ac-3"},
			expectedControls: []string{
				"ac-1",
				"ac-2",
			},
		},
		{
			name:        "valid-profile-multiple-imports",
			profilePath: "../../../test/unit/common/oscal/valid-profile-multiple-imports.yaml",
			include:     []string{},
			exclude:     []string{},
			expectedControls: []string{
				"ac-1",
				"ac-2",
				"ac-3",
				"s1.1.1",
				"s2.1.1",
			},
		},
		{
			name:             "valid-profile-exclude-all",
			profilePath:      "../../../test/unit/common/oscal/valid-profile-test-excludes.yaml",
			include:          []string{},
			exclude:          []string{},
			expectedControls: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.profilePath, tt.include, tt.exclude, tt.expectedControls)
		})
	}
}

func TestAddControl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		control  string
		include  []string
		exclude  []string
		expected bool
	}{
		{
			name:     "control-include-add-true",
			control:  "ac-1",
			include:  []string{"ac-1"},
			exclude:  []string{},
			expected: true,
		},
		{
			name:     "control-include-all-add-true",
			control:  "ac-1",
			include:  []string{},
			exclude:  []string{},
			expected: true,
		},
		{
			name:     "control-include-add-false",
			control:  "ac-2",
			include:  []string{"ac-1"},
			exclude:  []string{},
			expected: false,
		},
		{
			name:     "control-exclude-add-false",
			control:  "ac-1",
			include:  []string{},
			exclude:  []string{"ac-1"},
			expected: false,
		},
		{
			name:     "control-exclude-add-true",
			control:  "ac-2",
			include:  []string{},
			exclude:  []string{"ac-1"},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, oscal.AddControl(tt.control, tt.include, tt.exclude))
		})
	}
}
