package oscal_test

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

var (
	compdefValidMultiComponent           = "../../../test/unit/common/oscal/valid-multi-component.yaml"
	compdefValidMultiComponentPerControl = "../../../test/unit/common/oscal/valid-multi-component-per-control.yaml"
	source                               = "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.yaml"
	validProfileLocalCatalog             = "../../../test/unit/common/oscal/valid-profile.yaml"
	validProfileRemoteRev4               = "../../../test/unit/common/oscal/valid-profile-remote-rev4.yaml"
	validProfileNoControls               = "../../../test/unit/common/oscal/valid-profile-test-excludes.yaml"
	validSSP                             = "../../../test/unit/common/oscal/valid-ssp.yaml"
	validSSPNoComponents                 = "../../../test/unit/common/oscal/valid-ssp-no-components.yaml"
	validGeneratedSSP                    = "../../../test/unit/common/oscal/valid-generated-ssp.yaml"
)

func getComponentDefinition(t *testing.T, path string) *oscalTypes.ComponentDefinition {
	t.Helper()
	validComponentBytes := loadTestData(t, path)
	var validComponent oscalTypes.OscalCompleteSchema
	err := yaml.Unmarshal(validComponentBytes, &validComponent)
	require.NoError(t, err)
	return validComponent.ComponentDefinition
}

func getProfile(t *testing.T, path string) *oscalTypes.Profile {
	t.Helper()
	validProfileBytes := loadTestData(t, path)
	var validProfile oscalTypes.OscalCompleteSchema
	err := yaml.Unmarshal(validProfileBytes, &validProfile)
	require.NoError(t, err)
	return validProfile.Profile
}

func validateSSP(t *testing.T, ssp *oscal.SystemSecurityPlan) {
	t.Helper()
	dir := t.TempDir()
	modelPath := filepath.Join(dir, "ssp.yaml")
	defer os.Remove(modelPath)

	err := oscal.WriteOscalModelNew(modelPath, ssp)
	require.NoError(t, err)
}

func createSystemComponentMap(t *testing.T, ssp *oscal.SystemSecurityPlan) map[string]oscalTypes.SystemComponent {
	systemComponentMap := make(map[string]oscalTypes.SystemComponent)
	require.NotNil(t, ssp.Model)
	for _, sc := range ssp.Model.SystemImplementation.Components {
		systemComponentMap[sc.UUID] = sc
	}
	return systemComponentMap
}

// Tests that the SSP was generated, checking the control-implmentation.implemented-requirements and system-implementation.components links
func TestGenerateSystemSecurityPlan(t *testing.T) {

	t.Run("Simple generation of SSP - no components", func(t *testing.T) {
		validProfile := getProfile(t, validProfileLocalCatalog)

		ssp, err := oscal.GenerateSystemSecurityPlan("lula generate ssp <flags>", validProfileLocalCatalog, []string{"statement"}, validProfile)
		require.NoError(t, err)
		require.NotNil(t, ssp.Model)

		validateSSP(t, ssp)

		// Check the control-implementation.implemented-requirements links
		expectedControls := []string{"ac-1", "ac-2", "ac-3"}
		foundControls := make([]string, 0)
		for _, ir := range ssp.Model.ControlImplementation.ImplementedRequirements {
			foundControls = append(foundControls, ir.ControlId)
		}

		// All controls should be in the expectedControls list
		assert.ElementsMatch(t, expectedControls, foundControls)
	})

	t.Run("Simple generation of SSP - with component defn", func(t *testing.T) {
		validProfile := getProfile(t, validProfileRemoteRev4)
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponent)

		ssp, err := oscal.GenerateSystemSecurityPlan("lula generate ssp <flags>", validProfileRemoteRev4, []string{"statement"}, validProfile, validComponentDefn)
		require.NoError(t, err)
		require.NotNil(t, ssp.Model)

		validateSSP(t, ssp)

		// Check the control-implementation.implemented-requirements and system-implementation.components links
		expectedControls := []string{"ac-1", "ac-2", "ac-3"}
		foundControls := make([]string, 0)
		systemComponentMap := createSystemComponentMap(t, ssp)
		for _, ir := range ssp.Model.ControlImplementation.ImplementedRequirements {
			foundControls = append(foundControls, ir.ControlId)
			// All controls should have 1 component linked
			require.NotNil(t, ir.ByComponents)
			assert.Len(t, *ir.ByComponents, 1)
			for _, byComponent := range *ir.ByComponents {
				// Check that the component exists in the system-implementation.components
				_, ok := systemComponentMap[byComponent.ComponentUuid]
				assert.True(t, ok)
			}
		}

		// Check that only one component is specified in system-implementation.components
		require.Equal(t, 1, len(ssp.Model.SystemImplementation.Components))

		// All controls should be in the expectedControls list
		assert.ElementsMatch(t, expectedControls, foundControls)
	})

	t.Run("Generation of SSP using a profile with no controls", func(t *testing.T) {
		validProfile := getProfile(t, validProfileNoControls)

		_, err := oscal.GenerateSystemSecurityPlan("lula generate ssp <flags>", validProfileNoControls, []string{"statement"}, validProfile)
		require.Error(t, err)
	})
}

func TestMakeSSPDeterministic(t *testing.T) {
	t.Run("Make example SSP deterministic", func(t *testing.T) {
		validSSPBytes := loadTestData(t, validSSP)

		var validSSP oscalTypes.OscalCompleteSchema
		err := yaml.Unmarshal(validSSPBytes, &validSSP)
		require.NoError(t, err)

		ssp := oscal.SystemSecurityPlan{
			Model: validSSP.SystemSecurityPlan,
		}

		err = ssp.MakeDeterministic()
		require.NoError(t, err)

		// Check that system-implementation.components is sorted (by title)
		firstThreeComponents := ssp.Model.SystemImplementation.Components[0:3]
		expectedFirstThreeComponentsUuids := []string{"4938767c-dd8b-4ea4-b74a-fafffd48ac99", "795533ab-9427-4abe-820f-0b571bacfe6d", "fa39eb84-3014-46b4-b6bc-7da10527c262"}
		for i, component := range firstThreeComponents {
			assert.Equal(t, expectedFirstThreeComponentsUuids[i], component.UUID)
		}
	})
}

func TestCreateSourceControlsMap(t *testing.T) {
	t.Parallel()

	t.Run("Multiple control frameworks", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponent)

		sourceControlsMap := oscal.CreateSourceControlsMap(validComponentDefn)
		assert.Len(t, sourceControlsMap, 4) // Should return 4 frameworks

		// Check source values
		controlMap, ok := sourceControlsMap[source]
		require.True(t, ok)
		assert.Len(t, controlMap, 6) // source has 6 implemented requirements

		// Check only one component specifies ac-1
		ac_1, ok := controlMap["ac-1"]
		require.True(t, ok)
		assert.Len(t, ac_1, 1)
	})

	t.Run("Multiple Components per control", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponentPerControl)

		sourceControlsMap := oscal.CreateSourceControlsMap(validComponentDefn)
		assert.Len(t, sourceControlsMap, 1) // Should return 1 framework

		// Check source values
		controlMap, ok := sourceControlsMap[source]
		require.True(t, ok)
		assert.Len(t, controlMap, 3) // source has 3 implemented requirements

		// Check 2 components specify ac-1
		ac_1, ok := controlMap["ac-1"]
		require.True(t, ok)
		assert.Len(t, ac_1, 2)
	})
}

func TestRemapSourceToUUID(t *testing.T) {
	sourceMap := map[string]string{
		"https://raw.githubusercontent.com/GSA/fedramp-automation/refs/tags/fedramp-2.0.0-oscal-1.0.4/dist/content/rev5/baselines/json/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json": "foo",
		"https://raw.githubusercontent.com/defenseunicorns/lula/refs/heads/main/src/test/unit/common/validation/validation.opa.yaml":                                                                 "not-oscal",
		"not-a-link": "bar",
	}

	outMap := oscal.RemapSourceToUUID(sourceMap)

	assert.Equal(t, 1, len(outMap))
	v, ok := outMap["bc413ad0-23d7-4ff0-a7af-b48a03294873"]
	assert.True(t, ok)
	assert.Equal(t, "foo", v)
}

func TestMergeSystemSecurityPlanModels(t *testing.T) {
	t.Run("Merge two SSPs", func(t *testing.T) {
		original := &oscalTypes.SystemSecurityPlan{
			UUID: "original-uuid",
			ImportProfile: oscalTypes.ImportProfile{
				Href: "profile1.yaml",
			},
			SystemImplementation: oscalTypes.SystemImplementation{
				Components: []oscalTypes.SystemComponent{
					{
						UUID: "comp1",
					},
				},
				Users: []oscalTypes.SystemUser{
					{
						UUID: "user1",
					},
				},
			},
			ControlImplementation: oscalTypes.ControlImplementation{
				ImplementedRequirements: []oscalTypes.ImplementedRequirement{
					{
						ControlId: "ac-1",
						ByComponents: &[]oscalTypes.ByComponent{
							{
								ComponentUuid: "comp1",
							},
						},
					},
					{
						ControlId: "ac-2",
					},
				},
			},
		}

		latest := &oscalTypes.SystemSecurityPlan{
			UUID: "latest-uuid",
			ImportProfile: oscalTypes.ImportProfile{
				Href: "profile1.yaml",
			},
			SystemImplementation: oscalTypes.SystemImplementation{
				Components: []oscalTypes.SystemComponent{
					{
						UUID: "comp1",
					},
					{
						UUID: "comp2",
					},
				},
			},
			ControlImplementation: oscalTypes.ControlImplementation{
				ImplementedRequirements: []oscalTypes.ImplementedRequirement{
					{
						ControlId: "ac-1",
						ByComponents: &[]oscalTypes.ByComponent{
							{
								ComponentUuid: "comp1",
							},
							{
								ComponentUuid: "comp2",
							},
						},
					},
					{
						ControlId: "ac-2",
						ByComponents: &[]oscalTypes.ByComponent{
							{
								ComponentUuid: "comp2",
							},
						},
					},
				},
			},
		}

		merged, err := oscal.MergeSystemSecurityPlanModels(original, latest)
		require.NoError(t, err)

		// Check that the merged model has the expected values
		assert.Equal(t, 2, len(merged.SystemImplementation.Components))
		assert.Equal(t, 2, len(merged.ControlImplementation.ImplementedRequirements))
		assert.Equal(t, 1, len(merged.SystemImplementation.Users)) // Users should not change from original

		for _, ir := range merged.ControlImplementation.ImplementedRequirements {
			require.NotNil(t, ir.ByComponents)
			if ir.ControlId == "ac-1" {
				assert.Equal(t, 2, len(*ir.ByComponents))
			} else if ir.ControlId == "ac-2" {
				assert.Equal(t, 1, len(*ir.ByComponents))
			}
		}

		assert.NotEqual(t, "original-uuid", merged.UUID)
	})

	t.Run("Merge two SSPs with different profiles", func(t *testing.T) {
		original := &oscalTypes.SystemSecurityPlan{
			UUID: "original-uuid",
			ImportProfile: oscalTypes.ImportProfile{
				Href: "profile1.yaml",
			},
		}
		latest := &oscalTypes.SystemSecurityPlan{
			UUID: "latest-uuid",
			ImportProfile: oscalTypes.ImportProfile{
				Href: "profile2.yaml",
			},
		}
		_, err := oscal.MergeSystemSecurityPlanModels(original, latest)
		require.Error(t, err)
	})
}

func TestHandleExistingSSP(t *testing.T) {
	validSSPBytes := loadTestData(t, validGeneratedSSP)

	var validSSP oscalTypes.OscalCompleteSchema
	err := yaml.Unmarshal(validSSPBytes, &validSSP)
	require.NoError(t, err)

	t.Run("Handle Existing with no existing data", func(t *testing.T) {
		ssp := oscal.SystemSecurityPlan{}
		ssp.NewModel(validSSPBytes)

		tmpDir := t.TempDir()
		tmpFilePath := filepath.Join(tmpDir, "ssp.yaml")

		err := ssp.HandleExisting(tmpFilePath)
		require.NoError(t, err)

		// Check length of components are the same
		require.Equal(t, len(validSSP.SystemSecurityPlan.SystemImplementation.Components), len(ssp.Model.SystemImplementation.Components))
	})

	t.Run("Handle Existing with existing data", func(t *testing.T) {
		ssp := oscal.SystemSecurityPlan{}
		ssp.NewModel(validSSPBytes)

		err := ssp.HandleExisting(validSSPNoComponents)
		require.NoError(t, err)

		// Check length of components is 2
		require.Equal(t, 2, len(ssp.Model.SystemImplementation.Components))
	})
}
