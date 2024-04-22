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

func TestValidationStore_AddFromLink(t *testing.T) {
	tests := []struct {
		name    string
		link    oscalTypes_1_1_2.Link
		numIds  int
		wantErr bool
	}{
		{
			name: "multi-remote-validations: wildcard",
			link: oscalTypes_1_1_2.Link{
				Href:             "file://../../../../test/e2e/scenarios/remote-validations/multi-validations.yaml",
				ResourceFragment: "*",
			},
			numIds:  2,
			wantErr: false,
		},
		{
			name: "one from multi-remote-validation",
			link: oscalTypes_1_1_2.Link{
				Href:             "file://../../../../test/e2e/scenarios/remote-validations/multi-validations.yaml",
				ResourceFragment: "9d09b4fc-1a82-4434-9fbe-392935347a84",
			},
			numIds:  1,
			wantErr: false,
		},
		{
			name: "single validation",
			link: oscalTypes_1_1_2.Link{
				Href: "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml",
			},
			numIds:  1,
			wantErr: false,
		},
		{
			name: "invalid link",
			link: oscalTypes_1_1_2.Link{
				Href: "file://../../../../test/e2e/scenarios/remote-validations/invalid-link.opa.yaml",
			},
			numIds:  0,
			wantErr: true,
		},
		{
			name: "invalid resource fragment",
			link: oscalTypes_1_1_2.Link{
				Href:             "file://../../../../test/e2e/scenarios/remote-validations/multi-validations.yaml",
				ResourceFragment: "invalid-resource-fragment",
			},
			numIds:  0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validationstore.NewValidationStore()
			gotIds, err := v.AddFromLink(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidationStore.AddFromLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotIds) != tt.numIds {
				t.Errorf("ValidationStore.AddFromLink() = %v ids, want %v ids", len(gotIds), tt.numIds)
			}
		})
	}
}
