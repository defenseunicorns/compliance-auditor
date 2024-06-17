package requirementstore_test

import (
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
)

func TestNewRequirementStore(t *testing.T) {
	componentDefn := oscalTypes_1_1_2.ComponentDefinition{}
	r := requirementstore.NewRequirementStore(&componentDefn)
	if r == nil {
		t.Error("Expected a new RequirementStore, but got nil")
	}
}
