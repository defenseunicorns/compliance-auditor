package requirementstore_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"

	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
)

const (
	validCompDefMultiValidations = "../../../test/unit/common/oscal/valid-component-no-lula.yaml"
	controlImplementationSource  = "https://github.com/defenseunicorns/lula"
)

func TestNewRequirementStore(t *testing.T) {
	controlImplementations := []oscalTypes.ControlImplementationSet{}
	r := requirementstore.NewRequirementStore(&controlImplementations)
	if r == nil {
		t.Error("Expected a new RequirementStore, but got nil")
	}
}

func TestGenerateFindings(t *testing.T) {
	model := testhelpers.OscalFromPath(t, validCompDefMultiValidations)
	vs := validationstore.NewValidationStoreFromBackMatter(*model.ComponentDefinition.BackMatter)
	controlMap := oscal.FilterControlImplementations(model.ComponentDefinition)
	impls := controlMap[controlImplementationSource]
	rs := requirementstore.NewRequirementStore(&impls)

	findings := rs.GenerateFindings(vs)

	assert.Len(t, findings, 2)
	assert.Empty(t, findings["ID-1"].Remarks)
	assert.Equal(t, "not-satisfied", findings["ID-1"].Target.Status.State)
	assert.Equal(t, "fail", findings["ID-1"].Target.Status.Reason)
	assert.Empty(t, findings["ID-1"].Target.Status.Remarks)

	assert.Equal(t, "No Lula validations were defined for this control", findings["ID-2"].Remarks)
	assert.Equal(t, "not-satisfied", findings["ID-2"].Target.Status.State)
	assert.Equal(t, "other", findings["ID-2"].Target.Status.Reason)
	assert.Equal(t, "No Lula validations were defined for this control", findings["ID-2"].Target.Status.Remarks)
}
