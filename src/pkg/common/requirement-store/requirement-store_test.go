package requirementstore_test

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
)

func TestNewRequirementStore(t *testing.T) {
	controlImplementations := []oscalTypes.ControlImplementationSet{}
	r := requirementstore.NewRequirementStore(&controlImplementations)
	if r == nil {
		t.Error("Expected a new RequirementStore, but got nil")
	}
}
