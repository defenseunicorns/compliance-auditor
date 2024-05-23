package validationstore

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

type ValidationStore struct {
	backMatterMap map[string]string
	validationMap types.LulaValidationMap
}

// NewValidationStore creates a new validation store
func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		backMatterMap: make(map[string]string),
		validationMap: make(types.LulaValidationMap),
	}
}

// NewValidationStoreFromBackMatter creates a new validation store from a back matter
func NewValidationStoreFromBackMatter(backMatter oscalTypes_1_1_2.BackMatter) *ValidationStore {
	return &ValidationStore{
		backMatterMap: oscal.BackMatterToMap(backMatter),
		validationMap: make(types.LulaValidationMap),
	}
}

// AddValidation adds a validation to the store
func (v *ValidationStore) AddValidation(validation *common.Validation) (id string, err error) {
	if validation.Metadata.UUID == "" {
		validation.Metadata.UUID = uuid.NewUUID()
	}

	v.validationMap[validation.Metadata.UUID], err = validation.ToLulaValidation()

	if err != nil {
		return "", err
	}

	return validation.Metadata.UUID, nil
}

// AddLulaValidation adds a LulaValidation to the store
func (v *ValidationStore) AddLulaValidation(validation *types.LulaValidation, id string) {
	trimmedId := common.TrimIdPrefix(id)
	v.validationMap[trimmedId] = *validation
}

// GetLulaValidation gets the LulaValidation from the store
func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	trimmedId := common.TrimIdPrefix(id)

	if validation, ok := v.validationMap[trimmedId]; ok {
		return &validation, nil
	}

	if validationString, ok := v.backMatterMap[trimmedId]; ok {
		lulaValidation, err := common.ValidationFromString(validationString)
		if err != nil {
			return &lulaValidation, err
		}
		v.validationMap[trimmedId] = lulaValidation
		return &lulaValidation, nil
	}

	return validation, fmt.Errorf("validation #%s not found", trimmedId)
}

// Number of validations in the store
func (v *ValidationStore) Count() int {
	return len(v.validationMap)
}

// DryRun checks if the validations are performing execution actions
func (v *ValidationStore) DryRun() (executable bool, msg string) {
	executableValidations := make([]string, 0)
	for k, validation := range v.validationMap {
		if validation.Domain.IsExecutable() {
			executableValidations = append(executableValidations, k)
		}
	}
	if len(executableValidations) > 0 {
		return true, fmt.Sprintf("The following validations are executable: %v", executableValidations)
	}
	return false, "No validation is executable"
}

// RunValidations runs the validations in the store
func (v *ValidationStore) RunValidations(confirmExecution bool) {
	for k, validation := range v.validationMap {
		err := validation.Validate()
		if err != nil {
			message.Debugf("Error running validation %s: %v", k, err)
			// Update validation with failed results
		}
	}

}
