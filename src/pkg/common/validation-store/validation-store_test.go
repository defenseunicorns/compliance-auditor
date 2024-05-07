package validationstore_test

import (
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
)

const (
	validationPath = "../../../test/e2e/scenarios/remote-validations/validation.opa.yaml"
	componentPath  = "../../../test/e2e/scenarios/remote-validations/component-definition.yaml"
)

func generateValidation(t *testing.T) common.Validation {
	validationBytes, err := common.ReadFileToBytes(validationPath)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	var validation common.Validation
	err = validation.UnmarshalYaml(validationBytes)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	return validation
}

func TestNewValidationStore(t *testing.T) {
	v := validationstore.NewValidationStore()
	if v == nil {
		t.Error("Expected a new ValidationStore, but got nil")
	}
}

func TestNewValidationStoreFromBackMatter(t *testing.T) {
	backMatter := oscalTypes_1_1_2.BackMatter{}
	v := validationstore.NewValidationStoreFromBackMatter(backMatter)
	if v == nil {
		t.Error("Expected a new ValidationStore from back matter, but got nil")
	}
}

func TestAddValidation(t *testing.T) {
	validation := generateValidation(t)
	v := validationstore.NewValidationStore()

	id, err := v.AddValidation(&validation)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if id == "" {
		t.Error("Expected a non-empty ID, but got an empty string")
	}
}

func TestGetLulaValidation(t *testing.T) {
	validation := generateValidation(t)
	v := validationstore.NewValidationStore()
	id, _ := v.AddValidation(&validation)
	lulaValidation, err := v.GetLulaValidation(id)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if lulaValidation == nil {
		t.Error("Expected a LulaValidation, but got nil")
	}
}
