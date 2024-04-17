package validationstore

import (
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
)

const UUID_PREFIX = "#"
const WILDCARD = "*"

type ValidationStore struct {
	backMatterMap map[string]string
	validationMap types.LulaValidationMap
	hrefIdMap     map[string][]string
}

// NewValidationStore creates a new validation store
func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		backMatterMap: make(map[string]string),
		validationMap: make(types.LulaValidationMap),
		hrefIdMap:     make(map[string][]string),
	}
}

// NewValidationStoreFromBackMatter creates a new validation store from a back matter
func NewValidationStoreFromBackMatter(backMatter oscalTypes_1_1_2.BackMatter) *ValidationStore {
	return &ValidationStore{
		backMatterMap: oscal.BackMatterToMap(backMatter),
		validationMap: make(types.LulaValidationMap),
		hrefIdMap:     make(map[string][]string),
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

// GetLulaValidation gets the LulaValidation from the store
func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	trimmedId := TrimIdPrefix(id)

	if validation, ok := v.validationMap[trimmedId]; ok {
		return &validation, nil
	}

	if validationString, ok := v.backMatterMap[trimmedId]; ok {
		lulaValidation, err := common.ValidationFromString(validationString)
		if err != nil {
			return nil, err
		}
		v.validationMap[trimmedId] = lulaValidation
		return &lulaValidation, nil
	}

	return validation, fmt.Errorf("validation #%s not found", trimmedId)
}

// SetHrefIds sets the validation ids for a given href
func (v *ValidationStore) SetHrefIds(href string, ids []string) {
	v.hrefIdMap[href] = ids
}

// GetHrefIds gets the validation ids for a given href
func (v *ValidationStore) GetHrefIds(href string) (ids []string, err error) {
	if ids, ok := v.hrefIdMap[href]; ok {
		return ids, nil
	}
	return nil, fmt.Errorf("href #%s not found", href)
}

// TrimIdPrefix trims the id prefix from the given id

func TrimIdPrefix(id string) string {
	return strings.TrimPrefix(id, UUID_PREFIX)
}
