package oscal_test

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

var (
	compdefValidMultiComponent           = "../../../test/unit/common/oscal/valid-multi-component.yaml"
	compdefValidMultiComponentPerControl = "../../../test/unit/common/oscal/valid-multi-component-per-control.yaml"
	source                               = "https://raw.githubusercontent.com/usnistgov/oscal-content/main/nist.gov/SP800-53/rev5/yaml/NIST_SP-800-53_rev5_HIGH-baseline-resolved-profile_catalog.yaml"
	validSSP                             = "../../../test/unit/common/oscal/valid-ssp.yaml"
)

func getComponentDefinition(t *testing.T, path string) *oscalTypes.ComponentDefinition {
	t.Helper()
	validComponentBytes := loadTestData(t, path)
	var validComponent oscalTypes.OscalCompleteSchema
	err := yaml.Unmarshal(validComponentBytes, &validComponent)
	require.NoError(t, err)
	return validComponent.ComponentDefinition
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

	t.Run("Simple generation of SSP", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponentPerControl)

		ssp, err := oscal.GenerateSystemSecurityPlan("lula generate ssp <flags>", source, validComponentDefn)
		require.NoError(t, err)

		validateSSP(t, ssp)

		// Check the control-implementation.implemented-requirements and system-implementation.components links
		systemComponentMap := createSystemComponentMap(t, ssp)
		require.NotNil(t, ssp.Model)
		for _, ir := range ssp.Model.ControlImplementation.ImplementedRequirements {
			// All controls should have 2 components linked
			assert.Len(t, *ir.ByComponents, 2)
			for _, byComponent := range *ir.ByComponents {
				// Check that the component exists in the system-implementation.components
				_, ok := systemComponentMap[byComponent.ComponentUuid]
				assert.True(t, ok)
			}
		}
	})

	t.Run("Generation of SSP with mis-matched catalog source", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponent)

		_, err := oscal.GenerateSystemSecurityPlan("", "foo", validComponentDefn)
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

func TestCreateImplementedRequirementsByFramework(t *testing.T) {
	t.Parallel()

	t.Run("Multiple control frameworks", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponent)

		implementedRequirementMap := oscal.CreateImplementedRequirementsByFramework(validComponentDefn)
		assert.Len(t, implementedRequirementMap, 4) // Should return 4 frameworks

		// Check source values
		implmentedReqtSource, ok := implementedRequirementMap[source]
		require.True(t, ok)
		assert.Len(t, implmentedReqtSource, 6) // source has 6 implemented requirements

		// Check only one component specifies ac-1
		ac_1, ok := implmentedReqtSource["ac-1"]
		require.True(t, ok)
		assert.Len(t, *ac_1.ByComponents, 1)
	})

	t.Run("Multiple Components per control", func(t *testing.T) {
		validComponentDefn := getComponentDefinition(t, compdefValidMultiComponentPerControl)

		implementedRequirementMap := oscal.CreateImplementedRequirementsByFramework(validComponentDefn)
		assert.Len(t, implementedRequirementMap, 1) // Should return 1 framework

		// Check source values
		implmentedReqtSource, ok := implementedRequirementMap[source]
		require.True(t, ok)
		assert.Len(t, implmentedReqtSource, 3) // source has 3 implemented requirements

		// Check 2 components specify ac-1
		ac_1, ok := implmentedReqtSource["ac-1"]
		require.True(t, ok)
		assert.Len(t, *ac_1.ByComponents, 2)
	})
}
