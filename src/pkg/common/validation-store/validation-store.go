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
const WILD_CARD = "*"

type ValidationStore struct {
	backMatterMap map[string]string
	validationMap types.LulaValidationMap
	fileMap       map[string][]*types.LulaValidation
}

func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		backMatterMap: make(map[string]string),
		validationMap: make(types.LulaValidationMap),
		fileMap:       make(map[string][]*types.LulaValidation),
	}
}

func NewValidationStoreFromBackMatter(backMatter oscalTypes_1_1_2.BackMatter) *ValidationStore {
	return &ValidationStore{
		backMatterMap: oscal.BackMatterToMap(backMatter),
		validationMap: make(types.LulaValidationMap),
		fileMap:       make(map[string][]*types.LulaValidation),
	}
}

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

func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	trimmedId := trimIdPrefix(id)

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

func trimIdPrefix(id string) string {
	return strings.TrimPrefix(id, UUID_PREFIX)
}
